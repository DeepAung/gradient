package tasks

import (
	"database/sql"
	"errors"

	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/jmoiron/sqlx"
)

var (
	ErrTaskNotFound         = errors.New("task not found")
	ErrUniqueDisplayName    = errors.New("display name already exists")
	ErrUniqueUrlName        = errors.New("url name already exists")
	ErrInvalidTestcaseCount = errors.New("invalid testcase count")
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) types.TaskRepo {
	return &TaskRepo{
		db: db,
	}
}

func (r *TaskRepo) FindOneTask(id int) (types.Task, error) {
	var task types.Task
	err := r.db.Get(&task,
		`SELECT id, display_name, url_name, content_url, testcase_count
		FROM tasks 
		WHERE id = $1;`,
		id)
	if err == sql.ErrNoRows {
		return types.Task{}, ErrTaskNotFound
	}

	return task, err
}

// [startIndex, stopIndex)
func (r *TaskRepo) FindManyTasks(
	search string,
	onlyCompleted bool,
	startIndex, stopIndex int,
) ([]types.Task, error) {
	var tasks []types.Task
	err := r.db.Select(&tasks, `
		SELECT tasks.id, tasks.display_name, tasks.url_name, tasks.content_url, tasks.testcase_count
		FROM tasks
		INNER JOIN submissions
		ON (NOT $1) OR (submissions.task_id = tasks.id AND submissions.result_percent = 100)
		WHERE ($2 = '') OR ($2 % ANY(STRING_TO_ARRAY(url_name || ' ' || display_name, ' ')))
		GROUP BY tasks.id
		ORDER BY tasks.id
		LIMIT $3 OFFSET $4
	`, onlyCompleted, search, stopIndex-startIndex, startIndex)

	return tasks, err
}

func (r *TaskRepo) CreateTask(req types.CreateTaskReq) error {
	result, err := r.db.Exec(
		`INSERT INTO tasks (display_name, url_name, content_url, testcase_count) 
			VALUES ($1, $2, $3, $4)`,
		req.DisplayName, req.UrlName, req.ContentUrl, req.TestcaseCount)
	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "tasks_display_name_key"`:
			return ErrUniqueDisplayName
		case `pq: duplicate key value violates unique constraint "tasks_url_name_key"`:
			return ErrUniqueUrlName
		case `pq: new row for relation "tasks" violates check constraint "tasks_testcase_count_check"`:
			return ErrInvalidTestcaseCount
		default:
			return err
		}
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepo) UpdateTask(req types.Task) error {
	result, err := r.db.Exec(
		`UPDATE tasks SET display_name = $1, url_name = $2, content_url = $3, testcase_count = $4
		WHERE id = $5`,
		req.DisplayName, req.UrlName, req.ContentUrl, req.TestcaseCount, req.Id)
	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "tasks_display_name_key"`:
			return ErrUniqueDisplayName
		case `pq: duplicate key value violates unique constraint "tasks_url_name_key"`:
			return ErrUniqueUrlName
		case `pq: new row for relation "tasks" violates check constraint "tasks_testcase_count_check"`:
			return ErrInvalidTestcaseCount
		default:
			return err
		}
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *TaskRepo) DeleteTask(id int) error {
	result, err := r.db.Exec(`DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrTaskNotFound
	}
	return nil
}
