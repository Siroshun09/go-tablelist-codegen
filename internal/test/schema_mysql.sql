DROP TABLE IF EXISTS user_points;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users
(
    id         INT PRIMARY KEY AUTO_INCREMENT,
    email      VARCHAR(255) UNIQUE,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL
);
CREATE INDEX users_email_idx ON users (email);
CREATE INDEX users_name_idx ON users (name);

CREATE TABLE IF NOT EXISTS user_points
(
    user_id    INT UNIQUE,
    points     INT NOT NULL,
    created_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX user_points_user_id_idx ON user_points (user_id);
CREATE INDEX user_points_points_idx ON user_points (points);