package serve

import (
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/consul/api"
	"golang.org/x/sync/errgroup"

	consulClient "bitbucket.org/junglee_games/getsetgo/pandora/consul"
)

type Opt func(coordinator *Coordinator)

func WithSkippedServices(services []string) Opt {
	return func(coordinator *Coordinator) {
		coordinator.skippedServices = services
	}
}

type Coordinator struct {
	runners   map[string]Runner
	portIndex int
	consul    consulClient.Client

	// config
	skippedServices []string
}

func (c *Coordinator) Start(ctx context.Context, root string, opts ...Opt) error {
	for _, opt := range opts {
		opt(c)
	}

	consul, err := consulClient.NewClient(consulClient.ConsulConfig{})
	if err != nil {
		return err
	}

	c.runners = make(map[string]Runner)
	c.portIndex = 5000
	c.consul = consul

	if err := c.buildRunners(root); err != nil {
		return err
	}

	if err := c.registerConsul(); err != nil {
		return err
	}

	return c.run(ctx)
}

func (c *Coordinator) buildRunners(root string) error {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

scan:
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// no recursion plz
		if file.Name() == "pandora" {
			continue
		}

		for _, name := range c.skippedServices {
			if strings.EqualFold(name, file.Name()) {
				continue scan
			}
		}

		runner, err := NewRunner(file.Name(), path.Join(root, file.Name()), c.portIndex)
		if err != nil {
			return err
		}

		c.runners[file.Name()] = runner

		// Each service has a http and grpc server, so we need to increment by two
		c.portIndex += 2
	}

	return nil
}

func (c *Coordinator) run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, runner := range c.runners {
		func(r Runner) {
			g.Go(func() error {
				for changes := range r.Run(ctx) {
					pkg := filepath.Join("gitlab.com/jungleegames/backend", filepath.Dir(changes))

					for _, runner := range c.runners {
						if runner.HasDependency(pkg) {
							runner.Interrupt()
						}
					}
				}

				return nil
			})
		}(runner)
	}

	return g.Wait()
}

func (c *Coordinator) registerConsul() error {
	for svc, r := range c.runners {
		if err := c.consul.API().Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:      fmt.Sprintf("svc-%s-http", svc),
			Name:    fmt.Sprintf("svc-%s-http", svc),
			Port:    r.InitialPortIndex(),
			Address: "127.0.0.1",
		}); err != nil {
			return err
		}

		if err := c.consul.API().Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:      fmt.Sprintf("svc-%s-grpc", svc),
			Name:    fmt.Sprintf("svc-%s-grpc", svc),
			Port:    r.InitialPortIndex() + 1,
			Address: "127.0.0.1",
		}); err != nil {
			return err
		}
	}

	return nil
}

type Runner interface {
	Run(ctx context.Context) <-chan string
	Interrupt()
	InitialPortIndex() int
	HasDependency(pkg string) bool
}
