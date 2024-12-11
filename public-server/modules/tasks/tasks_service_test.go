package tasks

import (
	"testing"

	"github.com/DeepAung/gradient/public-server/config"
	"github.com/DeepAung/gradient/public-server/database"
	"github.com/DeepAung/gradient/public-server/modules/types"
	"github.com/DeepAung/gradient/public-server/pkg/asserts"
	"github.com/jmoiron/sqlx"
)

var (
	migrateSourceName = "../../migrations/migrate.sql"
	seedSourceName    = "../../migrations/seed.sql"
	cfg               *config.Config
	db                *sqlx.DB
	repo              types.TasksRepo
	svc               types.TasksSvc

	listTasks = []types.Task{
		{Id: 1, DisplayName: "Two Sum", UrlName: "two_sum", TestcaseCount: 10},
		{Id: 2, DisplayName: "Two Product", UrlName: "two_product", TestcaseCount: 10},
		{Id: 3, DisplayName: "Dijkstra", UrlName: "dijkstra", TestcaseCount: 1},
		{Id: 4, DisplayName: "Floyd Warshall", UrlName: "floyd_warshall", TestcaseCount: 1},
	}

	twoSumTask = listTasks[0]

	createTask = types.CreateTaskReq{
		DisplayName:   "New Task",
		UrlName:       "new_task",
		ContentUrl:    "new_task.com",
		TestcaseCount: 10,
	}

	updateTask = types.Task{
		Id:            9,
		DisplayName:   "New Task (updated)",
		UrlName:       "new_task (updated)",
		ContentUrl:    "new_task.com (updated)",
		TestcaseCount: 15,
	}
)

func init() {
	cfg = config.NewConfig("../../.env.dev")
	db = database.InitDB(cfg.App.DbUrl)
	database.RunSQL(db, migrateSourceName)
	database.RunSQL(db, seedSourceName)
	repo = NewTasksRepo(db)
	svc = NewTasksSvc(repo)
}

func TestGetTask(t *testing.T) {
	t.Run("id not found", func(t *testing.T) {
		task, err := svc.GetTask(1000)
		asserts.EqualError(t, err, ErrTaskNotFound)
		asserts.Equal(t, "task", task, types.Task{})
	})

	t.Run("normal get task", func(t *testing.T) {
		task, err := svc.GetTask(1)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "task", task, twoSumTask)
	})
}

func TestGetTasks(t *testing.T) {
	t.Run("normal get tasks", func(t *testing.T) {
		tasks, err := svc.GetTasks(1, "", false, 0, 100)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "tasks", tasks, listTasks)
	})

	t.Run("with search", func(t *testing.T) {
		tasks, err := svc.GetTasks(1, "two", false, 0, 100)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "tasks", tasks, listTasks[:2])
	})

	t.Run("only completed", func(t *testing.T) {
		tasks, err := svc.GetTasks(1, "", true, 0, 100)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "tasks", tasks, listTasks[:2])
	})

	t.Run("offset", func(t *testing.T) {
		tasks, err := svc.GetTasks(1, "", false, 1, 100)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "tasks", tasks, listTasks[1:])
	})

	t.Run("limit", func(t *testing.T) {
		tasks, err := svc.GetTasks(1, "", false, 0, 2)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "tasks", tasks, listTasks[0:2])
	})
}

func TestCreateTask(t *testing.T) {
	t.Run("unique display name", func(t *testing.T) {
		task := createTask
		task.DisplayName = twoSumTask.DisplayName
		err := svc.CreateTask(task)
		asserts.EqualError(t, err, ErrUniqueDisplayName)
	})

	t.Run("unique url name", func(t *testing.T) {
		task := createTask
		task.UrlName = twoSumTask.UrlName
		err := svc.CreateTask(types.CreateTaskReq{
			DisplayName:   "New Task",
			UrlName:       "two_sum",
			ContentUrl:    "new_task.com",
			TestcaseCount: 10,
		})
		asserts.EqualError(t, err, ErrUniqueUrlName)
	})

	t.Run("invalid testcase count", func(t *testing.T) {
		task := createTask
		task.TestcaseCount = 0
		err := svc.CreateTask(task)
		asserts.EqualError(t, err, ErrInvalidTestcaseCount)

		task.TestcaseCount = -10
		err = svc.CreateTask(task)
		asserts.EqualError(t, err, ErrInvalidTestcaseCount)
	})

	t.Run("normal create task", func(t *testing.T) {
		err := svc.CreateTask(createTask)
		asserts.EqualError(t, err, nil)
	})
}

func TestUpdateTask(t *testing.T) {
	t.Run("id not found", func(t *testing.T) {
		task := updateTask
		task.Id = 1000
		err := svc.UpdateTask(task)
		asserts.EqualError(t, err, ErrTaskNotFound)
	})

	t.Run("unique display name", func(t *testing.T) {
		task := updateTask
		task.DisplayName = twoSumTask.DisplayName
		err := svc.UpdateTask(task)
		asserts.EqualError(t, err, ErrUniqueDisplayName)
	})

	t.Run("unique url name", func(t *testing.T) {
		task := updateTask
		task.UrlName = twoSumTask.UrlName
		err := svc.UpdateTask(task)
		asserts.EqualError(t, err, ErrUniqueUrlName)
	})

	t.Run("invalid testcase count", func(t *testing.T) {
		task := updateTask
		task.TestcaseCount = -10
		err := svc.UpdateTask(task)
		asserts.EqualError(t, err, ErrInvalidTestcaseCount)
	})

	t.Run("normal update task", func(t *testing.T) {
		err := svc.UpdateTask(updateTask)
		asserts.EqualError(t, err, nil)

		task, err := svc.GetTask(updateTask.Id)
		asserts.EqualError(t, err, nil)
		asserts.Equal(t, "task", task, updateTask)
	})
}

func TestDeleteTask(t *testing.T) {
	t.Run("id not found", func(t *testing.T) {
		err := svc.DeleteTask(1000)
		asserts.EqualError(t, err, ErrTaskNotFound)
	})

	t.Run("normal delete task", func(t *testing.T) {
		err := svc.DeleteTask(updateTask.Id)
		asserts.EqualError(t, err, nil)
	})
}
