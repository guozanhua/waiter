package main

import (
	"./enet"
	"github.com/sauerbraten/jsonconf"
	"log"
	"runtime"
	"time"
)

var (
	// global enet host var (to call Flush() on)
	host enet.Host

	// global variable to indicate to the main loop that there are packets to be sent
	mustFlush = false

	// global server state
	state ServerState

	// global collection of clients
	clients = Clients{}

	// server configuration
	config Config

	positionPacketsToBroadcast chan PacketToBroadcast = make(chan PacketToBroadcast)
	otherPacketsToBroadcast    chan PacketToBroadcast = make(chan PacketToBroadcast)
)

func init() {
	config = Config{}

	err := jsonconf.ParseFile("config.json", &config)
	if err != nil {
		log.Fatalln(err)
	}

	state = ServerState{
		MasterMode:  MM_OPEN,
		GameMode:    GM_EFFIC,
		Map:         "hashi",
		TimeLeft:    TEN_MINUTES,
		NotGotItems: true,
		HasMaster:   false,
	}

	runtime.GOMAXPROCS(config.CPUCores)
}

func main() {
	var err error
	host, err = enet.StartServer(config.ListenAddress, config.ListenPort)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("server running on port", config.ListenPort)

	go countDown()

	positionBroadcaster := newBroadcaster(positionPacketsToBroadcast, true, func(cn ClientNumber, buf *[]byte) Packet { return Packet{} }, enet.PACKET_FLAG_NO_ALLOCATE, 0, 33*time.Millisecond)
	otherPacketsBroadcaster := newBroadcaster(otherPacketsToBroadcast, false, func(cn ClientNumber, buf *[]byte) Packet { return makePacket(N_CLIENT, cn, len(*buf)) }, enet.PACKET_FLAG_NO_ALLOCATE|enet.PACKET_FLAG_RELIABLE, 1, 33*time.Millisecond)

	go positionBroadcaster.run()
	go otherPacketsBroadcaster.run()

	for {
		event := host.Service(5)

		switch event.Type {
		case enet.EVENT_TYPE_CONNECT:
			log.Println("ENet: connected:", event.Peer.Address.String())
			client := addClient(event.Peer)
			err := event.Peer.SetData(&client.CN)
			if err != nil {
				log.Println("enet:", err)
			}
			client.sendServerInfo()

		case enet.EVENT_TYPE_DISCONNECT:
			log.Println("ENet: disconnected:", event.Peer.Address.String())
			client := clients[*(*ClientNumber)(event.Peer.Data)]
			client.leave()

		case enet.EVENT_TYPE_RECEIVE:
			// TODO: fix this maybe?
			if len(event.Packet.Data) == 0 {
				continue
			}

			go parsePacket(*(*ClientNumber)(event.Peer.Data), event.ChannelId, Packet{event.Packet.Data, 0})
		}

		if mustFlush {
			//log.Println("flushing")
			host.Flush()
			mustFlush = false
		}
	}
}
