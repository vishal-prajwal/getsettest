package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	"bitbucket.org/junglee_games/getsetgo/clients/aws"
	"github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host             string
	Username         string
	Password         string
	AppName          string
	Build            string
	SetDirect        bool
	SecretManagerKey string
}

func (c *Config) GetSecretKey() string {
	return c.SecretManagerKey
}
func (c *Config) SetSecret(s aws.Secrets) {
	c.Username = s.Username
	c.Password = s.Password
}

func GetMongoClient(cfg *Config) (*mongo.Client, error) {
	clientOptions := options.Client().SetAppName(cfg.AppName)
	clientOptions.SetDirect(cfg.SetDirect)
	clientOptions.SetMaxConnIdleTime(time.Minute * 10)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	hosts := strings.Split(cfg.Host, ",")
	if len(hosts) == 0 {
		return &mongo.Client{}, errors.New("missing db host")
	}

	clientOptions.SetHosts(hosts)

	if cfg.Build != "local" {
		clientOptions.SetAuth(options.Credential{
			AuthSource:    "admin",
			AuthMechanism: "SCRAM-SHA-256",
			Username:      cfg.Username,
			Password:      cfg.Password,
		})
	}

	nrMon := nrmongo.NewCommandMonitor(nil)
	clientOptions.SetMonitor(nrMon)

	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}
