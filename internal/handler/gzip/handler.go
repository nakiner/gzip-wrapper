package gzip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
)

type TarballService struct {
}

func (s *TarballService) Compress(ch chan *bytes.Buffer) (*bytes.Buffer, error) {
	// tar > gzip > buf
	archive := bytes.Buffer{}
	zr := gzip.NewWriter(&archive)
	tw := tar.NewWriter(zr)

	// produce tar
	if err := tw.Close(); err != nil {
		return nil, err
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return nil, err
	}
	//

	return &archive, nil
}
