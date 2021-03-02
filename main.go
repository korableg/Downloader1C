package main

import (
	"korableg/Downloader1C/args"
	"korableg/Downloader1C/downloader"
	"log"
)

func main() {

	login, err := args.Login()
	handleError(err)
	password, err := args.Password()
	handleError(err)
	path, err := args.Path()
	handleError(err)
	logFile, err := args.LogFile()
	handleError(err)
	defer logFile.Close()
	nicks, err := args.Nicks()
	handleError(err)
	startDate, err := args.StartDate()
	handleError(err)

	downloader.SetLogOutput(logFile)

	downldr := downloader.New(login, password, path, startDate, nicks)
	_, err = downldr.Get()
	if err != nil {
		log.Fatal(err)
	}

}

func handleError(err error) {
	if err != nil {
		args.Usage(err)
	}
}
