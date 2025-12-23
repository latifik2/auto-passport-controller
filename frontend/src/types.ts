export interface Infrastructure {
  infra_type: "VM" | "K8s";
  host: string;
  cluster: string;
  namespace: string;
}

export interface CommonPassport {
  service_type: string;
  infrastructure: Infrastructure;
  version: string;
  severity: "critical" | "important" | "normal";
}
