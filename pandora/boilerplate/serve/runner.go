package serve

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/tools/go/packages"

	"bitbucket.org/junglee_games/getsetgo/pandora/log"
)

func NewRunner(name, path string, initialPortIndex int) (Runner, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	reloadCh := make(chan string)

	return &runner{
		watcher:  watcher,
		reloadCh: reloadCh,

		name:             name,
		path:             path,
		initialPortIndex: initialPortIndex,
	}, nil
}

type runner struct {
	cmd     *exec.Cmd
	watcher *fsnotify.Watcher

	reloadCh chan string

	name             string
	path             string
	initialPortIndex int
	shutdown         bool

	dependencies []string
}

func (r *runner) InitialPortIndex() int {
	return r.initialPortIndex
}

func (r *runner) Run(ctx context.Context) <-chan string {
	changesCh := make(chan string)

	go r.watchForChanges(ctx)
	go r.buildDependencies()

	go debounce(time.Second*2, r.reloadCh, func(arg string) {
		r.Interrupt()
		r.buildDependencies()

		changesCh <- arg
	})

	go func() {
		<-ctx.Done()

		r.shutdown = true
		r.Interrupt()
	}()

	go func() {
		defer close(changesCh)

		for {
			if r.shutdown {
				return
			}

			path, err := r.build()
			if err != nil {
				log.
					WithError(err).
					Errorf("failed to build")

				time.Sleep(time.Second)

				continue
			}

			log.Infof("Starting service: %s", r.name)

			if err := r.run(path); err != nil {
				log.WithError(err).Errorf("failed to start service")
			}
		}
	}()

	return changesCh
}

func (r *runner) Interrupt() {
	log.Infof("Restarting service: %s", r.name)

	if err := r.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.WithError(err).Errorf("failed to send interrupt signal")
	}
}

func (r *runner) HasDependency(dep string) bool {
	for _, pkg := range r.dependencies {
		if pkg == dep {
			return true
		}
	}

	return false
}

func (r *runner) run(path string) error {
	cmd := exec.Command(path, "server")
	cmd.Dir = r.path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// todo: find a better way to map the existing env to the new process
	cmd.Env = []string{
		"LOG_DISABLE_REMOTE=true",

		fmt.Sprintf("NOMAD_ADDR_http=127.0.0.1:%d", r.initialPortIndex),
		fmt.Sprintf("NOMAD_ADDR_grpc=127.0.0.1:%d", r.initialPortIndex+1),
	}

	r.cmd = cmd

	return cmd.Run()
}

func (r *runner) build() (string, error) {
	out := fmt.Sprintf("/tmp/%s", r.name)

	cmd := exec.Command("/usr/local/bin/go", "build", "-o", out, "main.go")
	cmd.Dir = r.path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// todo: find a better way to map the existing env to the new process
	cmd.Env = []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("GOPATH=%s", os.Getenv("GOPATH")),
	}

	return out, cmd.Run()
}

func (r *runner) watchForChanges(ctx context.Context) {
	err := filepath.Walk(r.path, func(path string, info fs.FileInfo, err error) error {
		// We only care about go files
		if !info.IsDir() && filepath.Ext(path) != ".go" {
			return nil
		}

		log.Debugf("Watching file %s", path)

		if err := r.watcher.Add(path); err != nil {
			log.WithError(err).Warnf("failed to watch %s", path)
		}

		return nil
	})

	if err != nil {
		log.WithError(err).Errorf("failed to scan folders to watch")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		case e, ok := <-r.watcher.Errors:
			if !ok {
				return
			}

			log.WithError(e).Warnf("Error occurred while watching files")

		case e, ok := <-r.watcher.Events:
			if !ok {
				return
			}

			r.reloadCh <- e.Name
		}
	}
}

func (r *runner) buildDependencies() {
	startsAt := time.Now()
	defer log.Debugf("building dependencies for %s took %s\n", r.name, time.Since(startsAt).String())

	pkg := filepath.Join("gitlab.com/jungleegames/backend", r.path, "cmd")

	o, err := packages.Load(&packages.Config{Mode: packages.NeedImports}, pkg)
	if err != nil {
		panic(err)
	}

	var dependencies []string

	scannedPackages := make(map[string]bool)

	for _, p := range o {
		scanPackagesForImports(scannedPackages, &dependencies, p, r.name)
	}

	r.dependencies = dependencies
}

func debounce(interval time.Duration, input chan string, cb func(arg string)) {
	var item string

	timer := time.NewTimer(interval)

	for {
		select {
		case item = <-input:
			timer.Reset(interval)

		case <-timer.C:
			if item != "" {
				cb(item)
			}
		}
	}
}

func scanPackagesForImports(
	scannedPackages map[string]bool,
	dependencies *[]string,
	p *packages.Package,
	currentPackage string,
) {
	for packageName, pkg := range p.Imports {
		if _, ok := scannedPackages[packageName]; ok {
			continue
		}

		scannedPackages[packageName] = true

		if !strings.Contains(packageName, "gitlab.com/jungleegames/backend") {
			continue
		}

		if strings.Contains(packageName, currentPackage) {
			scanPackagesForImports(scannedPackages, dependencies, pkg, currentPackage)
		} else {
			*dependencies = append(*dependencies, packageName)
		}
	}
}
