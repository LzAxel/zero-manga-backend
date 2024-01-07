ALTER TABLE manga
    ADD COLUMN preview_file_id UUID NOT NULL UNIQUE,
    DROP COLUMN preview_url;