package winquit

import (
	"os"
	"time"

	"github.com/n1hility/winquit/pkg/winquit/win32"
	"github.com/sirupsen/logrus"
)

func RequestQuit(pid int) error {
	threads, err := win32.GetProcThreads(uint32(pid))
	if err != nil {
		return err
	}

	for _, thread := range threads {
		logrus.Debugf("Closing windows on thread %d", thread)
		win32.CloseThreadWindows(uint32(thread))
	}

	return nil
}

func QuitProcess(pid int, waitNicely time.Duration) error {
	_ = RequestQuit(pid)
	
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil
	}

	done := make(chan bool)
	
	go func() {
		proc.Wait()
		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-time.After(waitNicely):
	}

	return proc.Kill()
}