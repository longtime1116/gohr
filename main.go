package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

type FileModified map[string]time.Time
type FileInfo struct {
	name string
	mod  time.Time
}

func ReadDir() ([]FileInfo, error) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return nil, err
	}
	fi := make([]FileInfo, 0, len(files))
	for _, f := range files {
		fi = append(fi, FileInfo{f.Name(), f.ModTime()})
	}
	return fi, nil
}

func (m FileModified) register(files []FileInfo) int {
	for _, f := range files {
		m[f.name] = f.mod
	}
	return len(files)
}

func (m FileModified) update(fi *FileInfo) {
	m[fi.name] = fi.mod
}

func reload() {
	clear()

	fmt.Printf("Reloading...\n\n")
	// build
	out, err := exec.Command("go", "build", "-o", "main").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(out))
		//TODO: keep the previous result?
		return
	}
	// execute
	out, err = exec.Command("./main").CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("----- Results -----")
	fmt.Println(string(out))
}

func clear() {
	// clear the screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("[gohr running]")
}
func main() {
	clear()
	m := FileModified(make(map[string]time.Time))
	files, err := ReadDir()
	if err != nil {
		panic(err)
	}
	c := m.register(files)
	for {
		// TODO: get files recursively
		files, err := ReadDir()
		if err != nil {
			panic(err)
		}
		if c != len(files) {
			c = len(files)
			reload()
		}
		// check existing files
		for _, f := range files {
			// TODO: use argument
			if f.name == "main" {
				continue
			}
			if _, ok := m[f.name]; !ok || m[f.name] != f.mod {
				m.update(&f)
				reload()
			}
		}
	}
}
