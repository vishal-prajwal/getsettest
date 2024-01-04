package downloader

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

func DownloadFile(url string, destFile string) (path string, err error) {
	s := strconv.Itoa(rand.Intn(10000))
	dest := filepath.Join(os.TempDir(), s, destFile)
	err = ensureBaseDir(dest)
	if err != nil {
		err = errors.Wrapf(err, "Unable to open file %q", dest)
		return
	}
	file, err := os.Create(dest)
	if err != nil {
		err = errors.Wrapf(err, "Unable to open file %q", dest)
		return
	}
	defer file.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return dest, nil
}

func ensureBaseDir(fpath string) error {
	baseDir := path.Dir(fpath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(baseDir, 0755)
}
