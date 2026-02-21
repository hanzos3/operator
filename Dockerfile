FROM registry.access.redhat.com/ubi9/ubi-minimal:latest as build

RUN microdnf update -y --nodocs && microdnf install ca-certificates -y --nodocs

FROM registry.access.redhat.com/ubi9/ubi-micro:latest

ARG TAG

LABEL name="Hanzo S3" \
      vendor="Hanzo AI, Inc. <dev@hanzo.ai>" \
      maintainer="Hanzo AI, Inc. <dev@hanzo.ai>" \
      version="${TAG}" \
      release="${TAG}" \
      summary="Hanzo S3 Operator brings native support for S3-compatible object storage and encryption to Kubernetes." \
      description="Hanzo S3 is a high-performance, S3-compatible object storage system designed for AI infrastructure, private cloud environments, and mission-critical workloads. 100% open-source under AGPLv3."

# On RHEL the certificate bundle is located at:
# - /etc/pki/tls/certs/ca-bundle.crt (RHEL 6)
# - /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem (RHEL 7)
COPY --from=build /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem /etc/pki/ca-trust/extracted/pem/

COPY CREDITS /licenses/CREDITS
COPY LICENSE /licenses/LICENSE

COPY minio-operator /minio-operator

ENTRYPOINT ["/minio-operator"]
