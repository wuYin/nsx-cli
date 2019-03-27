package main

import (
	"errors"
	"nsx/registry"
	"nsx/server"
	"nsx/service"
	"reflect"
)

// 服务注册中心
type AdminRegistry struct {
	service2Addr map[string]string
}

func (r AdminRegistry) Register(uri, addr string) error {
	r.service2Addr[uri] = addr
	return nil
}

func (r AdminRegistry) UnRegister(uri, addr string) error {
	delete(r.service2Addr, uri)
	return nil
}

func (r AdminRegistry) GetService(uri string) ([]string, error) {
	addr, ok := r.service2Addr[uri]
	if !ok {
		return nil, errors.New(uri + " not registered yet.")
	}
	return []string{addr}, nil
}

// add 服务在另一程序中
const (
	ADD_SERVICE = "add-service"
)

type AddServiceInterface interface {
	Add(base int, diffs []int) int
}

type AddService struct{}

func (s AddService) Add(base int, diffs []int) int {
	for _, diff := range diffs {
		base += diff
	}
	return base
}

func main() {
	rs := AdminRegistry{
		service2Addr: make(map[string]string),
	}
	as := AddService{}
	services := []service.Service{
		{
			Uri:       service.SERVICE_ADMIN,
			Instance:  rs,
			Interface: reflect.TypeOf((*registry.Registry)(nil)).Elem(),
		},
		{
			Uri:       ADD_SERVICE,
			Instance:  as,
			Interface: reflect.TypeOf((*AddServiceInterface)(nil)).Elem(),
		},
	}

	// 使用默认注册中心
	// server.NewNsxServer("localhost:8080", []string{"localhost:8080"}, services, registry.REG_DEFAULT)
	// fmt.Println(rs.service2Addr) // map[admin-service:localhost:8080 add-service:localhost:8080]

	// 使用 zk 注册中心
	server.NewNsxServer("localhost:8080", []string{"127.0.0.1:2181"}, services, registry.REG_ZK)
	select {}
}
