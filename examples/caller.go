package main

import (
	"fmt"
	"nix-cli"
	"nix/service"
)

type FakeAddServiceProxy struct {
	Add func(base int, diff []int) int
}

func main() {
	c := client.NewCaller([]service.Service{
		{
			Uri:      "add-service",
			Instance: &FakeAddServiceProxy{},
		},
	})

	proxy, ok := c.GetService("add-service").(*FakeAddServiceProxy)
	if !ok {
		return
	}
	sum := proxy.Add(1, []int{1})
	fmt.Printf("1+1 = %d\n", sum)
}
