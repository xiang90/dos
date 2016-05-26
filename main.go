package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/xiang90/dos/api"
	"github.com/xiang90/dos/block"
	"github.com/xiang90/dos/group"
)

func main() {
	id := flag.Int("id", 1, "id")
	flag.Parse()

	m := block.NewManager()
	oh := &api.ObjectsHandler{Storage: m}
	ph := &api.PeerHandler{Storage: m}

	g := group.Group{Storage: m, Peers: make(map[int]*group.Peer)}
	if *id == 1 {
		g.Peers[2] = &group.Peer{URL: "http://127.0.0.1:8082/peer"}
	} else {
		g.Peers[1] = &group.Peer{URL: "http://127.0.0.1:8081/peer"}
	}

	go g.Sync()

	http.Handle("/objects/", oh)
	http.Handle("/peer/", ph)

	http.ListenAndServe(fmt.Sprintf("127.0.0.1:808%d", *id), nil)
}
