-- Filename: migrations/000003_add_entry_indexes.down.sql

DROP INDEX If EXISTS entries_name_idx;
DROP INDEX If EXISTS entries_level_idx;
DROP INDEX If EXISTS entries_mode_idx;