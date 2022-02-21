package goranker

type MemDBConfig struct {
	Type             string `json:"Type"`
	ConnectionString string `json:"ConnectionString"`
	Prefix           string `json:"Prefix"`
}
type DocumentDBConfig struct {
	Provider         string `json:"Provider"`
	ConnectionString string `json:"ConnectionString"`
	DatabaseName     string `json:"DatabaseName"`
	CollectionName   string `json:"CollectionName"`
}

type Config struct {
	TableSize  int               `json:"TableSize"`
	MemDB      *MemDBConfig      `json:"MemDB"`
	DocumentDB *DocumentDBConfig `json:"DocumentDB"`
}

var config *Config = nil

func SetConfig(conf *Config) {
	config = conf
}
