package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	// Import our internal packages
	"github.com/ruanxin/kyma-mcp-server/internal/k8s"
	"github.com/ruanxin/kyma-mcp-server/internal/k8s/services"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/handler"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/tools"
)

func main() {
	// 1. Initialize the Kubernetes client
	clientset, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	log.Println("Kubernetes client initialized.")

	// 2. Initialize the service layer
	podService := services.NewPodService(clientset)

	// 3. Initialize the handler layer
	podHandlers := handler.NewPodHandler(podService)

	// 4. Create the MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "kyma-mcp-server",
		Version: "v1.0.0",
	}, nil)

	// 5. Register all tools
	tools.RegisterPodTools(server, podHandlers)
	log.Println("MCP tools registered.")

	// 6. Run the server
	log.Println("Starting Kubernetes MCP server over stdio...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
