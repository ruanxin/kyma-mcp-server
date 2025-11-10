package types

type CreatePodInput struct {
	Name      string `json:"name" jsonschema:"the name for the new pod"`
	Image     string `json:"image" jsonschema:"the container image to use (e.g., 'nginx:latest')"`
	Namespace string `json:"namespace" jsonschema:"the namespace to create the pod in (e.g., 'default')"`
}

type CreatePodOutput struct {
	Status  string `json:"status" jsonschema:"the result of the operation (e.g., 'pod test created')"`
	PodName string `json:"podName" jsonschema:"the name of the created pod"`
	UID     string `json:"uid" jsonschema:"the unique ID of the created pod"`
}

// --- DeletePod Types ---

type DeletePodInput struct {
	Name      string `json:"name" jsonschema:"the name of the pod to delete"`
	Namespace string `json:"namespace" jsonschema:"the namespace of the pod to delete (e.g., 'default')"`
}

type DeletePodOutput struct {
	Status  string `json:"status" jsonschema:"the result of the delete operation"`
	PodName string `json:"podName" jsonschema:"the name of the pod that was deleted"`
}
