ALTER TABLE manga_tags
    DROP CONSTRAINT manga_tags_manga_id_fkey,
    ADD CONSTRAINT manga_tags_manga_id_fkey
        FOREIGN KEY (manga_id) REFERENCES manga(id),
    DROP CONSTRAINT manga_tags_tag_id_fkey,
    ADD CONSTRAINT manga_tags_tag_id_fkey
        FOREIGN KEY (tag_id) REFERENCES tag(id);