package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

//SetupDefaultConfig - setup the default config options that can be overridden via the config file
func SetupDefaultConfig() {
	viper.SetDefault("logging.level", "info")
}

/*SetupConfig - setup the configuration system */
func SetupConfig() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigName("0proxy")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

const (
	DeploymentDevelopment = 0
	DeploymentTestNet     = 1
	DeploymentMainNet     = 2
)

/*Config - all the config options passed from the command line*/
type Config struct {
	Port            int
	ChainID         string
	DeploymentMode  byte
	SignatureScheme string
	BlockWorker     string
}

/*Configuration of the system */
var Configuration Config

/*TestNet is the program running in TestNet mode? */
func TestNet() bool {
	return Configuration.DeploymentMode == DeploymentTestNet
}

/*Development - is the programming running in development mode? */
func Development() bool {
	return Configuration.DeploymentMode == DeploymentDevelopment
}
