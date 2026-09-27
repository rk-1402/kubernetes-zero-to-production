## 📘 Topic 72 — Top 30 System Design Questions Using Kubernetes

### Framework for any K8s system design question
1. **Clarify requirements** — scale, latency, availability, budget, regulation.
2. **Sketch high-level architecture** — clusters, network topology.
3. **Discuss compute strategy** — MNG/Autopilot/Fargate/Karpenter.
4. **Storage** — stateful vs stateless, database choice, PV strategy.
5. **Networking** — ingress, service mesh, cross-cluster, DNS.
6. **Security** — RBAC, network policies, secrets, image supply chain.
7. **Observability** — metrics, logs, traces, alerting.
8. **CI/CD & GitOps** — how do changes reach the cluster.
9. **Reliability** — HPA/VPA/CA, PDBs, multi-AZ/region, DR.
10. **Cost** — where the money goes, optimization levers.

Each question below sketches the answer around this frame.

---

### 🎯 1. Design a global e-commerce site on Kubernetes

**Requirements**: 50M DAU, ≤200 ms P95, active/active in 3 regions, 99.99% SLA.

**High-level**
- 3 regional EKS/GKE clusters + edge CDN (CloudFront/Cloud CDN).
- Global load balancer (Cloud LB / Route 53 latency + Global Accelerator).
- Regional databases with cross-region replication (Aurora Global / Cloud Spanner).
- Redis Cluster per region for cache.
- Kafka MSK / Pub/Sub for events.

**K8s**
- Compute: MNG for baseline + Karpenter for burst; separate node pool for GPUs (recommendations).
- Ingress: ALB / GCE Ingress + WAF.
- Service mesh (Istio) for retries, circuit breakers, mTLS.
- HPA on RPS custom metrics; CA/Karpenter for infra scale.

**Reliability**: PDBs, topology spread across AZs, chaos GameDay.

**Observability**: Central Grafana + Thanos across regions; distributed tracing via Tempo/Jaeger.

**Trade-offs**: Cross-region consistency (eventual) vs latency.

---

### 🎯 2. Multi-tenant SaaS platform

- Tenant isolation via **namespaces** with strict `ResourceQuota` + `LimitRange` + NetworkPolicy per tenant.
- Optional per-tenant node pool with taints for hostile tenants.
- Ingress via wildcard `*.tenant.example.com` + TLS.
- OPA/Kyverno policies enforcing labels `tenant=<id>`.
- Central logging tagged by namespace.

**Trade-off**: Namespace vs cluster-per-tenant — namespaces cheaper but noisier neighbor risk; clusters are stronger boundaries but 10× cost.

---

### 🎯 3. Real-time chat backend (WebSocket + presence)

- Stateful frontends (WebSocket) on StatefulSet with headless Service.
- Sticky sessions via `sessionAffinity: ClientIP` or ALB with cookies.
- Redis for presence (pub/sub).
- HPA on active-connections custom metric.
- Graceful drain via `preStop: sleep` + connection draining on LB.

---

### 🎯 4. CI/CD pipeline running on Kubernetes

- **Tekton** or **Argo Workflows** for pipelines; **GitHub Actions self-hosted runners** on K8s for CI.
- Build with **Kaniko** / **BuildKit** in-cluster (no Docker daemon).
- Cache in S3/GCS; artifacts in ECR/GAR.
- CD via **ArgoCD ApplicationSets** targeting multiple environments.
- Progressive delivery via **Argo Rollouts** or **Flagger** with metrics-driven analysis.

---

### 🎯 5. ML inference platform

- GPU node pool (A10G, L4) with taints.
- **KServe** or **Seldon** as serving framework.
- HPA on QPS (custom metrics) + VPA for right-sizing.
- Model artifacts on S3/GCS via a MinIO-backed cache.
- Autoscale to zero via KServe / Knative for cost.

---

### 🎯 6. Batch data-processing platform

- **Argo Workflows** / **Volcano** for DAGs.
- Spot-only Karpenter NodePool for cost.
- Kubernetes Jobs with `activeDeadlineSeconds` and `backoffLimit`.
- Results in S3/GCS; metadata in Postgres.

---

### 🎯 7. Financial trading system with sub-ms tail latencies

- Bare-metal or dedicated nodes; disable CPU frequency scaling.
- **CPU manager** with `static` policy for exclusive CPUs.
- HugePages, SR-IOV NICs.
- Isolated node pool with `Guaranteed` QoS pods.
- Not a fit for Fargate/Autopilot.

---

### 🎯 8. IoT ingestion platform (1M devices)

- MQTT broker (EMQX) on StatefulSet with LoadBalancer.
- Kafka for buffering.
- Karpenter with c-family nodes for CPU-bound consumers.
- Regional edge clusters via EKS Anywhere / GKE on-prem.

---

### 🎯 9. Video streaming CDN origin

- Nginx/Origin StatefulSet with local NVMe (`local-path` StorageClass).
- Kubernetes-native autoscaling on network throughput metrics.
- Global LB fronting HLS/DASH.
- Object storage tier (S3/GCS) for cold segments.

---

### 🎯 10. Design an in-cluster monitoring stack

- **kube-prometheus-stack** with 2 Prometheus replicas + Thanos sidecar.
- **Loki** for logs; **Tempo** for traces; **Grafana** for UI.
- **Alertmanager** with PagerDuty/Slack.
- Long-term storage: S3/GCS via Thanos + Loki chunks.

---

### 🎯 11. Design fleet-wide config management

- GitOps: ArgoCD ApplicationSets or Anthos Config Sync.
- Bases (Kustomize) + per-cluster overlays.
- Policy Controller / Kyverno for baseline enforcement.

---

### 🎯 12. Multi-region active-active database on K8s

- CockroachDB / YugabyteDB / TiDB StatefulSets with per-region PV.
- Anti-affinity across AZ; TopologySpread across regions via multi-cluster.
- Backup: incremental snapshots to S3/GCS.

**Note**: For most orgs, managed DB (Aurora Global, Spanner) is safer than self-managed on K8s.

---

### 🎯 13. Serverless functions platform

- **Knative Serving** or **OpenFaaS**.
- Scale-to-zero; concurrency-based scaling.
- Ingress via Gateway API / Istio.

---

### 🎯 14. Zero-trust internal APIs

- **Istio ambient / sidecar** with strict mTLS.
- **SPIFFE/SPIRE** identities.
- AuthorizationPolicy per service.
- NetworkPolicy default-deny.

---

### 🎯 15. HA control plane on bare-metal

- **kubeadm** with stacked etcd, 3 control-plane nodes, HAProxy VIP.
- Cluster API for provisioning.
- Regular etcd snapshots to S3-compatible object store.

---

### 🎯 16. Blue-green deployment for a legacy monolith

- Two Deployments (`app-blue`, `app-green`) behind one Service.
- Switch selector to shift traffic.
- Argo Rollouts BlueGreen strategy automates it.

---

### 🎯 17. Canary rollout with automated analysis

- **Argo Rollouts** `Rollout` with `AnalysisTemplate` querying Prometheus for `success_rate > 99%` and `p95_latency < 300ms`.
- Steps: 5% → 25% → 50% → 100%.
- Auto-abort on failure.

---

### 🎯 18. Design an Ingress at 1M RPS

- Multiple **ingress-nginx** DaemonSets or NGINX Plus / Envoy Gateway with HPA.
- LB with cross-zone.
- Layer-7 caching via Envoy.
- Anti-affinity so 1 node loss doesn't kill throughput.

---

### 🎯 19. Migrate monolith to K8s

- Strangler pattern: front with API gateway routing to legacy + new services.
- Data: dual-write then cut over.
- Feature flags for rollback.

---

### 🎯 20. Design workload isolation for compliance (PCI/HIPAA)

- Dedicated node pool with taints; only PCI Pods tolerate.
- Encrypted etcd; encrypted disks; Sealed Secrets / KMS.
- Restricted PSS; strict egress NetworkPolicies.
- Audit logs shipped to SIEM.

---

### 🎯 21. Cost-optimize a 500-node cluster

- Right-size with VPA recommendations.
- Karpenter with 80% spot mix + on-demand fallback.
- Consolidate underutilized nodes.
- Bin-pack via `PodTopologySpread` + priority classes.
- Delete idle namespaces; enforce `LimitRange`.

---

### 🎯 22. Disaster recovery strategy

- Velero cross-region backups (daily + hourly).
- Backup for GKE (workload + PV).
- etcd snapshots (self-managed) every 15 min → immutable S3.
- Quarterly restore drills.

---

### 🎯 23. Design multi-cluster service discovery

- **MCS ServiceExport/ServiceImport** or **Cilium ClusterMesh**.
- Global DNS via ExternalDNS.
- Cross-cluster mTLS via shared root CA.

---

### 🎯 24. Migrating from VM-based deployment to K8s

- Containerize with Dockerfiles / Buildpacks.
- Wrap in Deployment/StatefulSet.
- Introduce HPA gradually.
- Migrate config to ConfigMap / Secret.
- Start with Kubernetes as compute; keep DB external.

---

### 🎯 25. Design an internal developer platform (IDP)

- Backstage as portal.
- Golden path templates → Argo CD apps.
- Crossplane for infra self-service.
- Kyverno policies enforce guardrails.

---

### 🎯 26. Message queue on Kubernetes

- RabbitMQ or Kafka via community operator.
- StatefulSet + PVC per broker.
- Anti-affinity across zones.
- PDB `maxUnavailable: 1`.
- Cross-zone replication for durability.

---

### 🎯 27. GPU sharing platform for research

- NVIDIA GPU Operator.
- **Time-slicing** or **MIG** (Multi-Instance GPU) for A100/H100.
- Queue via `Kueue` / Volcano for fairness.
- Users request `nvidia.com/gpu: 1`.

---

### 🎯 28. Edge cluster (retail store) design

- K3s / KubeEdge lightweight nodes.
- Managed from central hub via ACM / Fleet.
- Local-only workloads; cloud sync for aggregates.
- Offline resilience with local etcd or SQLite backend.

---

### 🎯 29. Security incident response readiness

- Falco or Cilium Tetragon for runtime alerts.
- Automated Pod isolation via label change + NetworkPolicy.
- Central log store (SIEM) with 90-day retention.
- Runbook: quarantine node, evict Pod, snapshot for forensics.

---

### 🎯 30. Design an ephemeral dev environment per PR

- ArgoCD ApplicationSet with `pull request` generator → 1 namespace per PR.
- TTL controller to garbage-collect after N days.
- Ingress: `pr-<num>.dev.example.com` via cert-manager wildcard.
- Cost gate via ResourceQuota per namespace.

---

