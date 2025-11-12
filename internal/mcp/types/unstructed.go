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

// --- ListResources Types ---

type ListResourcesInput struct {
	ApiVersion    string `json:"apiVersion" jsonschema:"apiVersion of the resources (examples of valid apiVersion are: v1, apps/v1, networking.k8s.io/v1)"`
	Kind          string `json:"kind" jsonschema:"kind of the resources (examples of valid kind are: Pod, Service, Deployment, Ingress)"`
	Namespace     string `json:"namespace,omitempty" jsonschema:"Optional Namespace to retrieve the namespaced resources from (ignored in case of cluster scoped resources). If not provided, will list resources from all namespaces"`
	LabelSelector string `json:"labelSelector,omitempty" jsonschema:"Optional Kubernetes label selector (e.g. 'app=myapp,env=prod' or 'app in (myapp,yourapp)'), use this option when you want to filter the pods by label,pattern:([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]"`
}

type ResourceItem struct {
	Name      string                 `json:"name" jsonschema:"the name of the resource"`
	Namespace string                 `json:"namespace,omitempty" jsonschema:"the namespace of the resource"`
	Kind      string                 `json:"kind" jsonschema:"the kind of resource"`
	Labels    map[string]string      `json:"labels,omitempty" jsonschema:"labels on the resource"`
	CreatedAt string                 `json:"createdAt" jsonschema:"creation timestamp"`
	Status    map[string]interface{} `json:"status,omitempty" jsonschema:"complete status subresource as JSON"`
}

type ListResourcesOutput struct {
	ApiVersion string         `json:"apiVersion" jsonschema:"the apiVersion of resources that were listed"`
	Kind       string         `json:"kind" jsonschema:"the kind of resources that were listed"`
	Namespace  string         `json:"namespace,omitempty" jsonschema:"the namespace that was searched"`
	Items      []ResourceItem `json:"items" jsonschema:"list of resources found"`
	Count      int            `json:"count" jsonschema:"total number of resources found"`
}

// --- GetResource Types ---

type GetResourceInput struct {
	ApiVersion string `json:"apiVersion" jsonschema:"apiVersion of the resource (examples of valid apiVersion are: v1, apps/v1, networking.k8s.io/v1)"`
	Kind       string `json:"kind" jsonschema:"kind of the resource (examples of valid kind are: Pod, Service, Deployment, Ingress)"`
	Name       string `json:"name" jsonschema:"name of the resource to retrieve"`
	Namespace  string `json:"namespace,omitempty" jsonschema:"namespace of the resource (ignored for cluster-scoped resources)"`
}

type GetResourceOutput struct {
	ApiVersion string                 `json:"apiVersion" jsonschema:"the apiVersion of the retrieved resource"`
	Kind       string                 `json:"kind" jsonschema:"the kind of the retrieved resource"`
	Name       string                 `json:"name" jsonschema:"the name of the retrieved resource"`
	Namespace  string                 `json:"namespace,omitempty" jsonschema:"the namespace of the retrieved resource"`
	Labels     map[string]string      `json:"labels,omitempty" jsonschema:"labels on the resource"`
	CreatedAt  string                 `json:"createdAt" jsonschema:"creation timestamp"`
	Status     map[string]interface{} `json:"status,omitempty" jsonschema:"complete status subresource as JSON"`
	Manifest   map[string]interface{} `json:"manifest" jsonschema:"the complete resource manifest as JSON"`
}
