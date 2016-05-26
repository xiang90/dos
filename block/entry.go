package block

import (
	"encoding/binary"
	"io"
	"sync"
)

//
// - - - - - - - - - - - - - - - - -
// |   size  |  blob    |  crc     |
// | 8 bytes |  n bytes |  8 bytes |
// - - - - - - - - - - - - - - - - -
//

func entryEncode(w io.Writer, blob []byte) (int, error) {
	err := writeInt64(w, int64(len(blob)))
	if err != nil {
		return -1, err
	}
	_, err = w.Write(blob)
	if err != nil {
		return -1, err
	}
	return 8 + len(blob), nil
}

func entryDecode(r io.Reader) ([]byte, error) {
	size, err := readInt64(r)
	// avoid allocation
	b := make([]byte, size)
	_, err = r.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

var int64bufpool sync.Pool

func writeInt64(w io.Writer, n int64) error {
	var bufp *[]byte
	item := int64bufpool.Get()
	if item == nil {
		buf := make([]byte, 8)
		bufp = &buf
	} else {
		bufp = item.(*[]byte)
	}

	// http://golang.org/src/encoding/binary/binary.go
	binary.LittleEndian.PutUint64(*bufp, uint64(n))
	_, err := w.Write(*bufp)

	int64bufpool.Put(bufp)

	return err
}

func readInt64(r io.Reader) (int64, error) {
	var n int64
	err := binary.Read(r, binary.LittleEndian, &n)
	return n, err
}
