package repository

import (
	"context"
	"testing"
	"time"

	"github.com/zollidan/esmeralda/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&models.Task{}, &models.Game{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

// --- Generic Repository ---

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := New[models.Task](db)
	ctx := context.Background()

	task := &models.Task{ID: "t1", Date: "2025-01-01", Status: "pending", CreatedAt: time.Now()}
	if err := repo.Create(ctx, task); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := repo.FindByID(ctx, "t1")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if found.ID != "t1" {
		t.Errorf("ID = %q, want %q", found.ID, "t1")
	}
}

func TestRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := New[models.Task](db)
	ctx := context.Background()

	db.Create(&models.Task{ID: "t1", Date: "2025-01-01", Status: "pending"})
	db.Create(&models.Task{ID: "t2", Date: "2025-01-02", Status: "done"})

	tasks, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("len = %d, want 2", len(tasks))
	}
}

func TestRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := New[models.Task](db)
	ctx := context.Background()

	_, err := repo.FindByID(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing record")
	}
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := New[models.Task](db)
	ctx := context.Background()

	task := &models.Task{ID: "t1", Date: "2025-01-01", Status: "pending"}
	db.Create(task)

	task.Status = "done"
	if err := repo.Update(ctx, task); err != nil {
		t.Fatalf("update: %v", err)
	}

	found, _ := repo.FindByID(ctx, "t1")
	if found.Status != "done" {
		t.Errorf("Status = %q, want %q", found.Status, "done")
	}
}

func TestRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := New[models.Task](db)
	ctx := context.Background()

	db.Create(&models.Task{ID: "t1", Date: "2025-01-01", Status: "pending"})

	if err := repo.Delete(ctx, "t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := repo.FindByID(ctx, "t1")
	if err == nil {
		t.Fatal("expected record to be deleted")
	}
}

func TestRepository_CreateInBatches(t *testing.T) {
	db := setupTestDB(t)
	repo := New[models.Task](db)
	ctx := context.Background()

	tasks := []models.Task{
		{ID: "t1", Date: "2025-01-01", Status: "pending"},
		{ID: "t2", Date: "2025-01-02", Status: "pending"},
		{ID: "t3", Date: "2025-01-03", Status: "pending"},
	}

	if err := repo.CreateInBatches(ctx, tasks, 2); err != nil {
		t.Fatalf("create in batches: %v", err)
	}

	all, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("find all: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("len = %d, want 3", len(all))
	}
}

func TestRepository_DB(t *testing.T) {
	db := setupTestDB(t)
	repo := New[models.Task](db)
	if repo.DB() == nil {
		t.Fatal("DB() returned nil")
	}
}

// --- TaskRepository ---

func TestTaskRepository_GetAllOrdered(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()

	now := time.Now()
	db.Create(&models.Task{ID: "t1", Date: "2025-01-01", Status: "pending", CreatedAt: now.Add(-2 * time.Hour)})
	db.Create(&models.Task{ID: "t2", Date: "2025-01-02", Status: "done", CreatedAt: now.Add(-1 * time.Hour)})
	db.Create(&models.Task{ID: "t3", Date: "2025-01-03", Status: "pending", CreatedAt: now})

	tasks, err := repo.GetAllOrdered(ctx)
	if err != nil {
		t.Fatalf("get all ordered: %v", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("len = %d, want 3", len(tasks))
	}

	// Should be ordered desc by created_at: t3, t2, t1
	if tasks[0].ID != "t3" {
		t.Errorf("first task = %q, want t3", tasks[0].ID)
	}
	if tasks[1].ID != "t2" {
		t.Errorf("second task = %q, want t2", tasks[1].ID)
	}
	if tasks[2].ID != "t1" {
		t.Errorf("third task = %q, want t1", tasks[2].ID)
	}
}

func TestTaskRepository_GetAllOrdered_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()

	tasks, err := repo.GetAllOrdered(ctx)
	if err != nil {
		t.Fatalf("get all ordered: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("len = %d, want 0", len(tasks))
	}
}

func TestTaskRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)
	ctx := context.Background()

	db.Create(&models.Task{ID: "t1", Date: "2025-01-01", Status: "pending"})

	if err := repo.UpdateStatus(ctx, "t1", "done"); err != nil {
		t.Fatalf("update status: %v", err)
	}

	var task models.Task
	db.First(&task, "id = ?", "t1")
	if task.Status != "done" {
		t.Errorf("Status = %q, want %q", task.Status, "done")
	}
}

// --- GameRepository ---

func TestGameRepository_FindByTaskID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)
	ctx := context.Background()

	taskID := "task-1"
	otherTaskID := "task-2"

	db.Create(&models.Game{TaskID: &taskID, Day: 1, Month: 1, Year: 2025, HomeTeam: "A", AwayTeam: "B", League: "L1"})
	db.Create(&models.Game{TaskID: &taskID, Day: 2, Month: 1, Year: 2025, HomeTeam: "C", AwayTeam: "D", League: "L1"})
	db.Create(&models.Game{TaskID: &otherTaskID, Day: 3, Month: 1, Year: 2025, HomeTeam: "E", AwayTeam: "F", League: "L2"})

	games, err := repo.FindByTaskID(ctx, "task-1")
	if err != nil {
		t.Fatalf("find by task id: %v", err)
	}
	if len(games) != 2 {
		t.Errorf("len = %d, want 2", len(games))
	}
}

func TestGameRepository_FindByTaskID_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)
	ctx := context.Background()

	games, err := repo.FindByTaskID(ctx, "nonexistent")
	if err != nil {
		t.Fatalf("find by task id: %v", err)
	}
	if len(games) != 0 {
		t.Errorf("len = %d, want 0", len(games))
	}
}

func TestGameRepository_FindByDateRange(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)
	ctx := context.Background()

	games := []models.Game{
		{Day: 10, Month: 1, Year: 2025, HomeTeam: "A", AwayTeam: "B", League: "L"},
		{Day: 15, Month: 1, Year: 2025, HomeTeam: "C", AwayTeam: "D", League: "L"},
		{Day: 20, Month: 1, Year: 2025, HomeTeam: "E", AwayTeam: "F", League: "L"},
		{Day: 5, Month: 2, Year: 2025, HomeTeam: "G", AwayTeam: "H", League: "L"},
	}
	for i := range games {
		db.Create(&games[i])
	}

	// Range: Jan 10 - Jan 20
	result, err := repo.FindByDateRange(ctx, 2025, 1, 10, 2025, 1, 20)
	if err != nil {
		t.Fatalf("find by date range: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("len = %d, want 3", len(result))
	}
}

func TestGameRepository_FindByDateRange_SingleDay(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)
	ctx := context.Background()

	db.Create(&models.Game{Day: 15, Month: 3, Year: 2025, HomeTeam: "A", AwayTeam: "B", League: "L"})
	db.Create(&models.Game{Day: 16, Month: 3, Year: 2025, HomeTeam: "C", AwayTeam: "D", League: "L"})

	result, err := repo.FindByDateRange(ctx, 2025, 3, 15, 2025, 3, 15)
	if err != nil {
		t.Fatalf("find by date range: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("len = %d, want 1", len(result))
	}
}

func TestGameRepository_FindByDateRange_CrossMonth(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)
	ctx := context.Background()

	db.Create(&models.Game{Day: 28, Month: 1, Year: 2025, HomeTeam: "A", AwayTeam: "B", League: "L"})
	db.Create(&models.Game{Day: 5, Month: 2, Year: 2025, HomeTeam: "C", AwayTeam: "D", League: "L"})
	db.Create(&models.Game{Day: 15, Month: 2, Year: 2025, HomeTeam: "E", AwayTeam: "F", League: "L"})

	// Range: Jan 28 - Feb 5
	result, err := repo.FindByDateRange(ctx, 2025, 1, 28, 2025, 2, 5)
	if err != nil {
		t.Fatalf("find by date range: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("len = %d, want 2", len(result))
	}
}

func TestGameRepository_FindByDateRange_OrderedByDateTime(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGameRepository(db)
	ctx := context.Background()

	db.Create(&models.Game{Day: 2, Month: 1, Year: 2025, Time: "20:00", HomeTeam: "C", AwayTeam: "D", League: "L"})
	db.Create(&models.Game{Day: 1, Month: 1, Year: 2025, Time: "18:00", HomeTeam: "A", AwayTeam: "B", League: "L"})
	db.Create(&models.Game{Day: 1, Month: 1, Year: 2025, Time: "15:00", HomeTeam: "E", AwayTeam: "F", League: "L"})

	result, err := repo.FindByDateRange(ctx, 2025, 1, 1, 2025, 1, 2)
	if err != nil {
		t.Fatalf("find by date range: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}

	// Should be ordered: year, month, day, time
	if result[0].Time != "15:00" {
		t.Errorf("first game time = %q, want 15:00", result[0].Time)
	}
	if result[1].Time != "18:00" {
		t.Errorf("second game time = %q, want 18:00", result[1].Time)
	}
	if result[2].Time != "20:00" {
		t.Errorf("third game time = %q, want 20:00", result[2].Time)
	}
}
