package jsonrpconion

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cretz/bine/control"
)

func keyExist() bool {
	path := fmt.Sprintf("%s/%s/%s", cfgDataDir, cfgKeyDir, cfgKeyFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func loadKey() (*control.ED25519Key, error) {
	path := fmt.Sprintf("%s/%s/%s", cfgDataDir, cfgKeyDir, cfgKeyFile)
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, err := control.ED25519KeyFromBlob(string(buff))
	return key, err
}

func saveKey(blob []byte) error {
	enc := base64.StdEncoding.EncodeToString(blob)
	return ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", cfgDataDir, cfgKeyDir, cfgKeyFile), []byte(enc), 0600)
}
