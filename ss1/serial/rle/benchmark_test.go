package rle_test

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/inkyblackness/hacked/ss1/serial/rle"
)

func benchmarkRecompression(b *testing.B, size int, chanceLimit byte, nameSuffix string, seed int64) {
	input := make([]byte, size)
	reference := make([]byte, size)
	var chance [1]byte

	rand.Seed(seed) // nolint:gas

	b.ResetTimer()
	for run := 0; run < b.N; run++ {
		b.StopTimer()
		rand.Read(input) // nolint:gas
		copy(reference, input)
		for i := 0; i < size; i++ {
			rand.Read(chance[:]) // nolint:gas
			if chance[0] < chanceLimit {
				input[i] = chance[0]
			}
		}
		buf := bytes.NewBuffer(nil)
		b.StartTimer()

		err := rle.Compress(buf, input, reference)
		if err != nil {
			b.Errorf("Failed compression for %s in run %d of seed %v", nameSuffix, run, seed)
		}

		err = rle.Decompress(bytes.NewReader(buf.Bytes()), reference)
		if err != nil {
			b.Errorf("Failed decompression for %s in run %d of seed %v", nameSuffix, run, seed)
		}
		if !bytes.Equal(input, reference) {
			b.Errorf("Data mismatch in %s in run %d of seed %v", nameSuffix, run, seed)
		}
	}
}

func BenchmarkRecompression64KB_10(b *testing.B) {
	benchmarkRecompression(b, 1024*64, 10, "64KB,10", time.Now().UnixNano())
}

func BenchmarkRecompression64KB_128(b *testing.B) {
	benchmarkRecompression(b, 1024*64, 128, "64KB,128", time.Now().UnixNano())
}

func BenchmarkRecompression1KB_0(b *testing.B) {
	benchmarkRecompression(b, 1024, 0, "1KB,0", time.Now().UnixNano())
}

func BenchmarkRecompression1KB_10(b *testing.B) {
	benchmarkRecompression(b, 1024, 10, "1KB,10", time.Now().UnixNano())
}

func BenchmarkRecompression1KB_128(b *testing.B) {
	benchmarkRecompression(b, 1024, 128, "1KB,128", time.Now().UnixNano())
}

func BenchmarkRecompression1KB_255(b *testing.B) {
	benchmarkRecompression(b, 1024, 255, "1KB,255", time.Now().UnixNano())
}

func BenchmarkRecompression128B_10(b *testing.B) {
	benchmarkRecompression(b, 128, 10, "128B,10", time.Now().UnixNano())
}

func BenchmarkRecompression128B_128(b *testing.B) {
	benchmarkRecompression(b, 128, 128, "128B,128", time.Now().UnixNano())
}

func BenchmarkRecompression128B_255(b *testing.B) {
	benchmarkRecompression(b, 128, 255, "128B,255", time.Now().UnixNano())
}
