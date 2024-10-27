
CREATE TABLE "tasks" (
                         "id" bigserial PRIMARY KEY,
                         "title" varchar NOT NULL,
                         "text" text,
                         "created_at" timestamptz NOT NULL DEFAULT (now()),
                         "updated_at" timestamptz NOT NULL DEFAULT (now()),
                         "done" bool DEFAULT false
);