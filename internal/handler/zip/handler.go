package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
	"log"
	"sync"
)

type ZipperService struct {
}

func (s *ZipperService) Compress(target *bytes.Buffer, files chan *filer.File) {
	fmt.Println("start compressor")
	w := zip.NewWriter(target)
	var wg sync.WaitGroup
	defer w.Close()

	count := len(files)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case file := <-files:
				fmt.Println("start compress file")
				f, err := w.Create(file.Name)
				if err != nil {
					log.Fatal(err)
				}
				_, err = f.Write(file.Content.Bytes())
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("done compress file")
				count--
				if count < 1 {
					return
				}
			}
		}
	}()

	wg.Wait()
	fmt.Println("done compressor")
}
