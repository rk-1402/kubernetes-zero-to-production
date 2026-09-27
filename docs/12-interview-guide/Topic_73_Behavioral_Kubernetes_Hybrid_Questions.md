## 📘 Topic 73 — Behavioral + Kubernetes Hybrid Questions

Interviews often blend "tell me about a time…" with a Kubernetes twist. Use the **STAR framework** (Situation → Task → Action → Result) for each.

---

### 🎯 1. "Tell me about a time you diagnosed a production Kubernetes incident."

**STAR sketch**
- **S**: Payment service on EKS started returning 5xx during Black Friday peak.
- **T**: Restore <10 min; root-cause without a rollback.
- **A**: Checked HPA (maxed at 30), then `kubectl top pods` → memory near limit. `kubectl logs --previous` showed OOMKilled restart pattern. Scaled HPA `maxReplicas`, raised memory limit via hotfix Deployment patch, deployed. Followed up with VPA analysis + profiling to fix the leak.
- **R**: 8-min recovery; no failed transactions; leak fixed within 48 h; VPA + alerting added.
- **Lesson**: Ship VPA recommendations early; don't rely on average-case memory sizing.

---

### 🎯 2. "Describe a change you rolled back and why."

- Deployed a new ingress-nginx version; latency P99 spiked 2×.
- Rolled back via `helm rollback` in 3 min after Grafana alert.
- Traced to a default `worker_shutdown_timeout` change.
- Follow-up: added canary rollout with automated latency SLO analysis in Argo Rollouts.

---

### 🎯 3. "Tell me about a conflict with a teammate on architecture."

- Colleague wanted 1 giant cluster per environment; I argued for multiple by BU.
- Discussed blast radius, upgrade windows, quota conflicts.
- Compromise: one shared platform cluster + per-team workload clusters + fleet-wide GitOps.
- Result: onboarding time dropped from 2 weeks to 2 days.

---

### 🎯 4. "Time you influenced without authority."

- Champion for GitOps in a team using `kubectl apply` from laptops.
- Ran a lunch-and-learn, built a PoC ArgoCD instance in 1 day, wrote migration guide.
- Adopted org-wide within a quarter; audit trail from Git satisfied compliance.

---

### 🎯 5. "Time you made a mistake — how did you handle it?"

- Deleted a CRD without checking finalizers → cascading delete of prod DBs.
- Recovered from etcd snapshot within an hour; wrote a public postmortem.
- Introduced `kubectl` alias with confirmation prompt; added Kyverno policy blocking CRD deletes without label `sre-approved=true`.

---

### 🎯 6. "How do you keep up with Kubernetes changes?"

- KubeCon talks, Kubernetes changelog, KEPs, weekly release-team notes.
- Test each new minor in a Kind cluster before promoting.
- Run `kubent` in CI to detect deprecated APIs.

---

### 🎯 7. "Time you optimized cost."

- Migrated 500-node EKS from CA + MNG to Karpenter with 80/20 spot/on-demand.
- Introduced Namespace-level `ResourceQuota` and enforced default requests via LimitRange.
- Cut monthly compute by 42% ($180k → $104k) with no SLO regression.

---

### 🎯 8. "Time you had to onboard a junior."

- Structured 30/60/90 plan around Pods → Deployments → Services → Ingress → HPA.
- Paired on real tickets from day one; had them own a small operator.
- After 90 days they led a runbook for a common troubleshooting scenario.

---

### 🎯 9. "How do you approach security in Kubernetes?"

- Layered: image scanning (Trivy) → admission (Kyverno) → runtime (Falco) → network (NetworkPolicy) → identity (WLI/IRSA).
- Baseline PSS `restricted`; explicit exceptions via `pod-security.kubernetes.io/enforce=privileged` on tightly-controlled namespaces.
- Quarterly RBAC audits via `rakkess`.

---

### 🎯 10. "Time you designed for reliability under failure."

- Migrated stateful service to StatefulSet with 3 replicas across 3 AZs.
- Added PDB `maxUnavailable: 1`, topology spread, anti-affinity.
- Simulated zone loss with chaos experiment; hit SLO throughout.

---

### 🎯 11. "How do you decide between GKE and EKS for a new project?"

- Team's cloud skills, existing footprint, compliance boundaries.
- If ops maturity is low: GKE Autopilot.
- Deep AWS ecosystem dependence (SQS, Aurora, IAM ABAC): EKS.
- Multi-cloud strategy: use fleet abstractions to abstract differences.

---

### 🎯 12. "Time you communicated a complex K8s topic to non-technical stakeholders."

- Explained multi-region DR strategy to finance and legal using a car-insurance analogy.
- Got budget for cross-region backups and quarterly game-day.

---

### 🎯 13. "How do you handle disagreement on rollout speed?"

- Product wants weekly canary; SRE wants monthly.
- Meet in the middle: automated canary with SLO-based auto-abort.
- Data-driven review after 3 months; adjust cadence.

---

### 🎯 14. "Describe a time you missed a deadline."

- Underestimated ingress migration from ELB to ALB with cert-manager.
- Communicated early, replanned, delivered 5 days late.
- Built a checklist for future migrations to spot 3rd-party dependencies upfront.

---

### 🎯 15. "How do you decide when to use an Operator vs a Helm chart?"

- Helm for stateless, install-and-forget.
- Operator when you need day-2 behaviors: failover, backup, upgrade coordination.
- Prefer community-maintained operators; write custom only when domain-specific.

---

### 🎯 16. "Tell me about a K8s upgrade that went sideways."

- Upgrading GKE 1.24 → 1.25 broke a Helm chart using removed `PodSecurityPolicy` API.
- Discovered mid-upgrade; paused via blue-green; migrated chart to PSA + Kyverno; resumed.
- Now run `kubent` and canary each minor in dev.

---

### 🎯 17. "How do you drive operational excellence?"

- SLIs/SLOs per service tied to error budget.
- Weekly review of top toil sources; automate top 3 each quarter.
- Runbooks in Git, versioned with charts.

---

### 🎯 18. "How do you decide between Cluster Autoscaler and Karpenter on EKS?"

- CA is proven, works with ASG boundaries.
- Karpenter is faster, more flexible, cost-optimized, but younger.
- I default to Karpenter for new EKS clusters; retain CA on regulated clusters where change control is slower.

---

### 🎯 19. "Tell me about a time you led an on-call rotation improvement."

- On-call fatigue from 40+ pages/week.
- Categorized pages, killed 60% via better probes and alert tuning.
- Introduced SLO-based paging; page count dropped to <10/week.

---

### 🎯 20. "How do you approach production readiness reviews?"

- Checklist: probes, resources, PDB, HPA, RBAC, NetworkPolicy, image provenance, dashboards, runbook, chaos test result.
- Sign-off from SRE + Security before Ingress goes live.

---

### 🎯 21. "Handling ambiguity"

- Asked to 'move everything to Kubernetes' with no scope.
- Ran discovery: catalog apps, categorize by lift-and-shift vs refactor vs retire; produced a 6-month migration wave plan with metrics.
- Got exec buy-in; delivered wave 1 on time.

---

### 🎯 22. "Ownership under failure"

- Pager fires at 3 AM; feature I owned was root cause of DB overload.
- Woke myself, coordinated rollback, wrote postmortem, closed all action items within 2 sprints.

---

### 🎯 23. "Data-driven decision"

- Debated Istio vs Linkerd for mesh.
- Ran a bake-off measuring P99, memory footprint, incident count for 4 weeks; picked Linkerd for our latency profile with a documented list of features we gave up.

---

### 🎯 24. "How do you mentor engineers on K8s?"

- Curated learning path: CKAD → CKA → CKS.
- Weekly office hours; real-cluster labs for common scenarios.
- Rotate junior engineers into a shadow-on-call program.

---

### 🎯 25. "Long-term technical strategy"

- Started policy-as-code + progressive delivery + FinOps together as a 12-month platform bet.
- Measured DORA metrics quarterly; deploy frequency 4×, MTTR halved, cost/req dropped 30%.

---

