package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/handler"
)

func RegisterUnstructedTools(server *mcp.Server, handlers *handler.UnstructuredHandler) {

	mcp.AddTool(server, &mcp.Tool{
		Name:        "applyKubernetesResource",
		Description: "Applies any Kubernetes resource manifest (YAML or JSON) to the cluster.",
	}, handlers.HandleApplyKubernetesResource)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "listKubernetesResources",
		Description: "Lists Kubernetes resources of a given type (pods, deployments, services, etc.) with optional namespace and label filtering.",
	}, handlers.ListResources)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "getKubernetesResource",
		Description: "Gets a specific Kubernetes resource by name, apiVersion, and kind.",
	}, handlers.GetResource)
}
