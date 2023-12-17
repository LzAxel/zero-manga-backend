CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(30) UNIQUE NOT NULL,
    display_name VARCHAR(30),
    bio VARCHAR(300),
    email VARCHAR(40) UNIQUE NOT NULL,
    gender SMALLINT NOT NULL,
    type SMALLINT NOT NULL,
    password_hash bytea NOT NULL,
    online_at TIMESTAMP NOT NULL,
    registered_at TIMESTAMP NOT NULL
);

CREATE TABLE manga (
    id UUID PRIMARY KEY,
    title VARCHAR(100) UNIQUE NOT NULL,
    secondary_title VARCHAR(100),
    description VARCHAR(300),
    slug VARCHAR(100) UNIQUE NOT NULL,
    type SMALLINT NOT NULL,
    status SMALLINT NOT NULL,
    age_restrict SMALLINT NOT NULL,
    release_year INT NOT NULL
);

CREATE TABLE chapter (
    id UUID PRIMARY KEY,
    manga_id UUID NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
    title VARCHAR(100),
    number INT NOT NULL,
    volume INT NOT NULL,
    page_count INT NOT NULL,
    file_path VARCHAR(100) NOT NULL,
    uploader_id UUID NOT NULL REFERENCES users(id),
    uploaded_at TIMESTAMP NOT NULL
);