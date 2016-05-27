package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/xiang90/dos/api"
	"github.com/xiang90/dos/block"
	"github.com/xiang90/dos/frontend"
	"github.com/xiang90/dos/group"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "frontend" {
		bg := []string{"http://127.0.0.1:8081/blocks", "http://127.0.0.1:8082/blocks"}
		router := &frontend.Router{
			Groups: map[string][]string{"1": bg},
		}
		oh := &api.ObjectsHandler{
			Router: router,
		}

		http.Handle("/objects/", oh)
		http.ListenAndServe("127.0.0.1:9999", nil)
		return
	}

	id := flag.Int("id", 1, "id")
	flag.Parse()

	m := block.NewManager()
	bh := &api.BlocksHandler{Storage: m}
	ph := &api.PeerHandler{Storage: m}

	g := group.Group{Storage: m, Peers: make(map[int]*group.Peer)}
	if *id == 1 {
		g.Peers[2] = &group.Peer{URL: "http://127.0.0.1:8082/peer"}
	} else {
		g.Peers[1] = &group.Peer{URL: "http://127.0.0.1:8081/peer"}
	}

	go g.Sync()

	http.Handle("/blocks/", bh)
	http.Handle("/peer/", ph)

	http.ListenAndServe(fmt.Sprintf("127.0.0.1:808%d", *id), nil)
}
