package collector

import (
	"auto-passport/targets"
	"auto-passport/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"

	"k8s.io/client-go/kubernetes"
)

type AirflowCollector struct {
	K8sClientSet *kubernetes.Clientset
	Config       utils.Config
}

type AirflowVersion struct {
	GitVersion string `json:"git_version"`
	Version    string `json:"version"`
}

func (a AirflowCollector) GetStaticTargets() []targets.StaticTarget {
	var staticTargets []targets.StaticTarget
	hosts := a.Config.StaticTargets.Airflow

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

	pattern := ".*airflow.*"

	return GetK8sTargets(a, pattern, &a.Config)

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

		jsonVersion, errJson = utils.MakeApiCall(host, v1Url, b64Creds, "http")
		if errJson != nil {
			slog.Error(fmt.Sprintf("Error calling Airflow v1 API: %v. Falling back to Airflow v2 API", errJson))

			jsonVersion, errJson = utils.MakeApiCall(host, v2Url, b64Creds, "http")
			if errJson != nil {
				slog.Error(fmt.Sprintf("Error calling Airflow v2 API: %v", errJson))
				return []CommonPassport{}
			}
		}

		err := json.Unmarshal(jsonVersion, &airflowVersion)
		if err != nil {
			slog.Error(fmt.Sprintf("Error parsing Airflow version response: %v", err))
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

func (a AirflowCollector) GetK8sClientSet() *kubernetes.Clientset {
	return a.K8sClientSet
}
