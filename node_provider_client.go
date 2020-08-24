package main

import (
	"fmt"
	"net"
	"os"

	v1Core "k8s.io/api/core/v1"
	v1Meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type providerClientIf interface {
	listNodes() []node
}

type kubeProviderClient struct {
	c *kubernetes.Clientset
}

func newProviderClient() (p providerClientIf) {
	switch os.Getenv("OVERLAY_DATA_CLIENT") {
	case "kube-api", "":
		return newKubeProviderClient()
	default:
		pauseOnError(errNotImplemented, "newProviderClient")
	}
	return
}

func newKubeProviderClient() providerClientIf {
	client := newKubeClient()
	return &kubeProviderClient{client}
}

func (kC *kubeProviderClient) listNodes() []node {
	data, err := kC.c.CoreV1().Nodes().List(v1Meta.ListOptions{})
	pauseOnError(err, "Listing k8s nodes")
	nodes := []node{}
	for _, n := range data.Items {
		name := n.ObjectMeta.Name
		cidr := n.Spec.PodCIDR
		internalIPStr := ""
		for _, addr := range n.Status.Addresses {
			if addr.Type == v1Core.NodeInternalIP {
				internalIPStr = addr.Address
				break
			}
		}

		ip := net.ParseIP(internalIPStr)

		n := node{name: name, podCIDR: cidr, internalIP: ip}
		nodes = append(nodes, n)
	}
	return nodes
}

func newKubeClient() *kubernetes.Clientset {
	kubeAPIHost := fmt.Sprintf("https://%s:%s", os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT_HTTPS"))
	cfg := &rest.Config{
		Host:            kubeAPIHost,
		APIPath:         "/api",
		BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
	}
	client, err := kubernetes.NewForConfig(cfg)
	pauseOnError(err, "Creating K8s Client")
	return client
}
