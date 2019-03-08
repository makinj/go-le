// Package config uses viper to load configurations into an object
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Load uses Viper to load a configuration file into an object
func Load(cPath string, cName string, config interface{}) (err error) {
	// Setting up Viper with Config file Path= cPath and ConfigName=%s", cPath, cName)
	viper.SetConfigName(cName)
	viper.AddConfigPath(cPath)

	//Read Viper config
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Viper error reading config file %s%s: %s", cPath, cName, err)
	}

	//Unmarshal config
	err = viper.Unmarshal(config)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}
