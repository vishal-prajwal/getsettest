package config

import (
	"context"
	"fmt"
	"reflect"

	"bitbucket.org/junglee_games/getsetgo/aws"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type SM interface {
	GetFromSM(key string) (aws.Secrets, error)
}

type SMStruct interface {
	GetSecretKey() string
	SetSecret(aws.Secrets)
}

// more general load config function
func LoadConfig(env, path string, config interface{}) error {
	logger.Info(context.Background(), path)
	// =========================================================================
	// Configuration from file
	// Set the file name of the configurations file
	viper.SetConfigName(env)
	// Set the path to look for the configurations file
	viper.AddConfigPath(path)
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "Error reading config file")
	}
	if err := viper.Unmarshal(config); err != nil {
		return errors.Wrap(err, "Unable to decode into struct")
	}

	if env != "local" {
		a, err := aws.New()
		if err != nil {
			return errors.Wrap(err, "while starting aws for sm")
		}
		return LoadFromSM(a, config)
	}
	return nil
}

//to load config from secret manager just implement above interface and attach it to configuration , LoadConfig will automatically determines that this needs to be loaded from secret manager
func LoadFromSM(sm SM, config interface{}) error {
	smStructType := reflect.TypeOf((*SMStruct)(nil)).Elem()
	values := reflect.ValueOf(config)
	if values.Kind() == reflect.Ptr {
		values = values.Elem()
	}
	if values.Kind() != reflect.Struct {
		return nil
	}
	noOfFields := values.NumField()
	for i := 0; i < noOfFields; i++ {
		field := values.Field(i)
		if field.Type().Implements(smStructType) {
			if field.IsNil() {
				logger.Warn(context.Background(), fmt.Sprintf("%s is not set in config ", field.String()))
				continue
			}
			var s SMStruct = field.Interface().(SMStruct)
			v, err := sm.GetFromSM(s.GetSecretKey())
			if err != nil {
				return err
			}
			s.SetSecret(v)
		}
		err := LoadFromSM(sm, field.Interface())
		if err != nil {
			return fmt.Errorf("%s=>%s", field.String(), err.Error())
		}
	}
	return nil
}
