package compression_test

import (
	"flag"
	"io"
	"math/rand"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/inkyblackness/hacked/ss1/resource/lgres/internal/compression"
	"github.com/inkyblackness/hacked/ss1/serial"
)

// To run a profile, run it with the following command:
//     go test -bench . -args -cpuprofile=prof
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func rawData(size int) []byte {
	data := make([]byte, size)
	rand.Read(data) // nolint:gas
	return data
}

func initProfiling(b *testing.B, nameSuffix string) func() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile + "-" + b.Name() + nameSuffix)
		if err != nil {
			b.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		return func() { pprof.StopCPUProfile() }
	}
	return func() {}
}

func BenchmarkRawDataStorage(b *testing.B) {
	profileStop := initProfiling(b, "")
	defer profileStop()
	data := rawData(1024 * 1024)
	b.ResetTimer()
	for run := 0; run < b.N; run++ {
		encoder := serial.NewEncoder(serial.NewByteStore())
		_, _ = encoder.Write(data)
	}
}

func benchmarkCompression(b *testing.B, size int, nameSuffix string) {
	profileStop := initProfiling(b, nameSuffix)
	defer profileStop()
	data := rawData(size)
	b.ResetTimer()
	for run := 0; run < b.N; run++ {
		compressor := compression.NewCompressor(serial.NewByteStore())
		_, _ = compressor.Write(data)
		_ = compressor.Close()
	}
}

func BenchmarkCompression1KB(b *testing.B) {
	benchmarkCompression(b, 1024, "1KB")
}

func BenchmarkCompression16KB(b *testing.B) {
	benchmarkCompression(b, 1024*16, "16KB")
}

func BenchmarkCompression128KB(b *testing.B) {
	benchmarkCompression(b, 1024*128, "128KB")
}

func BenchmarkCompression512KB(b *testing.B) {
	benchmarkCompression(b, 1024*512, "512KB")
}

func BenchmarkCompression1024KB(b *testing.B) {
	benchmarkCompression(b, 1024*1024, "1024KB")
}

func benchmarkCompressionDecompression(b *testing.B, size int, nameSuffix string) {
	profileStop := initProfiling(b, nameSuffix)
	defer profileStop()
	data := rawData(size)
	output := make([]byte, len(data))
	b.ResetTimer()
	for run := 0; run < b.N; run++ {
		store := serial.NewByteStore()
		compressor := compression.NewCompressor(store)
		_, _ = compressor.Write(data)
		_ = compressor.Close()
		_, _ = store.Seek(0, io.SeekStart)
		decompressor := compression.NewDecompressor(store)
		_, _ = decompressor.Read(output)
	}
}

func BenchmarkCompressionDecompression1KB(b *testing.B) {
	benchmarkCompressionDecompression(b, 1024, "1KB")
}

func BenchmarkCompressionDecompression16KB(b *testing.B) {
	benchmarkCompressionDecompression(b, 1024*16, "16KB")
}

func BenchmarkCompressionDecompression128KB(b *testing.B) {
	benchmarkCompressionDecompression(b, 1024*128, "128KB")
}

func BenchmarkCompressionDecompression512KB(b *testing.B) {
	benchmarkCompressionDecompression(b, 1024*512, "512KB")
}

func BenchmarkCompressionDecompression1024KB(b *testing.B) {
	benchmarkCompressionDecompression(b, 1024*1024, "1024KB")
}
