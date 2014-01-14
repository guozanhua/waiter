package enet

/*
#cgo LDFLAGS: -lenet
#include <stdio.h>
#include <stdlib.h>
#include <enet/enet.h>

ENetHost * server;

int startServer(int port) {
	if (enet_initialize() != 0) {
		fprintf (stderr, "An error occurred while initializing ENet.\n");
		return 1;
	}
	atexit(enet_deinitialize);

	ENetAddress address;

	// Bind the server to the default localhost
	address.host = ENET_HOST_ANY;

	// Bind the server to port
	address.port = port;

	server = enet_host_create(&address, 2, 2, 0, 0);
	if (server == NULL) {
		fprintf(stderr, "An error occurred while trying to create an ENet server host.\n");
		exit(EXIT_FAILURE);
	}

	printf("server listening on 0.0.0.0:1234\n");
	return 0;
}

ENetEvent service(int timeout) {
	ENetEvent event;

	// Wait for an event (up to timeout milliseconds)
	int e = 0;

	do {
		e = enet_host_check_events(server, &event);
		if (e <= 0) {
			e = enet_host_service(server, &event, timeout);
		}
	} while (e < 0);

	return event;
}

void flush() {
	enet_host_flush(server);
}
*/
import "C"

import (
	"errors"
)

func StartServer(lport int) error {
	errCode := C.startServer(C.int(lport))
	if errCode != 0 {
		return errors.New("an error occured running the C code")
	}

	return nil
}

func Service(timeout int) Event {
	for {
		var cEvent C.ENetEvent = C.service(C.int(timeout))
		return eventFromCEvent(interface{}(&cEvent))
	}
}

func Flush() {
	C.flush()
}
