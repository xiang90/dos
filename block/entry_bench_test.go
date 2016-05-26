package block

import (
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkWriteEntry1MB(b *testing.B) { benchmarkWriteEntry(b, 1024*1024) }

func benchmarkWriteEntry(b *testing.B, size int) {
	f, err := ioutil.TempFile(os.TempDir(), "entry_bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(f.Name())

	data := make([]byte, size)
	for i := 0; i < len(data); i++ {
		data[i] = byte(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(int64(size))
	for i := 0; i < b.N; i++ {
		entryEncode(f, data)
	}
}
