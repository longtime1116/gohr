package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

// FileModified is a map with filename as key and modtime as value.
type FileModified map[string]time.Time

// FileInfo is the minimum required file information extracted.
type FileInfo struct {
	name string
	mod  time.Time
}

// ReadDir gets filename and mod time of all files under the current directory
// TODO: get files recursively
// FIXME: bad function name
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

func outfname() (string, error) {
	flag.Parse()
	if len(flag.Args()) == 0 {
		dirname, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return path.Base(dirname), nil
	}
	return flag.Args()[0], nil
}

func reload(bin string) {
	clear(bin)

	fmt.Printf("Reloading... ")
	// build
	out, err := exec.Command("go", "build", "-o", bin).CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(out))
		//TODO?: keep the previous result
		return
	}
	fmt.Printf("`" + bin + "` was built!\n\n")

	// execute
	out, err = exec.Command("./" + bin).CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("----- Results -----")
	fmt.Println(string(out))
}

func clear(bin string) {
	// clear the screen
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("$ gohr " + bin)
}
func main() {
	bin, err := outfname()
	if err != nil {
		panic(err)
	}

	m := FileModified(make(map[string]time.Time))
	files, err := ReadDir()
	if err != nil {
		panic(err)
	}
	c := m.register(files)

	reload(bin)
	for {
		files, err := ReadDir()
		if err != nil {
			panic(err)
		}
		// if files are removed or created, rebuild and rerun
		if c != len(files) {
			c = len(files)
			reload(bin)
		}
		// if a file modified, rebuild and rerun
		for _, f := range files {
			if f.name == bin {
				continue
			}
			if _, ok := m[f.name]; !ok || m[f.name] != f.mod {
				m.update(&f)
				reload(bin)
			}
		}
	}
}
