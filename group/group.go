package group

import (
	"log"
	"time"

	"github.com/xiang90/dos/block"
)

type Group struct {
	Storage *block.Manager

	Peers map[int]*Peer
}

func (g *Group) Sync() {
	syncInterval := 10 * time.Second
	log.Printf("start to sync with peers every %v", syncInterval)
	for {
		select {
		case <-time.After(syncInterval):
		}
		g.syncWithPeers()
	}
}

func (g *Group) syncWithPeers() {
	for _, p := range g.Peers {
		rbs := p.RecentBlocks()
		if rbs == nil {
			continue
		}
		var anyerr error
		for _, id := range rbs.Blocks {
			if g.Storage.Has(id) {
				continue
			}
			log.Printf("ask missing block %d from peer %s", id, p.URL)
			err := g.askBlock(p, id)
			if err != nil {
				anyerr = err
			}
			log.Printf("got missing block %d from peer %s", id, p.URL)
		}
		if anyerr == nil && rbs.Until > p.progress {
			p.progress = rbs.Until
			log.Printf("update progress of peer %s to %d", p.URL, p.progress)
		}
	}
}

func (g *Group) askBlock(p *Peer, id int) error {
	b, err := p.GetBlock(id)
	if err != nil {
		return err
	}
	err = g.Storage.Append(&block.Block{ID: id, Blob: b})
	if err != nil {
		// internal issue...
		panic(err)
	}
	return nil
}
