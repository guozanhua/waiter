package main

const (
	GLOBAL_AUTH_DOMAIN = ""
)

type Config struct {
	ListenAddress       string `json:"listen_address"`
	ListenPort          int    `json:"listen_port"`
	MasterServerAddress string `json:"master_server"`
	MasterServerPort    int    `json:"master_server_port"`

	ServerDescription       string   `json:"server_description"`
	ServerPassword          string   `json:"server_password"`
	ServerAuthDomains       []string `json:"server_auth_domains"`
	MaxClients              int      `json:"max_clients"`
	SendClientIPsViaExtinfo bool     `json:"send_client_ips_via_extinfo"`

	CPUCores int `json:"cpu_cores"`
}
