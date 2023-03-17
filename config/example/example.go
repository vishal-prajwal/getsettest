package main

import (
	"context"
	"encoding/json"

	"bitbucket.org/junglee_games/getsetgo/clients/aws"
	"bitbucket.org/junglee_games/getsetgo/config"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

type Config struct {
	CB *CBConfig
}

type CBConfig struct {
	URL      string
	Secret   string
	Password string
}

// implementing SMStruct interface as password should be taken from sm

func (c *CBConfig) GetSecretKey() string {
	return c.Secret
}
func (c *CBConfig) SetSecret(s aws.Secrets) {
	c.Password = s.Password
}

type SM struct {
}

func (sm *SM) GetFromSM(ctx context.Context, key string) (aws.Secrets, error) {
	return aws.Secrets{Username: "sm-username", Password: "sm-password", Token: "sm-token"}, nil
}

func main() {
	var c Config
	// giving env as local as I dont have aws secret implementing my own sm for it
	err := config.LoadConfig("local", ".", &c)
	if err != nil {
		logger.Panic(context.Background(), err.Error())
	}
	// loading a readed config from custom sm manager
	// Note: no need to call LoadFromSM if env is other than local its just for example here
	config.LoadFromSM(&SM{}, &c)
	bytes, _ := json.Marshal(&c)
	logger.Info(context.Background(), string(bytes))
}
