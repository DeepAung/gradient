BEGIN;

INSERT INTO "users" ("username", "email", "password", "is_admin")
  VALUES ('DeepAung', 'i.deepaung@gmail.com', '', FALSE),
  ('admin', 'admin@gmail.com', '', 1);

INSERT INTO "tasks" ("display_name", "url_name", "content_url", "testcase_count")
  VALUES ('Two Sum', 'two_sum', '', 10);

INSERT INTO "submissions" ("user_id", "task_id", "code", "results")
  VALUES (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', 'PPPPPPPPPP');

COMMIT;
