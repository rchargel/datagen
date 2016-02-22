package main

import "testing"

func BenchmarkTimeSeriesGenerateData(b *testing.B) {
	f := timeSeriesFile{
		totalTimeInHours:       1,
		timeStepInMilliseconds: 15000,
		minimum:                1000,
		maximum:                1300,
		valueType:              "double",
		numberOfEntities:       3,
		timeStepVariance:       1.5,
		valueVariance:          0.05,
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := make(chan string)
			go f.generateData(c)

			for _ = range c {

			}
		}
	})

}

func BenchmarkTimeSeriesGetStartTime(b *testing.B) {
	f := timeSeriesFile{
		totalTimeInHours: 6,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f.getStartTime()
		}
	})
}

func BenchmarkTimeSeriesFormat1(b *testing.B) {
	f := timeSeriesFile{
		totalTimeInHours: 6,
	}
	startTime := f.getStartTime()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			formatISO8601(startTime)
		}
	})
}
