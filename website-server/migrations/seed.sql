BEGIN;

INSERT INTO "users" ("username", "email", "password", "is_admin")
  VALUES ('DeepAung', 'i.deepaung@gmail.com', '$2a$10$OqveZpSgfd5KMU1Xeo6UUeCYMWYgz3kjkuZvrxEaVsujmnxI/P/oe', FALSE),
  ('admin', 'admin@gmail.com', '$2a$10$OqveZpSgfd5KMU1Xeo6UUeCYMWYgz3kjkuZvrxEaVsujmnxI/P/oeNUMB', TRUE);

INSERT INTO "tasks" ("display_name", "url_name", "content_url", "testcase_count")
  VALUES ('Sum of Pairs', 'sum_of_pairs', '', 10),
  ('Multiplication Madness', 'multiplication_madness', '', 10),
  ('Shortest Path in a Small Town', 'shortest_path_in_a_small_town', '', 1),
  ('Shortest Paths in a Large Town', 'shortest_paths_in_a_large_town', '', 1);

INSERT INTO "submissions" ("user_id", "task_id", "code", "language_index", "score", "max_time", "max_memory")
  VALUES (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', '3', 100, 0, 0),
  (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', '3', 100, 0, 0),
  (1, 1, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', '3', 100, 0, 0),
  (1, 1, 'print(123456)', '3', 0, 0, 0),
  (1, 1, 'println(123456)', '2', 0, 0, 0),
  (1, 2, 'for _ in range(len(int(input()))): print(int(input()) + int(input()))', '3', 100, 0, 0),
  (1, 2, 'print(123456)', '3', 0, 0, 0),
  (1, 3, 'print(123456)', '3', 0, 0, 0),
  (2, 1, 'print(123456)', '3', 0, 0, 0);

INSERT INTO "submission_results" ("submission_id", "time", "memory", "status")
  VALUES
    (1, 10, 18, 'P'), (1, 0, 20, 'P'), (1, 0, 17, 'P'), (1, 0, 0, 'P'), (1, 0, 0, 'P'), (1, 0, 0, 'P'), (1, 0, 0, 'P'), (1, 0, 0, 'P'), (1, 0, 0, 'P'), (1, 0, 0, 'P'),
    (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'), (2, 0, 0, 'P'),
    (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'), (3, 0, 0, 'P'),
    (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'), (4, 0, 0, '-'),
    (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'), (5, 0, 0, '-'),
    (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'), (6, 0, 0, 'P'),
    (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'), (7, 0, 0, '-'),
    (8, 0, 0, '-'),
    (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-'), (9, 0, 0, '-');

COMMIT;
