CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS tenants (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO tenants (id, name)
VALUES ('00000000-0000-4000-8000-000000000001', 'Local development')
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS hosts (
  id uuid PRIMARY KEY,
  tenant_id uuid NOT NULL REFERENCES tenants(id),
  name text,
  agent_version text NOT NULL,
  last_seen_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS hosts_tenant_seen_idx ON hosts (tenant_id, last_seen_at DESC);

CREATE TABLE IF NOT EXISTS metric_batches (
  batch_id uuid PRIMARY KEY,
  host_id uuid NOT NULL REFERENCES hosts(id),
  sequence bigint NOT NULL,
  accepted_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS metric_samples (
  id bigserial PRIMARY KEY,
  tenant_id uuid NOT NULL REFERENCES tenants(id),
  host_id uuid NOT NULL REFERENCES hosts(id),
  batch_id uuid NOT NULL REFERENCES metric_batches(batch_id),
  name text NOT NULL CHECK (name IN ('temperature_c', 'cpu_percent', 'memory_percent')),
  value double precision NOT NULL,
  unit text NOT NULL CHECK (unit IN ('celsius', 'percent')),
  source text,
  observed_at timestamptz NOT NULL,
  received_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS samples_host_time_idx ON metric_samples (tenant_id, host_id, observed_at DESC);

CREATE TABLE IF NOT EXISTS incidents (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id uuid NOT NULL REFERENCES tenants(id),
  host_id uuid NOT NULL REFERENCES hosts(id),
  rule_key text NOT NULL,
  severity text NOT NULL CHECK (severity IN ('warning', 'critical')),
  status text NOT NULL CHECK (status IN ('open', 'in_progress', 'resolved')),
  opened_at timestamptz NOT NULL DEFAULT now(),
  resolved_at timestamptz
);

CREATE TABLE IF NOT EXISTS audit_events (
  id bigserial PRIMARY KEY,
  tenant_id uuid NOT NULL REFERENCES tenants(id),
  actor_id uuid,
  action text NOT NULL,
  subject_type text NOT NULL,
  subject_id text NOT NULL,
  details jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);
