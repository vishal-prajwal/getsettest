package encryption

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"bitbucket.org/junglee_games/getsetgo/utils/files"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func ProtectPDFBytes(input []byte, password string) ([]byte, error) {
	// Create an in-memory reader for the input PDF bytes
	inputReader := bytes.NewReader(input)

	// Create an in-memory writer for the output PDF bytes
	var outputWriter bytes.Buffer

	// Protect the PDF using ProtectPdf function
	err := ProtectPdf(inputReader, &outputWriter, password, password)
	if err != nil {
		return nil, err
	}

	// Return the protected PDF bytes
	return outputWriter.Bytes(), nil
}

func ProtectPDFFile(path, filename, password string) error {
	file, err := os.Open(path + "/" + filename)
	if err != nil {
		return err
	}
	var outputWriter bytes.Buffer

	// Protect the PDF using ProtectPdf function
	err = ProtectPdf(file, &outputWriter, password, password)
	if err != nil {
		return err
	}
	return files.SaveFile(path, filename, outputWriter.String())
}

func ProtectPdf(inputFile io.ReadSeeker, outputWriter io.Writer, userPassword, ownerPassword string) error {
	// Create a new PDF configuration
	conf := model.NewDefaultConfiguration()

	// Set encryption settings
	conf.UserPW = userPassword
	conf.OwnerPW = ownerPassword
	conf.Permissions = model.PermissionsAll

	// Encrypt the PDF file
	err := api.Encrypt(inputFile, outputWriter, conf)
	if err != nil {
		return fmt.Errorf("error encrypting PDF: %v", err)
	}

	return nil
}
