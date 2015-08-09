package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func showPage(w http.ResponseWriter, request *http.Request) {
	path := request.URL.Path[1:]
	parts := strings.Split(path, "/")

	fmt.Fprintf(w, "<html><head>")
	fmt.Fprintf(w, "<title>"+parts[len(parts)-1]+"</title>")
	fmt.Fprintf(w, "<style>"+getCss(format)+"</style>")
	fmt.Fprintf(w, "<body><div style=\"display: table; width: 100%%\"><div style=\"display: table-row\">")
	fmt.Fprintf(w, "<div style=\"font-family: Courier; display: table-cell; width: 100%%;\"><a href=\"/\">home</a>&nbsp;&nbsp;"+breadcrumbHtml(path)+"</div>")
	fmt.Fprintf(w, "</div></div><hr />")

	ipath := path
	if ipath == "" {
		ipath = "."
	}
	if !strings.Contains(ipath, "..") {
		if isDirectory(ipath) {
			fmt.Fprintf(w, "<div style=\"font-family: Courier\"><ul>")
			entries, _ := ioutil.ReadDir(ipath)
			for _, entry := range entries {
				text := entry.Name()
				if entry.IsDir() {
					text += "/"
				}

				if hasMarkdownSuffix(entry.Name()) || entry.IsDir() {
					if !strings.HasPrefix(path, "/") && path != "" {
						path = "/" + path
					}
					fmt.Fprintf(w, "<li><a href=\""+path+"/"+entry.Name()+"\">"+text+"</a></li>")
				}
			}
			fmt.Fprintf(w, "</ul></div>")
		} else if hasMarkdownSuffix(path) && fileExists(path) {
			fmt.Fprintf(w, customMarkupHtml(getHtml(path)))
		} else {
			fmt.Fprintf(w, "File not found.")
		}
	} else {
		fmt.Fprintf(w, "File not found.")
	}
	fmt.Fprintf(w, "</body></html>")
}

func breadcrumbHtml(path string) string {
	html := ""
	partialPath := ""
	parts := strings.Split(path, "/")
	for _, part := range parts {
		partialPath += "/" + part
		html += " / <a href=\"" + partialPath + "\">" + part + "</a>"
	}
	return html
}
