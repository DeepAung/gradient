-- +goose Up
-- +goose StatementBegin

-- TODO: add go unit test that if we soft delete graident, the problem will not be visible anymore

CREATE TABLE "users" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL UNIQUE,
	email VARCHAR(255) NOT NULL UNIQUE,
	profile_url TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);

CREATE TABLE "socials" (
	name VARCHAR(255) NOT NULL PRIMARY KEY
);

CREATE TABLE "oauths" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	social_name VARCHAR(255) NOT NULL REFERENCES socials (name) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (user_id, social_name)
);

CREATE TABLE "tokens" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	access_token TEXT NOT NULL UNIQUE,
	refresh_token TEXT NOT NULL UNIQUE,
	expired_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "gradients" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	owner_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	name VARCHAR(255) NOT NULL UNIQUE,
	is_public BOOLEAN NOT NULL DEFAULT FALSE,
	banner_url TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);

CREATE TABLE "contests" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	gradient_id UUID NOT NULL REFERENCES gradients (id) ON DELETE CASCADE,
	start_at TIMESTAMP NOT NULL,
	end_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP,
	UNIQUE (gradient_id, name)
);

CREATE TYPE member_role AS ENUM ('owner', 'admin', 'participant');

CREATE TABLE "members" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	gradient_id UUID NOT NULL REFERENCES gradients (id) ON DELETE CASCADE,
	display_name VARCHAR(255) NOT NULL,
	role member_role NOT NULL DEFAULT 'participant',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP,
	UNIQUE (gradient_id, user_id),
	UNIQUE (gradient_id, display_name)
);

CREATE TABLE "problems" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	owner_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	gradient_id UUID NOT NULL REFERENCES gradients (id) ON DELETE CASCADE,
	contest_id UUID REFERENCES contests (id) ON DELETE CASCADE, -- NULL = this problem tied to gradient, not contest
	is_published BOOLEAN DEFAULT FALSE,
	problem_url TEXT NOT NULL,
	testcases_url TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);

CREATE TABLE "submissions" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	problem_id UUID NOT NULL REFERENCES problems (id) ON DELETE CASCADE,
	member_id UUID NOT NULL REFERENCES members (id) ON DELETE CASCADE,
	avg_time_ns BIGINT,
	avg_memory_bytes BIGINT,
	score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);

CREATE TYPE submission_result_type AS ENUM ('ok', 'server_error', 'runtime_error', 'timeout', 'memory_exceeded', 'wrong_answer');

CREATE TABLE "submission_results" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	submission_id UUID NOT NULL REFERENCES submissions (id) ON DELETE CASCADE,
	testcase_index INTEGER NOT NULL CHECK (testcase_index >= 0),
	time_ns BIGINT, -- Time in nano seconds
	memory_bytes BIGINT, -- Memory in bytes
	result_type submission_result_type NOT NULL,
	logs TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "tags" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP
);

CREATE TABLE "problems_tags" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	problem_id UUID NOT NULL REFERENCES problems (id) ON DELETE CASCADE,
	tag_id UUID NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE (problem_id, tag_id)
);

CREATE TABLE "pending_invites" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	gradient_id UUID NOT NULL REFERENCES gradients (id) ON DELETE CASCADE,
	inviter_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	invitee_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	role member_role NOT NULL DEFAULT 'participant',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expired_at TIMESTAMP NOT NULL
);

CREATE TYPE pending_join_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE "pending_joins" (
	id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
	gradient_id UUID NOT NULL REFERENCES gradients (id) ON DELETE CASCADE,
	status pending_join_status NOT NULL DEFAULT 'pending',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS "pending_joins";
DROP TYPE IF EXISTS pending_join_status;

DROP TABLE IF EXISTS "pending_invites";

DROP TABLE IF EXISTS "problems_tags";
DROP TABLE IF EXISTS "tags";

DROP TABLE IF EXISTS "submission_results";
DROP TYPE IF EXISTS submission_result_type;

DROP TABLE IF EXISTS "submissions";
DROP TABLE IF EXISTS "problems";

DROP TABLE IF EXISTS "members";
DROP TYPE IF EXISTS member_role;

DROP TABLE IF EXISTS "contests";
DROP TABLE IF EXISTS "gradients";

DROP TABLE IF EXISTS "tokens";
DROP TABLE IF EXISTS "oauths";

DROP TABLE IF EXISTS "socials";
DROP TABLE IF EXISTS "users";

-- +goose StatementEnd
