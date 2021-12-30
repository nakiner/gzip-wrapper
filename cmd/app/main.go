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
	"runtime/pprof"
	"sync"
	"time"
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
	//zipper := handler.NewCompressor().Zipper()
	tar := handler.NewCompressor().Tarball()

	count := 0
	var wg sync.WaitGroup
	archive := bytes.Buffer{}

	// push files
	files := gen.ListFiles("./export")
	count = len(files)
	byteCh := make(chan *filer.File, count)
	byteGenCh := make(chan *filer.File, count)
	for _, filename := range files {
		buf := filer.Buffer{}
		file := filer.File{
			Name:    filename,
			Content: &buf,
		}
		byteCh <- &file
	}

	var genWg sync.WaitGroup
	//init concurrent file generator
	for i := 0; i < count; i++ {
		wg.Add(1)
		genWg.Add(1)
		go func() {
			defer wg.Done()
			defer genWg.Done()
			byter := <-byteCh
			begin := time.Now()
			if err = gen.ReadFile(byter); err != nil {
				level.Error(logger).Log("err", err)
			} else {
				str := fmt.Sprintf("read %s finish, took %s, size: %d KB", byter.Name, time.Since(begin), byter.Content.Len()/1024)
				fmt.Println(str)
				byteGenCh <- byter
			}
		}()
	}

	go func() {
		defer close(byteCh)
		genWg.Wait()
		pr, _ := os.Create("read.pprof")
		pprof.WriteHeapProfile(pr)
		pr.Close()
	}()

	wg.Add(1)
	compr := make(chan bool, 1)
	go func() {
		defer wg.Done()
		defer close(byteGenCh)
		begin := time.Now()
		tar.Compress(count, &archive, byteGenCh)
		compr <- true
		defer func() {
			str := fmt.Sprintf("compressor finish, took %s, size: %d KB", time.Since(begin), archive.Len()/1024)
			fmt.Println(str)
			pc, _ := os.Create("compress.pprof")
			pprof.WriteHeapProfile(pc)
			pc.Close()
		}()
	}()

	// init archive uploader
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-compr
		size := archive.Len()
		begin := time.Now()
		if err = s3Repository.Upload(ctx, &archive, "archive.tar.gz"); err != nil {
			level.Error(logger).Log("err", err)
		}
		defer func() {
			str := fmt.Sprintf("upload finish, took %s, size: %d KB", time.Since(begin), size/1024)
			fmt.Println(str)
			pup, _ := os.Create("upload.pprof")
			pprof.WriteHeapProfile(pup)
			pup.Close()
		}()
	}()

	wg.Wait()
	os.Exit(0)
}
