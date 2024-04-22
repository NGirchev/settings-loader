package util

import (
	"github.com/spf13/viper"
	"path/filepath"
	"runtime"
)

func GetPath(path string) string {
	return filepath.Clean(filepath.Join(filepath.Dir(GetCurrentFilePath()), "../../", path))
}

func GetCurrentFilePath() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

func GetVar(path string) string {
	return ReplaceEnvVars(viper.GetString(path))
}
