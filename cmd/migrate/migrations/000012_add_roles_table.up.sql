CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description VARCHAR(255) NOT NULL
);

INSERT INTO roles (name, level, description) VALUES ('user', 1, 'User role, which can create, read, update and delete own data (posts, comments, etc.)');
INSERT INTO roles (name, level, description) VALUES ('moderator', 2, 'Moderator role, which can create, read, update and delete own data (posts, comments, etc.) and update users data (ban, delete)');
INSERT INTO roles (name, level, description) VALUES ('admin', 3, 'Admin role, which can create, read, update and delete any data (posts, comments, etc.)');
