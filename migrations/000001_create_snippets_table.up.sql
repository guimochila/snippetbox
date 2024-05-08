CREATE TABLE IF NOT EXISTS snippets (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  created DATE NOT NULL,
  expires DATE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_snippets_created ON snippets(created);
