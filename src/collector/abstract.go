package collector

import (
	"auto-passport/targets"

	"k8s.io/client-go/kubernetes"
)

type AbstractCollector interface {
	GetStaticTargets() []targets.StaticTarget
	GetDynamicTargets() []targets.DynamicTarget
	GetMetadata([]targets.Target) []CommonPassport
	GetK8sClientSet() *kubernetes.Clientset
}
