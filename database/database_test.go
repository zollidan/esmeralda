package database

import (
	"testing"

	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/models"
)

func TestInitDatabase(t *testing.T) {
	cfg := config.Config{
		DBURL: ":memory:",
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("InitDatabase panicked: %v", r)
		}
	}()

	database := InitDatabase(cfg)
	if database == nil {
		t.Fatal("database is nil after initialization")
	}

	for _, model := range []interface{}{&models.Files{}, &models.Task{}, &models.Parser{}, &models.User{}} {
		if !database.Migrator().HasTable(model) {
			t.Errorf("expected table for %T to exist", model)
		}
	}
}
