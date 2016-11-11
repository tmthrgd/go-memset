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

type size struct {
	name string
	l    int
}

var sizes = []size{
	{"32", 32},
	{"128", 128},
	{"1K", 1 * 1024},
	{"16K", 16 * 1024},
	{"128K", 128 * 1024},
	{"1M", 1024 * 1024},
	{"16M", 16 * 1024 * 1024},
	{"128M", 128 * 1024 * 1024},
	{"512M", 512 * 1024 * 1024},
}

func BenchmarkMemset(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			data := make([]byte, size.l)
			rand.Read(data)

			value := byte(rand.Intn(0x100))

			b.SetBytes(int64(size.l))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				Memset(data, value)
			}
		})
	}
}

func BenchmarkGoZero(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			data := make([]byte, size.l)
			rand.Read(data)

			b.SetBytes(int64(size.l))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for j := range data {
					data[j] = 0
				}
			}
		})
	}
}

func BenchmarkGoSet(b *testing.B) {
	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			data := make([]byte, size.l)
			rand.Read(data)

			value := byte(rand.Intn(0x100))

			b.SetBytes(int64(size.l))
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				for j := range data {
					data[j] = value
				}
			}
		})
	}
}
