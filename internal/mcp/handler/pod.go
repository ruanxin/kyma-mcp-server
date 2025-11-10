package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ruanxin/kyma-mcp-server/internal/k8s/services"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/types"
)

type PodHandler struct {
	podService *services.PodService
}

func NewPodHandler(podService *services.PodService) *PodHandler {
	return &PodHandler{podService: podService}
}

func (h *PodHandler) CreatePod(ctx context.Context, req *mcp.CallToolRequest, input types.CreatePodInput) (*mcp.CallToolResult, *types.CreatePodOutput, error) {
	log.Printf("Handler: Received createPod request: %+v", input)

	pod, err := h.podService.CreatePod(ctx, input.Namespace, input.Name, nil, map[string]interface{}{
		"containers": []map[string]interface{}{
			{
				"name":  input.Name,
				"image": input.Image,
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return nil, &types.CreatePodOutput{
		Status:  fmt.Sprintf("pod %s created", pod.Name),
		PodName: pod.Name,
		UID:     string(pod.UID),
	}, nil
}

func (h *PodHandler) DeletePod(ctx context.Context, req *mcp.CallToolRequest, input types.DeletePodInput) (*mcp.CallToolResult, *types.DeletePodOutput, error) {
	err := h.podService.DeletePod(ctx, input.Namespace, input.Name)
	if err != nil {
		return nil, nil, err
	}

	return nil, &types.DeletePodOutput{
		Status:  fmt.Sprintf("pod %s deleted", input.Name),
		PodName: input.Name,
	}, nil
}
