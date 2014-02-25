package main

import (
	"./enet"
	"log"
	"time"
)

type PacketToBroadcast struct {
	From    ClientNumber
	Payload Packet
}

type Broadcaster struct {
	PacketsToBroadcast chan PacketToBroadcast
	replacePrevious    bool
	clientPartPrefix   func(ClientNumber, *[]byte) Packet
	flags              enet.PacketFlag
	channel            uint8
	ticker             *time.Ticker
	groupedPackets     map[ClientNumber]*Packet
}

func newBroadcaster(packetsToBroadcast chan PacketToBroadcast, replacePrevious bool, interval time.Duration, clientPartPrefix func(ClientNumber, *[]byte) Packet, flags enet.PacketFlag, channel uint8) *Broadcaster {
	return &Broadcaster{
		PacketsToBroadcast: packetsToBroadcast,
		replacePrevious:    replacePrevious,
		clientPartPrefix:   clientPartPrefix,
		flags:              flags,
		channel:            channel,
		ticker:             time.NewTicker(interval),
		groupedPackets:     map[ClientNumber]*Packet{},
	}
}

func (b *Broadcaster) run() {
	for {
		select {
		case <-b.ticker.C:
			b.flush()

		case ptb := <-b.PacketsToBroadcast:
			p, ok := b.groupedPackets[ptb.From]
			if ok && !b.replacePrevious {
				p.putBytes(ptb.Payload.buf)
			} else {
				b.groupedPackets[ptb.From] = &ptb.Payload
			}
		}
	}
}

func (b *Broadcaster) flush() {
	masterPacket := Packet{}
	packetLengths := map[ClientNumber]int{}

	for _, client := range clients {
		if p, ok := b.groupedPackets[client.CN]; ok {
			if p.len() == 0 {
				continue
			}
			prefix := b.clientPartPrefix(client.CN, &(p.buf))
			masterPacket.putBytes(prefix.buf)
			masterPacket.putBytes(p.buf)
			packetLengths[client.CN] = prefix.len() + p.len()
			b.groupedPackets[client.CN] = &Packet{}
		}
	}

	if masterPacket.len() == 0 {
		return
	}

	masterPacket.putBytes(masterPacket.buf)
	log.Println(masterPacket)

	pos := 0
	for _, client := range clients {
		if !client.Joined {
			continue
		}

		length, ok := packetLengths[client.CN]
		if ok {
			pos += length
		}

		if masterPacket.len() == length*2 {
			// only the client's own packages are in the master packet
			continue
		}

		client.send(b.flags, b.channel, masterPacket.buf[pos:pos+(masterPacket.len()/2)-length])
	}
}
