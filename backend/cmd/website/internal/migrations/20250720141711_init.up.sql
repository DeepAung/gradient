CREATE TABLE "identities" (
	id UUID NOT NULL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	email VARCHAR(255) NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--bun:split
CREATE TABLE "tokens" (
	id UUID NOT NULL PRIMARY KEY,
	identity_id UUID NOT NULL REFERENCES identities (id),
	access_token VARCHAR(255) NOT NULL UNIQUE,
	refresh_token VARCHAR(255) NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMP NOT NULL
);

--bun:split
CREATE TABLE "classrooms" (
	id UUID NOT NULL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE,
	owner_id UUID NOT NULL REFERENCES identities (id),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--bun:split
CREATE TYPE person_role AS ENUM ('student', 'teacher');
CREATE TABLE "persons" (
	id UUID NOT NULL PRIMARY KEY,
	identity_id UUID NOT NULL REFERENCES identities (id),
	classroom_id UUID NOT NULL REFERENCES classrooms (id),
	display_name VARCHAR(255) NOT NULL,
	role person_role NOT NULL DEFAULT 'student',

	UNIQUE (classroom_id, display_name)
);

--bun:split
CREATE TABLE "tasks" (
	id UUID NOT NULL PRIMARY KEY,
	owner_id UUID NOT NULL REFERENCES identities (id),
	classroom_id UUID REFERENCES classrooms (id), -- if null, the task is global
	problem_url VARCHAR(255) NOT NULL,
	testcases_url VARCHAR(255) NOT NULL, -- e.g. storage.gradient.dev/{task-id}/testcases/
	due_at TIMESTAMP,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--bun:split
CREATE TABLE "submissions" (
	id UUID NOT NULL PRIMARY KEY,
	task_id UUID NOT NULL REFERENCES tasks (id),
	person_id UUID NOT NULL REFERENCES persons (id),
	score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--bun:split
CREATE TYPE submission_result_type AS ENUM ('ok', 'error', 'timeout', 'runtime_error', 'wrong_answer');
CREATE TABLE "submission_results" (
	id UUID NOT NULL PRIMARY KEY,
	submission_id UUID NOT NULL REFERENCES submissions (id),
	testcase_id INTEGER NOT NULL, -- uint
	result_type submission_result_type NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	logs TEXT
);

--bun:split
CREATE TABLE "tags" (
	id UUID NOT NULL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE "tasks_tags" (
	id UUID NOT NULL PRIMARY KEY,
	task_id UUID NOT NULL REFERENCES tasks (id),
	tag_id UUID NOT NULL REFERENCES tags (id),
	UNIQUE (task_id, tag_id)
);

--bun:split
CREATE TABLE "pending_invites" (
	id UUID NOT NULL PRIMARY KEY,
	classroom_id UUID NOT NULL REFERENCES classrooms (id),
	inviter_id UUID NOT NULL REFERENCES identities (id),
	invitee_id UUID NOT NULL REFERENCES identities (id),
	role person_role NOT NULL DEFAULT 'student',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TIMESTAMP NOT NULL
);

--bun:split
CREATE TYPE pending_join_status AS ENUM ('pending', 'approved', 'rejected');
CREATE TABLE "pending_joins" (
	id UUID NOT NULL PRIMARY KEY,
	identity_id UUID NOT NULL REFERENCES identities (id),
	classroom_id UUID NOT NULL REFERENCES classrooms (id),
	status pending_join_status NOT NULL DEFAULT 'pending',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
