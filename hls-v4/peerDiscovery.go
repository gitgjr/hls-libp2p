package main

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"sync"
	"time"
)

func connectDHT(node host.Host, ctx context.Context, config Config) *discovery.RoutingDiscovery {

	kademliaDHT, err := dht.New(ctx, node)
	if err != nil {
		panic(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		panic(err)
	}

	// Let's connect to the bootstrap nodes first. They will tell us about the
	// other nodes in the network.

	var wg sync.WaitGroup
	for _, peerAddr := range config.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := node.Connect(ctx, *peerinfo); err != nil {
				logger.Warn(err)
			} else {
				logger.Info("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	// We use a rendezvous point "exchange data" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	logger.Info("Announcing ourselves...")
	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)
	go func() { //persistently advertises the service through an Advertiser.
		for {
			d, err := routingDiscovery.Advertise(ctx, config.RendezvousString)
			if err != nil {
				panic(err.Error())
			} else {
				time.Sleep(d)
			}

		}
	}()
	logger.Debug("Successfully announced!")
	return routingDiscovery
}

//func searchPeer(routingDiscovery *discovery.RoutingDiscovery, ctx context.Context, config Config) {
//	logger.Debug("Searching for other peers...")
//	peerChan, err := routingDiscovery.FindPeers(ctx, config.RendezvousString)
//	if err != nil {
//		panic(err)
//	}
//	for peer := range peerChan {
//		if peer.ID == host.ID() {
//			continue
//		}
//		logger.Debug("Found peer:", peer)
//
//		logger.Debug("Connecting to:", peer)
//		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))
//
//		if err != nil {
//			logger.Warning("Connection failed:", err)
//			continue
//		} else {
//
//		}
//
//		logger.Info("Connected to:", peer)
//	}
//
//	select {}
//}
