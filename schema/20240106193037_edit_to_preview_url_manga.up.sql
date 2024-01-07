ALTER TABLE manga
    DROP COLUMN preview_file_id,
    ADD COLUMN preview_url TEXT NOT NULL;