package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hugolgst/rich-go/client"
)

func getIcon(command string) (string, string) {
	switch command {
	case "git":
		return "git", "git source control"
	case "curl":
		return "curl", "the downloady thing"
	case "go":
		return "golang", "go lang"
	case "rustup":
		fallthrough
	case "cargo":
		return "rust", "crab crab crab"
	}
	return "bash", "the bourne again shell"
}

var startTime = time.Now()

func updateStatus(command string) {
	parts := strings.Split(command, " ")
	icon, iconText := getIcon(parts[0])
	err := client.SetActivity(client.Activity{
		State:   fmt.Sprintf("Running `%s`", parts[0]),
		Details: command,
		Timestamps: &client.Timestamps{
			Start: &startTime,
		},
		LargeImage: icon,
		LargeText:  iconText,
		SmallImage: "bash",
		SmallText:  "the bourne again shell",
	})
	if err != nil {
		fmt.Println("couldn't publish discord rich status", err)
	}
}

func main() {
	err := client.Login("906585124303409213")
	if err != nil {
		panic(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		var line string
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					content, err := ioutil.ReadFile(event.Name)
					if err != nil {
						fmt.Println("error:", err)
					}
					hist := bytes.Split(content, []byte{'\n'})
					last := string(hist[len(hist)-2])
					if last == line {
						continue
					}
					line = last
					updateStatus(line)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				fmt.Println("error:", err)
			}
		}
	}()

	homedir, _ := os.UserHomeDir()
	err = watcher.Add(filepath.Join(homedir, ".bash_history"))
	if err != nil {
		panic(err)
	}
	<-done
}
