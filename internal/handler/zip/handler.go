package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
	"io"
	"log"
	"sync"
	"time"
)

type ZipperService struct {
}

func (s *ZipperService) Compress(count int, target *bytes.Buffer, files chan *filer.File) {
	fmt.Println("start compressor")
	w := zip.NewWriter(target)
	var wg sync.WaitGroup
	var mu sync.Mutex
	defer w.Close()

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(begin time.Time) {
			defer wg.Done()
			defer mu.Unlock()
			file := <-files
			size := file.Content.Len()
			mu.Lock()
			f, err := w.Create(file.Name)
			if err != nil {
				log.Fatal(err)
			}
			if _, err = io.Copy(f, file.Content); err != nil {
				log.Fatal(err)
			}
			defer func() {
				str := fmt.Sprintf("compress %s finish, took %s, size: %d KB", file.Name, time.Since(begin), size/1024)
				fmt.Println(str)
			}()
		}(time.Time{})
	}

	wg.Wait()
	fmt.Println("done compressor")
}
