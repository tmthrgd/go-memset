// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package memset

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func testCorrect(t *testing.T, size func(rand *rand.Rand) int, scale float64) {
	if err := quick.Check(func(data []byte, value byte) bool {
		Memset(data, value)

		for i := range data {
			if data[i] == value {
				continue
			}

			t.Errorf("Memset failed, expected %d at %d, got %d", value, i, data[i])
			return false
		}

		return true
	}, &quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			off := rand.Intn(32)

			data := make([]byte, off+size(rand))
			rand.Read(data[off:])
			args[0] = reflect.ValueOf(data[off:])

			args[1] = reflect.ValueOf(byte(rand.Intn(0x100)))
		},

		MaxCountScale: scale,
	}); err != nil {
		t.Error(err)
	}
}

func TestCorrect(t *testing.T) {
	testCorrect(t, func(rand *rand.Rand) int {
		return 1 + rand.Intn(128*1024)
	}, 20)
}

func TestCorrectHuge(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	testCorrect(t, func(rand *rand.Rand) int {
		return 32*1024*1024 + rand.Intn(96*1024*1024)
	}, 0.025)
}

func benchmarkMemset(b *testing.B, l int) {
	data := make([]byte, l)
	rand.Read(data)

	value := byte(rand.Intn(0x100))

	b.SetBytes(int64(l))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Memset(data, value)
	}
}

func BenchmarkMemset_32(b *testing.B) {
	benchmarkMemset(b, 32)
}

func BenchmarkMemset_128(b *testing.B) {
	benchmarkMemset(b, 128)
}

func BenchmarkMemset_1k(b *testing.B) {
	benchmarkMemset(b, 1*1024)
}

func BenchmarkMemset_16k(b *testing.B) {
	benchmarkMemset(b, 16*1024)
}

func BenchmarkMemset_128k(b *testing.B) {
	benchmarkMemset(b, 128*1024)
}

func BenchmarkMemset_1M(b *testing.B) {
	benchmarkMemset(b, 1024*1024)
}

func BenchmarkMemset_16M(b *testing.B) {
	benchmarkMemset(b, 16*1024*1024)
}

func BenchmarkMemset_128M(b *testing.B) {
	benchmarkMemset(b, 128*1024*1024)
}

func BenchmarkMemset_512M(b *testing.B) {
	benchmarkMemset(b, 512*1024*1024)
}

func benchmarkGoZero(b *testing.B, l int) {
	data := make([]byte, l)
	rand.Read(data)

	b.SetBytes(int64(l))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := range data {
			data[j] = 0
		}
	}
}

func BenchmarkGoZero_32(b *testing.B) {
	benchmarkGoZero(b, 32)
}

func BenchmarkGoZero_128(b *testing.B) {
	benchmarkGoZero(b, 128)
}

func BenchmarkGoZero_1k(b *testing.B) {
	benchmarkGoZero(b, 1*1024)
}

func BenchmarkGoZero_16k(b *testing.B) {
	benchmarkGoZero(b, 16*1024)
}

func BenchmarkGoZero_128k(b *testing.B) {
	benchmarkGoZero(b, 128*1024)
}

func BenchmarkGoZero_1M(b *testing.B) {
	benchmarkGoZero(b, 1024*1024)
}

func BenchmarkGoZero_16M(b *testing.B) {
	benchmarkGoZero(b, 16*1024*1024)
}

func BenchmarkGoZero_128M(b *testing.B) {
	benchmarkGoZero(b, 128*1024*1024)
}

func BenchmarkGoZero_512M(b *testing.B) {
	benchmarkGoZero(b, 512*1024*1024)
}

func benchmarkGoSet(b *testing.B, l int) {
	data := make([]byte, l)
	rand.Read(data)

	value := byte(rand.Intn(0x100))

	b.SetBytes(int64(l))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := range data {
			data[j] = value
		}
	}
}

func BenchmarkGoSet_32(b *testing.B) {
	benchmarkGoSet(b, 32)
}

func BenchmarkGoSet_128(b *testing.B) {
	benchmarkGoSet(b, 128)
}

func BenchmarkGoSet_1k(b *testing.B) {
	benchmarkGoSet(b, 1*1024)
}

func BenchmarkGoSet_16k(b *testing.B) {
	benchmarkGoSet(b, 16*1024)
}

func BenchmarkGoSet_128k(b *testing.B) {
	benchmarkGoSet(b, 128*1024)
}

func BenchmarkGoSet_1M(b *testing.B) {
	benchmarkGoSet(b, 1024*1024)
}

func BenchmarkGoSet_16M(b *testing.B) {
	benchmarkGoSet(b, 16*1024*1024)
}

func BenchmarkGoSet_128M(b *testing.B) {
	benchmarkGoSet(b, 128*1024*1024)
}

func BenchmarkGoSet_512M(b *testing.B) {
	benchmarkGoSet(b, 512*1024*1024)
}
