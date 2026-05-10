CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    email TEXT UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
CREATE INDEX IF NOT EXISTS users_name_idx ON users(name);

CREATE TABLE IF NOT EXISTS user_points (
    user_id INTEGER UNIQUE REFERENCES users(id),
    points INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS user_points_user_id_idx ON user_points(user_id);
CREATE INDEX IF NOT EXISTS user_points_points_idx ON user_points(points);
