package main

import (
	"fmt"
	"nsx-cli"
	"nsx/registry"
	"nsx/service"
)

type FakeAddServiceProxy struct {
	Add func(base int, diff []int) int
}

func main() {
	services := []service.Service{{Uri: "add-service", Instance: &FakeAddServiceProxy{}}}
	c := client.NewCaller(registry.REG_ZK, []string{"127.0.0.1:2181"}, services)

	proxy, ok := c.GetService("add-service").(*FakeAddServiceProxy)
	if !ok {
		return
	}
	sum := proxy.Add(1, []int{1})
	fmt.Printf("1+1 = %d\n", sum)
}
