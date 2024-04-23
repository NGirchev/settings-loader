//go:build manual
// +build manual

package main

import (
	"errors"
	"fmt"
	"github.com/ngirchev/settings-loader/cmd"
	"github.com/ngirchev/settings-loader/internal/api"
	"github.com/ngirchev/settings-loader/internal/util"
	"github.com/stretchr/testify/assert"
	"net"
	"net/rpc"
	"strings"
	"testing"
	"time"
)

// run db before
func TestRunAllE2E(t *testing.T) {
	// before all
	if err := cmd.InitConfig(); err != nil {
		panic(fmt.Sprintf("Couldnt initialize the config file: %s", err))
	}
	cmd.SetupLogging()
	appProps := cmd.BuildAppConf()

	host := "localhost" + appProps.ServerConf.BindAddress
	if checkPortAvailability(host) {
		go main()
		time.Sleep(1 * time.Second)
	}

	client, err := rpc.Dial("tcp", "localhost"+appProps.ServerConf.BindAddress)
	util.HandleError("client error", err)

	shouldReturnContentAndDefaultValuesWhenEmptyRequest(t, client)
	shouldReturnContentWhenFullRequest(t, client)
	shouldReturnEmptyContentWhenRequestedHashDifferent(t, client)
	shouldReturnErrorWhenRequestedContentNotFound(t, client)

	// after all
	err = client.Close()
	if err != nil {
		t.Fatalf("RPC call failed: %v", err)
	}
}

func checkPortAvailability(host string) bool {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return true
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return false
}

func startServer() {
	defer func() {
		if r := recover(); r != nil {
			var panicMessage string
			switch v := r.(type) {
			case error:
				panicMessage = v.Error()
			case string:
				panicMessage = v
			default:
				panicMessage = fmt.Sprintf("unknown panic: %v", v)
			}
			if !strings.HasSuffix(panicMessage, "bind: address already in use") {
				util.HandleError("Server can't start", errors.New(panicMessage))
			}
		}
	}()

	main()
}

func shouldReturnContentAndDefaultValuesWhenEmptyRequest(t *testing.T, client *rpc.Client) {
	// given
	req := api.Request{}
	var actual api.Response
	// when
	err := client.Call(cmd.LoadComponentMethod, req, &actual)
	// then
	if err != nil {
		t.Fatalf("RPC call failed: %v", err)
	}
	assert.Equal(t, "core", actual.Type)
	assert.Equal(t, "1.0.0", actual.Version)
	assert.Equal(t, 8080, len(actual.Content))
	assert.Equal(t, 16, len(actual.Hash))
}

func shouldReturnContentWhenFullRequest(t *testing.T, client *rpc.Client) {
	// given
	expectedHash := []byte{0x55, 0x76, 0x49, 0x34, 0x30, 0xd7, 0x96, 0xc1, 0x82, 0xac, 0x47, 0x4e, 0x68, 0x63, 0x12, 0xde}
	req := api.Request{Type: "core", Version: "1.0.0", Hash: expectedHash}
	var actual api.Response
	// when
	err := client.Call(cmd.LoadComponentMethod, req, &actual)
	// then
	if err != nil {
		t.Fatalf("RPC call failed: %v", err)
	}

	fmt.Println("Byte slice as Go array: ", fmt.Sprintf("%#v", actual.Hash))

	assert.Equal(t, "core", actual.Type)
	assert.Equal(t, "1.0.0", actual.Version)
	assert.Equal(t, 8080, len(actual.Content))
	assert.Equal(t, 16, len(actual.Hash))
}

func shouldReturnEmptyContentWhenRequestedHashDifferent(t *testing.T, client *rpc.Client) {
	// given
	req := api.Request{Type: "core", Version: "1.0.0", Hash: []byte{1}}
	var actual api.Response
	// when
	err := client.Call(cmd.LoadComponentMethod, req, &actual)
	// then
	if err != nil {
		t.Fatalf("RPC call failed: %v", err)
	}
	assert.Equal(t, "core", actual.Type)
	assert.Equal(t, "1.0.0", actual.Version)
	assert.Nil(t, actual.Content)
	assert.Equal(t, 16, len(actual.Hash))
}

func shouldReturnErrorWhenRequestedContentNotFound(t *testing.T, client *rpc.Client) {
	// given
	req := api.Request{Type: "ui", Version: "1.0.1"}
	var actual api.Response
	// when
	err := client.Call(cmd.LoadComponentMethod, req, &actual)
	// then
	assert.Error(t, err)
	assert.True(t, strings.HasSuffix(err.Error(), "/resources/ui/1.0.1.json: no such file or directory"),
		"Has string instead of expected "+err.Error())
}
