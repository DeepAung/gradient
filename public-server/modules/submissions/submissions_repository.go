package submissions

import (
	"database/sql"

	"github.com/DeepAung/gradient/public-server/modules/tasks"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/modules/users"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var (
	ErrSubmissionNotFound   = fiber.NewError(fiber.StatusBadRequest, "submission not found")
	ErrInvalidLanguage      = fiber.NewError(fiber.StatusBadRequest, "invalid language")
	ErrInvalidResultPercent = fiber.NewError(fiber.StatusBadRequest, "invalid result percent")
)

type submissionRepo struct {
	db *sqlx.DB
}

func NewSubmissionRepo(db *sqlx.DB) types.SubmissionsRepo {
	return &submissionRepo{
		db: db,
	}
}

func (r *submissionRepo) CreateSubmission(req types.CreateSubmissionReq) (int, error) {
	return r.createSubmissionWithDB(r.db, req)
}

func (r *submissionRepo) CanCreateSubmission(req types.CreateSubmissionReq) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	if _, err := r.createSubmissionWithDB(tx, req); err != nil {
		return err
	}

	if err := tx.Rollback(); err != nil {
		return err
	}
	return nil
}

func (r *submissionRepo) createSubmissionWithDB(
	db sqlx.Queryer,
	req types.CreateSubmissionReq,
) (int, error) {
	var id int
	err := sqlx.Get(db, &id,
		`INSERT INTO submissions (user_id, task_id, code, language, results, result_percent)
			VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		req.UserId,
		req.TaskId,
		req.Code,
		types.ProtoLanguageToString(req.Language),
		req.Results,
		req.ResultPercent,
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
		case `pq: new row for relation "submissions" violates check constraint "submissions_result_percent_check"`:
			return 0, ErrInvalidResultPercent
		default:
			return 0, err
		}
	}

	return id, nil
}

func (r *submissionRepo) FindOneSubmission(id int) (types.Submission, error) {
	var submission types.Submission
	err := r.db.Get(&submission,
		`SELECT id, user_id, task_id, code, language, results, result_percent 
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
	var submissions []types.Submission
	err := r.db.Select(&submissions,
		`SELECT id, user_id, task_id, code, language, results, result_percent 
		FROM submissions
		WHERE user_id = COALESCE($1, user_id) AND task_id = COALESCE($2, task_id)
		ORDER BY id DESC`,
		req.UserId, req.TaskId)

	return submissions, err
}
