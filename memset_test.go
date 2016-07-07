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

func memsetGo(data []byte, value byte) {
	if value == 0 {
		for i := range data {
			data[i] = 0
		}
	} else {
		for i := range data {
			data[i] = value
		}
	}
}

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

func benchmarkMemset(b *testing.B, value byte, l int) {
	data := make([]byte, l)
	rand.Read(data)

	b.SetBytes(int64(l))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Memset(data, value)
	}
}

func BenchmarkZero_32(b *testing.B) {
	benchmarkMemset(b, 0, 32)
}

func BenchmarkZero_128(b *testing.B) {
	benchmarkMemset(b, 0, 128)
}

func BenchmarkZero_1k(b *testing.B) {
	benchmarkMemset(b, 0, 1*1024)
}

func BenchmarkZero_16k(b *testing.B) {
	benchmarkMemset(b, 0, 16*1024)
}

func BenchmarkZero_128k(b *testing.B) {
	benchmarkMemset(b, 0, 128*1024)
}

func BenchmarkZero_1M(b *testing.B) {
	benchmarkMemset(b, 0, 1024*1024)
}

func BenchmarkZero_16M(b *testing.B) {
	benchmarkMemset(b, 0, 16*1024*1024)
}

func BenchmarkZero_128M(b *testing.B) {
	benchmarkMemset(b, 0, 128*1024*1024)
}

func BenchmarkZero_512M(b *testing.B) {
	benchmarkMemset(b, 0, 512*1024*1024)
}

func BenchmarkSet_32(b *testing.B) {
	benchmarkMemset(b, 0x55, 32)
}

func BenchmarkSet_128(b *testing.B) {
	benchmarkMemset(b, 0x55, 128)
}

func BenchmarkSet_1k(b *testing.B) {
	benchmarkMemset(b, 0x55, 1*1024)
}

func BenchmarkSet_16k(b *testing.B) {
	benchmarkMemset(b, 0x55, 16*1024)
}

func BenchmarkSet_128k(b *testing.B) {
	benchmarkMemset(b, 0x55, 128*1024)
}

func BenchmarkSet_1M(b *testing.B) {
	benchmarkMemset(b, 0x55, 1024*1024)
}

func BenchmarkSet_16M(b *testing.B) {
	benchmarkMemset(b, 0x55, 16*1024*1024)
}

func BenchmarkSet_128M(b *testing.B) {
	benchmarkMemset(b, 0x55, 128*1024*1024)
}

func BenchmarkSet_512M(b *testing.B) {
	benchmarkMemset(b, 0x55, 512*1024*1024)
}

func benchmarkMemsetGo(b *testing.B, value byte, l int) {
	data := make([]byte, l)
	rand.Read(data)

	b.SetBytes(int64(l))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		memsetGo(data, value)
	}
}

func BenchmarkZeroGo_32(b *testing.B) {
	benchmarkMemsetGo(b, 0, 32)
}

func BenchmarkZeroGo_128(b *testing.B) {
	benchmarkMemsetGo(b, 0, 128)
}

func BenchmarkZeroGo_1k(b *testing.B) {
	benchmarkMemsetGo(b, 0, 1*1024)
}

func BenchmarkZeroGo_16k(b *testing.B) {
	benchmarkMemsetGo(b, 0, 16*1024)
}

func BenchmarkZeroGo_128k(b *testing.B) {
	benchmarkMemsetGo(b, 0, 128*1024)
}

func BenchmarkZeroGo_1M(b *testing.B) {
	benchmarkMemsetGo(b, 0, 1024*1024)
}

func BenchmarkZeroGo_16M(b *testing.B) {
	benchmarkMemsetGo(b, 0, 16*1024*1024)
}

func BenchmarkZeroGo_128M(b *testing.B) {
	benchmarkMemsetGo(b, 0, 128*1024*1024)
}

func BenchmarkZeroGo_512M(b *testing.B) {
	benchmarkMemsetGo(b, 0, 512*1024*1024)
}

func BenchmarkSetGo_32(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 32)
}

func BenchmarkSetGo_128(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 128)
}

func BenchmarkSetGo_1k(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 1*1024)
}

func BenchmarkSetGo_16k(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 16*1024)
}

func BenchmarkSetGo_128k(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 128*1024)
}

func BenchmarkSetGo_1M(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 1024*1024)
}

func BenchmarkSetGo_16M(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 16*1024*1024)
}

func BenchmarkSetGo_128M(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 128*1024*1024)
}

func BenchmarkSetGo_512M(b *testing.B) {
	benchmarkMemsetGo(b, 0x55, 512*1024*1024)
}
