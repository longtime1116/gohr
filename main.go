package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// FileModified is a map with filename as key and modtime as value.
type FileModified map[string]time.Time

// FileInfo is the minimum required file information extracted.
type FileInfo struct {
	name string
	mod  time.Time
}

func (fi *FileInfo) String() string {
	return fmt.Sprintf("%v (%v)", fi.name, fi.mod)
}

// DirWalk gets filename and mod time of all files under the current directory recursively
func DirWalk(path string) ([]*FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	fi := make([]*FileInfo, 0, len(files))
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		if f.IsDir() {
			fi2, err := DirWalk(filepath.Join(path, f.Name()))
			if err != nil {
				return nil, err
			}
			fi = append(fi, fi2...)
			continue
		}
		fi = append(fi, &FileInfo{filepath.Join(path, f.Name()), f.ModTime()})
	}
	return fi, nil
}

func (m FileModified) register(files []*FileInfo) int {
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
func init() {
	testing.Init()
	flag.Usage = func() {
		txt := `Usage: gohr <output binary name>
When you ommit output binary name, the basename of current directory is used.`
		fmt.Fprintf(os.Stderr, "%s\n", txt)
	}
	flag.Parse()
}
func main() {
	bin, err := outfname()
	if err != nil {
		panic(err)
	}

	m := FileModified(make(map[string]time.Time))
	files, err := DirWalk("./")
	if err != nil {
		panic(err)
	}
	c := m.register(files)

	reload(bin)
	for {
		time.Sleep(500 * time.Millisecond)
		files, err := DirWalk("./")
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
				m.update(f)
				reload(bin)
			}
		}
	}
}
