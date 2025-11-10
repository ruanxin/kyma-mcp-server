package services

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceService struct {
	clnt *kubernetes.Clientset
}

func NewNamespaceService(clnt *kubernetes.Clientset) *NamespaceService {
	return &NamespaceService{clnt: clnt}
}

// EnsureNamespace checks if the given namespace exists in the cluster.
// If it does not exist, it creates the namespace.
func (n *NamespaceService) EnsureNamespace(ctx context.Context, ns string) error {
	if ns == "" {
		ns = metav1.NamespaceDefault
	}

	_, err := n.clnt.CoreV1().Namespaces().Get(ctx, ns, metav1.GetOptions{})
	if err == nil {
		// Namespace exists
		return nil
	}

	// Attempt to create the namespace
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}

	_, err = n.clnt.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace %s: %w", ns, err)
	}

	return nil
}
