package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
)

var releasesURL = "https://releases.1c.ru"
var loginURL = "https://login.1c.ru"

const projectHrefPrefix = "/project/"
const versionFilesHrefPrefix = "/version_files"
const fileServerHrefPrefix = "/public/file/get"

const tempFileSuffix = ".d1c"

var semaMaxConnections = make(chan struct{}, 10)
var logOutput = io.Writer(os.Stdout)

type FileToDownload struct {
	url  string
	path string
	name string
}

type Downloader struct {
	login      string
	password   string
	basePath   string
	nicks      map[string]bool
	httpClient *http.Client
	urlCh      chan *FileToDownload
	wg         sync.WaitGroup
	logger     *log.Logger
}

func New(login, password, basePath string, nicks map[string]bool) *Downloader {

	if len(basePath) > 0 && !os.IsPathSeparator(basePath[len(basePath)-1]) {
		basePath += string(os.PathSeparator)
	}

	dr := &Downloader{
		login:    login,
		password: password,
		basePath: basePath,
	}

	if nicks != nil {
		dr.nicks = make(map[string]bool, len(nicks))
		for k, v := range nicks {
			dr.nicks[projectHrefPrefix+strings.ToLower(k)] = v
		}
	}

	cj, _ := cookiejar.New(nil)
	dr.httpClient = &http.Client{
		Jar: cj,
	}

	dr.logger = log.New(logOutput, "", log.LstdFlags)

	return dr

}

func (dr *Downloader) Get() ([]os.FileInfo, error) {

	files := make([]os.FileInfo, 0)

	ticketUrl, err := dr.getURL()
	if err != nil {
		dr.logger.Println(err)
		return files, err
	}

	dr.urlCh = make(chan *FileToDownload, 10000)
	dr.wg.Add(1)
	go dr.findLinks(ticketUrl, dr.findProject)
	go func() { dr.wg.Wait(); close(dr.urlCh) }()

	for fileToDownload := range dr.urlCh {
		if fileInfo, ok := dr.downloadFile(fileToDownload); ok {
			files = append(files, fileInfo)
		}
	}

	dr.urlCh = nil

	return files, nil

}

func (dr *Downloader) getURL() (string, error) {

	type loginParams struct {
		Login       string `json:"login"`
		Password    string `json:"password"`
		ServiceNick string `json:"serviceNick"`
	}

	type ticket struct {
		Ticket string `json:"ticket"`
	}

	postBody, err := json.Marshal(
		loginParams{dr.login, dr.password, releasesURL})
	if err != nil {
		return "", err
	}

	acquireSemaConnections()
	resp, err := dr.httpClient.Post(
		loginURL+"/rest/public/ticket/get",
		"application/json",
		bytes.NewReader(postBody))
	releaseSemaConnections()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", string(responseBodyData))
	}

	var tick ticket
	err = json.Unmarshal(responseBodyData, &tick)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(loginURL+"/ticket/auth?token=%s", tick.Ticket), nil

}

func (dr *Downloader) findLinks(rawUrl string, f func(string, string)) {

	defer dr.wg.Done()

	acquireSemaConnections()
	resp, err := dr.httpClient.Get(rawUrl)
	releaseSemaConnections()

	if err != nil {
		dr.logger.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		dr.logger.Println(err)
		return
	}

	dr.eachNode(doc, rawUrl, f)

}

func (dr *Downloader) eachNode(node *html.Node, u string, f func(string, string)) {

	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				f(u, attr.Val)
				break
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		dr.eachNode(c, u, f)
	}

}

func (dr *Downloader) findProject(_, href string) {

	if (dr.nicks == nil && strings.HasPrefix(href, projectHrefPrefix)) || dr.nicks[strings.ToLower(href)] {
		dr.wg.Add(1)
		go dr.findLinks(releasesURL+href, dr.findVersion)
	}

}

func (dr *Downloader) findVersion(_, href string) {

	if strings.HasPrefix(href, versionFilesHrefPrefix) {
		dr.wg.Add(1)
		go dr.findLinks(releasesURL+href, dr.findToDownloadLink)
	}

}

func (dr *Downloader) findToDownloadLink(_, href string) {

	lowerHref := strings.ToLower(href)

	if strings.HasSuffix(lowerHref, "rar") ||
		strings.HasSuffix(lowerHref, "zip") ||
		strings.HasSuffix(lowerHref, "gz") ||
		strings.HasSuffix(lowerHref, "exe") ||
		strings.HasSuffix(lowerHref, "msi") ||
		strings.HasSuffix(lowerHref, "deb") ||
		strings.HasSuffix(lowerHref, "rpm") ||
		strings.HasSuffix(lowerHref, "epf") ||
		strings.HasSuffix(lowerHref, "erf") {

		dr.wg.Add(1)
		dr.findLinks(releasesURL+href, dr.findFileServerLink)

	} else if strings.HasSuffix(lowerHref, "txt") ||
		strings.HasSuffix(lowerHref, "pdf") ||
		strings.HasSuffix(lowerHref, "html") ||
		strings.HasSuffix(lowerHref, "htm") {

		dr.addFileToChannel(href, releasesURL+href)

	}

}

func (dr *Downloader) findFileServerLink(u, href string) {

	if strings.Contains(href, fileServerHrefPrefix) {
		dr.addFileToChannel(u, href)
	}

}

func (dr *Downloader) addFileToChannel(u, href string) {
	fileName, filePath, err := dr.fileNameFromUrl(u)
	if err == nil {
		fileToDownload := FileToDownload{
			url:  href,
			path: filePath,
			name: fileName,
		}
		dr.urlCh <- &fileToDownload
	} else {
		dr.logger.Println(err)
	}
}

func (dr *Downloader) fileNameFromUrl(rawUrl string) (string, string, error) {

	fileName := strings.Builder{}
	filePath := strings.Builder{}

	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", "", err
	}

	query, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return "", "", err
	}

	path := strings.Split(query.Get("path"), "\\")
	fileName.WriteString(path[len(path)-1])

	nick := query.Get("nick")
	ver := query.Get("ver")

	filePath.WriteString(nick)
	filePath.WriteRune(os.PathSeparator)
	filePath.WriteString(ver)
	filePath.WriteRune(os.PathSeparator)

	return fileName.String(), filePath.String(), nil
}

func (dr *Downloader) downloadFile(fileToDownload *FileToDownload) (os.FileInfo, bool) {

	fullpath := dr.basePath + fileToDownload.path + fileToDownload.name
	fileInfo, err := os.Stat(fullpath)
	if os.IsExist(err) {

		return fileInfo, true

	} else if os.IsNotExist(err) {

		acquireSemaConnections()
		resp, err := dr.httpClient.Get(fileToDownload.url)
		releaseSemaConnections()
		if err != nil {
			return nil, false
		}

		err = os.MkdirAll(dr.basePath+fileToDownload.path, os.ModeDir)
		if err != nil {
			return nil, false
		}

		f, err := os.Create(fullpath + tempFileSuffix)
		if err != nil {
			return nil, false
		}

		defer resp.Body.Close()

		_, err = io.Copy(f, resp.Body)
		if err != nil {
			return nil, false
		}
		f.Close()

		err = os.Rename(fullpath+tempFileSuffix, fullpath)
		if err != nil {
			return nil, false
		}

		fileInfo, err := os.Stat(fullpath)
		if err != nil {
			return nil, false
		}

		return fileInfo, true

	} else if err != nil {

		dr.logger.Println(err)

	}

	return nil, false

}

func acquireSemaConnections() {
	semaMaxConnections <- struct{}{}
}

func releaseSemaConnections() {
	_ = <-semaMaxConnections
}

func LogOutput() io.Writer {
	return logOutput
}

func SetLogOutput(out io.Writer) {
	logOutput = out
}
