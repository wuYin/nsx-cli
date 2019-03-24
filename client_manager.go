package client

import (
	"fmt"
	"nix/registry"
	"time"
	"tron"
)

type NixClientManager struct {
	clientManager *tron.ClientsManager
	loader        *NixServiceLoader
}

// 管理 uris 对应的 service 地址
func NewNixClientManager(centerAddr string, uris []string) *NixClientManager {
	clientManager := tron.NewClientsManager(tron.NewReconnectTaskManager(1*time.Second, 2))
	register := registry.NewDefaultRegistry(centerAddr)

	m := &NixClientManager{
		clientManager: clientManager,
	}
	loader := NewNixServiceLoader(uris, register, m.onReloadService)
	m.loader = loader

	return m
}

func (m *NixClientManager) SelectClient(uri string) (*tron.Client, error) {
	clients := m.clientManager.FindClients(uri, func(gid string, cli *tron.Client) bool {
		return false // 全选
	})
	return clients[0], nil
}

func (m *NixClientManager) onReloadService(uri string, addr string) {
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
