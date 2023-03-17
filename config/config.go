package config

import (
	"context"
	"fmt"
	"reflect"

	"bitbucket.org/junglee_games/getsetgo/clients/aws"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type SM interface {
	GetFromSM(ctx context.Context, key string) (aws.Secrets, error)
}

type SMStruct interface {
	GetSecretKey() string
	SetSecret(aws.Secrets)
}

func load(sm SM, s SMStruct) error {
	v, err := sm.GetFromSM(context.Background(), s.GetSecretKey())
	if err != nil {
		return fmt.Errorf(" secrete manager error : %s", err.Error())
	}
	s.SetSecret(v)
	return nil
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
		sm, err := aws.NewSecreteManager()
		if err != nil {
			return errors.Wrap(err, "while starting aws for sm")
		}
		return LoadFromSM(sm, config)
	}
	return nil
}

func getValue(a any) reflect.Value {
	value := reflect.ValueOf(a)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

//to load config from secret manager just implement above interface and attach it to configuration , LoadConfig will automatically determines that this needs to be loaded from secret manager
func LoadFromSM(sm SM, config interface{}) error {
	smStructType := reflect.TypeOf((*SMStruct)(nil)).Elem()
	value := getValue(config)
	if value.Kind() != reflect.Struct {
		return nil
	}
	noOfFields := value.NumField()
	for i := 0; i < noOfFields; i++ {
		field := value.Field(i)
		if field.Type().Implements(smStructType) {
			if field.IsNil() {
				logger.Warn(context.Background(), fmt.Sprintf("%s is not set in config ", field.String()))
				continue
			}
			var s SMStruct = field.Interface().(SMStruct)
			if err := load(sm, s); err != nil {
				return fmt.Errorf("%s => %s", field.String(), err.Error())
			}
		}
		err := LoadFromSM(sm, field.Interface())
		if err != nil {
			return fmt.Errorf("%s => %s", field.String(), err.Error())
		}
	}
	return nil
}
