package group

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/xiang90/dos/api"
)

type Peer struct {
	URL      string
	progress int
}

func (p *Peer) RecentBlocks() *api.RecentBlocks {
	resp, err := http.Get(p.URL + "/" + api.RecentBlocksPath + "/" + fmt.Sprint(p.progress))
	if err != nil {
		log.Printf("failed to get recent blocks from %s: %v", p.URL, err)
		return nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to get recent blocks from %s: %v", p.URL, err)
		return nil
	}
	rb := &api.RecentBlocks{}
	err = json.Unmarshal(b, rb)
	if err != nil {
		log.Printf("failed to get recent blocks from %s: %v", p.URL, err)
		return nil
	}
	return rb
}

func (p *Peer) GetBlock(id int) ([]byte, error) {
	resp, err := http.Get(p.URL + "/" + api.BlocksPath + "/" + fmt.Sprint(id))
	if err != nil {
		log.Printf("failed to get block %s from %s: %v", id, p.URL, err)
		return nil, nil
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to get block %s from %s: %v", id, p.URL, err)
		return nil, err
	}
	return b, nil
}
