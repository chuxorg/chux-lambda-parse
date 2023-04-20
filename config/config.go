// Description: This file contains the code to load the configuration file
// and unmarshal it into a struct.
// Author: Chuck Sailer
// Date: 2023-04-07
// Version: 1.0.0
// License: gpl-3.0
// Copyright: Chuck Sailer
// Credits: [Chuck Sailer]
// Maintainer: Chuck Sailer
// Status: Production
package config

import (
	bo "github.com/chuxorg/chux-models/config"
	"github.com/spf13/viper"
)

type DataStoreConfig struct {
	Target         string `mapstructure:"target"`
	URI            string `mapstructure:"uri"`
	Timeout        int    `mapstructure:"timeout"`
	DatabaseName   string `mapstructure:"databaseName"`
	CollectionName string `mapstructure:"collectionName"`
}

// LambdaConfig is the struct that holds the configuration
// for the data store.
// The struct is populated by the configuration file.
// The struct is used to initialize the data store.
// The struct is also used to initialize the logger.
//
// The configuration file is in YAML format.
//
//go:generate cp config.prod.yaml ../../
type LambdaConfig struct {
	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
	AWS struct {
		BucketName    string `mapstructure:"bucketName"`
		DownloadPath  string `mapstructure:"downloadPath"`
		ArchiveBucket string `mapstructure:"archiveBucket"`
		Profile       string `mapstructure:"profile"`
		Region        string `mapstructure:"region"`
		AccessKey     string `mapstructure:"accessKey"`
		SecretKey     string `mapstructure:"secretKey"`
	} `mapstructure:"aws"`
	Auth struct {
		IssuerURL string `mapstructure:"issuerUrl"`
		TokenURL  string `mapstructure:"tokenUrl"`
	} `mapstructure:"auth"`
	DataPath struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"data"`
	DataStores struct {
		// A map of data store configurations keyed by the data store name
		// e.g., "mongo" or "redis"
		DataStoreMap map[string]DataStoreConfig `mapstructure:"dataStore"`
	} `mapstructure:"dataStores"`
	Products []string `mapstructure:"productSources"`
}

func New() *LambdaConfig {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return cfg
}

func LoadConfig() (*LambdaConfig, error) {

	viper.SetConfigName("config.prod")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config LambdaConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func NewBizObjConfig(parserConfig *LambdaConfig) *bo.BizObjConfig {
	return &bo.BizObjConfig{
		Logging: struct {
			Level string `mapstructure:"level"`
		}{
			Level: parserConfig.Logging.Level,
		},
		DataStores: struct {
			DataStoreMap map[string]bo.DataStoreConfig `mapstructure:"dataStore"`
		}{
			DataStoreMap: ConvertDataStoreMap(parserConfig.DataStores.DataStoreMap),
		},
	}
}

func ConvertDataStoreMap(src map[string]DataStoreConfig) map[string]bo.DataStoreConfig {
	dst := make(map[string]bo.DataStoreConfig)
	for k, v := range src {
		dst[k] = bo.DataStoreConfig{
			Target:         v.Target,
			URI:            v.URI,
			Timeout:        v.Timeout,
			DatabaseName:   v.DatabaseName,
			CollectionName: v.CollectionName,
		}
	}
	return dst
}
