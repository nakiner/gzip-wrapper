package gzip

import (
	"archive/tar"
	"bytes"
	"fmt"
	gzip "github.com/klauspost/pgzip"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
	"log"
	"sync"
	"time"
)

type TarballService struct {
}

func (s *TarballService) Compress(count int, target *bytes.Buffer, files chan *filer.File) {
	zw, _ := gzip.NewWriterLevel(target, gzip.BestCompression)
	tw := tar.NewWriter(zw)
	var wg sync.WaitGroup
	var mu sync.Mutex

	defer zw.Close()
	defer tw.Close()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer mu.Unlock()
			file := <-files
			begin := time.Now()
			size := file.Content.Len()
			header := tar.Header{
				Name: file.Name,
				Size: int64(file.Content.Len()),
			}
			mu.Lock()
			if err := tw.WriteHeader(&header); err != nil {
				log.Fatal(err)
			}
			if _, err := file.Content.WriteTo(tw); err != nil {
				log.Fatal(err)
			}
			defer func() {
				str := fmt.Sprintf("compress %s finish, took %s, size: %d KB", file.Name, time.Since(begin), size/1024)
				fmt.Println(str)
			}()
		}()
	}

	wg.Wait()
}
