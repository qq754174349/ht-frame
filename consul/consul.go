package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/qq754174349/ht-frame/config"
	"log"
	"net"
	"strconv"
)

type AutoConfig struct{}

func (AutoConfig) Init(cfg *config.AppConfig) error {
	consul := cfg.Consul

	var consulConfig api.Config
	if consul.Addr == "" {
		consulConfig = *api.DefaultConfig()
	} else {
		consulConfig = api.Config{
			Address: consul.Addr,
		}
	}

	client, err := api.NewClient(&consulConfig)
	if err != nil {
		return err
	}

	// 2. 注册服务
	port, _ := strconv.Atoi(cfg.Web.Port)
	localIP := GetOutboundIP()
	reg := &api.AgentServiceRegistration{
		ID:   cfg.AppName,
		Name: cfg.AppName,
		Port: port,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", localIP, port),
			Interval: "10s",
		},
	}
	if err := client.Agent().ServiceRegister(reg); err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetOutboundIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
