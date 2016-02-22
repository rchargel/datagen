package main

import "testing"

func TestStaticGenerateData(t *testing.T) {
	s := staticFile{
		filePath:         "/test/path.csv",
		numberOfEntities: 30,
		possibleValues:   []string{"one", "two", "three", "four"},
	}

	c := make(chan string)

	go s.generateData(c)

	i := 0
	for _ = range c {
		i++
	}
	if uint16(i) == s.numberOfEntities {
		t.Logf("Generated %v items\n", i)
	} else {
		t.Errorf("Generated %v items, not %v\n", i, s.numberOfEntities)
	}
}

func BenchmarkStaticGenerateData(b *testing.B) {
	s := staticFile{
		filePath:         "/test/path.csv",
		numberOfEntities: 30,
		possibleValues:   []string{"one", "two", "three", "four"},
	}

	for i := 0; i < b.N; i++ {
		c := make(chan string)

		go s.generateData(c)

		i := 0
		for _ = range c {
			i++
		}
		if uint16(i) != s.numberOfEntities {
			b.Errorf("Generated %v items, not %v\n", i, s.numberOfEntities)
		}
	}
}

func BenchmarkParallelStaticRandValue(b *testing.B) {
	s := staticFile{
		filePath:         "/test/path.csv",
		numberOfEntities: 30,
		possibleValues:   []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "elevel", "twelve"},
	}
	l := len(s.possibleValues)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.randValue(1, l)
		}
	})
}
