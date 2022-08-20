package main

import (
	"bufio"
	"fmt"
	"github.com/libp2p/go-libp2p-core/network"
	"strings"
)

func handleStream(stream network.Stream) { //accept stream
	go echo(stream)
	// 'stream' will stay open until you close it (or the other side closes it).
}

func echo(stream network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	str, err := rw.ReadString('\n')
	if err != nil {
		stream.Close()
		return
	}
	args := strings.SplitN(strings.TrimSpace(str), " ", 2)
	if args[0] == "get" {
		err = sendfile(stream, args[1])
		if err == nil {
			fmt.Println("send successful")
		} else {
			fmt.Println("send fail:", err)
		}
	}
}
