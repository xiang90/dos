package block

type Block struct {
	ID   int
	Meta blockmeta
	Blob []byte
}
