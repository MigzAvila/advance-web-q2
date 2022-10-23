-- Filename new_migrations/000002_add_todos_indexes.down.sql
DROP INDEX IF EXISTS todo_title_idx;
DROP INDEX IF EXISTS todo_description_idx;
