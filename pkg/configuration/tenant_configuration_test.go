// This file is part of Hanzo S3 Operator
// Copyright (c) 2023 Hanzo AI, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package configuration

import (
	"reflect"
	"testing"

	miniov2 "github.com/minio/operator/pkg/apis/minio.min.io/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestEnvVarsToFileContent(t *testing.T) {
	type args struct {
		envVars []corev1.EnvVar
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Basic test case",
			args: args{
				envVars: []corev1.EnvVar{
					{
						Name:  "S3_UPDATE",
						Value: "on",
					},
				},
			},
			want: "export S3_UPDATE=\"on\"\n",
		},
		{
			name: "Two Vars test case",
			args: args{
				envVars: []corev1.EnvVar{
					{
						Name:  "S3_UPDATE",
						Value: "on",
					},
					{
						Name:  "S3_UPDATE_MINISIGN_PUBKEY",
						Value: "RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav",
					},
				},
			},
			want: `export S3_UPDATE="on"
export S3_UPDATE_MINISIGN_PUBKEY="RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := envVarsToFileContent(tt.args.envVars); got != tt.want {
				t.Errorf("envVarsToFileContent() = `%v`, want `%v`", got, tt.want)
			}
		})
	}
}

func TestGetTenantConfiguration(t *testing.T) {
	type args struct {
		tenant         *miniov2.Tenant
		cfgEnvExisting map[string]corev1.EnvVar
	}
	tests := []struct {
		name string
		args args
		want []corev1.EnvVar
	}{
		{
			name: "Defaulted Values",
			args: args{
				tenant:         &miniov2.Tenant{},
				cfgEnvExisting: nil,
			},
			want: []corev1.EnvVar{
				{
					Name:  "S3_ARGS",
					Value: "",
				},
				{
					Name:  "S3_PROMETHEUS_JOB_ID",
					Value: "minio-job",
				},
				{
					Name:  "S3_SERVER_URL",
					Value: "https://minio..svc.cluster.local:443",
				},
				{
					Name:  "S3_UPDATE",
					Value: "on",
				},
				{
					Name:  "S3_UPDATE_MINISIGN_PUBKEY",
					Value: "RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav",
				},
			},
		},
		{
			name: "Tenant has one env var",
			args: args{
				tenant: &miniov2.Tenant{
					Spec: miniov2.TenantSpec{
						Env: []corev1.EnvVar{
							{
								Name:  "TEST",
								Value: "value",
							},
						},
					},
				},
				cfgEnvExisting: nil,
			},
			want: []corev1.EnvVar{
				{
					Name:  "S3_ARGS",
					Value: "",
				},
				{
					Name:  "S3_PROMETHEUS_JOB_ID",
					Value: "minio-job",
				},
				{
					Name:  "S3_SERVER_URL",
					Value: "https://minio..svc.cluster.local:443",
				},
				{
					Name:  "S3_UPDATE",
					Value: "on",
				},
				{
					Name:  "S3_UPDATE_MINISIGN_PUBKEY",
					Value: "RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav",
				},
				{
					Name:  "TEST",
					Value: "value",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.tenant.EnsureDefaults()
			if got := buildTenantEnvs(tt.args.tenant, tt.args.cfgEnvExisting); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildTenantEnvs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConfEnvSecret(t *testing.T) {
	type args struct {
		secret *corev1.Secret
	}
	tests := []struct {
		name string
		args args
		want map[string]corev1.EnvVar
	}{
		{
			name: "Basic case",
			args: args{
				secret: &corev1.Secret{
					Data: map[string][]byte{"config.env": []byte(`export S3_ROOT_USER="minio"
export S3_ROOT_PASSWORD="minio123"
export S3_STORAGE_CLASS_STANDARD="EC:2"
export S3_BROWSER="on"`)},
				},
			},
			want: map[string]corev1.EnvVar{
				"S3_ROOT_USER": {
					Name:  "S3_ROOT_USER",
					Value: "minio",
				},
				"S3_ROOT_PASSWORD": {
					Name:  "S3_ROOT_PASSWORD",
					Value: "minio123",
				},
				"S3_STORAGE_CLASS_STANDARD": {
					Name:  "S3_STORAGE_CLASS_STANDARD",
					Value: "EC:2",
				},
				"S3_BROWSER": {
					Name:  "S3_BROWSER",
					Value: "on",
				},
			},
		},
		{
			name: "Basic case has tabs",
			args: args{
				secret: &corev1.Secret{
					Data: map[string][]byte{"config.env": []byte(`	export S3_ROOT_USER="minio"
	export S3_ROOT_PASSWORD="minio123"
	export S3_STORAGE_CLASS_STANDARD="EC:2"
	export S3_BROWSER="on"`)},
				},
			},
			want: map[string]corev1.EnvVar{
				"S3_ROOT_USER": {
					Name:  "S3_ROOT_USER",
					Value: "minio",
				},
				"S3_ROOT_PASSWORD": {
					Name:  "S3_ROOT_PASSWORD",
					Value: "minio123",
				},
				"S3_STORAGE_CLASS_STANDARD": {
					Name:  "S3_STORAGE_CLASS_STANDARD",
					Value: "EC:2",
				},
				"S3_BROWSER": {
					Name:  "S3_BROWSER",
					Value: "on",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseConfEnvSecret(tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfEnvSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFullTenantConfig(t *testing.T) {
	type args struct {
		tenant       *miniov2.Tenant
		configSecret *corev1.Secret
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty tenant with one env var",
			args: args{
				tenant: &miniov2.Tenant{
					Spec: miniov2.TenantSpec{
						Env: []corev1.EnvVar{
							{
								Name:  "TEST",
								Value: "value",
							},
						},
					},
				},
				configSecret: &corev1.Secret{
					Data: map[string][]byte{"config.env": []byte(`export S3_ROOT_USER="minio"
export S3_ROOT_PASSWORD="minio123"
export S3_STORAGE_CLASS_STANDARD="EC:2"
export S3_BROWSER="on"`)},
				},
			},
			want: `export S3_ARGS=""
export S3_BROWSER="on"
export S3_PROMETHEUS_JOB_ID="minio-job"
export S3_ROOT_PASSWORD="minio123"
export S3_ROOT_USER="minio"
export S3_SERVER_URL="https://minio..svc.cluster.local:443"
export S3_STORAGE_CLASS_STANDARD="EC:2"
export S3_UPDATE="on"
export S3_UPDATE_MINISIGN_PUBKEY="RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav"
export TEST="value"
`,
		},
		{
			name: "Empty tenant; with domains; one env var",
			args: args{
				tenant: &miniov2.Tenant{
					Spec: miniov2.TenantSpec{
						Env: []corev1.EnvVar{
							{
								Name:  "TEST",
								Value: "value",
							},
						},
						Features: &miniov2.Features{
							Domains: &miniov2.TenantDomains{
								Console: "http://console.minio",
							},
						},
					},
				},
				configSecret: &corev1.Secret{
					Data: map[string][]byte{"config.env": []byte(`export S3_ROOT_USER="minio"
export S3_ROOT_PASSWORD="minio123"
export S3_STORAGE_CLASS_STANDARD="EC:2"
export S3_BROWSER="on"`)},
				},
			},
			want: `export S3_ARGS=""
export S3_BROWSER="on"
export S3_BROWSER_REDIRECT_URL="http://console.minio"
export S3_PROMETHEUS_JOB_ID="minio-job"
export S3_ROOT_PASSWORD="minio123"
export S3_ROOT_USER="minio"
export S3_SERVER_URL="https://minio..svc.cluster.local:443"
export S3_STORAGE_CLASS_STANDARD="EC:2"
export S3_UPDATE="on"
export S3_UPDATE_MINISIGN_PUBKEY="RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav"
export TEST="value"
`,
		},
		{
			name: "One Pool Tenant; with domains; one env var",
			args: args{
				tenant: &miniov2.Tenant{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tenant",
						Namespace: "ns-x",
					},
					Spec: miniov2.TenantSpec{
						Env: []corev1.EnvVar{
							{
								Name:  "TEST",
								Value: "value",
							},
						},
						Features: &miniov2.Features{
							Domains: &miniov2.TenantDomains{
								Console: "http://console.minio",
							},
						},
						Pools: []miniov2.Pool{
							{
								Name:                "pool-0",
								Servers:             4,
								VolumesPerServer:    4,
								VolumeClaimTemplate: nil,
							},
						},
					},
				},
				configSecret: &corev1.Secret{
					Data: map[string][]byte{"config.env": []byte(`export S3_ROOT_USER="minio"
export S3_ROOT_PASSWORD="minio123"
export S3_STORAGE_CLASS_STANDARD="EC:2"
export S3_BROWSER="on"`)},
				},
			},
			want: `export S3_ARGS="https://tenant-pool-0-{0...3}.tenant-hl.ns-x.svc.cluster.local/export{0...3}"
export S3_BROWSER="on"
export S3_BROWSER_REDIRECT_URL="http://console.minio"
export S3_PROMETHEUS_JOB_ID="minio-job"
export S3_ROOT_PASSWORD="minio123"
export S3_ROOT_USER="minio"
export S3_SERVER_URL="https://minio.ns-x.svc.cluster.local:443"
export S3_STORAGE_CLASS_STANDARD="EC:2"
export S3_UPDATE="on"
export S3_UPDATE_MINISIGN_PUBKEY="RWTx5Zr1tiHQLwG9keckT0c45M3AGeHD6IvimQHpyRywVWGbP1aVSGav"
export TEST="value"
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.tenant.EnsureDefaults()
			if got, _, _ := GetFullTenantConfig(tt.args.tenant, tt.args.configSecret); got != tt.want {
				t.Errorf("GetFullTenantConfig() = `%v`, want `%v`", got, tt.want)
			}
		})
	}
}
