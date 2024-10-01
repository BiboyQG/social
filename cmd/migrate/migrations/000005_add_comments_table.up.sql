CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGSERIAL NOT NULL,
    user_id BIGSERIAL NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL default now(),
    updated_at TIMESTAMP(0) with time zone NOT NULL default now()
);