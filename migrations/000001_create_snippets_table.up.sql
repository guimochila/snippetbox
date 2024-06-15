CREATE TABLE IF NOT EXISTS snippets (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created DATE NOT NULL,
  expires DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_snippets_created ON snippets(created);

GRANT USAGE ON SCHEMA public TO web;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.snippets TO web;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE public.idx_snippets_created TO web;
