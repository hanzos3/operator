# Using PodSecurityPolicy for Hanzo S3 Pods 

This document explains how to apply `PodSecurityPolicy` to Hanzo S3 Pods created by the Hanzo S3 Operator. A Pod Security Policy is a cluster-level resource that controls security sensitive aspects of the pod specification. Read more in [Kubernetes PodSecurityPolicy Documentation](https://kubernetes.io/docs/concepts/policy/pod-security-policy/).

## Getting Started

Use the example to apply a custom `PodSecurityPolicy` to all the Hanzo S3 Pods created by the Operator.

```
kubectl create -f https://github.com/hanzos3/operator/tree/master/examples/tenant-pod-security-policy.yaml
```

This file creates a custom PodSecurityPolicy. Then it creates a `ClusterRole` attached to the `PodSecurityPolicy`. Finally a `ClusterRoleBinding` bounds the `ClusterRole` to a `ServiceAccount` which is added to all the Hanzo S3 Pods created by the Hanzo S3 Operator.
