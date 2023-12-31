CREATE TABLE IF NOT EXISTS users
(
    user_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_email ON users(email);

CREATE TABLE IF NOT EXISTS sources (
    source_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    feed_url VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS articles (
    article_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT,
    source_name VARCHAR(255),
    title VARCHAR(255) NOT NULL,
    link VARCHAR(255) NOT NULL UNIQUE,
    excerpt TEXT NOT NULL,
    image TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    posted_at TIMESTAMP
);

ALTER TABLE articles ADD CONSTRAINT fk_articles_user_id FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;
INSERT INTO users (user_name, email, password_hash) VALUES ('Bot', 'ihwuqih928u398dhfdfo;iweh', '2389uefhhr4e8934we56w3378sdeuibrw3345df');
INSERT INTO sources (name, feed_url) VALUES ('Habr', 'https://habr.com/ru/rss/hubs/go/articles/?fl=ru?with_tags=true:');