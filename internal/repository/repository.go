package repository

import (
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
)

type service struct {
}

type Repository interface {
	GenerateFile(file *filer.File, size int) error
}
