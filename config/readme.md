# config loader

`Its a more general config loader for go.`

## Params

1. env string : env in which program is running, if env is local then secret manager values will not be retrived.
2. configPath string : configPath is a path of directory where ($env).yml file is present.
3. config interface{} : its a address of config variable in which we want to load a config

## secret manager

`lets say we want to take a password from secret manager for ESConfig struct then ESConfig will implement SMStruct interface and rest will be taken care by LoadConfig`

checkout the example for more details
