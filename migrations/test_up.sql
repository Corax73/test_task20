CREATE TABLE IF NOT EXISTS "groups" (
    "id" serial PRIMARY KEY,
    "title" varchar NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS "songs" (
    "id" serial PRIMARY KEY,
    "group_id" int NOT NULL,
    "title" varchar NOT NULL DEFAULT '',
    "releaseDate" date NOT NULL DEFAULT 'now()',
    "text" text NOT NULL DEFAULT '',
    "link" varchar NOT NULL DEFAULT '',
    CONSTRAINT unique_group_song UNIQUE (group_id, title),
    FOREIGN KEY ("group_id") REFERENCES "groups" ("id")
);
CREATE INDEX ON "songs" ("group_id");
CREATE INDEX ON "songs" ("releaseDate");