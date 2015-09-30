package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/zgo/rrd"
)

var (
	addr = flag.String("http", ":4950", "")
	root = flag.String("root", "/var/lib/munin", "")
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	flag.Parse()
}

func handler(rw http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, "/") {
		f, err := http.Dir(*root).Open(req.URL.Path)
		if err != nil {
			fmt.Fprintln(rw, err)
			return
		}
		dirList(rw, f)
		return
	}

	v := new(rrd.RRD)
	v.Load(filepath.Join(*root, req.URL.Path))
	req.Header.Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(v)
}

func main() {
	log.Fatal(http.ListenAndServe(*addr, http.HandlerFunc(handler)))
}
