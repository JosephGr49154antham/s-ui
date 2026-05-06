package database

import (
	"os"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once sync.Once
	db   *gorm.DB
)

// GetDB returns the singleton database instance.
func GetDB() *gorm.DB {
	return db
}

// InitDB initialises the SQLite database at the given path,
// creating parent directories as needed.
func InitDB(dbPath string) error {
	var initErr error
	once.Do(func() {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o750); err != nil {
			initErr = err
			return
		}
		conn, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			initErr = err
			return
		}
		if err := conn.AutoMigrate(&User{}); err != nil {
			initErr = err
			return
		}
		db = conn
	})
	return initErr
}

// resetSingleton is used in tests to reset the package-level state.
func resetSingleton() {
	once = sync.Once{}
	db = nil
}
