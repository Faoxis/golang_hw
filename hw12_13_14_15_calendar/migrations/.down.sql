-- 001_create_events_table.down.sql
DROP INDEX IF EXISTS idx_events_start_time;
DROP TABLE IF EXISTS events;
