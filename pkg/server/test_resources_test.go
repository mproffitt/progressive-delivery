package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/fluxcd/flagger/pkg/apis/flagger/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/weaveworks/progressive-delivery/pkg/server"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/rand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func newCanary(
	ctx context.Context,
	t *testing.T,
	k client.Client,
	name, namespace string,
) v1beta1.Canary {
	resource := v1beta1.Canary{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.CanarySpec{
			Provider: "traefik",
			TargetRef: v1beta1.LocalObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       name,
			},
			SkipAnalysis: false,
			AutoscalerRef: &v1beta1.LocalObjectReference{
				APIVersion: "autoscaling/v2",
				Kind:       "HorizontalPodAutoscaler",
				Name:       name,
			},
			Service: v1beta1.CanaryService{
				Port:       80,
				TargetPort: intstr.FromInt(9999),
			},
			Analysis: &v1beta1.CanaryAnalysis{
				Iterations: 1,
				Interval:   "1m",
			},
		},
		Status: v1beta1.CanaryStatus{
			Phase:              v1beta1.CanaryPhaseSucceeded,
			FailedChecks:       0,
			CanaryWeight:       0,
			Iterations:         0,
			LastAppliedSpec:    "5978589476",
			LastPromotedSpec:   "5978589476",
			LastTransitionTime: metav1.NewTime(time.Now()),
			Conditions: []v1beta1.CanaryCondition{
				{
					LastUpdateTime:     metav1.NewTime(time.Now()),
					LastTransitionTime: metav1.NewTime(time.Now()),
					Message:            "Canary analysis completed successfully, promotion finished.",
					Reason:             "Succeeded",
					Status:             "True",
					Type:               v1beta1.PromotedType,
				},
			},
		},
	}

	err := k.Create(ctx, &resource)
	assert.NoError(t, err, "should be able to create canary: %s", resource.GetName())

	return resource
}

func newNamespace(ctx context.Context, t *testing.T, k client.Client) corev1.Namespace {
	ns := corev1.Namespace{}
	ns.Name = "kube-test-" + rand.String(5)

	err := k.Create(ctx, &ns)
	assert.NoError(t, err, "should be able to create namespace: %s", ns.GetName())

	return ns
}

func newDeployment(ctx context.Context, t *testing.T, k client.Client, name string, ns string) *appsv1.Deployment {
	dpl := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				server.LabelKustomizeName:      name,
				server.LabelKustomizeNamespace: ns,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "nginx",
						Image: "nginx",
					}},
				},
			},
		},
	}

	err := k.Create(ctx, dpl)
	assert.NoError(t, err, "should be able to create Deployment: %s", dpl.GetName())

	return dpl
}
