package main

import (
	"nsx/registry"
	"nsx/server"
	"nsx/service"
	"reflect"
)

// 服务注册中心
type AdminRegistry struct {
	service2Addr map[string]string
}

func (s AdminRegistry) RegisterService(uri string) error {
	s.service2Addr[uri] = "localhost:8080"
	return nil
}

func (s AdminRegistry) UnRegisterService(uri string) error {
	delete(s.service2Addr, uri)
	return nil
}

func (s AdminRegistry) GetService(uri string) (addr string) {
	return s.service2Addr[uri]
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

	server.NewNsxServer("localhost:8080", services)
	// fmt.Println(rs.service2Addr) // map[admin-service:localhost:8080 add-service:localhost:8080]
	select {}
}
