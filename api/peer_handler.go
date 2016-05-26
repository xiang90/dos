package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/xiang90/dos/block"
)

const (
	RecentBlocksPath = "recentblocks"
	BlocksPath       = "blocks"
)

type RecentBlocks struct {
	Until  int
	Blocks []int
}

type PeerHandler struct {
	Storage *block.Manager
}

func (ph *PeerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s [%s]", r.Method, r.URL, r.RemoteAddr)
	switch {
	case strings.HasPrefix(r.URL.Path, "/peer/"+BlocksPath+"/"):
		ph.serveBlock(w, r)
	case strings.HasPrefix(r.URL.Path, "/peer/"+RecentBlocksPath+"/"):
		ph.serveRecentBlocks(w, r)
	}
}

func (ph *PeerHandler) serveBlock(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Path[len("/peer/"+BlocksPath)+1:]
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := ph.Storage.Get(int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Write(b.Blob)
}

func (ph *PeerHandler) serveRecentBlocks(w http.ResponseWriter, r *http.Request) {
	sincestr := r.URL.Path[len("/peer/"+RecentBlocksPath)+1:]
	since, err := strconv.ParseInt(sincestr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	until, bs := ph.Storage.RecentBlocks(int(since))

	rb := &RecentBlocks{
		Until:  until,
		Blocks: bs,
	}

	b, err := json.Marshal(rb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
