-- Filename new_migrations/000001_add_todo_table.up.sql

CREATE TABLE IF NOT EXISTS todos (
    id bigserial primary key,
    create_at timestamp(0) without time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL,
    description text NOT NULL,
    complete boolean NOT NULL DEFAULT FALSE,
    version integer NOT NULL DEFAULT 1
)