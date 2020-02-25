package config

import (
	"github.com/tkanos/gonfig"
	"github.com/yossefazoulay/go_utils/utils"
)

func GetConfig(env string, configuration interface{}) {
	var configEnv = make(map[string]string)
	configEnv["dev"] = "./config/config.dev.json"
	configEnv["prod"] = "./config/config.prod.json"
	err := gonfig.GetConf(configEnv[env], configuration)
	utils.HandleError(err, "Cannot load/read config file")
}

