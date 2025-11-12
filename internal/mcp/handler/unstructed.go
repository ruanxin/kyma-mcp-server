package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ruanxin/kyma-mcp-server/internal/k8s/services"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type UnstructuredHandler struct {
	service *services.UnstructedService
}

// NewUnstructedHandler creates a new handler struct.
func NewUnstructedHandler(service *services.UnstructedService) *UnstructuredHandler {
	return &UnstructuredHandler{
		service: service,
	}
}

// HandleApplyKubernetesResource handles applying a generic manifest.
func (h *UnstructuredHandler) HandleApplyKubernetesResource(ctx context.Context, req *mcp.CallToolRequest, input types.ApplyKubernetesResourceInput) (
	*mcp.CallToolResult, types.ApplyKubernetesResourceOutput, error,
) {
	log.Printf("Handler: Received applyKubernetesResource request")

	appliedObj, err := h.service.ApplyManifest(ctx, input.Manifest)
	if err != nil {
		return nil, types.ApplyKubernetesResourceOutput{Status: "error: " + err.Error()}, nil
	}

	statusMessage := fmt.Sprintf("Successfully applied %s/%s", appliedObj.GetKind(), appliedObj.GetName())
	return nil, types.ApplyKubernetesResourceOutput{
		Status:    statusMessage,
		Kind:      appliedObj.GetKind(),
		Name:      appliedObj.GetName(),
		Namespace: appliedObj.GetNamespace(),
	}, nil
}

// ListResources handles listing Kubernetes resources of a given type
func (h *UnstructuredHandler) ListResources(ctx context.Context, req *mcp.CallToolRequest, input types.ListResourcesInput) (*mcp.CallToolResult, *types.ListResourcesOutput, error) {
	log.Printf("Handler: Received listResources request: %+v", input)

	result, gvk, err := h.service.ListResources(ctx, input.ApiVersion, input.Kind, input.Namespace, input.LabelSelector)
	if err != nil {
		return nil, nil, err
	}

	// Convert unstructured list to our output format
	items := make([]types.ResourceItem, 0, len(result.Items))
	for _, item := range result.Items {
		// Get full status subresource
		var statusObj map[string]interface{}
		if status, found, _ := unstructured.NestedMap(item.Object, "status"); found {
			statusObj = status
		}

		resourceItem := types.ResourceItem{
			Name:      item.GetName(),
			Namespace: item.GetNamespace(),
			Kind:      gvk.Kind,
			Labels:    item.GetLabels(),
			CreatedAt: item.GetCreationTimestamp().Format(time.RFC3339),
			Status:    statusObj,
		}
		items = append(items, resourceItem)
	}

	output := &types.ListResourcesOutput{
		ApiVersion: input.ApiVersion,
		Kind:       input.Kind,
		Namespace:  input.Namespace,
		Items:      items,
		Count:      len(items),
	}

	log.Printf("Handler: Successfully listed %d %s resources", output.Count, input.Kind)
	return nil, output, nil
}

// GetResource handles getting a specific Kubernetes resource by name
func (h *UnstructuredHandler) GetResource(ctx context.Context, req *mcp.CallToolRequest, input types.GetResourceInput) (*mcp.CallToolResult, *types.GetResourceOutput, error) {
	log.Printf("Handler: Received getResource request: %+v", input)

	result, gvk, err := h.service.GetResource(ctx, input.ApiVersion, input.Kind, input.Name, input.Namespace)
	if err != nil {
		return nil, nil, err
	}

	// Get full status subresource
	var statusObj map[string]interface{}
	if status, found, _ := unstructured.NestedMap(result.Object, "status"); found {
		statusObj = status
	}

	output := &types.GetResourceOutput{
		ApiVersion: gvk.GroupVersion().String(),
		Kind:       gvk.Kind,
		Name:       result.GetName(),
		Namespace:  result.GetNamespace(),
		Labels:     result.GetLabels(),
		CreatedAt:  result.GetCreationTimestamp().Format(time.RFC3339),
		Status:     statusObj,
		Manifest:   result.Object,
	}

	log.Printf("Handler: Successfully retrieved %s/%s resource %s", input.ApiVersion, input.Kind, input.Name)
	return nil, output, nil
}
