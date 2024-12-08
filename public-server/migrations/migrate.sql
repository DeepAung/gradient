BEGIN;

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "tokens" CASCADE;
DROP TABLE IF EXISTS "tasks" CASCADE;
DROP TABLE IF EXISTS "submissions" CASCADE;
DROP TYPE IF EXISTS LANGUAGE CASCADE;

CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "username" VARCHAR UNIQUE NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "picture_url" VARCHAR NOT NULL DEFAULT '',
  "is_admin" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE "tokens" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER REFERENCES "users" ("id") ON DELETE CASCADE,
  "access_token" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE "tasks" (
  "id" SERIAL PRIMARY KEY,
  "display_name" VARCHAR UNIQUE NOT NULL,
  "url_name" VARCHAR UNIQUE NOT NULL,
  "content_url" VARCHAR NOT NULL,
  "testcase_count" INTEGER NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TYPE LANGUAGE AS ENUM ('cpp', 'c', 'go', 'python');
CREATE TABLE "submissions" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER REFERENCES "users" ("id") ON DELETE CASCADE,
  "task_id" INTEGER REFERENCES "tasks" ("id") ON DELETE CASCADE,
  "code" VARCHAR NOT NULL,
  "language" LANGUAGE NOT NULL,
  "results" VARCHAR NOT NULL,
  "result_percent" REAL NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);

COMMIT;
