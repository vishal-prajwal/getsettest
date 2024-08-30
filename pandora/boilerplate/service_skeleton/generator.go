package skeleton

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	sampleService    = "sample"
	sampleImportPath = "bitbucket.org/junglee_games/getsetgo/pandora/boilerplate/service_skeleton"
	newImportPath    = "gitlab.com/jungleegames/backend/src"
)

func GenerateServiceSkeleton(serviceName string) {

	// get working directory
	dir, _ := os.Getwd()

	serviceFolder := filepath.Join(dir, serviceName)
	sample := filepath.Join(dir, "boilerplate/service_skeleton/sample")

	// copy the skeleton of new service from the sample service
	fmt.Println(sample)
	fmt.Println(serviceFolder)
	err := copyDir(sample, serviceFolder)
	if err != nil {
		panic(err)
	}

	// Go through the contents of files and replace sample filename or sample word from the files to new service

	// replacing for cmd folder
	cmd := filepath.Join(serviceFolder, "cmd")

	rootFile := filepath.Join(cmd, "root.go")
	replaceFileContent(rootFile, serviceName)

	serverFile := filepath.Join(cmd, "server.go")
	replaceFileContent(serverFile, serviceName)

	// replacing in all the files place at the root level of service structure
	envFile := filepath.Join(serviceFolder, ".env")
	replaceFileContent(envFile, serviceName)

	mainFile := filepath.Join(serviceFolder, "main.go")
	replaceFileContent(mainFile, serviceName)

	gitlabCIFile := filepath.Join(serviceFolder, ".gitlab-ci.yml")
	replaceFileContent(gitlabCIFile, serviceName)

	dockerFile := filepath.Join(serviceFolder, "/Dockerfile")
	replaceFileContent(dockerFile, serviceName)

	// replacing in provisioning folder
	provisioning := filepath.Join(serviceFolder, "/provisioning")

	stagFile := filepath.Join(provisioning, "staging.hcl")
	replaceFileContent(stagFile, serviceName)

	prodFile := filepath.Join(provisioning, "prod.hcl")
	replaceFileContent(prodFile, serviceName)

	// replacing in pkg folder

	pkg := filepath.Join(serviceFolder, "pkg")

	client := filepath.Join(pkg, "client")
	clientSampleFile := filepath.Join(client, "sample.go")
	clientFile := filepath.Join(client, serviceName+".go")

	if err := os.Rename(clientSampleFile, clientFile); err != nil {
		panic(err)
	}

	replaceFileContent(clientFile, serviceName)

	// replacing in internal folder
	internal := filepath.Join(serviceFolder, "internal")

	internalSampleFile := filepath.Join(internal, "sample.go")
	internalServiceFile := filepath.Join(internal, serviceName+".go")
	err = os.Rename(internalSampleFile, internalServiceFile)
	if err != nil {
		panic(err)
	}

	// replacing in http folder
	http := filepath.Join(internal, "http")

	httpServerFile := filepath.Join(http, "server.go")
	replaceFileContent(httpServerFile, serviceName)

	grpc := filepath.Join(internal, "grpc")
	serverFile = filepath.Join(grpc, "server.go")
	replaceFileContent(serverFile, serviceName)

	svc := filepath.Join(grpc, "svc")
	newFile := filepath.Join(svc, serviceName+".go")
	err = os.Rename(svc+"/sample.go", newFile)
	if err != nil {
		panic(err)
	}

	replaceFileContent(newFile, serviceName)

	fmt.Println(" ************ TODO's ****************")
	fmt.Println(" 1) Add your new db info in create_databases.sh file")
	fmt.Println(" 2) Create a new project on sentry.io and replace SENTRY_DSN value in both staging/production deployment files")
	fmt.Println(" 3) After generating required proto for your service, uncomment the code related to setting up grpc server and client")
	fmt.Println(" 4) Enough of generated code, go enjoy writing some code of your own now 😛")
	fmt.Println(" ************ END ****************")
}

func copyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		if err != nil {
			err = os.Chmod(dest, 0644)
		}
	}
	return
}

func copyDir(source string, dest string) (err error) { // create dest dir
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, obj := range objects {
		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = copyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		} else {
			// perform copy
			err = copyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func replaceFileContent(file, newServiceName string) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, sampleImportPath) {
			lines[i] = strings.ReplaceAll(lines[i], sampleImportPath, newImportPath)
		}

		if strings.Contains(line, sampleService) {
			lines[i] = strings.ReplaceAll(lines[i], sampleService, newServiceName)
		}
	}

	output := strings.Join(lines, "\n")

	// get the original file permissions
	info, err := os.Stat(file)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(file, []byte(output), info.Mode())
	if err != nil {
		panic(err)
	}
}
