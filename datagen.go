package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const yamlFormatDescription = `primaryKeyFileName:  The name of the primary key file
directory:           The path to the output directory
zipFileName:         The name of the zip file write in the output directory
numberOfEntities:    The total number of entities
totalTimeInHours:    The total number of hours of data to be output
files:               The list of output files
   - fileName:       The name of the file to output
     dataType:       One of "timeseries" or "static"
     values:         A list of comma separated values (static only) (eg: "[ one, two, three ]")
     valueType:      One of "double" or "long" (timeseries only)
     timeStepMillis: The number of approximate milliseconds between datum (timeseries only)
     minValue:       The minimum value range when generating data (timeseries only)
     maxValue:       The maximum value range when generating data (timeseries only)
     timeVariance:   The percentage of variance of the time steps (timeseries only)
     valueVariance:  The percentage of variance of the value (timeseries only)`

var yamlFile string
var outputYamlFormat bool

type dataGenerator interface {
	generateData(channel chan string)
	isStatic() bool
	getFilePath() string
}

func init() {
	const usage = "The YAML configuration file"
	flag.StringVar(&yamlFile, "yaml-file", "", usage)
	flag.StringVar(&yamlFile, "f", "", usage+" (shorthand)")
	flag.BoolVar(&outputYamlFormat, "h", false, "Describes the input YAML format")
}

func main() {
	flag.Parse()

	if outputYamlFormat {
		flag.Usage()
		fmt.Printf("The YAML input file:\n\n")
		fmt.Println(yamlFormatDescription)
		fmt.Printf("\n\n")
		os.Exit(0)
	}

	if len(yamlFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	start := time.Now()
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		panic(err)
	}
	dataset, err := NewDatasetFromYAML(string(data))

	if err != nil {
		panic(err)
	}
	err = writeDataFiles(dataset)
	if err != nil {
		panic(err)
	}

	elapsedTime := time.Since(start)
	fmt.Printf("Total Run Time: %s\n", elapsedTime)
}
