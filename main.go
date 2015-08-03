package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var port int
var format string

func main() {
	flag.StringVar(&format, "format", "simple", "Format for output.  Options are 'raw' or 'simple'.")
	flag.IntVar(&port, "port", 8101, "Port for HTTP server.  Defaults to 8101.")
	flag.Usage = showUsage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		showUsage()
		os.Exit(0)
	}

	if args[0] == "help" {
		showUsage()
		os.Exit(0)
	} else if args[0] == "render" {
		if len(args) != 3 {
			fmt.Println("invalid arguments - usage: cuyuna render [source file] [destination file]")
			os.Exit(1)
		}
		path := args[1]
		dest := args[2]
		if !fileExists(path) {
			fmt.Println("specified source file does not exist")
			os.Exit(1)
		}

		doRender(path, dest)
	} else if args[0] == "serve" {
		doServe()
	} else {
		fmt.Println("invalid command, type 'cuyuna help' for help")
		os.Exit(1)
	}
}

func showUsage() {
	fmt.Println("usage: cuyuna [command]")
	fmt.Println("  help          : Prints this menu.")
	fmt.Println("  render [file] : Renders HTML version of markdown file at the specified path.")
	fmt.Println("  serve [root]  : Starts a local HTTP server in the current working directory for viewing available markdown files.")
}

func doRender(path string, target string) {
	html := getHtml(path)
	if format == "raw" {
		html = "<html><body>" + html + "</body></html>"
	} else {
		html = customMarkupHtml(html)
		html = "<html><head><style>" + getCss(format) + "</style></head><body>" + html + "</body></html>"
	}
	data := []byte(html)
	err := ioutil.WriteFile(target, data, 0644)
	if err != nil {
		fmt.Println("an error occurred while writing the file")
		os.Exit(1)
	}
}

func doServe() {
	http.HandleFunc("/", showPage)

	fmt.Printf("Waiting for HTTP connections on port %d . . .\n", port)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}

func hasMarkdownSuffix(path string) bool {
	if strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown") {
		return true
	} else {
		return false
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
