# Waiter

A game server for Cube 2: Sauerbraten, written in Go.

## To Do

- implement most of the network events
- moar async
- Fork an observation server (connects to a game as spectator, broadcasts the game data to spectators with optional delay)
- Get global auth to work
- Implement local auth, too
- [Fix garbage collector (I think?) lags]

## Project Structure

| file                     | used for                                                                                                                                |
| ------------------------ | --------------------------------------------------------------------------------------------------------------------------------------- |
| auth.go                  | implementation of functions related to Sauer's authentication mechanism                                                                 |
| map_rotation.go          | contains sets of maps used for map rotations in the current mode                                                                        |
| server_state.go          | defines the `ServerState` type and methods on it which change the server state, e.g. `changeMap()`                                      |
| client.go                | everything related to a client (send methods for example)                                                                               |
| game_constants.go        | constants which are used by server and client (e.g. initial ammo depending on mode, game modes, privileges, etc.)                       |
| network_message_codes.go | lists all the network message codes used for communication                                                                              |
| timing.go                | stuff needed for timing of things; intermission handling                                                                                |
| config.go                | the configuration stuff which the config.json file is read into                                                                         |
| game_state.go            | a client's state in a game                                                                                                              |
| packet.go                | a Sauerbraten network event packet, plus methods to read and manipulate it                                                              |
| world_state.go           | functions to broadcast the world state (player event messages & player positions)                                                       |
| config.json              | the configuration file for the user to set up the server                                                                                |
| protocol.go              | processing of incoming packets; building of outgoing packets                                                                            |
| main.go                  | initialization of server (reading of config files, starting of needed goroutines) and server's main loop listening for incoming packets |

## License

This code is licensed under a BSD License:

Copyright (c) 2014 Alexander Willing. All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
