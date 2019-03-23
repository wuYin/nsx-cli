package client

import (
	"fmt"
	"net"
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
	conf := tron.NewConfig(16*1024, 16*1024, 100, 100, 1000, 5*time.Second)
	client := tron.NewClient(conn, conf, codec, func(cli *tron.Client, p *tron.Packet) {
		cli.NotifyReceived(p.Header.Seq, p.Data)
	})
	client.Run()

	g := tron.NewClientsGroup(uri, uri)
	m.clientManager.Add(g, client)
}

// 连接
func dial(addr string) (*net.TCPConn, error) {
	rAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp4", nil, rAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
