package main

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (t timeSeriesFile) generateData(channel chan string) {
	num := int(t.numberOfEntities)
	now := time.Now().UTC()
	startTime := t.getStartTime()
	ctime := startTime
	value := t.getStartValue()

	tsMin, tsMax := t.getTimeStepRange()
	vsMin, vsMax := t.getValueStepRange()
	for entityID := 1; entityID <= num; entityID++ {
		for ctime.Before(now) {
			millisToAdd := randInt(tsMin, tsMax)
			valToAdd := randFloat(vsMin, vsMax)
			ctime = ctime.Add(time.Duration(millisToAdd) * time.Millisecond)
			value = t.toTypedValue(value + valToAdd)
			channel <- formatJointString(entityID, ctime, value)
		}
		ctime = startTime
	}
	close(channel)
}

func (t timeSeriesFile) getStartTime() time.Time {
	now := time.Now()
	dur := time.Duration(-1 * time.Duration(t.totalTimeInHours) * time.Hour)
	return now.Add(dur).UTC()
}

func (t timeSeriesFile) getStartValue() float64 {
	return t.toTypedValue(t.minimum + ((t.maximum - t.minimum) / 2.0))
}

func (t timeSeriesFile) toTypedValue(value float64) float64 {
	if value > t.maximum {
		value = t.maximum
	}
	if value < t.minimum {
		value = t.minimum
	}
	if t.valueType == "long" {
		return math.Trunc(value)
	}
	return value
}

func (t timeSeriesFile) getTimeStepRange() (int, int) {
	timeStep := float64(t.timeStepInMilliseconds)
	tsMin := math.Min(1, timeStep-(t.timeStepVariance*timeStep))
	tsMax := timeStep + (t.timeStepVariance * timeStep)
	return int(tsMin), int(tsMax)
}

func (t timeSeriesFile) getValueStepRange() (float64, float64) {
	vrange := t.maximum - t.minimum
	variance := (vrange * t.valueVariance) / 2.0
	return 0 - (variance), variance
}

func formatJointString(eid int, t time.Time, value float64) string {
	return strings.Join([]string{strconv.Itoa(eid), formatISO8601(t), strconv.FormatFloat(value, 'f', 3, 64)}, ",")
}

func formatISO8601(t time.Time) string {
	return t.Format(time.RFC3339Nano)[:23] + "Z"
}

func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func randFloat(min, max float64) float64 {
	return min + (rand.Float64() * (max - min))
}

func (t timeSeriesFile) isStatic() bool {
	return false
}

func (t timeSeriesFile) getFilePath() string {
	return t.filePath
}
