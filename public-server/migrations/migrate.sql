BEGIN;

DROP EXTENSION IF EXISTS pg_trgm CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "tokens" CASCADE;
DROP TABLE IF EXISTS "tasks" CASCADE;
DROP TABLE IF EXISTS "users_tasks_info" CASCADE;
DROP TABLE IF EXISTS "submissions" CASCADE;
DROP TABLE IF EXISTS "evaluations" CASCADE;
DROP TYPE IF EXISTS LANGUAGE CASCADE;
DROP TRIGGER IF EXISTS update_score_and_solved_number ON submissions CASCADE;
DROP TRIGGER IF EXISTS update_max_time_memory ON evaluations CASCADE;
DROP FUNCTION IF EXISTS update_max_time_memory CASCADE;
DROP FUNCTION IF EXISTS update_score_and_solved_number CASCADE;

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

CREATE TABLE "submissions" (
  "id" SERIAL PRIMARY KEY,
  "user_id" INTEGER REFERENCES "users" ("id") ON DELETE CASCADE,
  "task_id" INTEGER REFERENCES "tasks" ("id") ON DELETE CASCADE,
  "code" VARCHAR NOT NULL,
  "language_index" INTEGER NOT NULL, -- refer to language index in proto file
  "score" REAL NOT NULL CHECK("score" >= 0),
  "max_time" INTEGER NOT NULL, -- micro seconds (ms)
  "max_memory" INTEGER NOT NULL, -- kilo bytes (kB)
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);

CREATE TABLE "evaluations" (
  "id" SERIAL PRIMARY KEY,
  "submission_id" INTEGER REFERENCES "submissions" ("id") ON DELETE CASCADE,
  "time" INTEGER NOT NULL, -- micro seconds (ms)
  "memory" INTEGER NOT NULL, -- kilo bytes (kB)
  "status" CHAR NOT NULL
);

CREATE OR REPLACE FUNCTION update_score_and_solved_number() RETURNS trigger AS $update_score_and_solved_number$
  DECLARE oldscore REAL := NULL;
  BEGIN

    SELECT score INTO oldscore FROM users_tasks_info WHERE user_id = NEW.user_id AND task_id = NEW.task_id;

    IF oldscore IS NULL THEN
      INSERT INTO users_tasks_info (user_id, task_id, score)
      VALUES (NEW.user_id, NEW.task_id, NEW.score);
      oldscore := 0; -- bz if oldscore IS NULL then oldscore != 100 will evaulate to false (https://stackoverflow.com/a/12108091)
    ELSE
      UPDATE users_tasks_info SET
      score = GREATEST(users_tasks_info.score, NEW.score);
    END IF;

    IF (NEW.score = 100 AND oldscore != 100) THEN
      UPDATE tasks SET solved_number = tasks.solved_number + 1
      WHERE tasks.id = NEW.task_id;
    END IF;

    RETURN NULL;
  END;
$update_score_and_solved_number$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_max_time_memory() RETURNS trigger AS $update_max_time_memory$
  BEGIN
    UPDATE submissions SET
    max_time = GREATEST(max_time, NEW.time),
    max_memory = GREATEST(max_memory, NEW.memory)
    WHERE submissions.id = NEW.submission_id;

    RETURN NULL;
  END;
$update_max_time_memory$ LANGUAGE plpgsql;

CREATE TRIGGER update_score_and_solved_number
AFTER INSERT ON submissions
  FOR EACH ROW EXECUTE FUNCTION update_score_and_solved_number();

CREATE TRIGGER update_max_time_memory
AFTER INSERT ON evaluations
  FOR EACH ROW EXECUTE FUNCTION update_max_time_memory();

COMMIT;
