package tasks

import (
	"context"
	"database/sql"
	"time"

	"github.com/DeepAung/gradient/website-server/modules/types"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

var (
	ErrTaskNotFound         = fiber.NewError(fiber.StatusBadRequest, "task not found")
	ErrUniqueDisplayName    = fiber.NewError(fiber.StatusBadRequest, "display name already exists")
	ErrUniqueUrlName        = fiber.NewError(fiber.StatusBadRequest, "url name already exists")
	ErrInvalidTestcaseCount = fiber.NewError(fiber.StatusBadRequest, "invalid testcase count")
)

type tasksRepoImpl struct {
	db      *sqlx.DB
	timeout time.Duration
}

func NewTasksRepo(db *sqlx.DB, timeout time.Duration) types.TasksRepo {
	return &tasksRepoImpl{
		db:      db,
		timeout: timeout,
	}
}

func (r *tasksRepoImpl) FindOneTask(userId, taskId int) (types.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var task types.Task
	err := r.db.GetContext(ctx, &task,
		`SELECT
			tasks.id,
			tasks.display_name,
			tasks.url_name,
			tasks.content_url,
			tasks.testcase_count,
			tasks.solved_number,
			COALESCE(info.score, 0) as score
		FROM tasks 
		LEFT JOIN users_tasks_info AS info
		ON info.user_id = $1 AND info.task_id = $2
		WHERE tasks.id = $2;`,
		userId, taskId)
	if err == sql.ErrNoRows {
		return types.Task{}, ErrTaskNotFound
	}

	return task, err
}

func (r *tasksRepoImpl) FindOneTaskTestcaseCount(taskId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var testcaseCount int
	err := r.db.GetContext(ctx, &testcaseCount,
		`SELECT tasks.testcase_count FROM tasks WHERE tasks.id = $1;`,
		taskId)
	if err == sql.ErrNoRows {
		return 0, ErrTaskNotFound
	}

	return testcaseCount, err
}

// [startIndex, stopIndex)
func (r *tasksRepoImpl) FindManyTasks(
	userId int,
	search string,
	onlyCompleted bool,
	startIndex, stopIndex int,
) ([]types.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var tasks []types.Task

	bindedStmt, args, err := r.db.BindNamed(`
		SELECT 
			tasks.id,
			tasks.display_name,
			tasks.url_name,
			tasks.content_url,
			tasks.testcase_count,
			tasks.solved_number,
			COALESCE(info.score, 0) as score
		FROM tasks
		LEFT JOIN users_tasks_info AS info
		ON info.user_id = :userId AND info.task_id = tasks.id
		WHERE (:search = '') OR (:search % ANY(STRING_TO_ARRAY(url_name || ' ' || display_name, ' ')))
		GROUP BY tasks.id, info.score
		HAVING NOT :onlyCompleted OR MAX(info.score) = 100
		ORDER BY tasks.id
		LIMIT :limit OFFSET :offset;`,
		map[string]interface{}{
			"userId":        userId,
			"onlyCompleted": onlyCompleted,
			"search":        search,
			"limit":         stopIndex - startIndex,
			"offset":        startIndex,
		})
	if err != nil {
		return nil, err
	}
	err = r.db.SelectContext(ctx, &tasks, bindedStmt, args...)

	return tasks, err
}

func (r *tasksRepoImpl) CreateTask(req types.CreateUpdateTaskReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx,
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

func (r *tasksRepoImpl) UpdateTask(id int, req types.CreateUpdateTaskReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET display_name = $1, url_name = $2, content_url = $3, testcase_count = $4
		WHERE id = $5`,
		req.DisplayName, req.UrlName, req.ContentUrl, req.TestcaseCount, id)
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

func (r *tasksRepoImpl) DeleteTask(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
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
