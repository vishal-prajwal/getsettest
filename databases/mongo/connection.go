package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(host, username, password, appname string) (*mongo.Client, error) {
	clientOptions := options.Client().SetAppName(appname)
	clientOptions.SetMaxConnIdleTime(time.Minute * 10)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	hosts := strings.Split(host, ",")
	if len(hosts) == 0 {
		return &mongo.Client{}, errors.New("missing db host")
	}

	clientOptions.SetHosts(hosts)

	clientOptions.SetAuth(options.Credential{
		AuthSource:    "admin",
		AuthMechanism: "SCRAM-SHA-256",
		Username:      username,
		Password:      password,
	})

	nrMon := nrmongo.NewCommandMonitor(nil)
	clientOptions.SetMonitor(nrMon)

	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}
