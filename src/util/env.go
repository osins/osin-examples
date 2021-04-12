package util

import (
	"fmt"

	"github.com/joho/godotenv"
)

func NewEnv() Env {
	return &env{}
}

type Env interface {
	Load(envfile string) error
}

type env struct {
}

func (s *env) Load(envfile string) error {
	fmt.Printf("init env start. path: %s\n", envfile)
	err := godotenv.Load(envfile)
	if err != nil {
		fmt.Printf("Error loading .env file: " + err.Error())
		return err
	}

	fmt.Printf("init env complete.")

	return nil
}
