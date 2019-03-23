package main

import (
	"fmt"
	"nix-cli"
	"nix/service"
)

type FakeAddServiceProxy struct {
	Add func(base, diff float64) float64
}

func main() {
	c := client.NewConsumer([]service.Service{
		{
			Uri:      "add-service",
			Instance: &FakeAddServiceProxy{},
		},
	})

	proxy, ok := c.GetService("add-service").(*FakeAddServiceProxy)
	if !ok {
		return
	}
	sum := proxy.Add(1, 1)
	fmt.Printf("1+1 = %f\n", sum)
}
