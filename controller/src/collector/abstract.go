package collector

import (
	"auto-passport/targets"

	"github.com/latifik2/auto-passport-controller/types"
	"k8s.io/client-go/kubernetes"
)

type AbstractCollector interface {
	GetStaticTargets() []targets.StaticTarget
	GetDynamicTargets() []targets.DynamicTarget
	GetMetadata([]targets.Target) []types.CommonPassport
	GetK8sClientSet() *kubernetes.Clientset
}
