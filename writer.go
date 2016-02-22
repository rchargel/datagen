package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

func writeDataFiles(d Dataset) error {
	if !exists(d.Directory) {
		fmt.Printf("Creating directory %v\n", d.Directory)
		err := os.Mkdir(d.Directory, os.ModePerm)
		if err != nil {
			return err
		}
	}

	generators := d.getFiles()

	zipChan := make(chan *os.File)

	wait := compressFiles(d.Directory, d.ZipFileName, zipChan)

	var wg sync.WaitGroup
	wg.Add(len(generators))
	for _, generator := range generators {
		go func(gen dataGenerator, dir string) {
			defer wg.Done()
			file, err := writeDataToFile(gen, dir)
			if err != nil {
				fmt.Printf("ERROR Writing %v: %v!\n", file, err.Error())
				panic(err)
			}
			zipChan <- file
		}(generator, d.Directory)
	}

	wg.Wait()
	close(zipChan)

	wait.Wait()
	return nil
}

func compressFiles(directory, zipFileName string, zipChan chan *os.File) *sync.WaitGroup {
	zipFilePath := path.Join(directory, zipFileName)
	if exists(zipFilePath) {
		fmt.Printf("File %v already exists, deleting file...\n", zipFileName)
		os.Remove(zipFilePath)
	}
	fmt.Printf("Creating file %v\n", zipFileName)

	zipFile, _ := os.Create(zipFilePath)
	var wg sync.WaitGroup
	wg.Add(1)

	zw := zip.NewWriter(zipFile)

	go func() {
		// Note the order (LIFO)
		defer wg.Done()       // 2. signal that we're done
		defer zipFile.Close() // 1. close the file

		var err error
		var fw io.Writer
		for file := range zipChan {
			// Loop until the channel is closed
			startTime := time.Now()
			zipName := "data/" + file.Name()[len(directory)+1:]
			if fw, err = zw.Create(zipName); err != nil {
				panic(err)
			}
			io.Copy(fw, file)
			if err = file.Close(); err != nil {
				panic(err)
			}
			os.Remove(file.Name())
			fmt.Printf("Compressed file %v in %v\n", file.Name(), time.Since(startTime))
		}
		zw.Close()
	}()
	return &wg
}

func writeDataToFile(gen dataGenerator, dir string) (*os.File, error) {
	start := time.Now()
	filepath := path.Join(dir, gen.getFilePath())
	var file *os.File
	var err error
	if exists(filepath) {
		fmt.Printf("File %v already exists, deleting old file\n", filepath)
		err = os.Remove(filepath)
		if err != nil {
			return file, err
		}
	}
	os.Remove(filepath)

	dataChan := make(chan string)
	file, err = os.Create(filepath)
	defer file.Close()
	if err != nil {
		return file, err
	}

	go gen.generateData(dataChan)

	for line := range dataChan {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			return file, err
		}
	}
	fmt.Printf("Finished writing file %v in %v\n", filepath, time.Since(start))
	return os.Open(filepath)
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
