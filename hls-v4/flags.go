package main

import (
	"flag"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"strings"

	maddr "github.com/multiformats/go-multiaddr"
)

// A new type we need for writing a custom flag parser
type addrList []maddr.Multiaddr

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := maddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

func StringsToAddrs(addrStrings []string) (maddrs []maddr.Multiaddr, err error) {
	for _, addrString := range addrStrings {
		addr, err := maddr.NewMultiaddr(addrString)
		if err != nil {
			return maddrs, err
		}
		maddrs = append(maddrs, addr)
	}
	return
}

type Config struct {
	RendezvousString string
	ProtocolID       string
	BootstrapPeers   addrList
}

func ParseFlags() (Config, error) {
	config := Config{}
	flag.StringVar(&config.RendezvousString, "rendezvous", "exchange data",
		"分组约定码，可以和你的朋友约定特殊字符串进行群聊。Unique string to identify group of nodes. Share this with your friends to let them connect with you.")
	flag.StringVar(&config.ProtocolID, "pid", "/sharefile/1.0.0", "Sets a protocol id for stream headers")
	flag.Var(&config.BootstrapPeers, "peer", "Adds a peer multiaddress to the bootstrap list")
	flag.Parse()
	if len(config.BootstrapPeers) == 0 {
		config.BootstrapPeers = dht.DefaultBootstrapPeers
	}

	return config, nil
}
