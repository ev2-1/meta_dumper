/*
proxy.go ©2023 derz <riley@e926.de>
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dragonfireclient/mt"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: proxy dial:port listen:port")
		os.Exit(1)
	}

	srvaddr, err := net.ResolveUDPAddr("udp", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	lc, err := net.ListenPacket("udp", os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer lc.Close()

	l := mt.Listen(lc)
	log.Printf("Listen %s; forwarding to %s\n", os.Args[2], srvaddr.String())

	for {
		clt, err := l.Accept()
		if err != nil {
			log.Printf("ERROR: Accept: %s\n", err)
			continue
		}

		log.Printf("Accept!\n")

		conn, err := net.DialUDP("udp", nil, srvaddr)
		if err != nil {
			clt.Close()
			continue
		}
		srv := mt.Connect(conn)

		log.Printf("Dial!\n")

		go srv2clt(srv, clt)
		go clt2srv(clt, srv)
	}
}

func srv2clt(srv, clt mt.Peer) {
	for {
		pkt, err := srv.Recv()
		if err != nil {
			log.Printf("Warn: Revc: %s\n", err)
		}

		if srvhandler(srv, clt, pkt.Cmd) {
			continue
		}

		// send
		_, err = clt.Send(pkt)
		if err != nil {
			log.Printf("ERROR: srv.Send: %s\n", err)
			return
		}
	}
}

func clt2srv(clt, srv mt.Peer) {
	for {
		pkt, err := clt.Recv()
		if err != nil {
			log.Printf("ERROR: clt Revc: %s\n", err)
			return
		}

		if clthandler(srv, clt, pkt.Cmd) {
			log.Printf("skip\n")
			continue
		}

		_, err = srv.Send(pkt)
		if err != nil {
			log.Printf("ERROR: clt.Send: %s\n", err)
			return
		}
	}
}

func clthandler(srv, clt mt.Peer, cmd mt.Cmd) bool {
	switch rlcmd := cmd.(type) {
	case *mt.ToSrvInit:
		srv.SendCmd(&mt.ToSrvInit{
			SerializeVer: 29,
			MinProtoVer:  39,
			MaxProtoVer:  39,
			PlayerName:   rlcmd.PlayerName,
		})

		return true
	}

	return false
}

func srvhandler(srv, clt mt.Peer, cmd mt.Cmd) bool {
	switch rlcmd := cmd.(type) {
	case *mt.ToCltItemDefs:
		log.Printf("got ItemDefs with %d entries\n", len(rlcmd.Defs))
		for k, v := range rlcmd.Defs {
			log.Printf("%d -> %#v\n", k, v)
		}

		// save them defs:
		f, err := os.OpenFile("itemdefs.json", os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			log.Fatalf("Error opening itemdefs.json: %s\n", err)
		}

		defer f.Close()
		f.Seek(0, 0)

		enc := json.NewEncoder(f)
		enc.Encode(rlcmd.Defs)

	case *mt.ToCltNodeDefs:
		log.Printf("got NodeDefs with %d entries\n", len(rlcmd.Defs))
		for k, v := range rlcmd.Defs {
			log.Printf("%d -> %#v\n", k, v)
		}

		// save them defs:
		f, err := os.OpenFile("nodedefs.json", os.O_CREATE|os.O_RDWR, 0777)
		if err != nil {
			log.Fatalf("Error opening nodedefs.json: %s\n", err)
		}

		defer f.Close()
		f.Seek(0, 0)

		enc := json.NewEncoder(f)
		enc.Encode(rlcmd.Defs)

	case *mt.ToCltHello:
		log.Printf("SV: %d; PV: %d; name: %s",
			rlcmd.SerializeVer, rlcmd.ProtoVer, rlcmd.Username,
		)

		//	case *mt.ToCltHello:
		//		rlcmd.			SerializeVer= 28,
		//rlcmd.					ProtoVer= 39,
		//rlcmd.				AuthMethods=  rlcmd.AuthMethods,
		//	rlcmd.			Username=     rlcmd.Username,
	}

	return false
}
