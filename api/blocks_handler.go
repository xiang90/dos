package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/xiang90/dos/block"
)

type BlocksHandler struct {
	Storage *block.Manager
}

func (bh *BlocksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Path[len("/blocks")+1:]
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		b, err := bh.Storage.Get(int(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Write(b.Blob)

	case "PUT":
		blob, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = bh.Storage.Append(&block.Block{ID: int(id), Blob: blob})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprint(id)))
	}
}
