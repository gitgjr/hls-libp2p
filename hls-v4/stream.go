package main

import (
	"github.com/libp2p/go-libp2p-quic-transport/integrationtests/stream"
	"path/filepath"
)

func getFile(stream stream.Stream, url string) {
	_, err := stream.Write([]byte("get" + url + "/n"))
	if err != nil {
		logger.Warn(err.Error())
		stream.Close()
		//contine
	}
	err = recvfile(stream, filepath.Base(url))
	if err != nil {
		logger.Info("recv successful")
	} else {
		logger.Info("recv failed")
	}
}
