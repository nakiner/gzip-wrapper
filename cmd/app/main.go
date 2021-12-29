package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/nakiner/gzip-wrapper/configs"
	"github.com/nakiner/gzip-wrapper/internal/handler"
	"github.com/nakiner/gzip-wrapper/internal/repository/filer"
	"github.com/nakiner/gzip-wrapper/internal/repository/generator/csv"
	"github.com/nakiner/gzip-wrapper/internal/store/s3"
	"github.com/nakiner/gzip-wrapper/tools/logging"
	"os"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg := configs.NewConfig()
	if err := cfg.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}

	// Print config
	if err := cfg.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %s", err)
		os.Exit(1)
	}
	ctx = logging.WithContext(ctx, logger)

	s3Client, err := s3.NewClient(ctx, &cfg.S3)
	if err != nil {
		level.Error(logger).Log("err", err, "msg", "failed to s3 client")
	}

	s3Repository := s3.NewRepository(s3Client)
	gen := csv.NewGenerator()
	zipper := handler.NewCompressor().Zipper()

	count := 10
	var wg sync.WaitGroup
	byteCh := make(chan *filer.File, count)
	byteGenCh := make(chan *filer.File, count)
	archive := bytes.Buffer{}

	//init concurrent file generator
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			byter := <-byteCh
			if err = gen.GenerateFile(byter); err != nil {
				level.Error(logger).Log("err", err)
			} else {
				byteGenCh <- byter
			}
		}()
	}

	// push files
	for i := 0; i < count; i++ {
		buf := filer.Buffer{}
		file := filer.File{
			Name:    fmt.Sprintf("text_%d.txt", i),
			Content: &buf,
		}
		byteCh <- &file
	}

	wg.Wait()
	close(byteCh)

	zipper.Compress(&archive, byteGenCh)
	close(byteGenCh)

	// init archive uploader
	if err = s3Repository.Upload(ctx, &archive, "archive.zip"); err != nil {
		level.Error(logger).Log("err", err)
	}

	os.Exit(0)
}
