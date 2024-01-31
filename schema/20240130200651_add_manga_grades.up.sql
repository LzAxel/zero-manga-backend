CREATE TABLE grade(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    manga_id UUID REFERENCES manga(id) NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    grade SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL
);