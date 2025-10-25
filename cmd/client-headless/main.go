package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/typetrait/pingo/internal/networking"
	"github.com/typetrait/pingo/internal/packet/clientbound"
	"github.com/typetrait/pingo/internal/packet/serverbound"
)

var (
	ErrUnknownPacketType = errors.New("unknown packet type")
)

func main() {
	conn, err := net.Dial("tcp4", "localhost:7777")
	if err != nil {
		log.Fatal(err)
	}

	_ = sendPacket(conn, &serverbound.Handshake{
		ProtocolVersion: 0,
	})

	pkt, err := readPacket(conn)
	if err != nil {
		log.Fatal(err)
	}

	switch pkt.(type) {
	case *clientbound.Handshake:
		log.Println("handshake ack received")
	}

	for {
		fmt.Print("action (host/join): ")
		var opt string
		_, err = fmt.Scanf("%s", &opt)
		if err != nil {
			log.Fatal(err)
		}

		switch opt {
		case "host":
			log.Println("requesting match creation")
			createMatchPacket := &serverbound.CreateMatch{}
			err = sendPacket(conn, createMatchPacket)
			if err != nil {
				log.Fatal(err)
			}

			p, err := readPacket(conn)
			if err != nil {
				log.Fatal(err)
			}

			matchCreatedPkt, ok := p.(*clientbound.MatchCreated)
			if !ok {
				log.Fatal("unexpected packet")
			}
			fmt.Printf("match created: %v\n", matchCreatedPkt.MatchID)

			p, err = readPacket(conn)
			if err != nil {
				log.Fatal(err)
			}

			_, ok = p.(*clientbound.Play)
			if !ok {
				log.Fatal("unexpected packet")
			}

			handlePlay(conn)
			os.Exit(1)

		case "join":
			log.Println("requesting match join")

			var matchID, playerName string
			fmt.Print("match id: ")
			_, err = fmt.Scanf("%s", &matchID)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print("player name: ")
			_, err = fmt.Scanf("%s", &playerName)
			if err != nil {
				log.Fatal(err)
			}

			joinMatchPacket := &serverbound.JoinMatch{
				MatchID:    matchID,
				PlayerName: playerName,
			}
			err = sendPacket(conn, joinMatchPacket)
			if err != nil {
				log.Fatal(err)
			}

			p, err := readPacket(conn)
			if err != nil {
				log.Fatal(err)
			}

			_, ok := p.(*clientbound.Play)
			if !ok {
				log.Fatal("unexpected packet")
			}

			handlePlay(conn)
			os.Exit(1)

		default:
			log.Println("unknown action")
			continue
		}
	}
}

func handlePlay(conn net.Conn) {
	for {
		p, err := readPacket(conn)
		if err != nil {
			return
		}

		gameState, ok := p.(*clientbound.GameState)
		if !ok {
			log.Fatal("unexpected packet")
		}

		fmt.Printf("receiving game state\n")
		fmt.Printf("p1 ypos = %f\n", gameState.PlayerOnePos.Y)
		fmt.Printf("p2 ypos = %f\n", gameState.PlayerTwoPos.Y)
		fmt.Printf("ball xpos = %f\n", gameState.BallPos.X)
		fmt.Printf("ball ypos = %f\n", gameState.BallPos.Y)
	}
}

func sendPacket(conn net.Conn, p networking.Packet) error {
	p.Write(conn)
	return nil
}

func readPacket(conn net.Conn) (networking.Packet, error) {
	readBuf := make([]byte, 1)
	n, err := conn.Read(readBuf)
	if err != nil || n != 1 {
		return nil, fmt.Errorf("could not read from connection: %w", err)
	}

	packetID := readBuf[0]

	fmt.Printf("packet id = %d\n", packetID)

	switch packetID {
	case clientbound.S2CHandshake:
		pkt := clientbound.Handshake{}
		pkt.Read(conn)
		return &pkt, nil
	case clientbound.S2CMatchCreated:
		pkt := clientbound.MatchCreated{}
		pkt.Read(conn)
		return &pkt, nil
	case clientbound.S2CGameState:
		pkt := clientbound.GameState{}
		pkt.Read(conn)
		return &pkt, nil
	case clientbound.S2CPlay:
		pkt := clientbound.Play{}
		pkt.Read(conn)
		return &pkt, nil
	default:
		return nil, ErrUnknownPacketType
	}
}
