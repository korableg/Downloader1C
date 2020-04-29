package args

import (
	"errors"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {

	_, err := Login()
	if err != nil && !errors.Is(err, ErrLoginEmpty) {
		t.Error(err)
	}

}

func TestPassword(t *testing.T) {

	_, err := Password()
	if err != nil && err != ErrPasswordEmpty {
		t.Error(err)
	}

}

func TestPath(t *testing.T) {

	_, err := Path()
	if err != nil && !os.IsNotExist(err) {
		t.Error(err)
	}

}

func TestNicks(t *testing.T) {

	_, err := Nicks()
	if err != nil {
		t.Error(err)
	}

}

func TestLogFile(t *testing.T) {

	logfile, err := LogFile()
	if err != nil {
		t.Error(err)
	} else {
		err = logfile.Close()
		if err != nil {
			t.Error(err)
		} else {
			err = os.Remove(logfile.Name())
			if err != nil {
				t.Error(err)
			}
		}
	}

}

func TestStartDate(t *testing.T) {

	_, err := StartDate()
	if err != nil {
		t.Error(err)
	}

}

func TestInstance(t *testing.T) {

	_, err := Instance()
	if err != nil {
		t.Error(err)
	}

}
