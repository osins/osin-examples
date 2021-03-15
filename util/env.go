package util

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

var (
	_, f, _, _ = runtime.Caller(0)
	BasePATH   = filepath.Dir(f)
	ENVFile    = BasePATH + "/../.env"
)

func NewEnv() Env {
	return &env{}
}

type Env interface {
	Load() error
}

type env struct {
}

func (s *env) Load() error {
	fmt.Printf("init env start. path: %s\n", ENVFile)
	err := godotenv.Load(ENVFile)
	if err != nil {
		fmt.Printf("Error loading .env file: " + err.Error())
		return err
	}

	fmt.Printf("init env complete.")

	return nil
}
