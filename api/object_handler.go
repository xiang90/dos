package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/xiang90/dos/block"
)

type ObjectsHandler struct {
	Storage *block.Manager
}

func (oh *ObjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		idstr := r.URL.Path[len("/objects")+1:]
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		b, err := oh.Storage.Get(int(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("got an object [id: %d, size: %d]", id, len(b.Blob))
		w.Write(b.Blob)

	case "POST":
		blob, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// fixme...
		id := rand.Int()
		err = oh.Storage.Append(&block.Block{ID: id, Blob: blob})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprint(id)))
		log.Printf("created an object [id: %d, size: %d]", id, len(blob))
	}
}
