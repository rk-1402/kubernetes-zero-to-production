# ☸️ Kubernetes: Zero to Production

<div align="center">

<img src="assets/kubernetes.gif" width="90" alt="Kubernetes Animation"/>

### 🚀 Learn Kubernetes from basics to production — without pretending YAML is always fun.

![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?logo=kubernetes&logoColor=white)
![GKE](https://img.shields.io/badge/GKE-4285F4?logo=googlecloud&logoColor=white)
![EKS](https://img.shields.io/badge/EKS-FF9900?logo=amazonaws&logoColor=white)
![GitOps](https://img.shields.io/badge/GitOps-ArgoCD-EF7B4D?logo=argo&logoColor=white)
![Status](https://img.shields.io/badge/Status-Work%20in%20Progress-yellow)

</div>

---

## 🌈 About This Repository

**Kubernetes: Zero to Production** is a practical study guide covering Kubernetes from fundamentals to production-level deployments.

It is designed for:

- ☁️ DevOps Engineers
- 🛠️ Cloud Engineers
- ⚙️ Platform Engineers
- 🔥 SREs
- 🎯 Kubernetes interview preparation
- 🧪 Hands-on learning

The final goal is to convert all chapters into a colorful, structured Kubernetes PDF study guide.

> Because reading 600 pages of dry documentation voluntarily should probably earn a certification by itself. 😄

---

# 🗺️ Learning Path

```text
Fundamentals
    ↓
Workloads
    ↓
Configuration
    ↓
Networking
    ↓
Storage
    ↓
Scheduling
    ↓
Security
    ↓
Observability
    ↓
Scaling
    ↓
Advanced Kubernetes
    ↓
GKE + EKS
    ↓
Troubleshooting
    ↓
Production Labs
```

---

# 📚 Modules

## 🧱 Module 1 — Foundations

- Kubernetes history and concepts
- Containers vs VMs vs Kubernetes
- Control Plane and Worker Nodes
- API Server, Scheduler, etcd
- kubelet and kube-proxy
- Minikube, Kind, kubeadm
- GKE and EKS
- kubectl and kubeconfig

---

## 🚀 Module 2 — Workloads

- Pods
- Init Containers
- Sidecars
- ReplicaSets
- Deployments
- StatefulSets
- DaemonSets
- Jobs and CronJobs
- Rolling updates and rollback

---

## ⚙️ Module 3 — Configuration

- Labels and Selectors
- Annotations
- ConfigMaps
- Secrets
- Namespaces
- Resource Requests and Limits
- Environment Variables
- Downward API

---

## 🌐 Module 4 — Networking

- Kubernetes Networking Model
- ClusterIP
- NodePort
- LoadBalancer
- Headless Services
- Ingress
- Network Policies
- CoreDNS
- Calico
- Cilium
- AWS VPC CNI
- GKE Dataplane V2
- Istio and Linkerd

---

## 💾 Module 5 — Storage

- Volumes
- PV and PVC
- Storage Classes
- Dynamic Provisioning
- CSI Drivers
- EBS
- EFS
- GCE Persistent Disk
- Filestore
- StatefulSet storage patterns

---

## 🎯 Module 6 — Scheduling

- Scheduler internals
- Node Selectors
- Affinity / Anti-Affinity
- Taints and Tolerations
- Topology Spread Constraints
- Priority and Preemption
- Custom Schedulers

---

## 🔐 Module 7 — Security

- Authentication and Authorization
- RBAC
- Service Accounts
- Pod Security Admission
- Security Contexts
- Sealed Secrets
- External Secrets
- GKE Workload Identity
- EKS IRSA
- OPA Gatekeeper
- Kyverno

---

## 📊 Module 8 — Observability

- kubectl logs
- Kubernetes Events
- Prometheus
- Grafana
- Loki
- CloudWatch
- GKE Cloud Logging
- Metrics Server
- Liveness Probes
- Readiness Probes
- Startup Probes
- Distributed Tracing

---

## 📈 Module 9 — Scaling & Reliability

- HPA
- VPA
- Cluster Autoscaler
- Karpenter
- GKE Node Auto-Provisioning
- Pod Disruption Budgets
- Multi-Zone HA
- Multi-Region architecture

---

## 🧠 Module 10 — Advanced Kubernetes

- CRDs
- Operators
- Helm
- Kustomize
- ArgoCD
- Flux
- Admission Controllers
- API Aggregation
- etcd Backup and Restore
- Multi-Cluster Management

> Congratulations. You have now entered the part where Kubernetes starts asking *you* interview questions. 😄

---

# ☁️ Module 11 — GKE & EKS

### 🔵 GKE

- GKE Standard
- GKE Autopilot
- Workload Identity
- Managed Prometheus
- GCE Ingress
- Cluster upgrades

### 🟠 EKS

- Managed Node Groups
- Self-Managed Nodes
- Fargate
- IRSA
- Karpenter
- ALB Controller
- Cluster upgrades

### ⚔️ Comparisons

- GKE Standard vs Autopilot
- EKS Node Groups vs Fargate
- GKE Workload Identity vs EKS IRSA

---

# 🎤 Module 12 — Interview Preparation

Includes:

- 🔥 50+ troubleshooting scenarios
- 🏗️ 30+ Kubernetes system design questions
- 💬 Behavioral + Kubernetes questions

Common scenarios:

```text
CrashLoopBackOff
ImagePullBackOff
OOMKilled
Pending Pods
Node NotReady
DNS Failures
PVC Pending
FailedScheduling
Ingress Issues
Service Connectivity Issues
```

---

# 🧪 Module 13 — Capstone Projects

10 production-style projects combining multiple Kubernetes concepts.

### 🚀 Projects

1. 3-Tier Application — Minikube, GKE and EKS
2. GitHub Actions → Docker → Kubernetes CI/CD
3. HPA + Cluster Autoscaler + Karpenter
4. Ingress + TLS + Let's Encrypt
5. Prometheus + Grafana + Loki + Alertmanager
6. PostgreSQL StatefulSet + Backup/Restore
7. Kubernetes Security Hardening
8. GitOps with ArgoCD
9. Custom Kubernetes Operator
10. Kubernetes Disaster Recovery

Every lab includes:

```text
Objective
Architecture
Prerequisites
Setup
YAML
Commands
Verification
Troubleshooting
Cleanup
Bonus Challenges
```

---

# 🛠️ Repository Structure

```text
kubernetes-zero-to-production/
│
├── README.md
├── assets/
│
├── docs/
│   ├── 01-foundations/
│   ├── 02-workloads/
│   ├── 03-configuration/
│   ├── 04-networking/
│   ├── 05-storage/
│   ├── 06-scheduling/
│   ├── 07-security/
│   ├── 08-observability/
│   ├── 09-scaling-reliability/
│   ├── 10-advanced-topics/
│   ├── 11-managed-kubernetes/
│   ├── 12-interview-guide/
│   └── 13-capstone-projects/
│
├── manifests/
├── diagrams/
├── scripts/
└── releases/
```

---

# 💻 Commands You Will Use

```bash
kubectl get pods
kubectl describe pod
kubectl logs
kubectl exec
kubectl apply
kubectl delete
kubectl rollout
kubectl scale
kubectl top
kubectl debug
kubectl auth
kubectl config
```

And yes, eventually:

```bash
kubectl get pods
```

will become muscle memory. 😄

---

# 📕 Final Kubernetes Study Guide

The completed repository will eventually become:

```text
Kubernetes: Zero to Production
Complete Kubernetes Study Guide
```

The PDF will contain:

- 🎨 Colorful diagrams
- 🧠 Concepts
- 💻 kubectl commands
- 📄 YAML examples
- ☁️ GKE and EKS
- 🔐 Security
- 📊 Monitoring
- 🛠️ Troubleshooting
- 🎤 Interview questions
- 🧪 Production labs

Final PDF location:

```text
releases/kubernetes-zero-to-production.pdf
```

---

# 🚀 Clone This Repository

```bash
git clone https://github.com/rk-1402/kubernetes-zero-to-production.git
```

```bash
cd kubernetes-zero-to-production
```

Start here:

```text
docs/01-foundations/
```

---

# 🤝 Contributions

Found an error?

Found a better Kubernetes command?

Found another YAML indentation problem that somehow consumed two hours?

Contributions are welcome. 😄

Fork → Branch → Commit → Pull Request.

---

# ⭐ Support

If this repository helps you understand Kubernetes, consider giving it a ⭐.

It helps other engineers discover the project.

---

# ⚠️ Cloud Cost Warning

GKE, EKS, Load Balancers, NAT Gateways, persistent disks, and other cloud resources can generate charges.

Clean up your lab resources after testing.

Your cloud provider will happily keep billing even if you forgot the cluster exists. 😄

---

## 👨‍💻 Author

**Ravi Kumar Gupta**

`DevOps` • `AWS` • `GCP` • `Kubernetes` • `Linux` • `Cloud Infrastructure`

---

> ### ☸️ Learn it. Break it. Debug it. Automate it. Run it in production.
