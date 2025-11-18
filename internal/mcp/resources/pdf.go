package resources

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ruanxin/kyma-mcp-server/internal/mcp/handler"
)

func RegisterPDFResources(server *mcp.Server, handler *handler.PDFHandler) {
	server.AddResource(
		&mcp.Resource{
			Name:     "API Gateway Kyma Module Documentation",
			MIMEType: "application/pdf",
			URI:      "embedded:api-gateway.pdf",
		},
		handler.HandlePDFResource,
	)
	server.AddResource(
		&mcp.Resource{
			Name:     "Telemetry Kyma Module Documentation",
			MIMEType: "application/pdf",
			URI:      "embedded:telemetry-modules.pdf",
		},
		handler.HandlePDFResource,
	)
}
