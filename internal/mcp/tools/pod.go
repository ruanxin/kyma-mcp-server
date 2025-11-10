package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/handler"
)

// RegisterTools adds all available tools to the MCP server.
func RegisterPodTools(server *mcp.Server, handlers *handler.PodHandler) {

	// --- Tool 1: Create Pod ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "createKubernetesPod",
		Description: "Creates a new pod in a Kubernetes cluster",
	}, handlers.CreatePod)

	// --- Tool 3: Delete Pod ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "deleteKubernetesPod",
		Description: "Deletes a specific pod from a Kubernetes namespace",
	}, handlers.DeletePod)
}
