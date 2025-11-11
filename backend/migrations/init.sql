CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE columns (
                         id SERIAL PRIMARY KEY,
                         title VARCHAR(100) NOT NULL
);

CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL,
                       description TEXT,
                       column_id INT REFERENCES columns(id)
);
