package client

import (
	"fmt"
	"nsx/registry"
	"time"
)

type NsxServiceLoader struct {
	uris      []string
	uri2Addrs map[string][]string
	registry  registry.Registry
	onReload  OnReload
}

type OnReload func(uri string, mewServoceAddrs []string)

func NewNsxServiceLoader(uris []string, registry registry.Registry, onReload OnReload) *NsxServiceLoader {
	l := &NsxServiceLoader{
		uris:      uris,
		uri2Addrs: make(map[string][]string),
		registry:  registry,
		onReload:  onReload,
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
		addrs, err := l.registry.GetService(uri)
		if err != nil || len(addrs) == 0 {
			fmt.Printf("load %s failed: %v\n", uri, err)
			continue
		}

		l.uri2Addrs[uri] = addrs // 更新
		l.onReload(uri, addrs)
		time.Sleep(1 * time.Second)
	}
}
