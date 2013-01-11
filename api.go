// 10 january 2013
package main

import (
	"net/http"
	"net/url"
	"io"
	"strings"
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

func post(action string, format string, MIMEtype string, query ...string) *http.Response {
	if len(query) % 2 == 1 {
		panic("query sent to post not set of key/value pairs (odd number of arguments)")
	}
	req, err := http.NewRequest("POST", apiaddr, buildQueryBody(action, format, query...))
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
