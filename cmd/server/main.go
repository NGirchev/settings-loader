package main

import (
	"fmt"
	"github.com/ngirchev/settings-loader/cmd"
	"github.com/ngirchev/settings-loader/internal/api"
	"github.com/ngirchev/settings-loader/internal/domain/postgresql"
	"github.com/ngirchev/settings-loader/internal/service"
	"github.com/ngirchev/settings-loader/internal/service/json"
	"github.com/ngirchev/settings-loader/internal/util"
)

func main() {
	if err := cmd.InitConfig(); err != nil {
		panic(fmt.Sprintf("Couldnt initialize the config file: %s", err))
	}
	cmd.SetupLogging()

	appProps := cmd.BuildAppConf()
	dbProps := cmd.BuildPostgreSQLConf()

	rpcServer := api.NewRpcServer(createLoader(appProps, dbProps))
	rpcServer.Start(appProps.ServerConf)
}

func createLoader(appConfig cmd.AppConf, dbConf postgresql.DBConf) *api.LoaderController {
	var hasher service.IHasher
	if appConfig.Hash == "md5" {
		hasher = service.NewMD5Hasher()
	} else if appConfig.Hash == "sha256" {
		hasher = service.NewSHA256Hasher()
	}

	repo := postgresql.NewSettingsRepo(dbConf)
	loaderService := service.NewLoaderService(hasher, json.NewJsonReader(), repo, util.GetPath(appConfig.Path))
	return api.NewLoaderController(loaderService)
}
