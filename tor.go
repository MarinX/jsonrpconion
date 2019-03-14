package jsonrpconion

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/cretz/bine/tor"
)

var torProcess *tor.Tor
var onceStart sync.Once
var onceClose sync.Once

// getTorProcess returns Tor process with default config
func getTorProcess() (*tor.Tor, error) {
	var err error
	onceStart.Do(func() {
		os.Mkdir(".data", 0700)
		file, _ := os.Create(".data/rpctorrc")
		if file != nil {
			file.Close()
		}

		dir, _ := filepath.Abs(".data")
		torrc, _ := filepath.Abs(".data/rpctorrc")

		torProcess, err = tor.Start(nil, &tor.StartConf{
			DataDir:   dir,
			TorrcFile: torrc,
		})
	})
	return torProcess, err
}

// closeTorProcess closes Tor connection
func closeTorProcess() error {
	var err error
	onceClose.Do(func() {
		if torProcess != nil {
			err = torProcess.Close()
		}
	})
	return err
}
