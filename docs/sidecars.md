# Configuring Sidecars for a Tenant

This document explains how to enable configure sidecars for your Hanzo S3 Tenant.

Sidecars are containers that run in the same pod as the Hanzo S3 container, this makes it so they run together on the same machine and have the ability to community with each other over `localhost`.

## Getting Started

### Prerequisites

- Hanzo S3 Operator up and running as explained in the [document here](https://github.com/hanzos3/operator#operator-setup).

## Sidecars Configuration

Sidecars Configuration is a part of Tenant yaml. 

The following example configures a warp container to run in the same pod as the Hanzo S3 pod.

```yaml
...
  sideCars:
    containers:
      - name: warp
        image: ghcr.io/hanzos3/warp:v0.3.21
        args:
          - client
        ports:
          - containerPort: 7761
            name: http
            protocol: TCP
```

**Note:** the Hanzo S3 Service for the tenant won't expose the ports added in the sidecar. It's up to the user to expose these ports with their own services.

A complete list of values is available [here](tenant_crd.adoc##sidecars) in the API reference.