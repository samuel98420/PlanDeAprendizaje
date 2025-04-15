package config

import (
	"path/filepath"
	"runtime"
)

func GetDBPath() string {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	return filepath.Join(projectRoot, "internal", "database", "tasks.db")
}
