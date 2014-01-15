package enet

/*
#cgo LDFLAGS: -lenet
#include <enet/enet.h>
*/
import "C"

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"unsafe"
)

type PeerState uint

const (
	PEER_STATE_DISCONNECTED             PeerState = C.ENET_PEER_STATE_DISCONNECTED
	PEER_STATE_CONNECTING               PeerState = C.ENET_PEER_STATE_CONNECTING
	PEER_STATE_ACKNOWLEDGING_CONNECT    PeerState = C.ENET_PEER_STATE_ACKNOWLEDGING_CONNECT
	PEER_STATE_CONNECTION_PENDING       PeerState = C.ENET_PEER_STATE_CONNECTION_PENDING
	PEER_STATE_CONNECTION_SUCCEEDED     PeerState = C.ENET_PEER_STATE_CONNECTION_SUCCEEDED
	PEER_STATE_CONNECTED                PeerState = C.ENET_PEER_STATE_CONNECTED
	PEER_STATE_DISCONNECT_LATER         PeerState = C.ENET_PEER_STATE_DISCONNECT_LATER
	PEER_STATE_DISCONNECTING            PeerState = C.ENET_PEER_STATE_DISCONNECTING
	PEER_STATE_ACKNOWLEDGING_DISCONNECT PeerState = C.ENET_PEER_STATE_ACKNOWLEDGING_DISCONNECT
	PEER_STATE_ZOMBIE                   PeerState = C.ENET_PEER_STATE_ZOMBIE
)

type Peer struct {
	sync.Mutex
	Address net.UDPAddr
	State   PeerState
	Data    unsafe.Pointer
	cPeer   *C.ENetPeer
}

func peerFromCPeer(cPeer *C.ENetPeer) Peer {
	if cPeer == nil {
		return Peer{}
	}

	ip := uint32(cPeer.address.host)

	p := Peer{
		Address: net.UDPAddr{
			IP:   net.IPv4(byte((ip<<24)>>24), byte((ip<<16)>>8), byte((ip<<8)>>16), byte(ip>>24)),
			Port: int(cPeer.address.port),
		},
		State: PeerState(cPeer.state),
		Data:  cPeer.data,
		cPeer: cPeer,
	}

	return p
}

func (p *Peer) SetData(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("error setting peer data: invalid type" + reflect.TypeOf(v).String())
	}

	p.cPeer.data = unsafe.Pointer(rv.Pointer())
	p.Data = p.cPeer.data

	return nil
}

func (p *Peer) Send(payload []byte, flags PacketFlag, channel uint8) {
	if len(payload) <= 0 {
		return
	}

	packet := C.enet_packet_create(unsafe.Pointer(&payload[0]), C.size_t(len(payload)), C.enet_uint32(flags))

	p.Lock()
	C.enet_peer_send(p.cPeer, C.enet_uint8(channel), packet)
	p.Unlock()
}

func (p *Peer) Disconnect(reason uint32) {
	C.enet_peer_disconnect(p.cPeer, C.enet_uint32(reason))
}
