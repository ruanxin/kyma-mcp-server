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
}
