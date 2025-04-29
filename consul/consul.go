package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/qq754174349/ht-frame/autoconfigure"
	baseConfig "github.com/qq754174349/ht-frame/config"
	log "github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/web"
	"github.com/spf13/viper"
	"net"
	"strconv"
	"time"
)

var config *Consul

type AutoConfig struct{}

type Consul struct {
	Consul Config
}

type Config struct {
	Addr string
}

func init() {
	err := autoconfigure.Register(AutoConfig{})
	if err != nil {
		log.Fatal("consul 自动配置注册失败")
	}
}

func (AutoConfig) Init() error {
	config = &Consul{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatal("配置文件格式错误")
	}
	return StartConsulAutoRegister(config)
}

func StartConsulAutoRegister(cfg *Consul) error {
	port, err := strconv.Atoi(web.GetConfig().Port)
	if err != nil {
		return fmt.Errorf("端口转换失败: %v", err)
	}

	localIP := GetOutboundIP()
	reg := &api.AgentServiceRegistration{
		ID:   baseConfig.GetAppCfg().AppName,
		Name: baseConfig.GetAppCfg().AppName,
		Port: port,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", localIP, port),
			Interval: "20s",
			Timeout:  "3s",
		},
	}

	consulAddr := cfg.Consul.Addr
	if consulAddr == "" {
		consulAddr = "127.0.0.1:8500"
	}

	client, err := api.NewClient(&api.Config{Address: consulAddr})
	if err != nil {
		return fmt.Errorf("创建 Consul 客户端失败: %v", err)
	}

	go monitorConsulAndRegister(client, reg)
	return nil
}

func monitorConsulAndRegister(client *api.Client, reg *api.AgentServiceRegistration) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		_, err := client.Agent().Self()
		if err != nil {
			log.Warn("[Consul] 不可用，等待恢复...")
			continue
		}

		services, err := client.Agent().Services()
		if err != nil {
			log.Warn("[Consul] 获取服务列表失败：", err)
			continue
		}

		if _, ok := services[reg.ID]; ok {
			// 已注册
			continue
		}

		// 注册服务
		if err := client.Agent().ServiceRegister(reg); err != nil {
			log.Warn("[Consul] 注册失败：", err)
		} else {
			log.Info("[Consul] 注册成功")
		}
	}
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal("获取本地 IP 失败：", err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}
