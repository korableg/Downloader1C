package args

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	ErrLoginEmpty    = errors.New("Login not filled in")
	ErrPasswordEmpty = errors.New("Password not filled in")
	ErrNotDir        = errors.New("The path is not a directory")

	login        string
	password     string
	path         string
	nicksRaw     string
	logPath      string
	startDateRaw string
	instance     string

	fs *flag.FlagSet
)

func init() {

	fs = flag.NewFlagSet("downloader", flag.ExitOnError)

	fs.StringVar(&login, "login", "", "Login to releases.1c.ru (Required)")
	fs.StringVar(&password, "password", "", "Password to releases.1c.ru (Required)")
	fs.StringVar(&path, "path", "."+string(os.PathSeparator), "The directory to save the archives")
	fs.StringVar(&nicksRaw, "nicks", "", "Comma separated string (example: platform83, EnterpriseERP20, hrm)")
	fs.StringVar(&logPath, "log", "./downloader.log", "Path to log file")
	fs.StringVar(&startDateRaw, "startdate", "", "Minimum release date (example: 01.01.2020)")
	fs.StringVar(&instance, "instance", "Downloader1C", "Instance name")

	i := 1
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		i = 2
	}

	fs.Parse(os.Args[i:])

}

func Login() (string, error) {
	if len(login) == 0 {
		return "", ErrLoginEmpty
	}
	return login, nil
}

func Password() (string, error) {
	if len(password) == 0 {
		return "", ErrPasswordEmpty
	}
	return password, nil
}

func Path() (string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !stat.IsDir() {
		return "", ErrNotDir
	}
	return path, nil
}

func Nicks() (map[string]bool, error) {
	var nicksM map[string]bool
	if len(nicksRaw) > 0 {
		nicksS := strings.Split(nicksRaw, ",")
		nicksM = make(map[string]bool, 0)
		for _, nick := range nicksS {
			nicksM[strings.Trim(nick, " ")] = true
		}
	}
	return nicksM, nil
}

func LogFile() (*os.File, error) {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return nil, err
	}
	return logFile, err
}

func StartDate() (time.Time, error) {
	if startDateRaw == "" {
		return time.Unix(0, 0), nil
	}
	return time.Parse("02.01.2006", startDateRaw)
}

func Instance() (string, error) {
	return instance, nil
}

func Usage(err error) {
	fmt.Println("© Downloader 1C by Dmitry Titov\nGithub: github.com/korableg/Downloader1C, E-mail: titov-de@yandex.ru, 2020")
	fmt.Println("Error:", err)
	fmt.Println("Help:")
	fs.PrintDefaults()
	os.Exit(2)
}
