CREATE TABLE
    tags (
        id SERIAL PRIMARY KEY
    );

CREATE TABLE
    features (
        id SERIAL PRIMARY KEY
    );

CREATE TABLE
    banners (
        id SERIAL PRIMARY KEY,
        feature_id INT REFERENCES features (id) ON DELETE CASCADE,
        title TEXT NOT NULL,
        text TEXT NOT NULL,
        url TEXT NOT NULL,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP NOT NULL DEFAULT now(),
        updated_at TIMESTAMP NOT NULL DEFAULT now()
    );
CREATE TABLE
    bannertags (
        tag_id INT REFERENCES tags (id) ON DELETE CASCADE,
        banner_id INT REFERENCES banners (id) ON DELETE CASCADE,
        feature_id INT REFERENCES features (id) ON DELETE CASCADE,
        UNIQUE (tag_id, banner_id)
    );

CREATE TABLE
    users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(30) UNIQUE NOT NULL,
        token VARCHAR(36) UNIQUE NOT NULL
    );

CREATE TABLE
    admins (
        id SERIAL PRIMARY KEY,
        name VARCHAR(30) UNIQUE NOT NULL,
        token VARCHAR(36) UNIQUE NOT NULL
    );
