package main

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/multiformats/go-multiaddr"
)

// NewNode initialize a node
func NewNode(listen multiaddr.Multiaddr) host.Host {
	h, err := libp2p.New(
		libp2p.ListenAddrs(listen),
	)
	if err != nil {
		panic(err)
	}
	return h
}

func setPermission(permission Permission) {
	logger.Info("Input the operation of peer, upload and download")
	fmt.Scanln(&permission.Upload, &permission.Download)
}
