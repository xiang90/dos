package block

import "os"

type bloblog struct {
	f *os.File
	// offset of tail
	tailoff int64
}

func (bl *bloblog) append(blob []byte) error {
	_, err := bl.f.Seek(bl.tailoff, os.SEEK_SET)
	if err != nil {
		return err
	}
	var n int
	n, err = entryEncode(bl.f, blob)
	if err == nil {
		bl.tailoff += int64(n)
	}
	return err
}

func (bl *bloblog) readAt(offset int64) ([]byte, error) {
	_, err := bl.f.Seek(offset, os.SEEK_SET)
	if err != nil {
		return nil, err
	}
	return entryDecode(bl.f)
}

func (bl *bloblog) trylock() bool {
	return true
}

func (bl *bloblog) unlock() {

}

func (bl *bloblog) tail() int64 {
	return bl.tailoff
}
