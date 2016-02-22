package main

import (
	"fmt"
	"os"
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
	numFiles := len(generators)

	doneChan := make(chan string)

	for _, gen := range generators {
		go writeDataToFile(gen, doneChan)
	}

	c := 0
	for reply := range doneChan {
		c++
		if len(reply) > 0 {
			fmt.Println(reply)
		}
		if c == numFiles {
			close(doneChan)
			break
		}
	}

	return nil
}

func writeDataToFile(gen dataGenerator, doneChan chan string) error {
	start := time.Now()
	filepath := gen.getFilePath()
	if exists(filepath) {
		fmt.Printf("File %v already exists, deleting old file\n", filepath)
		err := os.Remove(filepath)
		if err != nil {
			doneChan <- fmt.Sprintf("Error writing file %v: %v", filepath, err.Error())
			return err
		}
	}
	os.Remove(filepath)

	dataChan := make(chan string)
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		doneChan <- err.Error()
		return err
	}

	go gen.generateData(dataChan)

	for line := range dataChan {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			doneChan <- fmt.Sprintf("Error writing file %v: %v", filepath, err.Error())
			return err
		}
	}
	fmt.Printf("Finished writing file %v in %v\n", filepath, time.Since(start))
	doneChan <- ""
	return nil
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
