// 11 january 2013
package main

import (
	"net/http"
	"encoding/json"
)

// TODO better error handling

type loginResult struct {
	Result		string	`json:"result"`
	LGUserID		uint32	`json:"lguserid"`
	LGUsername	string	`json:"lgusername"`
	Token		string	`json:"token"`			// returned by step 1
	LGToken		string	`json:"lgtoken"`			// returned by step 2
	CookiePrefix	string	`json:"cookieprefix"`
	SessionID		string	`json:"sessionid"`
}

var loginCookie *http.Cookie

func login(username string, password string) {
	var result struct {
		R	loginResult	`json:"login"`
	}

	// step 1
	resp := post("login",
		"lgname", username,
		"lgpassword", password)
	d := json.NewDecoder(resp.Body)
	err := d.Decode(&result)
	if err != nil {
		panic(err)
	}
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		panic("zero cookies or more than one cookie returned")
	}
	loginCookie = cookies[0]
	resp.Body.Close()

	// do we need step 2?
	switch result.R.Result {
	case "Success":
		// TODO complain about possibly outdated or insecure MediaWiki?
		return		// no
	case "NeedToken":
		// yes; continue to step 2
	default:
		panic("received error " + result.R.Result + " during login (1/2)")
	}

	// step 2
	resp = post("login",
		"lgname", username,
		"lgpassword", password,
		"lgtoken", result.R.Token)
	d = json.NewDecoder(resp.Body)
	err = d.Decode(&result)
	if err != nil {
		panic(err)
	}
	if result.R.Result != "Success" {
		panic("received error " + result.R.Result + " during login (2/2)")
	}
	resp.Body.Close()
	// logged in!
}

func logout() {
	post("logout")
}

/*
func main() {
	login("andlabs", "fake password")
	defer logout()
}
*/
