CREATE TABLE IF NOT EXISTS "songs" (
    "id" serial PRIMARY KEY,
    "group_id" int NOT NULL,
    "title" varchar NOT NULL DEFAULT '',
    "release_date" date NOT NULL DEFAULT 'now()',
    "text" text NOT NULL DEFAULT '',
    "link" varchar NOT NULL DEFAULT '',
    CONSTRAINT unique_group_song UNIQUE (group_id, title),
    FOREIGN KEY ("group_id") REFERENCES "groups" ("id")
);
CREATE INDEX ON "songs" ("group_id");
CREATE INDEX ON "songs" ("release_date");