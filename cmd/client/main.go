package main

import (
	"fmt"
	"github.com/ngirchev/settings-loader/cmd"
	"github.com/ngirchev/settings-loader/internal/api"
	"github.com/ngirchev/settings-loader/internal/util"
	log "github.com/sirupsen/logrus"
	"net/rpc"
)

func main() {
	if err := cmd.InitConfig(); err != nil {
		panic(fmt.Sprintf("Couldnt initialize the config file: %s", err))
	}
	cmd.SetupLogging()
	client, err := rpc.Dial("tcp", "localhost"+cmd.BuildAppConf().ServerConf.BindAddress)
	util.HandleError("Client error: ", err)

	var response api.Response
	req := &api.Request{Version: "2.0.0"}

	err = client.Call(cmd.LoadComponentMethod, req, &response)
	util.HandleError("Client error: ", err)
	log.Println(response)
}
