import 'dotenv/config';
import express, { type NextFunction, type Request, type Response } from 'express';
import { createHash, timingSafeEqual } from 'node:crypto';
import { createClient } from 'redis';
import { z } from 'zod';
import { pool } from './store.js';

const app = express();
const port = Number(process.env.PORT ?? 8080);
const maxBatch = Number(process.env.MAX_BATCH_SIZE ?? 250);
const redis = createClient({ url: process.env.REDIS_URL });

app.disable('x-powered-by');
app.use(express.json({ limit: '256kb', strict: true }));

const metric = z.object({
  name: z.enum(['temperature_c', 'cpu_percent', 'memory_percent']),
  value: z.number().finite(),
  unit: z.enum(['celsius', 'percent']),
  source: z.string().min(1).max(100).optional(),
  observedAt: z.string().datetime({ offset: true })
});
const batchSchema = z.object({
  hostId: z.string().uuid(),
  agentVersion: z.string().min(1).max(40),
  sequence: z.number().int().nonnegative(),
  batchId: z.string().uuid(),
  metrics: z.array(metric).min(1).max(maxBatch)
});

function authorized(req: Request): boolean {
  const configured = process.env.AGENT_TOKEN;
  const supplied = req.header('authorization')?.replace(/^Bearer\s+/i, '');
  if (!configured || !supplied) return false;
  const expectedHash = createHash('sha256').update(configured).digest();
  const suppliedHash = createHash('sha256').update(supplied).digest();
  return timingSafeEqual(expectedHash, suppliedHash);
}

function requireAgent(req: Request, res: Response, next: NextFunction): void {
  if (!authorized(req)) {
    res.status(401).json({ error: 'unauthorized' });
    return;
  }
  next();
}

app.get('/health', (_req, res) => res.json({ status: 'ok' }));
app.get('/ready', async (_req, res) => {
  try {
    await pool.query('SELECT 1');
    if (!redis.isReady) throw new Error('Redis is not ready');
    res.json({ status: 'ready' });
  } catch {
    res.status(503).json({ status: 'not_ready' });
  }
});

app.post('/v1/telemetry', requireAgent, async (req, res) => {
  const parsed = batchSchema.safeParse(req.body);
  if (!parsed.success) {
    res.status(400).json({ error: 'invalid_payload', details: parsed.error.flatten() });
    return;
  }
  const batch = parsed.data;
  const client = await pool.connect();
  try {
    await client.query('BEGIN');
    const host = await client.query(
      `INSERT INTO hosts (id, tenant_id, agent_version, last_seen_at)
       VALUES ($1, '00000000-0000-4000-8000-000000000001', $2, now())
       ON CONFLICT (id) DO UPDATE SET agent_version = EXCLUDED.agent_version, last_seen_at = now()
       RETURNING id`, [batch.hostId, batch.agentVersion]
    );
    if (host.rowCount !== 1) throw new Error('Host upsert failed');
    const inserted = await client.query(
      `INSERT INTO metric_batches (batch_id, host_id, sequence)
       VALUES ($1, $2, $3) ON CONFLICT (batch_id) DO NOTHING RETURNING batch_id`,
      [batch.batchId, batch.hostId, batch.sequence]
    );
    if (inserted.rowCount) {
      for (const sample of batch.metrics) {
        await client.query(
          `INSERT INTO metric_samples (tenant_id, host_id, batch_id, name, value, unit, source, observed_at)
           VALUES ('00000000-0000-4000-8000-000000000001', $1, $2, $3, $4, $5, $6, $7)`,
          [batch.hostId, batch.batchId, sample.name, sample.value, sample.unit, sample.source ?? null, sample.observedAt]
        );
      }
    }
    await client.query('COMMIT');
    res.status(202).json({ accepted: true, duplicate: inserted.rowCount === 0 });
  } catch (error) {
    await client.query('ROLLBACK');
    console.error('telemetry_ingest_failed', error);
    res.status(500).json({ error: 'ingestion_failed' });
  } finally {
    client.release();
  }
});

app.use((err: Error, _req: Request, res: Response, _next: NextFunction) => {
  console.error('request_failed', err.message);
  res.status(400).json({ error: 'invalid_request' });
});

async function start(): Promise<void> {
  await redis.connect();
  app.listen(port, () => console.log(`Servitor API listening on ${port}`));
}

start().catch((error) => {
  console.error('startup_failed', error);
  process.exit(1);
});
