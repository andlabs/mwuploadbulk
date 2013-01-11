// 10 january 2013
package main

import (
	"net/http"
	"net/url"
"io";"os";"fmt"
)

// TODO better error handling

const apiaddr = "http://tcrf.net/api.php"
const format = "json"

func buildURL(action string, query ...string) string {
	if len(query) % 2 == 1 {
		panic("odd number of arguments passed to buildQuery")
	}
	u, err := url.Parse(apiaddr)
	if err != nil {
		panic(err)
	}
	v := u.Query()
	v.Set("format", format)			// TODO Add instead of Set for all of these?
	v.Set("action", action)
	for i := 0; i < len(query); i += 2 {
		v.Set(query[i], query[i + 1])
	}
	u.RawQuery = v.Encode()
	return u.String()
}

func post(action string, cookie *http.Cookie, query ...string) *http.Response {
	if len(query) % 2 == 1 {
		panic("query sent to post not set of key/value pairs (odd number of arguments)")
	}
	req, err := http.NewRequest("POST", buildURL(action, query...), nil)
	if err != nil {
		panic(err)
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

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
