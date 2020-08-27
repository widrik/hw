package config

import (
	"errors"
)

var (
	ErrFilePathEmpty = errors.New("file path is empty")
	ErrReadFile      = errors.New("can't read file")
)

type Logging struct {
	File  string
	Level string
}

type GRPCServer struct {
	Host string
	Port string
}
