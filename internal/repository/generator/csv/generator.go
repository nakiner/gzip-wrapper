package csv

import (
	"github.com/nakiner/gzip-wrapper/internal/repository"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
	"io"
	"math/rand"
	"os"
	"path/filepath"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type service struct{}

func NewGenerator() repository.Repository {
	return &service{}
}

func (s *service) GenerateFile(w *filer.File, size int) error {
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

func (s *service) ReadFile(w *filer.File) error {
	data, err := os.Open(w.Name)
	defer data.Close()
	if err != nil {
		return err
	}
	if _, err = io.Copy(w.Content, data); err != nil {
		return err
	}
	return nil
}

func (s *service) ListFiles(dir string) []string {
	var res []string
	filepath.Walk(dir, func(file string, fi os.FileInfo, _ error) error {
		if !fi.IsDir() {
			res = append(res, filepath.FromSlash(file))
		}
		return nil
	})

	return res
}
