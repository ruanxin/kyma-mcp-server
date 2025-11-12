package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
)

// UnstructedService holds the logic for interacting with Kubernetes.
type UnstructedService struct {
	clnt          kubernetes.Interface
	dynamicClient dynamic.Interface
	mapper        *restmapper.DeferredDiscoveryRESTMapper
}

// NewUnstructedService creates a new service.
func NewUnstructedService(clnt kubernetes.Interface, dynamicClient dynamic.Interface, mapper *restmapper.DeferredDiscoveryRESTMapper) *UnstructedService {
	return &UnstructedService{
		clnt:          clnt,
		dynamicClient: dynamicClient,
		mapper:        mapper,
	}
}

func convertToUnstructured(manifest string) (*unstructured.Unstructured, *schema.GroupVersionKind, error) {
	// 1. Decode the manifest string into an unstructured object
	var obj unstructured.Unstructured
	// Use NewDecodingSerializer to handle both JSON and YAML
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode([]byte(manifest), nil, &obj)
	if err != nil {
		log.Printf("Service: Error decoding manifest: %v", err)
		return nil, nil, err
	}
	return &obj, gvk, nil
}

// ApplyManifest is a generic function that applies any Kubernetes YAML/JSON.
func (s *UnstructedService) ApplyManifest(ctx context.Context, manifest string) (*unstructured.Unstructured, error) {
	log.Println("Service: Attempting to apply manifest...")

	obj, gvk, err := convertToUnstructured(manifest)
	if err != nil {
		log.Printf("Service: Error converting to unstructured: %v", err)
		return nil, err
	}

	// 2. Find the GVR (GroupVersionResource) for this GVK (GroupVersionKind)
	mapping, err := s.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		log.Printf("Service: Error mapping GVK to GVR: %v", err)
		return nil, err
	}

	// 3. Get the correct dynamic resource client
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// This is a namespaced resource
		if obj.GetNamespace() == "" {
			// Default to 'default' namespace if not specified
			obj.SetNamespace("default")
		}
		dr = s.dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// This is a cluster-scoped resource
		dr = s.dynamicClient.Resource(mapping.Resource)
	}

	// 4. Apply the resource (Server-Side Apply)
	// We use "mcp-server-apply" as the FieldManager
	applyOptions := metav1.ApplyOptions{
		FieldManager: "mcp-server-apply",
		Force:        true, // Force ownership
	}

	appliedObj, err := dr.Apply(ctx, obj.GetName(), obj, applyOptions)
	if err != nil {
		log.Printf("Service: Error applying resource: %v", err)
		return nil, err
	}

	log.Printf("Service: Successfully applied %s/%s", appliedObj.GetKind(), appliedObj.GetName())
	return appliedObj, nil
}

// ListResources lists Kubernetes resources of a given type
func (s *UnstructedService) ListResources(ctx context.Context, apiVersion, kind, namespace, labelSelector string) (*unstructured.UnstructuredList, *schema.GroupVersionKind, error) {
	log.Printf("Service: Listing %s/%s resources in namespace %s with selector %s", apiVersion, kind, namespace, labelSelector)

	// Parse the apiVersion to get group and version
	var group, version string
	if apiVersion == "v1" {
		group = ""
		version = "v1"
	} else {
		// Split by "/" to get group and version (e.g., "apps/v1" -> group="apps", version="v1")
		parts := strings.Split(apiVersion, "/")
		if len(parts) == 2 {
			group = parts[0]
			version = parts[1]
		} else {
			return nil, nil, fmt.Errorf("invalid apiVersion format: %s", apiVersion)
		}
	}

	gvk := schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}

	// Get the GVR mapping
	mapping, err := s.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		log.Printf("Service: Error mapping GVK to GVR: %v", err)
		return nil, nil, err
	}

	// Get the correct dynamic resource client
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if namespace == "" {
			// List across all namespaces
			dr = s.dynamicClient.Resource(mapping.Resource)
		} else {
			// List in specific namespace
			dr = s.dynamicClient.Resource(mapping.Resource).Namespace(namespace)
		}
	} else {
		// This is a cluster-scoped resource
		dr = s.dynamicClient.Resource(mapping.Resource)
	}

	// Set up list options
	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	// List the resources
	result, err := dr.List(ctx, listOptions)
	if err != nil {
		log.Printf("Service: Error listing resources: %v", err)
		return nil, nil, err
	}

	log.Printf("Service: Successfully listed %d %s resources", len(result.Items), kind)
	return result, &gvk, nil
}

// GetResource gets a specific Kubernetes resource by name
func (s *UnstructedService) GetResource(ctx context.Context, apiVersion, kind, name, namespace string) (*unstructured.Unstructured, *schema.GroupVersionKind, error) {
	log.Printf("Service: Getting %s/%s resource %s in namespace %s", apiVersion, kind, name, namespace)

	// Parse the apiVersion to get group and version
	var group, version string
	if apiVersion == "v1" {
		group = ""
		version = "v1"
	} else {
		// Split by "/" to get group and version (e.g., "apps/v1" -> group="apps", version="v1")
		parts := strings.Split(apiVersion, "/")
		if len(parts) == 2 {
			group = parts[0]
			version = parts[1]
		} else {
			return nil, nil, fmt.Errorf("invalid apiVersion format: %s", apiVersion)
		}
	}

	gvk := schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}

	// Get the GVR mapping
	mapping, err := s.mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		log.Printf("Service: Error mapping GVK to GVR: %v", err)
		return nil, nil, err
	}

	// Get the correct dynamic resource client
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if namespace == "" {
			// Default to 'default' namespace if not specified for namespaced resources
			namespace = "default"
		}
		dr = s.dynamicClient.Resource(mapping.Resource).Namespace(namespace)
	} else {
		// This is a cluster-scoped resource
		dr = s.dynamicClient.Resource(mapping.Resource)
	}

	// Get the resource
	result, err := dr.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Service: Error getting resource: %v", err)
		return nil, nil, err
	}

	log.Printf("Service: Successfully retrieved %s/%s resource %s", apiVersion, kind, name)
	return result, &gvk, nil
}
