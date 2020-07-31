package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func prepareTestDirWalk(workDir string) {
	// file just in the tmp dir
	file := filepath.Join(workDir, "file.txt")
	if err := ioutil.WriteFile(file, []byte("test"), 0666); err != nil {
		log.Fatal(err)
	}

	// directory and a file in it
	dir := workDir + "/dir"
	if err := os.Mkdir(dir, 0777); err != nil {
		log.Fatal(err)
	}
	file = filepath.Join(dir, "file_in_dir.txt")
	if err := ioutil.WriteFile(file, []byte("test"), 0666); err != nil {
		log.Fatal(err)
	}

	// hidden directory and file in it
	hiddenDir, err := ioutil.TempDir(workDir, ".hidden_dir")
	if err != nil {
		log.Fatal(err)
	}
	file = filepath.Join(hiddenDir, "file_in_hidden_dir.txt")
	if err := ioutil.WriteFile(file, []byte("test"), 0666); err != nil {
		log.Fatal(err)
	}

	// hidden file
	file = filepath.Join(workDir, ".hidden_file.txt")
	if err := ioutil.WriteFile(file, []byte("test"), 0666); err != nil {
		log.Fatal(err)
	}
}
func TestDirWalk(t *testing.T) {
	cur, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// create tmp dir
	workDir, err := ioutil.TempDir(cur, "TestDirWalk")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(workDir)

	// move to tmp dir and create necessary files, then execute DirWalk
	os.Chdir(workDir)
	prepareTestDirWalk("./")
	now := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fis, err := DirWalk("./")
	for _, fi := range fis {
		fi.mod = now
	}
	ans := []*FileInfo{
		{"dir/file_in_dir.txt", now},
		{"file.txt", now},
	}
	if !reflect.DeepEqual(fis, ans) {
		t.Errorf("TestDirWalk() failed.\nanswer: %v\nresult: %v", ans, fis)
	}
}
