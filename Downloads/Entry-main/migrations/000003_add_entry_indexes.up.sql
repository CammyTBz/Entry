-- Filename: migrations/000003_add_entry_indexes.up.sql

CREATE INDEX IF NOT EXISTS entries_name_idx ON entries USING GIN(to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS entries_level_idx ON entries USING GIN(to_tsvector('simple', level));
CREATE INDEX IF NOT EXISTS entries_mode_idx ON entries USING GIN(mode);