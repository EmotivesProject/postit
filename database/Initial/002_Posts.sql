\c postit_db;

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    username VARCHAR(128) NOT NULL,
    message TEXT default null,
    image_path TEXT default null,
    longitude decimal default null,
    latitude decimal default null,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    active BOOLEAN default false,

    CONSTRAINT fk_username_on_posts
    FOREIGN KEY(username)
    REFERENCES users(username)
);
