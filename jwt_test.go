package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/net/publicsuffix"
)

func TestLogin(t *testing.T) {
	url := "http://localhost:8000/signin"

	var jsonStr = []byte(`{"username":"user1","password":"password1"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Status:%s", resp.Status)
	}
	SetCookie, exists := resp.Header["Set-Cookie"]
	if !exists {
		t.Fatal("missing: Set-Cookie")
	}
	if strings.Contains(SetCookie[0], "token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InVzZXIxIiwiZXhwIjoxNTc5NjQ4MDkzfQ") {
		t.Error("Token doesn't match expected")
	}
}

func TestWelcomeFail(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8000/welcome", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("Status:%s (expected 401 Unauthorized)", resp.Status)
	}

}

func TestSigninWelcome(t *testing.T) {
	u, err := url.Parse("http://localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	// create client with Jar
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Jar: jar,
	}

	// signin
	var jsonStr = []byte(`{"username":"user1","password":"password1"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/signin", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	t.Log("After 1st request:")
	for _, cookie := range jar.Cookies(u) {
		t.Logf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	req, err = http.NewRequest("GET", "http://localhost:8000/welcome", nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Status:%s", resp.Status)
	}

	t.Log("After 2st request:")
	for _, cookie := range jar.Cookies(u) {
		t.Logf("  %s: %s\n", cookie.Name, cookie.Value)
	}

}
