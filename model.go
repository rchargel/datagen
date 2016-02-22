package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// Dataset is a wrapper around a configuration for a dataset export
type Dataset struct {
	Directory          string   `yaml:"directory"`
	ZipFileName        string   `yaml:"zipFileName"`
	PrimaryKeyFileName string   `yaml:"pkFileName"`
	NumberOfEntities   uint16   `yaml:"numberOfEntities"`
	TotalTimeInHours   uint16   `yaml:"totalTimeInHours"`
	Files              []DsFile `yaml:"files"`
}

// DsFile is a wrapper around the configuration necessary to generate a file.
type DsFile struct {
	FileName               string   `yaml:"fileName"`
	DataType               string   `yaml:"dataType"`
	ValueType              string   `yaml:"valueType"`
	TimeStepInMilliseconds uint32   `yaml:"timeStepMillis"`
	PossibleValues         []string `yaml:"values"`
	Minimum                float64  `yaml:"minValue"`
	Maximum                float64  `yaml:"maxValue"`
	TimeStepVariance       float64  `yaml:"timeVariance"`
	ValueVariance          float64  `yaml:"valueVariance"`
}

type primaryKeyFile struct {
	filePath         string
	numberOfEntities uint16
}

type staticFile struct {
	filePath         string
	numberOfEntities uint16
	possibleValues   []string
}

type timeSeriesFile struct {
	filePath               string
	numberOfEntities       uint16
	valueType              string
	totalTimeInHours       uint16
	timeStepInMilliseconds uint32
	minimum                float64
	maximum                float64
	timeStepVariance       float64
	valueVariance          float64
}

// NewDatasetFromYAML creates a new dataset configuration from a yaml file.
func NewDatasetFromYAML(data string) (Dataset, error) {
	dataset := Dataset{}
	err := yaml.Unmarshal([]byte(data), &dataset)
	return dataset, err
}

// String the toString() method for the dataset.
func (d Dataset) String() string {
	return fmt.Sprintf(`Dataset {
	Directory:          "%v"
	PrimaryKeyFileName: "%v"
	NumberOfEntities:   "%v"
	TotalTimeInHours:   "%v"
	Files:              %v
}`, d.Directory, d.PrimaryKeyFileName, d.NumberOfEntities, d.TotalTimeInHours, d.Files)
}

// String the toString() method for the dataset file.
func (f DsFile) String() string {
	return fmt.Sprintf(`
		File {
			FileName:               "%v"
			DataType:               "%v"
			ValueType:              "%v"
			TimeStepInMilliseconds: "%v"
			PossibleValues:         %v
			Minimum:                "%v"
			Maximum:                "%v"
			TimeStepVariance:       "%v"
			ValueVariance:          "%v"
		}`, f.FileName, f.DataType, f.ValueType, f.TimeStepInMilliseconds, f.PossibleValues, f.Minimum, f.Maximum, f.TimeStepVariance, f.ValueVariance)
}

func (d Dataset) getFiles() []dataGenerator {
	l := make([]dataGenerator, len(d.Files)+1, len(d.Files)+1)

	l[0] = primaryKeyFile{
		filePath:         d.PrimaryKeyFileName,
		numberOfEntities: d.NumberOfEntities,
	}
	for i, f := range d.Files {
		l[i+1] = f.toGeneratedFile(d)
	}
	return l
}

func (f DsFile) toGeneratedFile(d Dataset) dataGenerator {
	fp := f.FileName

	if f.DataType == "static" {
		return staticFile{
			filePath:         fp,
			numberOfEntities: d.NumberOfEntities,
			possibleValues:   f.PossibleValues,
		}
	}
	return timeSeriesFile{
		filePath:               fp,
		numberOfEntities:       d.NumberOfEntities,
		valueType:              f.ValueType,
		totalTimeInHours:       d.TotalTimeInHours,
		timeStepInMilliseconds: f.TimeStepInMilliseconds,
		minimum:                f.Minimum,
		maximum:                f.Maximum,
		timeStepVariance:       f.TimeStepVariance,
		valueVariance:          f.ValueVariance,
	}
}
