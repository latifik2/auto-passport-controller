package collector

import "auto-passport/targets"

type AbstractCollector interface {
	GetStaticTargets() []targets.StaticTarget
	GetDynamicTargets() []targets.DynamicTarget
	GetMetadata([]targets.Target) []CommonPassport
}

type CommonPassport struct {
	ServiceType    string         // тип сервиса - Airflow, Metabase, Superset и т.д.
	Infrastructure Infrastructure // где сервис развренут - K8s, VM
	Version        string         // версия сервиса
	Severity       string         // уровень важности сервиса - критичный, важный, неважный
}

type Infrastructure struct {
	InfrastructureType string // тип инфраструктуры - K8s, VM
	Host               string // хост, если тип VM
	Cluster            string // кластер, если тип K8s
	Namespace          string // неймспейс, если тип K8s

	// TODO: добавить информацию о ресурах, которые выделены для сервиса
}
