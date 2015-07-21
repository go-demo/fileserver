package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fileserver"
	app.Usage = "An simple file server."
	app.Author = "Lyric"
	app.Email = "tiannianshou@gmail.com"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 5230,
			Usage: "listen port",
		},
	}

	app.Action = func(c *cli.Context) {
		args := c.Args()
		if len(args) == 0 {
			fmt.Println("ERROR:Please specify a static directory.")
			os.Exit(1)
		}
		Run(args[0], c.Int("port"))
	}

	app.RunAndExitOnError()
}

type MyHandler struct {
	FilePath string
}

func (m *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		http.ServeFile(w, r, filepath.Join(m.FilePath, "index.html"))
		return
	}
	fileName := m.FilePath + path
	f, err := os.Stat(fileName)
	if err != nil {
		fmt.Fprintln(w, "404 Not found.")
		return
	}
	if f.IsDir() {
		http.ServeFile(w, r, filepath.Join(fileName, "index.html"))
	} else {
		http.ServeFile(w, r, fileName)
	}
}

func Run(dir string, port int) {
	handler := &MyHandler{
		FilePath: dir,
	}
	fmt.Printf("===> Server is running at %d port.\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
