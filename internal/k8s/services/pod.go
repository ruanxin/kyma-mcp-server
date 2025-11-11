package services

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodService struct {
	clnt kubernetes.Interface
}

func NewPodService(clnt kubernetes.Interface) *PodService {
	return &PodService{clnt: clnt}
}

func (p *PodService) GetPod(ctx context.Context, ns, name string) (*corev1.Pod, error) {
	return p.clnt.CoreV1().Pods(ns).Get(ctx, name, metav1.GetOptions{})
}

func (p *PodService) DeletePod(ctx context.Context, ns, name string) error {
	return p.clnt.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{})
}

func (p *PodService) CreatePod(ctx context.Context, ns, name string, metadata map[string]interface{}, spec map[string]interface{}) (*corev1.Pod, error) {
	if ns == "" {
		ns = "default"
	}
	nsService := NewNamespaceService(p.clnt)

	if err := nsService.EnsureNamespace(ctx, ns); err != nil {
		return nil, fmt.Errorf("namespace ensure failed: %w", err)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	// Populate metadata
	if metadata != nil {
		if labels, ok := metadata["labels"].(map[string]string); ok {
			pod.ObjectMeta.Labels = labels
		}
		if annotations, ok := metadata["annotations"].(map[string]string); ok {
			pod.ObjectMeta.Annotations = annotations
		}
	}

	// Populate PodSpec
	if spec != nil {
		if containers, ok := spec["containers"].([]map[string]interface{}); ok {
			for _, c := range containers {
				pod.Spec.Containers = append(pod.Spec.Containers, corev1.Container{
					Name:  c["name"].(string),
					Image: c["image"].(string),
				})
			}
		}
	}

	if pod.ObjectMeta.Name == "" && len(pod.Spec.Containers) > 0 {
		pod.ObjectMeta.Name = fmt.Sprintf("%s-%d", pod.Spec.Containers[0].Name, time.Now().Unix())
	}

	return p.clnt.CoreV1().Pods(ns).Create(ctx, pod, metav1.CreateOptions{})
}
