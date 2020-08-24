package main

import (
	"net"
)

type node struct {
	internalIP net.IP
	podCIDR    string
	name       string
}

type dataProviderIf interface {
	getNodeByName(string) *node
	getNodeMatchingPodCIDR(string) *node
}

type dataProviderT struct {
	c                   providerClientIf
	nodeNameLookup      map[string]*node
	podCIDRPrefixLookup map[string]*node
}

func newNodeDataProvider() dataProviderIf {
	client := newProviderClient()
	p := dataProviderT{c: client}
	p.buildCaches()
	return &p
}

func (p *dataProviderT) getNodeByName(name string) *node {
	node, ok := p.nodeNameLookup[name]
	if ok {
		return node
	}
	p.buildCaches()
	return p.nodeNameLookup[name]
}

func (p *dataProviderT) getNodeMatchingPodCIDR(cidr string) *node {
	maskedCidr := applyCIDRMask(cidr)
	node, ok := p.podCIDRPrefixLookup[maskedCidr]
	if ok {
		return node
	}

	// if not present, list, rebuild cache
	p.buildCaches()

	// lookup
	node = p.podCIDRPrefixLookup[maskedCidr]
	// if this lookup fails, so be it
	return node
}

func (p *dataProviderT) buildCaches() {
	nodes := p.c.listNodes()
	if p.nodeNameLookup == nil {
		p.nodeNameLookup = map[string]*node{}
	}
	if p.podCIDRPrefixLookup == nil {
		p.podCIDRPrefixLookup = map[string]*node{}
	}

	for i, n := range nodes {
		p.nodeNameLookup[n.name] = &nodes[i]
		maskedCIDR := applyCIDRMask(n.podCIDR)
		p.podCIDRPrefixLookup[maskedCIDR] = &nodes[i]
	}
}
