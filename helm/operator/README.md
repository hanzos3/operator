# Hanzo S3 ![license](https://img.shields.io/badge/license-AGPL%20V3-blue)

[Hanzo S3](https://s3.hanzo.ai) is a high-performance, S3-compatible object storage system released under GNU AGPLv3 or later.
Use Hanzo S3 to build high-performance infrastructure for machine learning, analytics, and application data workloads.

For more detailed documentation please visit [here](https://s3.hanzo.ai/docs)

Introduction
------------

This chart bootstraps the Hanzo S3 Operator on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

Configure Helm repo
--------------------

```bash
helm repo add hanzos3 https://hanzos3.github.io/operator
```

Installing the Chart
--------------------

Install this chart using:

```bash
helm install \
  --namespace minio-operator \
  --create-namespace \
  minio-operator hanzos3/operator
```

The command deploys the Hanzo S3 Operator on the Kubernetes cluster in the default configuration.

Creating a Tenant
-----------------

Once the Hanzo S3 Operator Chart is successfully installed, create a tenant using:

```bash
helm install --namespace tenant-ns \
  --create-namespace tenant hanzos3/tenant
```

This creates a 4-node tenant (cluster). To change the default values, take a look at various [values.yaml](https://github.com/hanzos3/operator/blob/main/helm/tenant/values.yaml).
