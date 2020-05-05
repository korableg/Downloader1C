// +build windows

package main

import (
	dargs "Downloader1C/args"
	"Downloader1C/downloader"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"sync"
	"time"
)

type service struct{}

var elog debug.Log
var mutex sync.Mutex
var err error
var d = 30 * time.Hour
var ticker *time.Ticker

var login, password, path string
var startDate time.Time
var nicks map[string]bool

type elogWriterAdapter struct{}

func (l *elogWriterAdapter) Write(p []byte) (int, error) {
	err := elog.Info(1, string(p))
	return len(p), err
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	downloader.SetLogOutput(&elogWriterAdapter{})

	login, err = dargs.Login()
	if err != nil {
		elog.Error(1, err.Error())
		handleError(err)
	}
	password, err = dargs.Password()
	if err != nil {
		elog.Error(1, err.Error())
		handleError(err)
	}
	path, err = dargs.Path()
	if err != nil {
		elog.Error(1, err.Error())
		handleError(err)
	}
	startDate, err = dargs.StartDate()
	if err != nil {
		elog.Error(1, err.Error())
		handleError(err)
	}
	nicks, err = dargs.Nicks()
	if err != nil {
		elog.Error(1, err.Error())
		handleError(err)
	}

	changes <- svc.Status{State: svc.StartPending}

	ticker = time.NewTicker(d)
	tick := ticker.C

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	mutex.Lock()
	go getFiles()

loop:
	for {
		select {
		case <-tick:
			mutex.Lock()
			go getFiles()
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				ticker.Stop()
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				ticker = time.NewTicker(d)
				tick = ticker.C
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string, isDebug bool) {
	var err error
	if isDebug {
		elog = debug.New(name)
	} else {
		elog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	if isDebug {
		run = debug.Run
	}
	err = run(name, &service{})
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}

func getFiles() {
	defer mutex.Unlock()
	dwnldr := downloader.New(login, password, path, startDate, nicks)
	_, err = dwnldr.Get()
	if err != nil {
		elog.Error(1, err.Error())
	}
}
