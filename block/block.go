package block

var (
	// 1MB per block
	MaxSize = 1 * 1024 * 1024
)

type Block struct {
	ID   int
	Meta blockmeta
	Blob []byte
}
