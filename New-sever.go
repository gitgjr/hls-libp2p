package main

import (
	"bytes"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peerstore"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"net/http"
)

func newHost(listen multiaddr.Multiaddr) host.Host {
	h, err := libp2p.New(
		libp2p.ListenAddrs(listen),
	)
	if err != nil {
		panic(err)
	}
	return h
}
func main() {
	m1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10000")
	m2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10001")
	srvHost := newHost(m1)
	clientHost := newHost(m2)
	defer srvHost.Close()
	defer clientHost.Close()

	srvHost.Peerstore().AddAddrs(clientHost.ID(), clientHost.Addrs(), peerstore.PermanentAddrTTL)
	clientHost.Peerstore().AddAddrs(srvHost.ID(), srvHost.Addrs(), peerstore.PermanentAddrTTL)

	listener, err := gostream.Listen(srvHost, "/testiti-test")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	go func() {
		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			resp := fmt.Sprintf("Hi %s!", body)
			w.Write([]byte(resp))
		})
		server := &http.Server{}
		server.Serve(listener)
	}()

	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost, p2phttp.ProtocolOption("/testiti-test")))
	client := &http.Client{Transport: tr}

	buf := bytes.NewBufferString("Hector")
	res, err := client.Post(fmt.Sprintf("libp2p://%s/hello", srvHost.ID().Pretty()), "text/plain", buf)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	text, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	if string(text) != "Hi Hector!" {
		fmt.Println("expected Hi Hector! but got %s", text)
	}

	fmt.Println(string(text))
}
