package main

import (
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	// Import our internal packages
	"github.com/ruanxin/kyma-mcp-server/internal/k8s"
	"github.com/ruanxin/kyma-mcp-server/internal/k8s/services"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/handler"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/tools"
)

func main() {
	// 1. Initialize the Kubernetes client
	clientset, dynamicClient, mapper, err := k8s.NewInClusterClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}
	log.Println("Kubernetes client initialized.")

	// 2. Initialize the service layer
	podService := services.NewPodService(clientset)
	unstructuredService := services.NewUnstructedService(clientset, dynamicClient, mapper)

	// 3. Initialize the handler layer
	podHandlers := handler.NewPodHandler(podService)
	unstructuredHandlers := handler.NewUnstructedHandler(unstructuredService)

	// 4. Create the MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "kyma-mcp-server",
		Version: "v1.0.0",
	}, nil)

	// 5. Register all tools
	tools.RegisterPodTools(server, podHandlers)
	tools.RegisterUnstructedTools(server, unstructuredHandlers)
	log.Println("MCP tools registered.")

	// 6. Run the server
	handler := mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{},
	)

	// 7. Register the handler with the default http server
	http.Handle("/mcp", handler) // Your server will be available at the /mcp endpoint

	// 8. Run the standard Go HTTP server on port 8080
	log.Println("Starting Kubernetes MCP server on :8080/mcp ...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
