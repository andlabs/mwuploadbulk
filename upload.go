// 11 january 2013
package main

import (
	"encoding/xml"
	"encoding/json"
	"io/ioutil"
	"crypto/sha1"
	"path/filepath"
	"encoding/hex"
	"bytes"
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
	Result		string	`json:"result"`
	ImageInfo	struct {
		SHA1	string	`json:"sha1"`
	}	`json:"imageinfo"`
}

const uploadMIME = "multipart/form-data"

func upload(filename string, description string, editToken string) {
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

	var result struct {
		R	uploadResult	`json:"upload"`
	}

	resp := post_multipart("upload", "json", b,
		"filename", outname,
		"token", editToken,
		"text", description,
		"ignorewarnings", "")		// TODO make it an option
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&result)
	if err != nil {
		panic(err)
	}
	if result.R.Result != "Success" {
		panic("received error " + result.R.Result + " during upload")
	}
	tosha1, err := hex.DecodeString(result.R.ImageInfo.SHA1)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(tosha1, sum.Sum(nil)) {
		panic("SHA-1 mismatch")
	}
	// TODO check if description sent successfully
}
