package zipper

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrInvalidGZipFile        = fmt.Errorf("invalid gzip file")
	ErrReadingGZipFileContent = fmt.Errorf("error reading gzip file")
)

// GZip takes the content of the file and compress it and return new name with .gz and compressed content.
func GZip(name string, content string) (string, string, error) {
	reader := strings.NewReader(content)
	buffer := new(bytes.Buffer)
	// Create a gzip reader
	gzipWriter := gzip.NewWriter(buffer)

	// Copy the original file to the gzip file
	_, err := io.Copy(gzipWriter, reader)
	if err != nil {
		return "", "", errors.Wrap(ErrReadingGZipFileContent, err.Error())
	}

	gzipWriter.Close()

	return name + ".gz", buffer.String(), nil
}

//GunZip will uncompress the content and returns the new name with no .gz and uncompressed content
func GunZip(name string, content string) (string, string, error) {
	if len(name) < 3 || name[len(name)-3:] != ".gz" {
		return "", "", ErrInvalidGZipFile
	}
	gzipFileReader := strings.NewReader(content)
	// Create a gzip reader
	gzipReader, err := gzip.NewReader(gzipFileReader)
	if err != nil {
		return "", "", errors.Wrap(ErrReadingGZipFileContent, err.Error())
	}
	defer gzipReader.Close()

	bytes, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return "", "", errors.Wrap(ErrReadingGZipFileContent, err.Error())
	}

	return name[:len(name)-3], string(bytes), nil
}
