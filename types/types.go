package types

type CommonPassport struct {
	ServiceType    string         `json:"service_type"`   // тип сервиса - Airflow, Metabase, Superset и т.д.
	Infrastructure Infrastructure `json:"infrastructure"` // где сервис развренут - K8s, VM
	Version        string         `json:"version"`        // версия сервиса
	Severity       string         `json:"severity"`       // уровень важности сервиса - критичный, важный, неважный
}

type Infrastructure struct {
	InfrastructureType string `json:"infra_type"` // тип инфраструктуры - K8s, VM
	Host               string `json:"host"`       // хост, если тип VM
	Cluster            string `json:"cluster"`    // кластер, если тип K8s
	Namespace          string `json:"namespace"`  // неймспейс, если тип K8s

	// TODO: добавить информацию о ресурах, которые выделены для сервиса
}
