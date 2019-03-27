package client

import (
	"fmt"
	"nsx/registry"
	"time"
	"tron"
)

type NsxClientManager struct {
	clientManager *tron.ClientsManager
	loader        *NsxServiceLoader
}

// 管理 uris 对应的 service 地址
func NewNsxClientManager(regType int, regServers []string, uris []string) *NsxClientManager {
	clientManager := tron.NewClientsManager(tron.NewReconnectTaskManager(1*time.Second, 2))
	var regCenter registry.Registry
	switch regType {
	case registry.REG_DEFAULT:
		regCenter = registry.NewDefaultRegistry(regServers[0])
	case registry.REG_ZK:
		regCenter = registry.NewZKRegistry(regServers)
	}
	m := &NsxClientManager{
		clientManager: clientManager,
	}
	loader := NewNsxServiceLoader(uris, regCenter, m.onReloadService)
	m.loader = loader

	return m
}

func (m *NsxClientManager) SelectClient(uri string) (*tron.Client, error) {
	clients := m.clientManager.FindClients(uri, func(gid string, cli *tron.Client) bool {
		return false // 全选
	})
	return clients[0], nil
}

func (m *NsxClientManager) onReloadService(uri string, newServiceAddrs []string) {
	for _, addr := range newServiceAddrs {
		if addr == "" {
			fmt.Println("empty addr")
			return
		}

		conn, err := dial(addr)
		if err != nil {
			fmt.Printf("dial %s failed: %v\n", addr, err)
			return
		}

		codec := NewClientCodec()
		conf := tron.NewDefaultConf(1 * time.Minute)
		client := tron.NewClient(conn, conf, codec, func(cli *tron.Client, p *tron.Packet) {
			cli.NotifyReceived(p.Seq(), p.Data)
		})
		client.ReadWriteAndHandle()

		g := tron.NewClientsGroup(uri, uri)
		m.clientManager.Add(g, client)
	}
}
