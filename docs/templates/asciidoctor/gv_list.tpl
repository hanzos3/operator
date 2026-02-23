{{- define "gvList" -}}
{{- $groupVersions := . -}}

// Generated documentation. Please do not edit.
:anchor_prefix: k8s-api

[id="{p}-api-reference"]
== API Reference

:minio-image: https://github.com/hanzos3/operator/pkgs[ghcr.io/hanzos3/s3:RELEASE.2025-04-08T15-41-24Z]
:kes-image: https://github.com/hanzos3/kes/pkgs[ghcr.io/hanzos3/kes:2025-03-12T09-35-18Z]
:mc-image: https://github.com/hanzos3/mc/pkgs[ghcr.io/hanzos3/mc:RELEASE.2024-10-02T08-27-28Z]

{{ range $groupVersions }}
{{ template "gvDetails" . }}
{{ end }}

{{- end -}}
