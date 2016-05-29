package object

import (
	"io"
	"math/rand"
	"time"

	"github.com/xiang90/dos/block"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Object struct {
	Reader    io.ReadCloser
	metablock *block.Block
}

func NewObjectFromBlock(b *block.Block) *Object {
	return &Object{
		Reader:    newBlockr(b),
		metablock: b,
	}
}

func (o *Object) NextBlock() (*block.Block, error) {
	buf := make([]byte, maxBlobSize)
	n, err := o.Reader.Read(buf)
	if err == io.EOF {
		// simple case: object is smaller than max block size.
		// metablock and datablock is the same!
		if o.metablock == nil {
			o.metablock = makeDataBlock(buf[:n])
			return o.metablock, io.EOF
		}
		panic("unimplemented")
	}
	panic("unimplemented")
}

func (o *Object) ID() int {
	return o.metablock.ID
}
