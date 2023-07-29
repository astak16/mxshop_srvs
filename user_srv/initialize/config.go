package initialize

import (
	"mxshow_srvs/user_srv/global"

	"github.com/spf13/viper"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	debug := GetEnvInfo("MXSHOP_DEBUG")
	if debug {
		viper.SetConfigName("config-debug")
	} else {
		viper.SetConfigName("config-pro")
	}

	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
