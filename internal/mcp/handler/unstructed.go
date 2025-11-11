package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ruanxin/kyma-mcp-server/internal/k8s/services"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/types"
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
