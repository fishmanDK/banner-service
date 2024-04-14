CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    hash_password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL,

    CONSTRAINT chk_role CHECK (role IN ('user', 'admin'))
);

CREATE TABLE banners (
    id SERIAL PRIMARY KEY,
    content JSONB NOT NULL,
    status BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE features (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE banner_tags (
    banner_id INT NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (banner_id, tag_id),
    FOREIGN KEY (banner_id) REFERENCES banners(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);

CREATE TABLE banner_features (
    banner_id INTEGER UNIQUE,
    feature_id INTEGER UNIQUE,
    PRIMARY KEY (banner_id, feature_id),
    FOREIGN KEY (banner_id) REFERENCES banners(id),
    FOREIGN KEY (feature_id) REFERENCES features(id)
);
