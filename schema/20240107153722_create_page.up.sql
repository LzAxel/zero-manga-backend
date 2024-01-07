CREATE TABLE page (
    id UUID PRIMARY KEY,
    chapter_id UUID NOT NULL REFERENCES chapter(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    number INT NOT NULL,
    height INT NOT NULL,
    width INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);