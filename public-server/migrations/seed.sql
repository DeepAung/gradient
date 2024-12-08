BEGIN;

INSERT INTO "users" ("username", "email", "password", "is_admin")
  VALUES ('DeepAung', 'i.deepaung@gmail.com', '$2a$10$OqveZpSgfd5KMU1Xeo6UUeCYMWYgz3kjkuZvrxEaVsujmnxI/P/oe', FALSE),
  ('admin', 'admin@gmail.com', '$2a$10$OqveZpSgfd5KMU1Xeo6UUeCYMWYgz3kjkuZvrxEaVsujmnxI/P/oeNUMB', TRUE);

INSERT INTO "tasks" ("display_name", "url_name", "content_url", "testcase_count")
  VALUES ('Two Sum', 'two_sum', '', 10),
  ('Two Product', 'two_product', '', 10),
  ('Dijkstra', 'dijkstra', '', 1),
  ('Floyd Warshall', 'floyd_warshall', '', 1);

INSERT INTO "submissions" ("user_id", "task_id", "code", "language", "results", "result_percent")
  VALUES (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', 'python', 'PPPPPPPPPP', 100),
  (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', 'python', 'PPPPPPPPPP', 100),
  (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', 'python', 'PPPPPPPPPP', 100),
  (1, 1, 'print(123456)', 'python', '----------', 0),
  (1, 1, 'println(123456)', 'go', '----------', 0),
  (1, 2, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', 'python', 'PPPPPPPPPP', 100),
  (1, 2, 'print(123456)', 'python', '----------', 0),
  (1, 3, 'print(123456)', 'python', '----------', 0);

COMMIT;
