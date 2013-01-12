// 11 january 2013
package main

import (
	"encoding/xml"
	"encoding/json"
	"io/ioutil"
	"crypto/sha1"
	"path/filepath"
"fmt"
)

// I would use JSON for this but MediaWiki does the following:
// 	{ "query": { "pages": { "[the actual page ID]": { ... } } }
// this doesn't work with the struct tag system of encoding/json
// they do provide an option to return a list of page IDs, indexpageids, but that still doesn't help when you're using struct tags
// cespare in #go-nuts suggested using a map, but map[string]editTokenResult didn't work (because I was doing it wrong, which I only found out AFTER WRITING ALL THIS BELOW) and map[string]interface{} would just lead to more work, so I'll just use XML and end it now
type editTokenResult struct {
	EditToken		string	`xml:"edittoken,attr"`
}

func getEditToken(filename string) string {
	var result struct {
		P	[]editTokenResult	`xml:"query>pages>page"`
	}

	resp := post("query", "xml", queryMIME,
		"prop", "info",
		"intoken", "edit",
		"titles", "File:" + filename)
	defer resp.Body.Close()
	d := xml.NewDecoder(resp.Body)
	err := d.Decode(&result)
	if err != nil {
		panic(err)
	}
	if len(result.P) != 1 {
		panic("zero pages or more than one page returned when getting edit token")
	}
	return result.P[0].EditToken
}

type uploadResult struct {
	// ...
}

const uploadMIME = "multipart/form-data"

func upload(filename string, editToken string) {
	outname := filepath.Base(filename)
	if outname == "" {
		panic("filename not given to upload")
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	sum := sha1.New()
	n, err := sum.Write(b)
	if err != nil {
		panic(err)
	}
	if n < len(b) {
		panic("short write to SHA-1 sum without error")
	}

	resp := post_multipart("upload", "json", b,
		"filename", outname,
		"token", editToken)
	defer resp.Body.Close()
/*	d := json.NewDecoder(resp.Body)
	err = d.Decode(&result)
	if err != nil {
		panic(err)
	}
*/	_=json.NewDecoder
	b,err=ioutil.ReadAll(resp.Body)
	fmt.Printf("%v %s\n", err, b)
}
