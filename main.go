package main

import (
	"Downloader1C/downloader"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	var nicksM map[string]bool

	pLogin := flag.String("login", "", "Login to releases.1c.ru (Required)" )
	pPassword := flag.String("password", "", "Password to releases.1c.ru (Required)")
	pPath := flag.String("path", "./", "Path to save file")
	pNicks := flag.String("nicks", "", "Comma separated string (example: platform83, EnterpriseERP20, hrm)")
	pLogPath := flag.String("log", "./downloader.log", "Path to log file")

	flag.Parse()

	checkRequiredStringFlag(pLogin, fmt.Errorf("%s", "Login not filled in"))
	checkRequiredStringFlag(pPassword, fmt.Errorf("%s", "Password not filled in"))

	logFile, err := os.OpenFile(*pLogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	downloader.SetLogOutput(logFile)

	nicksS := strings.Split(*pNicks, ",")
	if len(nicksS) > 0 {
		nicksM = make(map[string]bool, 0)
		for _, nick := range nicksS {
			nicksM[strings.Trim(nick, " ")] = true
		}
	}

	downldr := downloader.New(*pLogin, *pPassword, *pPath, nicksM)
	_, err = downldr.Get()
	if err != nil {
		log.Fatal(err)
	}

}

func checkRequiredStringFlag(flagPtr *string, err error) {
	if len(*flagPtr) == 0 {
		fmt.Println(err)
		flag.Usage()
		os.Exit(2)
	}
}