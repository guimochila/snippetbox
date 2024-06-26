CREATE TABLE IF NOT EXISTS sessions (
  token CHAR(43) PRIMARY KEY,
  data BYTEA NOT NULL,
  expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.sessions TO web;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.sessions_expiry_idx TO web;
