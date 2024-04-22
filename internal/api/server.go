package api

import (
	"github.com/ngirchev/settings-loader/internal/util"
	"net"
	"net/rpc"
)

type RpcServer struct {
	loaderController *LoaderController
}

func NewRpcServer(loaderController *LoaderController) *RpcServer {
	return &RpcServer{loaderController: loaderController}
}

type ServerProps struct {
	BindAddress string
}

func (s *RpcServer) Start(config ServerProps) {
	err := rpc.Register(s.loaderController)
	util.HandleError("Register controller error:", err)

	listener, err := net.Listen("tcp", config.BindAddress)
	util.HandleError("Listen error:", err)

	defer func(listener net.Listener) {
		err := listener.Close()
		util.HandleError("listener close error", err)
	}(listener)

	rpc.Accept(listener)
}
