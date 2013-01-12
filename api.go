// 10 january 2013
package main

import (
	"net/http"
	"net/url"
	"io"
	"strings"
	"mime/multipart"
	"bytes"
//"os";"fmt"
)

// TODO better error handling

var apiaddr string

const queryMIME = "application/x-www-form-urlencoded"

// [22:46] <kevlar> see the implementation of PostForm for how you need to encode and set the POST body.
// which just posts the encoded queries with MIME type queryMIME
func buildQueryBody(action string, format string, query ...string) io.Reader {
	if len(query) % 2 == 1 {
		panic("odd number of arguments passed to buildQueryBody")
	}
	v := url.Values{}
	v.Set("format", format)			// TODO Add instead of Set for all of these?
	v.Set("action", action)
	for i := 0; i < len(query); i += 2 {
		v.Set(query[i], query[i + 1])
	}
	return strings.NewReader(v.Encode())
}

func buildMultipartBody(action string, format string, filecontents []byte, query ...string) (MIMEtype string, r io.Reader) {
	if len(query) % 2 == 1 {
		panic("odd number of arguments passed to buildQueryBody")
	}
	b := new(bytes.Buffer)
	v := multipart.NewWriter(b)
	err := v.WriteField("format", format)
	if err != nil {
		panic(err)
	}
	err = v.WriteField("action", action)
	if err != nil {
		panic(err)
	}
	filename := ""
	for i := 0; i < len(query); i += 2 {
		err = v.WriteField(query[i], query[i + 1])
		if err != nil {
			panic(err)
		}
		if query[i] == "filename" {
			filename = query[i + 1]
		}
	}
	if filename == "" {
		panic("filename not provided")
	}
	fw, err := v.CreateFormFile("file", filename)
	if err != nil {
		panic(err)
	}
	n, err := fw.Write(filecontents)
	if err != nil {
		panic(err)
	} else if n < len(filecontents) {
		panic("short writte shoving in file but no error returned")
	}
	v.Close()
	return v.FormDataContentType(), b
}

func dopost(MIMEtype string, r io.Reader) *http.Response {
	req, err := http.NewRequest("POST", apiaddr, r)
	if err != nil {
		panic(err)
	}
	if loginCookie != nil {
		req.AddCookie(loginCookie)
	}
	req.Header.Set("Content-Type", MIMEtype)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

// TODO abstract MIMEtype away from here

func post(action string, format string, MIMEtype string, query ...string) *http.Response {
	return dopost(MIMEtype, buildQueryBody(action, format, query...))
}

func post_multipart(action string, format string, filecontents []byte, query...string) *http.Response {
	return dopost(buildMultipartBody(action, format, filecontents, query...))
}

/*
func main() {
	resp := post("login", nil,
		"lgname", "Andlabs",
		"lgpassword", "fake password")
	fmt.Printf("%+v\n", resp)
	io.Copy(os.Stdout, resp.Body)
	fmt.Println()
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		panic("too many cookies")
	}
	resp = post("login", cookies[0],
		"lgname", "Andlabs",
		"lgpassword", "fake password")
	fmt.Printf("%+v\n", resp)
	io.Copy(os.Stdout, resp.Body)
	fmt.Println()
}
*/
