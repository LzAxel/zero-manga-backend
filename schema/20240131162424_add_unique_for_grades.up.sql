ALTER TABLE grade
    ADD CONSTRAINT grade_id_manga_id_unique UNIQUE (manga_id, user_id);