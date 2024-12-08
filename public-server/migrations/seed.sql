BEGIN;

INSERT INTO "users" ("username", "email", "password", "is_admin")
  VALUES ('DeepAung', 'i.deepaung@gmail.com', '$2a$10$OqveZpSgfd5KMU1Xeo6UUeCYMWYgz3kjkuZvrxEaVsujmnxI/P/oe', FALSE),
  ('admin', 'admin@gmail.com', '$2a$10$OqveZpSgfd5KMU1Xeo6UUeCYMWYgz3kjkuZvrxEaVsujmnxI/P/oe', TRUE);

INSERT INTO "tasks" ("display_name", "url_name", "content_url", "testcase_count")
  VALUES ('Two Sum', 'two_sum', '', 10);

INSERT INTO "submissions" ("user_id", "task_id", "code", "language", "results")
  VALUES (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', 'python', 'PPPPPPPPPP');

COMMIT;
