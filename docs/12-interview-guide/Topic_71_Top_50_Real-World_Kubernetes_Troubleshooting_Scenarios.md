## 📘 Topic 71 — Top 50 Real-World Kubernetes Troubleshooting Scenarios

**Reference material**
- https://kubernetes.io/docs/tasks/debug/
- https://kubernetes.io/docs/tasks/debug/debug-application/
- https://kubernetes.io/docs/tasks/debug/debug-cluster/
- https://cloud.google.com/kubernetes-engine/docs/troubleshooting
- https://docs.aws.amazon.com/eks/latest/userguide/troubleshooting.html

Each scenario follows: **Symptoms → Root causes → Diagnostic commands → Resolution → Prevention.**

---

### 🚨 1. Pod stuck in `Pending`

**Symptoms**: `kubectl get pods` shows `Pending`; no container starts.

**Common causes**
- Insufficient node capacity (CPU/mem/GPU).
- PVC not bound.
- No node matches nodeSelector / taints / affinity.
- ResourceQuota exceeded in namespace.
- Cluster autoscaler at max size.

**Diagnose**
```bash
kubectl describe pod <pod>
kubectl get events --sort-by=.lastTimestamp
kubectl top nodes
kubectl get pvc
kubectl get quota -n <ns>
```
Look for `FailedScheduling` and the exact reason.

**Fix**
- Add nodes / scale node pool.
- Fix nodeSelector or tolerations.
- Create the referenced StorageClass / PV.
- Raise or delete blocking ResourceQuota.

**Prevent**
- Set sensible requests; monitor cluster utilization; alert on `unschedulable_pods_total`.

---

### 🚨 2. `CrashLoopBackOff`

**Symptoms**: Pod restarts repeatedly; `RESTARTS` climbs.

**Causes**
- App crashes on startup (config, missing env, DB unreachable).
- Failing liveness probe.
- OOM (see #4).
- Wrong command/args.

**Diagnose**
```bash
kubectl logs <pod> --previous
kubectl describe pod <pod>
kubectl get pod <pod> -o yaml | grep -A5 lastState
```

**Fix**
- Fix bug or config; restart.
- Raise `initialDelaySeconds` on livenessProbe.
- Add missing ConfigMap/Secret.

**Prevent**
- Add readiness probe to avoid routing during warmup.
- Structured logging + startup diagnostics.

---

### 🚨 3. `ImagePullBackOff` / `ErrImagePull`

**Symptoms**: Pod shows `ImagePullBackOff`.

**Causes**
- Wrong image name/tag.
- Private registry, missing `imagePullSecret`.
- Rate-limited by Docker Hub.
- Node has no network route to registry.

**Diagnose**
```bash
kubectl describe pod <pod> | grep -A5 Events
```

**Fix**
- Correct image / re-tag.
- Create Secret + reference in Pod:
  ```bash
  kubectl create secret docker-registry regcred --docker-server=... --docker-username=... --docker-password=... --docker-email=...
  ```
- Use ECR/GAR with IRSA/WLI so no secret needed.

**Prevent**
- Mirror upstream images to your own registry.
- Digest-pinning (`image@sha256:...`) for immutability.

---

### 🚨 4. `OOMKilled`

**Symptoms**: `kubectl describe pod` → `Last State: Terminated, Reason: OOMKilled, ExitCode: 137`.

**Causes**
- Container's memory usage exceeded `resources.limits.memory`.
- Memory leak.
- Underestimated JVM/Go/Node heap.

**Diagnose**
```bash
kubectl top pod <pod> --containers
kubectl describe pod <pod>
```
Compare `limits.memory` to observed usage.

**Fix**
- Raise limits (make sure node has capacity).
- Tune GC / heap flags.
- Fix leak.

**Prevent**
- VPA in `Off` recommender mode for right-sizing.
- Grafana alert on memory-usage / limit ratio > 80%.

---

### 🚨 5. Node `NotReady`

**Symptoms**: `kubectl get nodes` shows `NotReady`.

**Causes**
- kubelet down or unhealthy.
- CNI plugin failure.
- Disk pressure / PID pressure.
- Network partition.
- Docker/containerd hung.

**Diagnose**
```bash
kubectl describe node <node>
# on node:
journalctl -u kubelet
crictl ps
df -h
```

**Fix**
- Restart kubelet.
- Reinstall CNI.
- Reboot the node (GKE/EKS: replace it).

**Prevent**
- Monitor kubelet + container-runtime health; auto-replace unhealthy nodes.

---

### 🚨 6. DNS resolution failures inside cluster

**Symptoms**: `nslookup kubernetes.default` fails; apps get `no such host`.

**Causes**
- CoreDNS Pods down / OOM.
- NetworkPolicy blocking DNS to `kube-system`.
- Wrong `dnsPolicy` on Pod.
- Node-local DNS cache issue.

**Diagnose**
```bash
kubectl -n kube-system get pods -l k8s-app=kube-dns
kubectl -n kube-system logs -l k8s-app=kube-dns
kubectl exec -it netshoot -- dig kubernetes.default
```

**Fix**
- Scale CoreDNS; raise memory limits.
- Add NetworkPolicy allowing UDP/TCP 53 to `kube-system`.
- Restart node-local DNS DaemonSet.

**Prevent**
- Run CoreDNS with PDB and topology spread; monitor `coredns_dns_request_duration_seconds`.

---

### 7. Service returns 503 / connection refused

**Causes**
- Endpoints empty (no ready Pods).
- Selector mismatch.
- Pod readinessProbe failing.
- kube-proxy dead.

**Diagnose**
```bash
kubectl get endpoints <svc>
kubectl get pods -l <selector> -o wide
kubectl describe pod <pod> | grep -A3 Readiness
```

**Fix**: Align labels; fix probe; scale up healthy Pods.

**Prevent**: CI test that the selector matches at least one label.

---

### 8. `Terminating` Pod stuck

**Causes**: Finalizer not removed, node lost, PV won't detach.

**Diagnose**
```bash
kubectl get pod <p> -o json | jq .metadata.finalizers
kubectl describe pod <p>
```

**Fix**
- `kubectl patch pod <p> --type=merge -p '{"metadata":{"finalizers":[]}}'`
- Force delete: `kubectl delete pod <p> --grace-period=0 --force`
- Repair the underlying controller.

**Prevent**: Idempotent finalizer logic in operators.

---

### 9. PVC stuck `Pending`

**Causes**: No StorageClass, StorageClass missing provisioner, quota, no matching PV.

**Diagnose**
```bash
kubectl describe pvc <pvc>
kubectl get sc
kubectl get pv
```

**Fix**: Create SC; annotate default SC; ensure CSI driver Pods healthy.

**Prevent**: `is-default-class: true` annotation on cluster's default SC.

---

### 10. Pod can't mount PVC (multi-attach error)

**Causes**: RWO volume attached to another node; StatefulSet scaled down mid-rollout.

**Diagnose**
```bash
kubectl describe pod <p>
kubectl get volumeattachments
```

**Fix**: Wait for detach; force delete the old attachment; use RWX volume (EFS/Filestore).

**Prevent**: Only RWX for multi-Pod workloads.

---

### 11. HPA never scales

**Causes**: metrics-server not installed; missing CPU/memory requests; wrong target ref.

**Diagnose**
```bash
kubectl get hpa <name>
kubectl top pods -l app=<x>
kubectl describe hpa <name>
```

**Fix**: Install metrics-server; add requests; check `scaleTargetRef` kind/name.

**Prevent**: Alert on `hpa_scaler_metric_stale`.

---

### 12. Cluster Autoscaler never adds nodes

**Causes**: Node group at max; no matching instance type; Pod requests exceed largest node; missing tolerations on new nodes.

**Diagnose**: CA logs (`kubectl logs -n kube-system deploy/cluster-autoscaler`), events on Pending Pod.

**Fix**: Raise max; add larger instance type; fix nodeSelector.

**Prevent**: Regularly reconcile `max` size vs peak demand.

---

### 13. HPA and VPA fighting

**Symptoms**: Pod restarts loop, replicas thrash.

**Fix**: Use VPA in `Off` mode (recommendations only) with HPA on CPU; or VPA + HPA on custom metric.

**Prevent**: Documented autoscaling policy per workload.

---

### 14. `Unauthorized` from `kubectl`

**Causes**: Expired token, wrong context, kubeconfig missing `exec` plugin (`aws eks get-token`, `gke-gcloud-auth-plugin`).

**Fix**: `aws eks update-kubeconfig`, `gcloud container clusters get-credentials`, refresh token.

**Prevent**: SSO with short-lived tokens; rotate CI kubeconfigs.

---

### 15. `Forbidden` on a resource

**Cause**: RBAC missing.

**Diagnose**
```bash
kubectl auth can-i <verb> <resource> -n <ns>
kubectl auth can-i <verb> <resource> --as=<user>
```

**Fix**: Add Role/ClusterRole binding.

**Prevent**: RBAC lint in CI (e.g., `rbac-lookup`, `polaris`).

---

### 16. Ingress returns 404 or 502

**Causes**: Missing `IngressClass`, wrong `pathType`, backend Service selector mismatch, LB annotation typo, TLS cert missing.

**Diagnose**: `kubectl describe ingress`; check controller logs; `curl -H "Host: ..."` to LB IP.

**Fix**: Set `ingressClassName`; verify backend Service targets ready endpoints.

**Prevent**: Synthetic checks per Ingress route.

---

### 17. TLS cert not issued (cert-manager)

**Causes**: DNS-01 propagation slow, HTTP-01 firewall, ClusterIssuer misconfig.

**Diagnose**: `kubectl describe certificate/challenge/order`.

**Fix**: Fix DNS/records; open port 80 to LB; check ACME account key.

---

### 18. NetworkPolicy blocking wanted traffic

**Symptoms**: Timeouts within cluster.

**Diagnose**: `kubectl describe networkpolicy`; test with `kubectl exec busybox -- nc -zv target 80`.

**Fix**: Add egress rule; whitelist DNS to `kube-system`.

**Prevent**: Default deny + explicit allow patterns; policy tests with `netpol-analyzer`.

---

### 19. StatefulSet Pod stuck `0/1 Ready`

**Causes**: PVC pending, ordinal peer down, headless Service missing.

**Diagnose**: `kubectl get pvc`; verify headless Service.

**Fix**: Create Service `clusterIP: None`; ensure peers healthy.

---

### 20. `evicted` Pods due to node pressure

**Causes**: Disk / memory / PID pressure.

**Diagnose**: `kubectl describe node`; look for `MemoryPressure=True`.

**Fix**: Increase node size; add ephemeral-storage requests; clean image cache.

**Prevent**: Node-local ephemeral-storage requests + eviction thresholds.

---

### 21. `too many open files`

**Fix**: Raise `ulimits` via node OS or LimitRange; on GKE/EKS use custom AMI/OS config.

---

### 22. Application drops requests during rolling update

**Causes**: No `preStop` hook + `terminationGracePeriodSeconds`; no readinessProbe.

**Fix**:
```yaml
lifecycle:
  preStop:
    exec:
      command: ["/bin/sh", "-c", "sleep 15"]
terminationGracePeriodSeconds: 60
```
Add readinessProbe so the Service removes the Pod before it exits.

---

### 23. Rolling update stuck

**Causes**: New Pods failing readiness; PDB blocking eviction.

**Fix**: `kubectl rollout status`; `kubectl rollout undo`; lower PDB minAvailable temporarily.

---

### 24. `deployment.apps` stuck at old ReplicaSet

**Cause**: Deployment paused or `progressDeadlineSeconds` exceeded.

**Fix**: `kubectl rollout resume deploy/<x>`; check logs of failing pods.

---

### 25. Node clock skew breaks TLS

**Symptoms**: `x509: certificate is not yet valid`.

**Fix**: Install/enable NTP; on EKS/GKE the nodes ship with time sync.

---

### 26. `Too many pods` on node

**Cause**: `--max-pods` limit hit (EKS VPC CNI computes based on ENI/IP limits).

**Fix**: Use larger instance; enable prefix delegation on VPC CNI (`ENABLE_PREFIX_DELEGATION=true`); or switch CNI.

---

### 27. Cluster suddenly slow to schedule

**Causes**: apiserver overloaded, etcd slow disk, huge list requests, watch storm.

**Diagnose**: `apiserver_request_duration_seconds`, `etcd_disk_wal_fsync_duration_seconds`, `apiserver_flowcontrol_dispatched_requests_total`.

**Fix**: Add FlowSchemas, tune informers, upgrade disk class.

---

### 28. Pod's outbound traffic times out

**Causes**: NetworkPolicy egress, NAT gateway saturated, cloud firewall.

**Fix**: Allow egress; add NAT capacity; check cloud VPC routes.

---

### 29. Node cannot join cluster (kubeadm)

**Cause**: Bad token, apiserver unreachable, wrong CA hash.

**Fix**: `kubeadm token create --print-join-command`.

---

### 30. Controller-manager not electing leader

**Cause**: Clock skew, Lease API blocked, etcd write latency.

**Fix**: Repair time sync; raise `--leader-elect-lease-duration`.

---

### 31. `no route to host` between Pods on different nodes

**Cause**: CNI misconfig, MTU mismatch, cloud firewall.

**Fix**: Ensure CNI matches cluster CIDR; lower Pod MTU (`1450` on GCP with encapsulation).

---

### 32. Container image huge, slow start

**Fix**: Multi-stage builds; distroless base; layer caching; use image streaming (GKE Image Streaming, EKS via lazy pulling).

---

### 33. Secrets not updated in running Pod

**Cause**: Env vars from Secret are captured at Pod start.

**Fix**: Mount as volume (auto-refresh) or use Reloader / Stakater.

---

### 34. Job never completes

**Cause**: Pod OOMKilled retried past `backoffLimit`; wrong `restartPolicy`.

**Fix**: Set correct `restartPolicy: OnFailure`, adjust `backoffLimit` and `activeDeadlineSeconds`.

---

### 35. CronJob missing runs

**Cause**: `concurrencyPolicy: Forbid` + long-running job; `startingDeadlineSeconds` too small.

**Fix**: Adjust concurrency; use `Replace` if idempotent.

---

### 36. Persistent volume filled up

**Fix**: Expand volume (`kubectl edit pvc` → `resources.requests.storage`, needs `allowVolumeExpansion: true` on SC).

---

### 37. Autopilot rejects a Pod

**Cause**: Privileged, hostPath, missing requests, or violating compute-class ratios.

**Fix**: Remove offending fields; add requests; choose correct compute class.

---

### 38. Fargate Pod stuck Pending on EKS

**Cause**: No matching Fargate Profile.

**Fix**: Create Fargate Profile matching namespace + labels.

---

### 39. GKE Ingress not creating LB

**Cause**: `ingressClass` wrong (`gce` vs `nginx`); missing `NEG` annotation.

**Fix**: `kubernetes.io/ingress.class: gce`; add `cloud.google.com/neg: '{"ingress": true}'` on Service.

---

### 40. EKS ALB Ingress: `no target group`

**Cause**: AWS Load Balancer Controller not installed or IRSA missing perms.

**Fix**: Install controller with correct IAM policy.

---

### 41. Metrics-server cannot reach kubelet

**Fix**: `--kubelet-insecure-tls`, `--kubelet-preferred-address-types=InternalIP`.

---

### 42. Prometheus scrape targets down

**Cause**: NetworkPolicy blocks Prometheus namespace.

**Fix**: Allow ingress from `monitoring`.

---

### 43. GKE node upgrade disrupts stateful app

**Fix**: Use PDB + blue-green upgrades; drain deadline `--pdb-timeout-seconds`.

---

### 44. Karpenter kills node hosting critical Pod

**Fix**: Annotate Pod `karpenter.sh/do-not-disrupt: "true"`; add PDB.

---

### 45. IngressClass conflict

**Cause**: Two controllers watching same class.

**Fix**: Distinct `IngressClass` per controller; set `--ingress-class` flag.

---

### 46. `kubectl exec` hangs

**Cause**: Long-lived TCP through LB with idle timeout; SPDY blocked.

**Fix**: Use `--tty=false` for scripts; enable websocket support in LB.

---

### 47. Cluster upgrade blocked by deprecated API

**Fix**: Run `kubent`; migrate manifests to current API version.

---

### 48. `too many requests` (429) from apiserver

**Cause**: APF throttling; controller doing unbounded lists.

**Fix**: Add FlowSchema at `workload-high`; move to informer-based caching.

---

### 49. Pod gets wrong ServiceAccount

**Cause**: `automountServiceAccountToken: false` on SA/Pod.

**Fix**: Enable automount or use projected token.

---

### 50. Workload works in dev, fails in prod (same YAML)

**Common differentiators**
- Different Kubernetes minor version (API deprecations).
- Prod PDB / ResourceQuota / NetworkPolicy stricter.
- Prod uses Autopilot / Fargate — enforced admission constraints.
- Prod IRSA/WLI missing → SDK falls back to node role that fails.
- Prod uses PSA `restricted` — dev uses `privileged`.

**Fix**: Use identical baseline (same CRDs, same admission policies, same node OS) via GitOps.

**Prevent**: Prod-parity dev environment; policy-as-code enforcement in CI.

---

