CREATE TABLE users (
    id  SERIAL PRIMARY KEY,
    name TEXT NOT UNIQUE,
    password TEXT NOT NULL,
    date DATE NOT NULL
);
