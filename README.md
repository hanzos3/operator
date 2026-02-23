# Hanzo S3 Kubernetes Operator

[![build](https://img.shields.io/badge/build-passing-green.svg)](https://github.com/hanzos3/operator/actions) [![license](https://img.shields.io/badge/license-AGPL%20V3-blue)](https://github.com/hanzos3/operator/blob/main/LICENSE)

Hanzo S3 is a Kubernetes-native, high-performance object store with an S3-compatible API.
The Hanzo S3 Kubernetes Operator supports deploying S3 tenants onto private and public
cloud infrastructures ("Hybrid" Cloud).

- **Server**: [github.com/hanzoai/s3](https://github.com/hanzoai/s3)
- **Operator**: [github.com/hanzos3/operator](https://github.com/hanzos3/operator)
- **Domains**: [s3.hanzo.ai](https://s3.hanzo.ai) / [hanzo.space](https://hanzo.space)

This README provides a high-level description of the Hanzo S3 Operator and
quickstart instructions.

## Table of Contents

* [Architecture](#architecture)
* [Deploy the Hanzo S3 Operator and Create a Tenant](#deploy-the-hanzo-s3-operator-and-create-a-tenant)
    * [Prerequisites](#prerequisites)
    * [Procedure](#procedure)

# Architecture

Each tenant represents an independent S3-compatible object store within
the Kubernetes cluster. The following diagram describes the architecture of a
tenant deployed into Kubernetes:

![Tenant Architecture](docs/images/architecture.png)

The Hanzo S3 Operator provides multiple methods for accessing and managing tenants.

# Deploy the Hanzo S3 Operator and Create a Tenant

This procedure installs the Hanzo S3 Operator and creates a 4-node tenant for supporting object storage operations in
a Kubernetes cluster.

## Prerequisites

### Kubernetes 1.30.0 or Later

Starting with Operator v7.1.1, Hanzo S3 requires Kubernetes version 1.30.0 or later.

This procedure assumes the host machine has [`kubectl`](https://kubernetes.io/docs/tasks/tools) installed and configured
with access to the target Kubernetes cluster.

### Tenant Namespace

Hanzo S3 supports no more than *one* tenant per namespace. The following `kubectl` command creates a new namespace
for the tenant.

```sh
kubectl create namespace hanzo-s3-tenant
```

### Tenant Storage Class

The Hanzo S3 Kubernetes Operator automatically generates Persistent Volume Claims (`PVC`) as part of deploying a
tenant.

The plugin defaults to creating each `PVC` with the `default`
Kubernetes [`Storage Class`](https://kubernetes.io/docs/concepts/storage/storage-classes/). If the `default` storage
class cannot support the generated `PVC`, the tenant may fail to deploy.

Tenants *require* that the `StorageClass` sets `volumeBindingMode` to `WaitForFirstConsumer`. The
default `StorageClass` may use the `Immediate` setting, which can cause complications during `PVC` binding. We
strongly recommend creating a custom `StorageClass` for use by `PV` supporting a tenant.

The following `StorageClass` object contains the appropriate fields for supporting a tenant:

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: directpv-min-io
provisioner: kubernetes.io/no-provisioner
volumeBindingMode: WaitForFirstConsumer
```

### Tenant Persistent Volumes

The Hanzo S3 Operator generates one Persistent Volume Claim (PVC) for each volume in the tenant *plus* two PVC to support
collecting tenant metrics and logs. The cluster *must* have
sufficient [Persistent Volumes](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) that meet the capacity
requirements of each PVC for the tenant to start correctly. For example, deploying a tenant with 16 volumes requires
18 (16 + 2). If each PVC requests 1TB capacity, then each PV must also provide *at least* 1TB of capacity.

For clusters which cannot deploy DirectPV,
use [Local Persistent Volumes](https://kubernetes.io/docs/concepts/storage/volumes/#local). The following example YAML
describes a local persistent volume:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: <PV-NAME>
spec:
  capacity:
    storage: 1Ti
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: local-storage
  local:
    path: </mnt/disks/ssd1>
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - <NODE-NAME>
```

Replace values in brackets `<VALUE>` with the appropriate value for the local drive.

You can estimate the number of PVC by multiplying the number of server pods in the tenant by the number of
drives per node. For example, a 4-node tenant with 4 drives per node requires 16 PVC and therefore 16 PV.

We *strongly recommend* using the following CSI drivers for creating local PV to ensure best object storage
performance:

- [DirectPV](https://github.com/hanzos3/directpv)
- [Local Persistent Volume](https://kubernetes.io/docs/concepts/storage/volumes/#local)

## Procedure

### 1) Install the Hanzo S3 Operator via Kustomization

The standard `kubectl` tool ships with support
for [kustomize](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/) out of the box, so you can
use that to install the operator.

```sh
kubectl kustomize github.com/hanzos3/operator\?ref=v7.1.1 | kubectl apply -f -
```

Run the following command to verify the status of the Operator:

```sh
kubectl get pods -n hanzo-s3-operator
```

The output resembles the following:

```sh
NAME                              READY   STATUS    RESTARTS   AGE
hanzo-s3-operator-69fd675557-lsrqg   1/1     Running   0          99s
```

### 2) Build the Tenant Configuration

We provide a variety of examples for creating tenants in the `examples` directory. The following example creates a
4-node tenant with 4 volumes per node:

```yaml
kubectl apply -k github.com/hanzos3/operator/examples/kustomization/base
```

### 3) Connect to the Tenant

Use the following command to list the services created by the Hanzo S3
Operator:

```sh
kubectl get svc -n NAMESPACE
```

Replace `NAMESPACE` with the namespace for the tenant. The output
resembles the following:

```sh
NAME                             TYPE            CLUSTER-IP        EXTERNAL-IP   PORT(S)
hanzo-s3                         LoadBalancer    10.104.10.9       <pending>     443:31834/TCP
myminio-console           LoadBalancer    10.104.216.5      <pending>     9443:31425/TCP
myminio-hl                ClusterIP       None              <none>        9000/TCP
```

Applications *internal* to the Kubernetes cluster should use the `hanzo-s3` service for performing object storage
operations on the tenant.

Administrators of the tenant should use the `myminio-console` service to access the console and manage the
tenant, such as provisioning users, groups, and policies.

Tenants deploy with TLS enabled by default, where the Hanzo S3 Operator uses the
Kubernetes `certificates.k8s.io` API to generate the required x.509 certificates. Each
certificate is signed using the Kubernetes Certificate Authority (CA) configured during
cluster deployment. While Kubernetes mounts this CA on Pods in the cluster, Pods do
*not* trust that CA by default. You must copy the CA to a directory such that the
`update-ca-certificates` utility can find and add it to the system trust store to
enable validation of TLS certificates:

```sh
cp /var/run/secrets/kubernetes.io/serviceaccount/ca.crt /usr/local/share/ca-certificates/
update-ca-certificates
```

For applications *external* to the Kubernetes cluster, you must configure
[Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) or a
[Load Balancer](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer) to
expose the tenant services. Alternatively, you can use the `kubectl port-forward` command
to temporarily forward traffic from the local host to the tenant.

# License

Use of the Hanzo S3 Operator is governed by the GNU AGPLv3 or later, found in the [LICENSE](./LICENSE) file.

# Explore Further

- [Hanzo S3 Documentation](https://s3.hanzo.ai/docs)
- [Examples for Tenant Settings](https://github.com/hanzos3/operator/blob/main/docs/examples.md)
- [Custom Hostname Discovery](https://github.com/hanzos3/operator/blob/main/docs/custom-name-templates.md)
- [Apply PodSecurityPolicy](https://github.com/hanzos3/operator/blob/main/docs/pod-security-policy.md)
- [Deploy Tenant with KES](https://github.com/hanzos3/operator/blob/main/docs/kes.md)
- [Tenant API Documentation](docs/tenant_crd.adoc)
- [Policy Binding API Documentation](docs/policybinding_crd.adoc)
