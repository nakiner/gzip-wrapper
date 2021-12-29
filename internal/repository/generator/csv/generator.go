package csv

import (
	"github.com/nakiner/gzip-wrapper/internal/repository"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type service struct{}

func NewGenerator() repository.Repository {
	return &service{}
}

func (s *service) GenerateFile(w *filer.File) error {
	size := 1 * 1024 // 1 mB
	for i := 0; i < size; i++ {
		if _, err := w.Content.Write(RandStringBytes(1024)); err != nil {
			return err
		}
	}
	return nil
}

func RandStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
