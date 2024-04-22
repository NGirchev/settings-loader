package cmd

import (
	"fmt"
	"github.com/ngirchev/settings-loader/internal/api"
	"github.com/ngirchev/settings-loader/internal/domain/postgresql"
	"github.com/ngirchev/settings-loader/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LoadComponentMethod = "LoaderController.LoadComponent"
)

func InitConfig() error {
	configPath := util.GetPath("configs")
	log.Debugf("Config path:" + configPath)
	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("yml")
	return viper.ReadInConfig()
}

func SetupLogging() {
	logLevel := util.GetVar("app.logLevel")

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

type AppConf struct {
	ServerConf api.ServerProps
	Hash       string
	Path       string
}

func BuildPostgreSQLConf() postgresql.DBConf {
	return postgresql.DBConf{
		Host:     util.GetVar("db.host"),
		Port:     util.GetVar("db.port"),
		Username: util.GetVar("db.username"),
		Password: util.GetVar("db.password"),
		DBName:   util.GetVar("db.dbname"),
		SSLMode:  util.GetVar("db.sslMode"),
	}
}

func BuildAppConf() AppConf {
	conf := AppConf{
		ServerConf: api.ServerProps{
			BindAddress: util.GetVar("app.api.bindAddress"),
		},
		Hash: util.GetVar("app.hash"),
		Path: util.GetVar("app.path"),
	}
	validate(conf)
	return conf
}

func validate(conf AppConf) {
	if conf.Hash != "md5" &&
		conf.Hash != "sha256" {
		panic(fmt.Sprintf("Couldn't initialize the conf file, because hash function=%s is invalid",
			conf.Hash))
	}
}
