package client

import (
	"fmt"
	"nsx/registry"
	"time"
)

type NsxServiceLoader struct {
	uris     []string
	uri2Addr map[string]string
	registry registry.Registry
	onReload OnReload
}

type OnReload func(uri string, addr string)

func NewNsxServiceLoader(uris []string, registry registry.Registry, onReload OnReload) *NsxServiceLoader {
	l := &NsxServiceLoader{
		uris:     uris,
		uri2Addr: make(map[string]string),
		registry: registry,
		onReload: onReload,
	}

	l.RefreshAddr()
	go func() {
		t := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-t.C:
				l.RefreshAddr()
			}
		}
	}()
	return l
}

// 刷新拉取最新的可用服务地址
func (l *NsxServiceLoader) RefreshAddr() {
	for _, uri := range l.uris {
		addr := l.registry.GetService(uri)
		if addr == "" {
			fmt.Printf("[addr is empty]: %s\n", uri)
			continue
		}

		l.uri2Addr[uri] = addr // 更新
		l.onReload(uri, addr)
	}
}
