-- Filename new_migrations/000002_add_todos_indexes.up.sql

CREATE INDEX IF NOT EXISTS todo_title_idx ON todos USING GIN(to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS todo_description_idx ON todos USING GIN(to_tsvector('simple', description));