package collector

import (
	"auto-passport/targets"
	"auto-passport/utils"
	"context"
	"fmt"
	"log/slog"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetK8sClientSet() *kubernetes.Clientset {

	config, err := rest.InClusterConfig()
	if err != nil {
		slog.Error("Error creating in-cluster config", slog.Any("error", err))
		return nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Error creating Kubernetes client", slog.Any("error", err))
		return nil
	}

	return clientset
}

func GetK8sTargets(ac AbstractCollector, pattern string, config *utils.Config) []targets.DynamicTarget {

	var dynamicTargets []targets.DynamicTarget
	clientset := ac.GetK8sClientSet()

	if clientset == nil {
		slog.Error("K8s clientset is nill. Kubernetes service discovery won`t work")
		return []targets.DynamicTarget{}
		// panic("K8s clientset is nill. Terminating programm...")
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		slog.Error("Error listing namespaces", slog.Any("error", err))
		return []targets.DynamicTarget{}
	}

	slog.Info(fmt.Sprintf("Looking for namespaces, that match regex '%s'", pattern))
	re := regexp.MustCompile(pattern)

	for _, ns := range namespaces.Items {

		if isMatched := re.MatchString(ns.Name); !isMatched {
			continue
		}

		slog.Info(fmt.Sprintf("Namespace matched: %s", ns.Name))

		ingresses, err := clientset.NetworkingV1().Ingresses(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			slog.Error("Error listing ingresses", slog.Any("error", err))
			continue
		}

		for _, ingress := range ingresses.Items {
			for _, rule := range ingress.Spec.Rules {
				host := rule.Host
				slog.Info(fmt.Sprintf("Found Ingress host: %s in namespace: %s", host, ns.Name))
				dynamicTargets = append(dynamicTargets, targets.DynamicTarget{
					Host:      host,
					Cluster:   config.Cluster,
					Namespace: ns.Name,
				})
			}
		}

		// достать из пакета utlits/config.go паеременную CLUSTER, котораяв в свою очередь будет прочитана функцией ReadConfig

	}
	return dynamicTargets
}
