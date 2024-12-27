package submissions

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DeepAung/gradient/public-server/modules/tasks"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/modules/users"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var (
	ErrSubmissionNotFound = fiber.NewError(fiber.StatusBadRequest, "submission not found")
	ErrInvalidLanguage    = fiber.NewError(fiber.StatusBadRequest, "invalid language")
	ErrInvalidScore       = fiber.NewError(fiber.StatusBadRequest, "invalid score")
)

type submissionRepo struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewSubmissionRepo(
	db *sqlx.DB,
	timeout time.Duration,
) types.SubmissionsRepo {
	return &submissionRepo{
		db:      db,
		timeout: timeout,
	}
}

func (r *submissionRepo) CreateSubmission(req types.CreateSubmissionReq) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, err
	}

	id, err := r.createSubmissionWithDB(ctx, tx, req)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *submissionRepo) CanCreateSubmission(req types.CreateSubmissionReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if _, err := r.createSubmissionWithDB(ctx, tx, req); err != nil {
		return err
	}

	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}

type mydb interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
}

func (r *submissionRepo) createSubmissionWithDB(
	ctx context.Context,
	db mydb,
	req types.CreateSubmissionReq,
) (int, error) {
	var id int

	err := sqlx.GetContext(ctx, db, &id,
		`INSERT INTO submissions ( user_id, task_id, code, language_index, score)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id`,
		req.UserId, req.TaskId, req.Code, req.LanguageIndex, req.Score,
	)
	if err != nil {
		switch err.Error() {
		case sql.ErrNoRows.Error():
			return 0, ErrSubmissionNotFound
		case `pq: insert or update on table "submissions" violates foreign key constraint "submissions_user_id_fkey"`:
			return 0, users.ErrUserNotFound
		case `pq: insert or update on table "submissions" violates foreign key constraint "submissions_task_id_fkey"`:
			return 0, tasks.ErrTaskNotFound
		case `pq: invalid input value for enum language: ""`:
			return 0, ErrInvalidLanguage
		case `pq: new row for relation "submissions" violates check constraint "submissions_score_check"`:
			return 0, ErrInvalidScore
		default:
			return 0, err
		}
	}

	for _, evaluation := range req.Evaluations {
		result, err := db.ExecContext(ctx,
			`INSERT INTO evaluations (submission_id, time, memory, status)
			VALUES ($1, $2, $3, $4)`,
			evaluation.SubmissionId,
			evaluation.TimeMicroSeconds,
			evaluation.MemoryKiloBytes,
			evaluation.Status,
		)
		if err != nil {
			return 0, err
		}
		n, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, errors.New("invalid sql schema?")
		}
	}

	return id, nil
}

// TODO: join things up
func (r *submissionRepo) FindOneSubmission(id int) (types.Submission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var submission types.Submission
	err := r.db.GetContext(ctx, &submission,
		`SELECT id, user_id, task_id, code, language_index, score
		FROM submissions WHERE id = $1`,
		id)
	if err == sql.ErrNoRows {
		return types.Submission{}, ErrSubmissionNotFound
	}

	return submission, err
}

func (r *submissionRepo) FindManySubmissions(
	req types.GetSubmissionsReq,
) ([]types.Submission, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var submissions []types.Submission
	err := r.db.SelectContext(ctx, &submissions,
		`SELECT id, user_id, task_id, code, language_index, score
		FROM submissions
		WHERE user_id = COALESCE($1, user_id) AND task_id = COALESCE($2, task_id)
		ORDER BY id DESC`,
		req.UserId, req.TaskId)

	return submissions, err
}
