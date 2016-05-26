package block

type metaStore struct {
	backend map[int]blockmeta
}

func (ms *metaStore) get(blockid int) (blockmeta, bool) {
	bm, ok := ms.backend[blockid]
	return bm, ok
}

func (ms *metaStore) put(blockid int, bm blockmeta) {
	ms.backend[blockid] = bm
}
