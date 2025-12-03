package collector

import (
	"auto-passport/targets"
	"auto-passport/utils"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type AirflowCollector struct{}

type AirflowVersion struct {
	GitVersion string `json:"git_version"`
	Version    string `json:"version"`
}

func (a AirflowCollector) GetStaticTargets() []targets.StaticTarget {
	var staticTargets []targets.StaticTarget
	hosts := utils.ReadConfig("airflow")

	if hosts == nil {
		slog.Error("Unsupported conifg section Airflow")
		return []targets.StaticTarget{}
	}

	if len(hosts) == 0 {
		slog.Info("No static targets found for Airflow in config")
		return []targets.StaticTarget{}
	}

	for _, host := range hosts {
		staticTargets = append(staticTargets, targets.StaticTarget{Host: host})
	}
	return staticTargets
}

func (a AirflowCollector) GetDynamicTargets() []targets.DynamicTarget {

	var dynamicTargets []targets.DynamicTarget

	config, err := rest.InClusterConfig()
	if err != nil {
		slog.Error("Error creating in-cluster config", slog.Any("error", err))
		return []targets.DynamicTarget{}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.Error("Error creating Kubernetes client", slog.Any("error", err))
		return []targets.DynamicTarget{}
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		slog.Error("Error listing namespaces", slog.Any("error", err))
		return []targets.DynamicTarget{}
	}

	slog.Info("Find namespaces, that match regex '.*airflow.*'")
	pattern := ".*airflow.*"
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
					Cluster:   config.Host,
					Namespace: ns.Name,
				})
			}
		}

	}

	return dynamicTargets

}

func (a AirflowCollector) GetMetadata(airflowTargets []targets.Target) []CommonPassport {
	// Implement logic to generate metadata for the given targets

	v1Url := "/api/v1/version"
	v2Url := "/api/v2/version"

	username, password := utils.GetCredentials("airflow")

	creds := username + ":" + password
	b64Creds := base64.StdEncoding.EncodeToString([]byte(creds))

	var jsonVersion []byte
	var errJson error

	var airflowVersion AirflowVersion

	var passports []CommonPassport
	var infrastructure Infrastructure

	var host string

	for _, afTarget := range airflowTargets {
		switch t := afTarget.(type) {
		case targets.StaticTarget:
			host = t.Host
			infrastructure = Infrastructure{
				InfrastructureType: "VM",
				Host:               host,
			}
		case targets.DynamicTarget:
			host = t.Host
			infrastructure = Infrastructure{
				InfrastructureType: "K8s",
				Cluster:            t.Cluster,
				Namespace:          t.Namespace,
				Host:               host,
			}
		}

		jsonVersion, errJson = utils.MakeApiCall(host+v1Url, b64Creds)
		if errJson != nil {
			log.Printf("Error calling Airflow v1 API: %v", errJson)
			log.Printf("Falling back to Airflow v2 API")

			jsonVersion, errJson = utils.MakeApiCall(host+v2Url, b64Creds)
			if errJson != nil {
				log.Printf("Error calling Airflow v2 API: %v", errJson)
				return []CommonPassport{}
			}
		}

		err := json.Unmarshal(jsonVersion, &airflowVersion)
		if err != nil {
			log.Printf("Error parsing Airflow version response: %v", err)
			return []CommonPassport{}
		}

		passport := CommonPassport{
			ServiceType:    "Airflow",
			Infrastructure: infrastructure,
			Version:        airflowVersion.Version,
			Severity:       "critical",
		}
		passports = append(passports, passport)
	}

	return passports
}
