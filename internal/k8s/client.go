package k8s

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func NewInClusterClient() (kubernetes.Interface, dynamic.Interface, *restmapper.DeferredDiscoveryRESTMapper, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	// 1. Standard typed clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	// 2. Dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	// 3b. Wrap the basic client in a memory cache
	cachedDiscoveryClient := memory.NewMemCacheClient(discoveryClient)

	// 3c. Pass the caching client (not the basic one) to the mapper
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)
	return clientset, dynamicClient, mapper, nil
}
