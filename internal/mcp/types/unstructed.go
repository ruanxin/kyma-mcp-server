package types

// --- ApplyKubernetesResource Types ---
type ApplyKubernetesResourceInput struct {
	Manifest string `json:"manifest" jsonschema:"The complete Kubernetes resource manifest as a YAML or JSON string"`
}

type ApplyKubernetesResourceOutput struct {
	Status    string `json:"status" jsonschema:"The result of the apply operation"`
	Kind      string `json:"kind" jsonschema:"The kind of the applied resource (e.g., 'Deployment')"`
	Name      string `json:"name" jsonschema:"The name of the applied resource"`
	Namespace string `json:"namespace" jsonschema:"The namespace of the applied resource (if applicable)"`
}
