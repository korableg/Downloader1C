package downloader

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func init() {
	SetLogOutput(bytes.NewBuffer(nil))
}

func TestNewDownloader(t *testing.T) {
	nicks := make(map[string]bool, 0)
	nicks["erp"] = true
	New("test", "test", "/", nil)
	New("test", "test", "/", nicks)
}

func TestGet(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/login/rest/public/ticket/get" {
			fmt.Fprint(w, "{\"ticket\": \"Hello\"}")
		} else if r.RequestURI == "/login/ticket/auth?token=Hello" {
			fmt.Fprintln(w, "<a href=\"/project/test1\"/>")
			fmt.Fprintln(w, "<a href=\"/project/test2\"/>")
			fmt.Fprintln(w, "<a href=\"/project/test3\"/>")
		} else if strings.Contains(r.RequestURI, "project/test1") {
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test1&ver=1.0\"/>")
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test1&ver=1.1\"/>")
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test1&ver=1.2\"/>")
		} else if strings.Contains(r.RequestURI, "project/test2") {
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test2&ver=1.0\"/>")
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test2&ver=1.1\"/>")
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test2&ver=1.2\"/>")
		} else if strings.Contains(r.RequestURI, "project/test3") {
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test3&ver=1.0\"/>")
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test3&ver=1.1\"/>")
			fmt.Fprintln(w, "<a href=\"/version_files?nick=test3&ver=1.2\"/>")
		} else if r.URL.Path == "/releases/version_files" {
			query, err := url.ParseQuery(r.URL.RawQuery)
			if err != nil {
				log.Fatal(err)
			}

			nick := query.Get("nick")
			ver := query.Get("ver")
			ver = strings.Replace(ver, ".", "_", -1)

			fmt.Fprintf(w, "<a href=\"/version_file?%s&path=%s\\%s\\Readme.txt\"/>\n",
				r.URL.RawQuery, nick, ver)
			fmt.Fprintf(w, "<a href=\"/version_file?%s&path=%s\\%s\\release.exe\"/>\n",
				r.URL.RawQuery, nick, ver)
		} else if strings.HasSuffix(r.RequestURI, ".exe") {
			fmt.Fprintf(w, "<a href=\"%s/public/file/get/test\"/>", releasesURL)
		} else if strings.HasSuffix(r.RequestURI, ".txt") {
			fmt.Fprintln(w, "Hello! i'm test")
		} else if strings.Contains(r.RequestURI, "/public/file/get/test") {
			fmt.Fprintln(w, "Hello! i'm test")
		}

	}))

	defer ts.Close()

	nicks := make(map[string]bool, 0)
	nicks["test1"] = true
	//nicks["test3"] = true

	releasesURL_bak := releasesURL
	releasesURL = ts.URL + "/releases"

	loginURL_bak := loginURL
	loginURL = ts.URL + "/login"

	defer func() { releasesURL = releasesURL_bak; loginURL = loginURL_bak }()

	downldr := New("test", "test", "./", nicks)

	files, err := downldr.Get()
	if err != nil {
		t.Error(err)
	}
	if len(files) < 6 {
		t.Errorf("files must be 6")
	}
	os.RemoveAll("./test1")
	os.RemoveAll("./test2")
	os.RemoveAll("./test3")

	downldr = New("test", "test", "./", nil)

	files, err = downldr.Get()
	if err != nil {
		t.Error(err)
	}
	if len(files) < 18 {
		t.Errorf("files must be 18")
	}
	os.RemoveAll("./test1")
	os.RemoveAll("./test2")
	os.RemoveAll("./test3")
}

func TestBadLogin(t *testing.T) {

	downldr := New("test", "test", "/", nil)
	_, err := downldr.Get()

	if !(strings.Contains(err.Error(), "Incorrect login or password") ||
		strings.Contains(err.Error(), "Too many failed attempts")) {
		t.Error("Test bad login :(")
	}

}
