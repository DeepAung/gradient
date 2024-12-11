BEGIN;

DROP EXTENSION IF EXISTS pg_trgm CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "tokens" CASCADE;
DROP TABLE IF EXISTS "tasks" CASCADE;
DROP TABLE IF EXISTS "users_tasks_info" CASCADE;
DROP TABLE IF EXISTS "submissions" CASCADE;
DROP TYPE IF EXISTS LANGUAGE CASCADE;
DROP TRIGGER IF EXISTS after_submissions_insert ON submissions CASCADE;
DROP FUNCTION IF EXISTS after_submissions_insert CASCADE;

CREATE EXTENSION pg_trgm;

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
  "testcase_count" INTEGER NOT NULL CHECK("testcase_count" > 0),
  "solved_number" INTEGER NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE "users_tasks_info" (
  "user_id" INTEGER REFERENCES "users" ("id") ON DELETE CASCADE,
  "task_id" INTEGER REFERENCES "tasks" ("id") ON DELETE CASCADE,
  "score" REAL NOT NULL DEFAULT 0,
  PRIMARY KEY ("user_id", "task_id")
);

CREATE TYPE LANGUAGE AS ENUM ('cpp', 'c', 'go', 'python');
CREATE TABLE "submissions" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER REFERENCES "users" ("id") ON DELETE CASCADE,
  "task_id" INTEGER REFERENCES "tasks" ("id") ON DELETE CASCADE,
  "code" VARCHAR, -- NOT NULL,
  "language" LANGUAGE NOT NULL,
  "results" VARCHAR NOT NULL,
  "result_percent" REAL NOT NULL CHECK("result_percent" >= 0),
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE OR REPLACE FUNCTION after_submissions_insert() RETURNS trigger AS $after_submissions_insert$
  DECLARE oldscore REAL := NULL;
  BEGIN

    SELECT score INTO oldscore FROM users_tasks_info WHERE user_id = NEW.user_id AND task_id = NEW.task_id;

    IF oldscore IS NULL THEN
      INSERT INTO users_tasks_info (user_id, task_id, score)
      VALUES (NEW.user_id, NEW.task_id, NEW.result_percent);
      oldscore := 0; -- bz if oldscore IS NULL then oldscore != 100 will evaulate to false (https://stackoverflow.com/a/12108091)
    ELSE
      UPDATE users_tasks_info SET
      score = GREATEST(users_tasks_info.score, NEW.result_percent);
    END IF;

    IF (NEW.result_percent = 100 AND oldscore != 100) THEN
      UPDATE tasks SET solved_number = tasks.solved_number + 1
      WHERE tasks.id = NEW.task_id;
    END IF;

    RETURN NULL;
  END;
$after_submissions_insert$ LANGUAGE plpgsql;

CREATE TRIGGER after_submissions_insert
AFTER INSERT ON submissions
  FOR EACH ROW EXECUTE FUNCTION after_submissions_insert();

COMMIT;
