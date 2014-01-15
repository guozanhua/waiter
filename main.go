package main

import (
	"./enet"
	"github.com/sauerbraten/jsonconf"
	"log"
)

type ServerState struct {
	MasterMode  MasterMode
	GameMode    GameMode
	Map         string
	TimeLeft    int32 // in milliseconds
	NotGotItems bool
}

var (
	// global enet host var (to call Flush() on)
	host enet.Host

	// global server state
	state ServerState

	// global collection of clients
	clients = map[ClientNumber]*Client{}

	// server configuration
	config Config
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
		TimeLeft:    600000,
		NotGotItems: true,
	}
}

func main() {
	var err error
	host, err = enet.StartServer(config.ListenPort)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("server running on port", config.ListenPort)

	for {
		event := host.Service(5)

		switch event.Type {
		case enet.EVENT_TYPE_CONNECT:
			client := addClient(event.Peer)
			err := event.Peer.SetData(&client.CN)
			if err != nil {
				log.Println("enet:", err)
			}
			client.sendServerInfo()

		case enet.EVENT_TYPE_DISCONNECT:
			client := clients[*(*ClientNumber)(event.Peer.Data)]
			client.leave()
			log.Println("ENet: disconnected:", event.Peer.Address.String())

		case enet.EVENT_TYPE_RECEIVE:
			// TODO: fix this maybe?
			if len(event.Packet.Data) == 0 {
				continue
			}
			//log.Println("got:", event.Packet.Data)
			parsePacket(*(*ClientNumber)(event.Peer.Data), event.ChannelId, Packet{event.Packet.Data, 0})
		}

		log.Println("flushing")
		host.Flush()
	}
}
