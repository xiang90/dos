package object

import (
	"io"

	"github.com/xiang90/dos/block"
)

type blockr struct {
	firstbk *block.Block

	// currently reading block
	curbk  *block.Block
	offset int
	// how many blocks we have read
	count int
	// is the current reading block the last block?
	eofbk bool
	more  chan *block.Block
}

func newBlockr(first *block.Block) *blockr {
	hd := blobHeader{}
	hd.decode(first.Blob[0:4])
	if hd.typ != blobdata {
		panic("metablob not implmented")
	}

	return &blockr{
		firstbk: first,

		curbk:  first,
		offset: 4,
		eofbk:  true,
	}
}

func (br *blockr) Read(p []byte) (int, error) {
	n := copy(p, br.curbk.Blob[br.offset:])
	br.offset += n
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func (br *blockr) Close() error {
	return nil
}
