package frontend

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/xiang90/dos/block"
	"github.com/xiang90/dos/object"
)

type Router struct {
	Groups map[string][]string
}

func (r *Router) Distribute(o *object.Object) error {
	// only handle the simplest case for now.
	// object < max block size.
	block, err := o.NextBlock()
	if err != io.EOF {
		panic(err)
	}
	ok := r.PutBlock(block)
	if ok {
		return nil
	}
	return errors.New("router: failed to distribute object to block group")
}

func (r *Router) Get(id int) (*object.Object, error) {
	// todo: support dynamic groups
	// fixme: do not hard code group name.
	gurls := r.Groups["1"]

	var err error
	for _, url := range gurls {
		var resp *http.Response
		// again, now we only handle small block case
		resp, err = http.Get(url + "/" + fmt.Sprintf("%d", id))
		if err == nil {
			return &object.Object{
				Reader: resp.Body,
			}, nil
		}
		log.Printf("failed to get block %d from %s", id, url)
	}
	return nil, err
}

func (r *Router) PutBlock(b *block.Block) bool {
	// todo: support dynamic groups
	// fixme: do not hard code group name.
	gurls := r.Groups["1"]
	buf := bytes.NewBuffer(b.Blob)

	// define replication requirement.
	var ok bool
	for _, url := range gurls {
		l := url + "/" + fmt.Sprintf("%d", b.ID)
		req, err := http.NewRequest(http.MethodPut, l, buf)
		if err != nil {
			panic(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("failed to put block %d to %s", b.ID, l)
			continue
		}
		ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			log.Printf("failed to put block %d to %s (bad status code: %d)", b.ID, l, http.StatusCreated)
			continue
		}
		log.Printf("put block %d to %s", b.ID, l)
		ok = true
	}
	return ok
}
