Here's a breakdown of what KTHW teaches you and what it deliberately leaves out:

---

## What Kubernetes the Hard Way teaches you

**The core bootstrapping process** — you manually stand up every component that tools like `kubeadm` or managed cloud services (GKE, EKS, AKS) do invisibly for you. Specifically:

**Infrastructure & TLS**
- Provisioning VMs and a jumpbox, then wiring them together on the same network
- Running your own Certificate Authority (CA) and generating TLS certs for every component: the API server, kubelet, scheduler, controller-manager, etcd, and admin
- Understanding *why* each component needs its own cert and what SANs/CN values matter

**The control plane internals**
- Setting up **etcd** from scratch — the distributed key-value store that holds all cluster state
- Bootstrapping the **kube-apiserver**, **kube-controller-manager**, and **kube-scheduler** as systemd units
- Writing kubeconfig files by hand, so you deeply understand how authentication and context-switching work

**Worker nodes**
- Setting up **containerd** (the container runtime) directly, not through Docker
- Configuring **CNI plugins** for pod networking
- Bootstrapping **kubelet** and **kube-proxy** as systemd services
- Manually adding static pod network routes between nodes

**What you deeply internalize:** the trust chain (who signs what cert, why), how components authenticate to each other, and that "Kubernetes" is really ~6 separate processes that cooperate.

---

## What KTHW deliberately doesn't cover

KTHW is a single-node control plane, no-automation learning exercise. After finishing it, the major Kubernetes topics still ahead of you are:

**Day-2 operations & production concerns**
- High-availability control planes (etcd clusters, multiple API servers behind a load balancer)
- Cluster upgrades without downtime
- etcd backup, restore, and disaster recovery

**Workload management**
- Deployments, ReplicaSets, StatefulSets, DaemonSets, Jobs/CronJobs
- Rolling updates, rollback strategies, and Pod Disruption Budgets
- Horizontal Pod Autoscaler (HPA) and Vertical Pod Autoscaler (VPA)
- Resource requests/limits and QoS classes

**Networking (the deep end)**
- CNI plugin choices (Calico, Cilium, Flannel) and their trade-offs
- NetworkPolicy (pod-level firewall rules)
- Services: ClusterIP, NodePort, LoadBalancer, Headless
- Ingress controllers and the Gateway API
- DNS (CoreDNS configuration and customization)
- Service meshes (Istio, Linkerd)

**Storage**
- PersistentVolumes, PersistentVolumeClaims, StorageClasses
- Dynamic provisioning
- CSI drivers
- StatefulSet storage patterns

**Security**
- RBAC (roles, rolebindings, service accounts — KTHW barely touches this)
- Pod Security Admission / Pod Security Standards
- Secrets management (Sealed Secrets, External Secrets, Vault integration)
- Network policies, OPA/Gatekeeper, or Kyverno for policy enforcement
- Image scanning and supply chain security

**Observability**
- Logging pipelines (Fluent Bit, Loki)
- Metrics (Prometheus, kube-state-metrics, metrics-server)
- Tracing (OpenTelemetry)
- Dashboards (Grafana)

**The ecosystem / real-world usage**
- Helm (package management)
- Operators and Custom Resource Definitions (CRDs)
- GitOps with Argo CD or Flux
- CI/CD integration
- Multi-cluster and multi-tenancy patterns
- Cloud-managed Kubernetes (EKS, GKE, AKS) and their managed add-ons

---

**The bottom line:** KTHW gives you an unusually solid mental model of *how Kubernetes works under the hood* — most people who skip it treat the control plane as a black box forever. But it covers almost none of the day-to-day operational work of actually *running* workloads on Kubernetes. Everything in the second list is what separates a Kubernetes administrator from someone who just finished the tutorial.