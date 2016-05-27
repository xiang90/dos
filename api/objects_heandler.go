package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/xiang90/dos/frontend"
	"github.com/xiang90/dos/object"
)

type ObjectsHandler struct {
	Router *frontend.Router
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
		obj, err := oh.Router.Get(int(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.Copy(w, obj.Reader)
		obj.Reader.Close()

	case "POST":
		obj := &object.Object{
			Reader: r.Body,
		}
		err := oh.Router.Distribute(obj)
		r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprint(obj.ID())))
	}
}
