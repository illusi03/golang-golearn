package main

import (
	"os"
	"strings"

	"github.com/illusi03/golearn/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}
	cmd.InitApi(&cmd.ApiConfig{
		Port:   viper.GetInt("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	})
}
