# 📘 Topic 51 — Horizontal Pod Autoscaler (HPA)

**Official Documentation Reference**
- HPA concept: https://kubernetes.io/docs/concepts/workloads/autoscaling/horizontal-pod-autoscale/
- HPA walkthrough: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/
- GKE HPA: https://cloud.google.com/kubernetes-engine/docs/concepts/horizontalpodautoscaler
- EKS HPA: https://docs.aws.amazon.com/eks/latest/userguide/horizontal-pod-autoscaler.html

---

## 🔹 Definition
The **Horizontal Pod Autoscaler (HPA)** automatically scales the number of pod replicas in a Deployment, ReplicaSet, or StatefulSet based on observed CPU utilization, memory utilization, or custom/external metrics.

---

## 🔹 About

### How it works
1. HPA controller (in kube-controller-manager) queries metrics every **15 seconds** (default `--horizontal-pod-autoscaler-sync-period`).
2. Calculates desired replicas: `desired = ceil(current × (currentMetric / targetMetric))`.
3. Scales up/down within `minReplicas` and `maxReplicas`.

### Metric types
| Type | Source | Example |
|---|---|---|
| **Resource** | Metrics Server | CPU utilization, memory utilization |
| **Pods** | Custom metrics API | `requests_per_second` per pod |
| **Object** | Custom metrics API | `queue_depth` on a Service |
| **External** | External metrics API | SQS queue length, Pub/Sub backlog |

### Scaling behavior (K8s 1.23+ stable)
```yaml
behavior:
  scaleUp:
    stabilizationWindowSeconds: 0        # scale up immediately
    policies:
    - type: Percent
      value: 100                          # double pods in one step
      periodSeconds: 60
  scaleDown:
    stabilizationWindowSeconds: 300       # wait 5 min before scaling down
    policies:
    - type: Percent
      value: 10                           # remove max 10% per minute
      periodSeconds: 60
```

### Key fields
| Field | Purpose |
|---|---|
| `minReplicas` | Floor (default 1) |
| `maxReplicas` | Ceiling |
| `metrics[].type` | Resource / Pods / Object / External |
| `metrics[].resource.target.type` | `Utilization` (%) or `AverageValue` |
| `behavior` | Fine-grained scale-up/down policies |

### Requirements
- **Metrics Server** must be installed (for resource metrics).
- **Prometheus Adapter** for custom metrics.
- **Stackdriver/CloudWatch Adapter** for external metrics.
- Pods **must have `resources.requests`** — HPA calculates utilization as `actual / request`.

### HPA + VPA interaction
Running both on the same Deployment is tricky — HPA changes replicas, VPA changes requests. Best practice: use HPA for scaling, VPA in `Off` mode for recommendations only.

---

## 🏗️ Diagrams

### HPA control loop
```mermaid
flowchart LR
  MS[Metrics Server / Prometheus Adapter]
  HPA[HPA Controller<br/>every 15s]
  DEPLOY[Deployment<br/>replicas: N]
  MS -->|current CPU| HPA
  HPA -->|desired = ceil(current × actual/target)| DEPLOY
  DEPLOY -->|scale| PODS[Pods]
```

### Multi-metric HPA
```mermaid
flowchart TB
  HPA[HPA]
  HPA --> M1[Resource: CPU 70%]
  HPA --> M2[Pods: RPS 1000]
  HPA --> M3[External: SQS depth < 100]
  HPA --> MAX[Take max desired across all metrics]
  MAX --> SCALE[Scale to max(desired)]
```

---

## 💻 Real-World Examples

### Example 1 — Simple CPU-based HPA
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata: { name: web-hpa }
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: web
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target: { type: Utilization, averageUtilization: 70 }
```

### Example 2 — Multi-metric (CPU + custom RPS)
```yaml
metrics:
- type: Resource
  resource: { name: cpu, target: { type: Utilization, averageUtilization: 70 } }
- type: Pods
  pods:
    metric: { name: http_requests_per_second }
    target: { type: AverageValue, averageValue: "1000" }
```

### Example 3 — External metric (SQS queue)
```yaml
metrics:
- type: External
  external:
    metric:
      name: sqs_queue_length
      selector: { matchLabels: { queue: orders } }
    target: { type: AverageValue, averageValue: "50" }
```

### Example 4 — Behavior tuning
```yaml
behavior:
  scaleUp:
    stabilizationWindowSeconds: 0
    policies:
    - { type: Pods, value: 4, periodSeconds: 60 }     # add max 4 pods/min
    - { type: Percent, value: 100, periodSeconds: 60 } # or double, whichever is larger
    selectPolicy: Max
  scaleDown:
    stabilizationWindowSeconds: 300
    policies:
    - { type: Percent, value: 10, periodSeconds: 60 }
```

### Example 5 — Imperative HPA
```bash
kubectl autoscale deployment web --cpu-percent=70 --min=2 --max=20
kubectl get hpa
kubectl describe hpa web-hpa
```

---

## 💻 Related Commands

| Command | Purpose | Example |
|---|---|---|
| `kubectl autoscale deploy` | Create HPA imperatively | `kubectl autoscale deploy web --cpu-percent=70 --min=2 --max=20` |
| `kubectl get hpa` | Status + current replicas | `kubectl get hpa` |
| `kubectl describe hpa` | Metrics + events + conditions | `kubectl describe hpa web-hpa` |
| `kubectl top pods` | Current CPU/mem | `kubectl top pods -l app=web` |
| `kubectl delete hpa` | Remove | `kubectl delete hpa web-hpa` |

---

## 🚨 Scenarios

### Scenario 1 — HPA shows `<unknown>` for CPU
- **Cause**: Metrics Server not installed or pods missing `resources.requests`.
- **Fix**: install Metrics Server; set CPU requests on all target pods.

### Scenario 2 — HPA scales too aggressively → thrashing
- **Cause**: `stabilizationWindowSeconds: 0` for scale-down; traffic is bursty.
- **Fix**: set `scaleDown.stabilizationWindowSeconds: 300`.

### Scenario 3 — HPA never scales down
- **Cause**: `stabilizationWindowSeconds` too long; or utilization stays above target.
- **Fix**: verify actual utilization with `kubectl top`; tune target percentage.

### Scenario 4 — Custom metric HPA stuck at minReplicas
- **Cause**: Prometheus Adapter not configured for that metric; or metric name mismatch.
- **Fix**: `kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1` → verify metric exists.

---

## ☁️ GKE Perspective
- HPA works on Autopilot and Standard.
- **GKE custom metrics**: use Stackdriver Adapter or GKE Managed Prometheus + Prometheus Adapter.
- **Autopilot**: min pod resources may affect utilization calculations.
- **Multidimensional Pod Autoscaler (MPA)**: GKE-specific CRD combining HPA + VPA.

## 🟧 EKS Perspective
- Metrics Server installed via EKS add-on or Helm.
- **Custom metrics**: Prometheus Adapter or KEDA.
- **External metrics**: CloudWatch Adapter or KEDA with SQS/Kinesis triggers.
- **Karpenter**: complements HPA — HPA adds pods, Karpenter adds nodes.

---

## 🧪 Labs

### Lab 51.1 — CPU-based HPA (Beginner, 30 min)
```bash
kubectl create deploy web --image=registry.k8s.io/hpa-example --requests='cpu=200m'
kubectl expose deploy web --port=80
kubectl autoscale deploy web --cpu-percent=50 --min=1 --max=10
# Load test
kubectl run load --image=busybox --rm -it -- sh -c "while true; do wget -qO- http://web; done"
# Watch
kubectl get hpa web -w
```

### Lab 51.2 — Custom metric HPA with Prometheus (Intermediate, 90 min)
1. Deploy app exposing `http_requests_total`. 2. Install Prometheus + Prometheus Adapter. 3. Create HPA on `http_requests_per_second`. 4. Load test → observe scaling.

### Lab 51.3 — External metric HPA with SQS on EKS (Advanced, 90 min)
1. Create SQS queue. 2. Install KEDA or CloudWatch Adapter. 3. Create HPA on SQS `ApproximateNumberOfMessages`. 4. Push 1000 messages → pods scale up. 5. Drain queue → pods scale down.

---

## 🎯 Interview Deliverable Answer (60–90 s)
"**What**: HPA automatically scales pod replicas based on metrics — CPU, memory, custom (RPS), or external (queue depth). It queries metrics every 15 seconds and calculates desired replicas using the formula `ceil(current × actual/target)`.  
**Why**: Manual scaling doesn't respond to traffic. HPA matches capacity to demand — scaling up under load and down during quiet periods to save cost.  
**How**: I always use `autoscaling/v2` for multi-metric support. CPU target at 70%, plus a custom metric like RPS for traffic-driven services, plus external metrics (SQS depth) for queue workers. I tune `behavior` — fast scale-up (0s stabilization), slow scale-down (300s window, 10%/min) to avoid thrashing. Pods must have `resources.requests` or HPA can't calculate utilization.  
**Example**: On Black Friday our API HPA scaled from 5 to 45 replicas in under 3 minutes based on CPU. Cluster Autoscaler added 8 nodes to accommodate. After traffic subsided, the 300-second cooldown and 10%/min scale-down policy brought us back to 5 replicas over 20 minutes — no manual intervention."

---

## ❓ Interview Questions

**Q1. HPA formula?** `desired = ceil(current × actual/target)`. **Q2. What does HPA need to work?** Metrics Server (or custom API) + `resources.requests` on pods. **Q3. `<unknown>` in HPA — why?** Metrics Server missing or no requests set. **Q4. Metric types?** Resource, Pods, Object, External. **Q5. Multi-metric: which wins?** HPA takes the max desired across all metrics. **Q6. How to prevent thrashing?** `behavior.scaleDown.stabilizationWindowSeconds`. **Q7. HPA + VPA together?** Risky — HPA scales replicas, VPA changes requests (which changes utilization). Use VPA in Off mode. **Q8. Custom metrics source?** Prometheus Adapter or KEDA exposing via `custom.metrics.k8s.io`. **Q9. External metrics example?** SQS queue depth, Pub/Sub backlog, CloudWatch metric. **Q10. HPA on StatefulSet?** Supported but consider ordered scaling semantics and PVC lifecycle.

