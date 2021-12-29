package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
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
		go func() {
			defer wg.Done()
			defer mu.Unlock()
			file := <-files
			begin := time.Now()
			size := file.Content.Len()
			mu.Lock()
			f, err := w.Create(file.Name)
			if err != nil {
				log.Fatal(err)
			}
			if _, err = file.Content.WriteTo(f); err != nil {
				log.Fatal(err)
			}
			defer func() {
				str := fmt.Sprintf("compress %s finish, took %s, size: %d KB", file.Name, time.Since(begin), size/1024)
				fmt.Println(str)
			}()
		}()
	}

	wg.Wait()
	fmt.Println("done compressor")
}
