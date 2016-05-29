package object

import (
	"encoding/binary"
	"math/rand"

	"github.com/xiang90/dos/block"
)

const (
	version1 = 1

	blobmeta = 1
	blobdata = 2
)

var (
	maxBlobSize = block.MaxSize - 4
)

type blobHeader struct {
	version uint16
	typ     uint16
}

func (bh *blobHeader) encode(b []byte) {
	binary.LittleEndian.PutUint16(b, bh.version)
	binary.LittleEndian.PutUint16(b[2:], bh.typ)
}

func (bh *blobHeader) decode(b []byte) {
	bh.version = binary.LittleEndian.Uint16(b)
	bh.typ = binary.LittleEndian.Uint16(b[2:])
}

func (bh *blobHeader) size() int {
	return 4
}

func makeDataBlock(data []byte) *block.Block {
	bh := &blobHeader{version: version1, typ: blobdata}
	blob := make([]byte, len(data)+bh.size())
	bh.encode(blob)
	copy(blob[4:], data)

	bk := &block.Block{
		ID:   rand.Int(),
		Blob: blob,
	}
	return bk
}
