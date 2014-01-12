package main

const (
	GLOBAL_AUTH_DOMAIN = ""
)

type Config struct {
	ServerDescription string   `json:"server_description"`
	MaxClients        int      `json:"max_clients"`
	MasterServer      string   `json:"master_server"`
	MasterServerPort  int      `json:"master_server_port"`
	ServerAuthDomains []string `json:"server_auth_domains"`
}
