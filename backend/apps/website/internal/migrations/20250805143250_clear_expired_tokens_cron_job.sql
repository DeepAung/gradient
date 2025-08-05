-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_cron;

SELECT cron.schedule(
	'delete-expired-tokens-job',
	'0 3 * * *',
	$$DELETE FROM tokens WHERE expired_at < NOW();$$
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT cron.unschedule('delete-expired-tokens-job');

DROP EXTENSION IF EXISTS pg_cron;
-- +goose StatementEnd
