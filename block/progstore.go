package block

type progressStore struct {
	blocks []int
}

func (ps *progressStore) put(blockid int) {
	ps.blocks = append(ps.blocks, blockid)
}

func (ps *progressStore) recent(since int) (int, []int) {
	return len(ps.blocks), ps.blocks[since:]
}
