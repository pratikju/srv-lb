package clb

import (
	"github.com/benschw/consul-clb-go/dns"
	"github.com/benschw/consul-clb-go/randomclb"
	"github.com/benschw/consul-clb-go/roundrobinclb"
)

type LoadBalancerType int

const (
	Random     LoadBalancerType = iota
	RoundRobin LoadBalancerType = iota
)

type LoadBalancer interface {
	GetAddress(name string) (dns.Address, error)
}

func NewClb(address string, port string, lbType LoadBalancerType) LoadBalancer {
	switch lbType {
	case RoundRobin:
		return NewRoundRobinClb(address, port)
	case Random:
		return NewRandomClb(address, port)
	}
	return nil
}

func NewRoundRobinClb(address string, port string) *roundrobinclb.RoundRobinClb {
	return roundrobinclb.NewRoundRobinClb(address, port)
}

func NewRandomClb(address string, port string) *randomclb.RandomClb {
	return randomclb.NewRandomClb(address, port)
}