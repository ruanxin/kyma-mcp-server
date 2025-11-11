package services

import (
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestConvertToUnstructured(t *testing.T) {
	tests := []struct {
		name          string
		manifest      string
		expectedGVK   *schema.GroupVersionKind
		expectedName  string
		expectedNS    string
		expectedKind  string
		shouldError   bool
		errorContains string
	}{
		{
			name: "valid pod manifest YAML",
			manifest: `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: default
spec:
  containers:
  - name: nginx
    image: nginx:latest`,
			expectedGVK: &schema.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "Pod",
			},
			expectedName: "test-pod",
			expectedNS:   "default",
			expectedKind: "Pod",
			shouldError:  false,
		},
		{
			name: "valid deployment manifest YAML",
			manifest: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-deployment
  namespace: test-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - name: nginx
        image: nginx:latest`,
			expectedGVK: &schema.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
			expectedName: "test-deployment",
			expectedNS:   "test-ns",
			expectedKind: "Deployment",
			shouldError:  false,
		},
		{
			name: "valid service manifest JSON",
			manifest: `{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "test-service",
    "namespace": "kube-system"
  },
  "spec": {
    "ports": [
      {
        "port": 80,
        "targetPort": 8080
      }
    ],
    "selector": {
      "app": "test"
    }
  }
}`,
			expectedGVK: &schema.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "Service",
			},
			expectedName: "test-service",
			expectedNS:   "kube-system",
			expectedKind: "Service",
			shouldError:  false,
		},
		{
			name: "cluster-scoped resource (ClusterRole)",
			manifest: `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: test-clusterrole
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]`,
			expectedGVK: &schema.GroupVersionKind{
				Group:   "rbac.authorization.k8s.io",
				Version: "v1",
				Kind:    "ClusterRole",
			},
			expectedName: "test-clusterrole",
			expectedNS:   "", // ClusterRole is cluster-scoped
			expectedKind: "ClusterRole",
			shouldError:  false,
		},
		{
			name:          "invalid YAML syntax",
			manifest:      `apiVersion: v1\nkind: Pod\nmetadata:\n  name: test-pod\n  invalid: [unclosed`,
			shouldError:   true,
			errorContains: "",
		},
		{
			name:          "empty manifest",
			manifest:      "",
			shouldError:   true,
			errorContains: "",
		},
		{
			name: "missing apiVersion",
			manifest: `kind: Pod
metadata:
  name: test-pod`,
			expectedGVK: &schema.GroupVersionKind{
				Group:   "",
				Version: "",
				Kind:    "Pod",
			},
			expectedName: "test-pod",
			expectedKind: "Pod",
			shouldError:  false, // The decoder might still work
		},
		{
			name: "missing kind",
			manifest: `apiVersion: v1
metadata:
  name: test-resource`,
			expectedGVK: &schema.GroupVersionKind{
				Group:   "",
				Version: "v1",
				Kind:    "",
			},
			expectedName: "test-resource",
			shouldError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj, gvk, err := convertToUnstructured(tt.manifest)

			if tt.shouldError {
				if err == nil {
					t.Errorf("convertToUnstructured() expected error but got none")
				}
				if tt.errorContains != "" && err != nil {
					if !contains(err.Error(), tt.errorContains) {
						t.Errorf("convertToUnstructured() error = %v, expected to contain %v", err, tt.errorContains)
					}
				}
				return
			}

			if err != nil {
				t.Errorf("convertToUnstructured() error = %v, expected no error", err)
				return
			}

			if obj == nil {
				t.Error("convertToUnstructured() returned nil object")
				return
			}

			if gvk == nil {
				t.Error("convertToUnstructured() returned nil GVK")
				return
			}

			// Verify GVK
			if tt.expectedGVK != nil {
				if gvk.Group != tt.expectedGVK.Group {
					t.Errorf("convertToUnstructured() GVK Group = %v, expected %v", gvk.Group, tt.expectedGVK.Group)
				}
				if gvk.Version != tt.expectedGVK.Version {
					t.Errorf("convertToUnstructured() GVK Version = %v, expected %v", gvk.Version, tt.expectedGVK.Version)
				}
				if gvk.Kind != tt.expectedGVK.Kind {
					t.Errorf("convertToUnstructured() GVK Kind = %v, expected %v", gvk.Kind, tt.expectedGVK.Kind)
				}
			}

			// Verify object metadata
			if tt.expectedName != "" {
				if obj.GetName() != tt.expectedName {
					t.Errorf("convertToUnstructured() object name = %v, expected %v", obj.GetName(), tt.expectedName)
				}
			}

			if tt.expectedNS != "" {
				if obj.GetNamespace() != tt.expectedNS {
					t.Errorf("convertToUnstructured() object namespace = %v, expected %v", obj.GetNamespace(), tt.expectedNS)
				}
			}

			if tt.expectedKind != "" {
				if obj.GetKind() != tt.expectedKind {
					t.Errorf("convertToUnstructured() object kind = %v, expected %v", obj.GetKind(), tt.expectedKind)
				}
			}

			// Verify that we got a proper unstructured object
			if _, ok := obj.Object["metadata"]; !ok && tt.expectedName != "" {
				t.Error("convertToUnstructured() object missing metadata field")
			}
		})
	}
}

func TestConvertToUnstructured_ObjectFields(t *testing.T) {
	manifest := `apiVersion: v1
kind: Pod
metadata:
  name: test-pod
  namespace: test-ns
  labels:
    app: test
    version: "1.0"
  annotations:
    description: "test pod"
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80`

	obj, gvk, err := convertToUnstructured(manifest)
	if err != nil {
		t.Fatalf("convertToUnstructured() error = %v", err)
	}

	// Test that we can access nested fields
	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil {
		t.Errorf("Error accessing nested containers: %v", err)
	}
	if !found {
		t.Error("containers field not found in spec")
	}
	if len(containers) != 1 {
		t.Errorf("Expected 1 container, got %d", len(containers))
	}

	// Test labels
	labels := obj.GetLabels()
	if labels["app"] != "test" {
		t.Errorf("Expected label app=test, got %v", labels["app"])
	}
	if labels["version"] != "1.0" {
		t.Errorf("Expected label version=1.0, got %v", labels["version"])
	}

	// Test annotations
	annotations := obj.GetAnnotations()
	if annotations["description"] != "test pod" {
		t.Errorf("Expected annotation description='test pod', got %v", annotations["description"])
	}

	// Verify GVK is correctly set
	if gvk.Kind != "Pod" || gvk.Version != "v1" || gvk.Group != "" {
		t.Errorf("Unexpected GVK: %+v", gvk)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) > 0 && indexOfString(s, substr) >= 0))
}

// Simple substring search
func indexOfString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
