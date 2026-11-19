import 'dotenv/config';
import { readFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import { pool } from './store.js';

try {
  const path = fileURLToPath(new URL('./migrations/001_initial.sql', import.meta.url));
  const sql = await readFile(path, 'utf8');
  await pool.query(sql);
  console.log('Migration applied');
} finally {
  await pool.end();
}
