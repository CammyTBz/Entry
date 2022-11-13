-- Filename: migrations/000002_add_entry_check_constraint.down.sql

ALTER TABLE entries DROP CONSTRAINT IF EXISTS mode_length_check;