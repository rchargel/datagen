package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

var sizeTypes = [4]string{"B", "KB", "MB", "GB"}

func writeDataFiles(d Dataset) error {
	if !exists(d.Directory) {
		log.Printf("Creating directory %v", d.Directory)
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
			log.Printf("Started writing file %v", gen.getFilePath())
			file, err := writeDataToFile(gen, dir)
			if err != nil {
				log.Printf("ERROR Writing %v: %v!", file, err.Error())
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
		log.Printf("File %v already exists", zipFileName)
		log.Printf("Deleting file %v", zipFileName)
		os.Remove(zipFilePath)
	}
	log.Printf("Creating new file %v", zipFileName)

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
			log.Printf("Compressed file %v in %v", file.Name(), time.Since(startTime))
		}
		zw.Close()
		log.Printf("Zip file %v (%v) completed", zipFile.Name(), getFileSize(zipFile))
	}()
	return &wg
}

func writeDataToFile(gen dataGenerator, dir string) (*os.File, error) {
	start := time.Now()
	filepath := path.Join(dir, gen.getFilePath())
	var file *os.File
	var err error
	if exists(filepath) {
		log.Printf("File %v already exists", filepath)
		log.Printf("Deleting file %v", filepath)
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

	file, err = os.Open(filepath)
	log.Printf("Finished writing file %v (%v) in %v\n", filepath, getFileSize(file), time.Since(start))
	return file, err
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getFileSize(file *os.File) string {
	fileInfo, _ := file.Stat()
	size := float64(fileInfo.Size())
	sizeIdx := 0
	for size >= 1024 && sizeIdx < len(sizeTypes) {
		size /= 1024
		sizeIdx++
	}

	intsize := uint32(size * 10)
	size = float64(intsize) / 10.0

	return fmt.Sprintf("%g%v", size, sizeTypes[sizeIdx])
}
