package helper

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/evilsocket/opensnitch/daemon/log"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FindFileByName(path string) string {
	matches, err := filepath.Glob(path + "*")

	if err != nil {
		return ""
	}

	if len(matches) != 0 {
		return matches[0]
	}

	return ""
}

func FileDownload(URL, domain string, fileName string) error {
	url, _ := url.Parse(URL)
	url.Scheme = "http"
	if url.Host == "" {
		url.Host = domain
	}

	log.Info("Downloading file from: ", url.String())

	response, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Bad response code during download request")
	}

	os.Mkdir(path.Dir(fileName), 0777)

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
