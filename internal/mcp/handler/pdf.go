package handler

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PDFHandler struct {
}

func NewPDFHandler() *PDFHandler {
	return &PDFHandler{}
}

func (p *PDFHandler) HandlePDFResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	if req == nil || req.Params.URI == "" {
		return nil, mcp.ResourceNotFoundError("<empty>")
	}

	// Expect URIs like embedded:api-gateway.pdf or embedded:telemetry-modules.pdf
	fileName := strings.TrimPrefix(req.Params.URI, "embedded:")
	if fileName == req.Params.URI { // prefix not found
		return nil, mcp.ResourceNotFoundError(req.Params.URI)
	}

	full := filepath.Join(".", "resources", fileName)
	data, err := os.ReadFile(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, mcp.ResourceNotFoundError(req.Params.URI)
		}
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/pdf",
				Blob:     data,
			},
		},
	}, nil
}
