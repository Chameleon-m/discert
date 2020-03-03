package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/Chameleon-m/discert/internal/app/apiserver"
	"log"
	//"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	//EnvType := os.Getenv("ENV_TYPE) // prod,stage,test,dev

	config := apiserver.NewConfig()
	// TODO https://github.com/pelletier/go-toml OR ENV PARAMS
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
