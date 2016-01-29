package main

import (
	"log"
	"net/rpc/jsonrpc"

	"github.com/bryanl/doit-provider-tag"
	"github.com/natefinch/pie"
)

const (
	pluginName = "tag"
)

func main() {
	log.SetPrefix("[doit-provider-tag] ")

	p := pie.NewProvider()
	if err := p.RegisterName(pluginName, &doittag.PluginAPI{}); err != nil {
		log.Fatalf("failed to register plugin: %s", err)
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
}
