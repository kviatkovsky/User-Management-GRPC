CREATE TABLE IF NOT EXISTS users
(
    id        uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email     TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL
);