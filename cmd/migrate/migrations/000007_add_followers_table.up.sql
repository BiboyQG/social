CREATE TABLE IF NOT EXISTS followers (
    user_id BIGSERIAL NOT NULL,
    follower_id BIGSERIAL NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, follower_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE
);
