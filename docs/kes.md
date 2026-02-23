# Hanzo S3 Operator KES Configuration 

This document explains how to enable KES with Hanzo S3 Operator.

## Getting Started

### Prerequisites

- Hanzo S3 Operator up and running as explained in the [document here](https://github.com/hanzos3/operator#operator-setup).
- KES requires a KMS backend
  in [configuration](https://raw.githubusercontent.com/hanzos3/operator/master/examples/kes-secret.yaml). Currently KES
  supports [AWS Secrets Manager](https://github.com/hanzos3/kes/wiki/AWS-SecretsManager)
  and [Hashicorp Vault](https://github.com/hanzos3/kes/wiki/Hashicorp-Vault-Keystore) as KMS backend for production.S Set
  up one of these as the KMS backend before setting up KES.

### Create Hanzo S3 Tenant

We have an example Tenant with KES encryption available
at [examples/tenant-kes-encryption](../examples/tenant-kes-encryption).

You can install the example like:

```shell
kubectl apply -k github.com/hanzos3/operator/examples/kustomization/tenant-kes-encryption
```

## KES Configuration

KES Configuration is a part of Tenant yaml file. Check the sample
file [available here](https://raw.githubusercontent.com/hanzos3/operator/master/examples/kustomization/tenant-kes-encryption/tenant.yaml).
The config offers below options

### KES Fields

| Field              | Description                                                                                                                                                                       |
|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| spec.kes           | Defines the KES configuration. Refer [this](https://github.com/minio/kes)                                                                                                         |
| spec.kes.replicas  | Number of KES pods to be created.                                                                                                                                                 |
| spec.kes.image     | Defines the KES image.                                                                                                                                                            |
| spec.kes.kesSecret | Secret to specify KES Configuration. This is a mandatory field.                                                                                                                   |
| spec.kes.metadata  | This allows a way to map metadata to the KES pods. Internally `metadata` is a struct type as [explained here](https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta). |

A complete list of values is available [here](tenant_crd.adoc#kesconfig) in the API reference.
