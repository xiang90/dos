package block

import (
	"errors"
	"io/ioutil"
	"os"
)

var (
	ErrBlockNotFound = errors.New("blockgroup: block not found")
	ErrUnavailable   = errors.New("blockgroup: not available")
)

type Manager struct {
	logs map[int]*bloblog

	ms *metaStore
	ps *progressStore
}

func NewManager() *Manager {
	f, err := ioutil.TempFile(os.TempDir(), "bloglog")
	if err != nil {
		panic(err)
	}
	bl := &bloblog{f: f}
	ms := &metaStore{
		backend: make(map[int]blockmeta),
	}
	ps := &progressStore{}

	return &Manager{logs: map[int]*bloblog{0: bl}, ms: ms, ps: ps}
}

func (m *Manager) Append(b *Block) error {
	for k, l := range m.logs {
		if l.trylock() {
			tail := l.tail()
			err := l.append(b.Blob)
			l.unlock()
			if err != nil {
				continue
			}
			meta := blockmeta{
				Log:       k,
				LogOffset: tail,
			}
			m.ms.put(b.ID, meta)
			m.ps.put(b.ID)
			return nil
		}
	}
	return ErrUnavailable
}

func (m *Manager) Get(id int) (*Block, error) {
	meta, ok := m.ms.get(id)
	if !ok {
		return nil, ErrBlockNotFound
	}
	b, err := m.logs[0].readAt(meta.LogOffset)
	if err != nil {
		return nil, err
	}
	return &Block{ID: id, Meta: meta, Blob: b}, nil
}

func (m *Manager) Has(id int) bool {
	_, ok := m.ms.get(id)
	return ok
}

func (m *Manager) RecentBlocks(since int) (int, []int) {
	return m.ps.recent(since)
}
