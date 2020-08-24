package main

import (
	"net"
	"testing"
)

var mockNodes = []node{
	{internalIP: net.ParseIP("192.168.200.1"), podCIDR: "10.244.0.0/24", name: "master"},
	{internalIP: net.ParseIP("192.168.200.2"), podCIDR: "10.244.1.0/24", name: "worker1"},
	{internalIP: net.ParseIP("192.168.200.3"), podCIDR: "10.244.3.0/24", name: "worker2"},
}

type mockClient struct {
	nodes []node
}

func (c *mockClient) listNodes() []node {
	return c.nodes
}

func makeMockClient(nodes []node) providerClientIf {
	return &mockClient{nodes}
}

func TestDataProviderBuildCaches(t *testing.T) {
	c := makeMockClient(mockNodes)
	provider := dataProviderT{c: c}
	provider.buildCaches()

	for _, n := range mockNodes {
		if provider.nodeNameLookup[n.name].internalIP.String() != n.internalIP.String() {
			t.Errorf("Expected %s -> Got %s", n.internalIP, provider.nodeNameLookup[n.name].internalIP)
		}

		masked := applyCIDRMask(n.podCIDR)
		if provider.podCIDRPrefixLookup[masked].name != n.name {
			t.Errorf("Expected %s -> Got %s", n.name, provider.podCIDRPrefixLookup[masked].name)
		}
	}
}

func TestGetNodeMatchingPodCIDR(t *testing.T) {
	c := makeMockClient(mockNodes)
	provider := dataProviderT{c: c}

	specs := []testSpec{
		{"10.244.0.32/24", "master"},
		{"10.244.1.92/24", "worker1"},
		{"10.244.3.7/24", "worker2"},
	}

	for _, spec := range specs {
		destinationNode := provider.getNodeMatchingPodCIDR(spec.input)
		if destinationNode.name != spec.expected {
			t.Errorf("Expected %s -> Got %s", spec.expected, destinationNode.name)
		}
	}
}

func TestGetNodeByName(t *testing.T) {
	c := makeMockClient(mockNodes)
	provider := dataProviderT{c: c}

	specs := []testSpec{
		{"master", "192.168.200.1"},
		{"worker1", "192.168.200.2"},
		{"worker2", "192.168.200.3"},
	}

	for _, spec := range specs {
		node := provider.getNodeByName(spec.input)
		if node.internalIP.String() != spec.expected {
			t.Errorf("Expected %s -> Got %s", spec.expected, node.internalIP)
		}
	}
}
