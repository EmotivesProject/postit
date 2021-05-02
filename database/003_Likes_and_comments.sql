\c postit_db;

CREATE TABLE likes (
    id SERIAL PRIMARY KEY,
	post_id INT NOT NULL,
    username VARCHAR(128) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    active BOOLEAN default false,

    CONSTRAINT fk_username_on_like
    FOREIGN KEY(username)
    REFERENCES users(username),

    CONSTRAINT fk_post_id_on_like
    FOREIGN KEY(post_id)
    REFERENCES posts(id),

    CONSTRAINT user_and_like_contraint UNIQUE (post_id, username)
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
	post_id INT NOT NULL,
    username VARCHAR(128) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    active BOOLEAN default false,

    CONSTRAINT fk_username_on_comments
    FOREIGN KEY(username)
    REFERENCES users(username),

    CONSTRAINT fk_post_id_on_comments
    FOREIGN KEY(post_id)
    REFERENCES posts(id)
);
