package model

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

var result int

/******** Benchmark access to a map with explicit mutex ********/

type MutexMap struct {
	rwm *sync.RWMutex
	m   map[string]int
}

func (mm *MutexMap) get(k string) (i int) {
	mm.rwm.RLock()
	i = mm.m[k]
	mm.rwm.RUnlock()
	return
}

func benchmarkMutexMap(size int, b *testing.B) {
	mm := &MutexMap{
		rwm: &sync.RWMutex{},
		m:   make(map[string]int, size),
	}
	for i := 0; i < size; i++ {
		mm.m[strconv.Itoa(i)] = i
	}

	for n := 0; n < b.N; n++ {
		r := rand.Intn(size)
		result = mm.get(strconv.Itoa(r))
	}
}

func BenchmarkMutexMap10(b *testing.B)   { benchmarkMutexMap(9, b) }
func BenchmarkMutexMap50(b *testing.B)   { benchmarkMutexMap(49, b) }
func BenchmarkMutexMap100(b *testing.B)  { benchmarkMutexMap(99, b) }
func BenchmarkMutexMap1000(b *testing.B) { benchmarkMutexMap(999, b) }

/******** Benchmark access to a sync map ********/

func benchmarkSyncMap(size int, b *testing.B) {
	m := &sync.Map{}
	for i := 0; i < size; i++ {
		m.Store(strconv.Itoa(i), i)
	}

	for n := 0; n < b.N; n++ {
		r := rand.Intn(size)
		temp, _ := m.Load(strconv.Itoa(r))
		result = temp.(int)
	}
}

func BenchmarkSyncMap10(b *testing.B)   { benchmarkSyncMap(9, b) }
func BenchmarkSyncMap50(b *testing.B)   { benchmarkSyncMap(49, b) }
func BenchmarkSyncMap100(b *testing.B)  { benchmarkSyncMap(99, b) }
func BenchmarkSyncMap1000(b *testing.B) { benchmarkSyncMap(999, b) }

/******** Benchmark cases ********/

func check(s string) int {
	switch s {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	}
	return 0
}

func benchmarkCase(b *testing.B) {
	for n := 0; n < b.N; n++ {
		r := rand.Intn(10)
		result = check(strconv.Itoa(r))
	}
}

func BenchmarkCase(b *testing.B) { benchmarkCase(b) }
