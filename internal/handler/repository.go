package handler

import (
	"bytes"
	zp "github.com/nakiner/gzip-wrapper/internal/handler/zip"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
)

type Archive struct {
	FileCount int
}

type service struct {
}

type Repository interface {
	Zipper() Compressor
	//Tarball() Compressor
}

type Compressor interface {
	Compress(target *bytes.Buffer, files chan *filer.File)
}

func NewCompressor() Repository {
	return &service{}
}

func (s *service) Zipper() Compressor {
	return &zp.ZipperService{}
}

//func (s *service) Tarball() Compressor {
//	return &gzip.TarballService{}
//}
