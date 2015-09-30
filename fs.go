package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

// https://github.com/golang/go/blob/master/src/net/http/fs.go

func dirList(w http.ResponseWriter, f http.File) {
	dirs, err := f.Readdir(-1)
	if err != nil {
		// TODO: log err.Error() to the Server.ErrorLog, once it's possible
		// for a handler to get at its Server via the ResponseWriter. See
		// Issue 12438.
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}
	sort.Sort(byName(dirs))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		} else if !strings.HasSuffix(name, ".rrd") {
			continue
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), htmlReplacer.Replace(name))
	}
	fmt.Fprintf(w, "</pre>\n")
}

type byName []os.FileInfo

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)
