package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"math/rand"
)

var logger = log.Logger("rendezvous")

type Permission struct {
	Upload   bool
	Download bool
}

//TODO permission
func main() {
	url := "E:\\go\\libp2p-chatroom\\videoFile\\origin.m3u8"
	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel("rendezvous", "info")

	config, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	//Input the operation of peer
	permission := Permission{}
	setPermission(permission)

	//create node
	port := rand.Intn(1000)
	severAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	Node := NewNode(severAddr)
	logger.Info("Host created. We are:", Node.ID())
	Node.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)

	ctx := context.Background()
	routingDiscovery := connectDHT(Node, ctx, config) //Connect to DHT

	logger.Debug("Searching for other peers...")
	peerChan, err := routingDiscovery.FindPeers(ctx, config.RendezvousString)
	for peer := range peerChan {
		if peer.ID == Node.ID() {
			continue
		}
		logger.Debug("Found peer:", peer)

		logger.Debug("Connecting to:", peer)
		stream, err := Node.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			logger.Warn("Connection failed:", err)
			continue
		} else { //Key part,operation after the connection of peers
			//if permission.Upload==true{
			//	go sendData()
			//}
			if permission.Download == true {
				go getFile(stream, url)
			}
		}

		logger.Info("Connected to:", peer)
	}

	select {}
}
