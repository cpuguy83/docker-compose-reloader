package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/jaschaephraim/lrserver"
	"gopkg.in/fsnotify.v1"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "watch",
			Value: &cli.StringSlice{"."},
		},
		cli.StringFlag{
			Name: "service",
		},
	}

	app.Action = run
	app.Run(os.Args)
}

func run(ctx *cli.Context) {
	watchDirs := ctx.StringSlice("watch")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error setting up file watcher: %v", err)
	}
	for _, d := range watchDirs {
		watcher.Add(d)
	}
	defer watcher.Close()

	lr, err := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error starting live reload: %v", err)
	}
	go lr.ListenAndServe()

	composeBin, err := exec.LookPath("docker-compose")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not find docker-compose path")
		os.Exit(1)
	}

	var cmd *exec.Cmd
	var t <-chan time.Time

	var ignore bool
MainLoop:
	for {
		select {
		case e := <-watcher.Events:
			select {
			case <-t:
				ignore = false
			default:
				if ignore {
					continue MainLoop
				}

				match, _ := filepath.Match("*.sw*", e.Name)
				if match || strings.HasSuffix(e.Name, "~") {
					continue MainLoop
				}
			}

			if ctx.String("service") == "" {
				cmd = exec.Command(composeBin, "up", "-d")
			} else {
				cmd = exec.Command(composeBin, "up", "-d", ctx.String("service"))
			}

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			ignore = true
			time.Sleep(1 * time.Second)
			lr.Reload("")
			t = time.After(3 * time.Second)
		case err := <-watcher.Errors:
			fmt.Fprintf(os.Stderr, "error watching files: %v", err)
			os.Exit(1)
		}
	}
}
