// 11 january 2013
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"bufio"
	"code.google.com/p/gopass"
	"io/ioutil"
)

// general TODOs:
// - better error handling all around
// - glob support
// - delay between files

var description string
var startingPath string

func doUpload(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}
	if info.IsDir() {
		if path == startingPath {
			// if we just skip the starting path the Walk() will end immediately; we don't want that
			return nil
		} else {		// TODO add recursive option
			return filepath.SkipDir
		}
	}
	outname := filepath.Base(path)
	et := getEditToken(outname)
	upload(path, description, et)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: %s API-URL files...\n", os.Args[0])
		os.Exit(1)
	}
	apiaddr = os.Args[1]
	fmt.Printf("Enter username: ")
	username, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {	// TODO EOF
		panic(err)
	}
	username = username[:len(username) - 1]	// strip \n
	password, err := gopass.GetPass("Enter password: ")
	if err != nil {	// TODO EOF
		panic(err)
	}
	login(username, password)
	defer logout()

	// TODO grab description from file
	fmt.Println("Enter description, ending with EOF.")
	desc, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	description = string(desc)

	for i := 2; i < len(os.Args); i++ {
		// TODO will this break on a filename instead of a directory name?
		startingPath = os.Args[i]
		filepath.Walk(startingPath, doUpload)
	}
}
