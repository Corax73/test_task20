CREATE TABLE IF NOT EXISTS "groups" (
    "id" serial PRIMARY KEY,
    "title" varchar NOT NULL UNIQUE
);