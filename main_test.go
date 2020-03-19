package main

import (
	"testing"

	cmp "github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	fake "k8s.io/client-go/kubernetes/fake"
)

func getSecretManagerConfig(secretManager string) secretManagerConfig {
	var smCfg secretManagerConfig
	smCfg.aws.config.enabled = false
	smCfg.gcp.config.enabled = false
	smCfg.vault.config.enabled = false

	switch secretManager {
	case "aws":
		smCfg.aws.config.enabled = true
		smCfg.aws.config.region = "us-west-2"
		smCfg.aws.config.roleARN = "arn:aws:iam::user:role/secretmanger"
		smCfg.aws.config.secretName = "test-aws-secret"
		smCfg.aws.config.previousVersion = "true"
	case "gcp":
		smCfg.gcp.config.enabled = true
		smCfg.gcp.config.projectID = "project-x"
		smCfg.gcp.config.secretName = "gcp-test-secret"
		smCfg.gcp.config.secretVersion = "5"
		smCfg.gcp.config.serviceAccountKeySecretName = "gcp-credentials"
	case "vault-k8s":
		smCfg.vault.config.enabled = true
		smCfg.vault.config.addr = "https://vault:8200"
		smCfg.vault.config.path = "/secret/data/top-secret"
		smCfg.vault.config.role = "x-role"
		smCfg.vault.config.tlsSecretName = "vault-tls"
		smCfg.vault.config.tokenPath = "/tmp/key"
		smCfg.vault.config.backend = "kubernetes"
		smCfg.vault.config.useSecretNamesAsKeys = true
	case "vault-gcp":
		smCfg.vault.config.enabled = true
		smCfg.vault.config.addr = "https://vault:8200"
		smCfg.vault.config.path = "/secret/data/top-secret"
		smCfg.vault.config.role = "x-role"
		smCfg.vault.config.tlsSecretName = "vault-tls"
		smCfg.vault.config.tokenPath = "/tmp/key"
		smCfg.vault.config.backend = "gcp"
		smCfg.vault.config.useSecretNamesAsKeys = true
		smCfg.vault.config.gcpServiceAccountKeySecretName = "vault-sa-gcp-creds"
	case "vault-secret-version":
		smCfg.vault.config.enabled = true
		smCfg.vault.config.addr = "https://vault:8200"
		smCfg.vault.config.path = "/secret/data/top-secret"
		smCfg.vault.config.role = "x-role"
		smCfg.vault.config.tlsSecretName = "vault-tls"
		smCfg.vault.config.tokenPath = "/tmp/key"
		smCfg.vault.config.backend = "kubernetes"
		smCfg.vault.config.useSecretNamesAsKeys = true
		smCfg.vault.config.version = "2"
	default:
		return smCfg
	}
	return smCfg
}

func Test_mutatingWebhook_mutateContainers(t *testing.T) {
	// arguments passed to the webhook
	type args struct {
		containers          []corev1.Container
		podSpec             *corev1.PodSpec
		ns                  string
		secretManagerConfig secretManagerConfig
	}

	type fields struct {
		k8sClient kubernetes.Interface
	}

	tests := []struct {
		name             string
		fields           fields
		args             args
		mutated          bool
		wantErr          bool
		wantedContainers []corev1.Container
	}{
		{
			name: "Will not mutate container without enabled aws gcp or vault annotation",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "MyContainer",
						Image:   "some-image",
						Command: []string{"/bin/bash"},
						Args:    nil,
						Env: []corev1.EnvVar{
							{Name: "SOME_VARIABLE", Value: "non-of-your-business"},
						},
					},
				},
				secretManagerConfig: getSecretManagerConfig(""),
			},
			mutated: false,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{

					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/bin/bash"},
					Args:    nil,
					Env: []corev1.EnvVar{
						{Name: "SOME_VARIABLE", Value: "non-of-your-business"},
					},
				},
			},
		}, {
			name: "Will mutate container for AWS",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "MyContainer",
						Image:   "some-image",
						Command: []string{"/app"},
						Args:    nil,
						Env: []corev1.EnvVar{
							{Name: "SOME_VARIABLE", Value: "non-of-your-business"},
						},
					},
				},
				secretManagerConfig: getSecretManagerConfig("aws"),
			},
			mutated: true,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{
					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args:    []string{"/app"},
					Env: []corev1.EnvVar{
						{Name: "SOME_VARIABLE", Value: "non-of-your-business"},
						{Name: "SECRET_NAME", Value: "test-aws-secret"},
						{Name: "REGION", Value: "us-west-2"},
						{Name: "ROLE_ARN", Value: "arn:aws:iam::user:role/secretmanger"},
						{Name: "PREVIOUS_VERSION", Value: "true"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "secrets-consumer-env", MountPath: "/secrets-consumer"},
					},
				},
			},
		}, {
			name: "Will mutate container for GCP",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "MyContainer",
						Image:   "some-image",
						Command: []string{"/app"},
						Args:    nil,
						Env: []corev1.EnvVar{
							{Name: "HOST", Value: "127.0.0.1"},
						},
					},
				},
				secretManagerConfig: getSecretManagerConfig("gcp"),
			},
			mutated: true,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{
					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args:    []string{"/app"},
					Env: []corev1.EnvVar{
						{Name: "HOST", Value: "127.0.0.1"},
						{Name: "SECRET_NAME", Value: "gcp-test-secret"},
						{Name: "PROJECT_ID", Value: "project-x"},
						{Name: "SECRET_VERSION", Value: "5"},
						{
							Name:  "GOOGLE_APPLICATION_CREDENTIALS",
							Value: "/var/run/secret/cloud.google.com/service-account.json",
						},
					},
					VolumeMounts: []v1.VolumeMount{
						{Name: "google-cloud-key", MountPath: "/var/run/secret/cloud.google.com"},
						{Name: "secrets-consumer-env", MountPath: "/secrets-consumer"},
					},
				},
			},
		}, {
			name: "Will mutate container for Vault with Kubernetes backend",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "MyContainer",
						Image:   "some-image",
						Command: []string{"/app"},
						Args:    nil,
						Env: []corev1.EnvVar{
							{Name: "HOST", Value: "127.0.0.1"},
						},
					},
				},
				secretManagerConfig: getSecretManagerConfig("vault-k8s"),
			},
			mutated: true,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{
					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args:    []string{"/app"},
					Env: []corev1.EnvVar{
						{Name: "HOST", Value: "127.0.0.1"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_PATH", Value: "/secret/data/top-secret"},
						{Name: "VAULT_ROLE", Value: "x-role"},
						{Name: "TOKEN_PATH", Value: "/tmp/key"},
						{Name: "VAULT_USE_SECRET_NAMES_AS_KEYS", Value: "true"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/ca.pem"},
					},
					VolumeMounts: []v1.VolumeMount{
						{Name: "vault-tls", MountPath: "/etc/tls/ca.pem", SubPath: "ca.pem"},
						{Name: "secrets-consumer-env", MountPath: "/secrets-consumer"},
					},
				},
			},
		}, {
			name: "Will mutate container for Vault with GCP backend",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "MyContainer",
						Image:   "some-image",
						Command: []string{"/app"},
						Args:    nil,
						Env: []corev1.EnvVar{
							{Name: "HOST", Value: "127.0.0.1"},
						},
					},
				},
				secretManagerConfig: getSecretManagerConfig("vault-gcp"),
			},
			mutated: true,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{
					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args:    []string{"/app"},
					Env: []corev1.EnvVar{
						{Name: "HOST", Value: "127.0.0.1"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_PATH", Value: "/secret/data/top-secret"},
						{Name: "VAULT_ROLE", Value: "x-role"},
						{Name: "VAULT_BACKEND", Value: "gcp"},
						{
							Name:  "GOOGLE_APPLICATION_CREDENTIALS",
							Value: "/var/run/secret/cloud.google.com/service-account.json",
						},
						{Name: "TOKEN_PATH", Value: "/tmp/key"},
						{Name: "VAULT_USE_SECRET_NAMES_AS_KEYS", Value: "true"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/ca.pem"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "google-cloud-key", MountPath: "/var/run/secret/cloud.google.com"},
						{Name: "vault-tls", MountPath: "/etc/tls/ca.pem", SubPath: "ca.pem"},
						{Name: "secrets-consumer-env", MountPath: "/secrets-consumer"},
					},
				},
			},
		}, {
			name: "Will mutate container for Vault with a prefix",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "MyContainer",
						Image:   "some-image",
						Command: []string{"/app"},
						Args:    nil,
						Env: []corev1.EnvVar{
							{Name: "API_KEY", Value: "vault:API_KEY"},
						},
					},
				},
				secretManagerConfig: getSecretManagerConfig("vault-gcp"),
			},
			mutated: true,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{
					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args:    []string{"/app"},
					Env: []corev1.EnvVar{
						{Name: "API_KEY", Value: "vault:API_KEY"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_PATH", Value: "/secret/data/top-secret"},
						{Name: "VAULT_ROLE", Value: "x-role"},
						{Name: "VAULT_BACKEND", Value: "gcp"},
						{
							Name:  "GOOGLE_APPLICATION_CREDENTIALS",
							Value: "/var/run/secret/cloud.google.com/service-account.json",
						},
						{Name: "TOKEN_PATH", Value: "/tmp/key"},
						{Name: "VAULT_USE_SECRET_NAMES_AS_KEYS", Value: "true"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/ca.pem"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "google-cloud-key", MountPath: "/var/run/secret/cloud.google.com"},
						{Name: "vault-tls", MountPath: "/etc/tls/ca.pem", SubPath: "ca.pem"},
						{Name: "secrets-consumer-env", MountPath: "/secrets-consumer"},
					},
				},
			},
		}, {
			name: "Will mutate container for Vault with a prefix",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "MyContainer",
						Image:   "some-image",
						Command: []string{"/app"},
						Args:    nil,
						Env: []corev1.EnvVar{
							{Name: "API_KEY", Value: "vault:API_KEY"},
						},
					},
				},
				secretManagerConfig: getSecretManagerConfig("vault-secret-version"),
			},
			mutated: true,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{
					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args:    []string{"/app"},
					Env: []corev1.EnvVar{
						{Name: "API_KEY", Value: "vault:API_KEY"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_PATH", Value: "/secret/data/top-secret"},
						{Name: "VAULT_ROLE", Value: "x-role"},
						{Name: "TOKEN_PATH", Value: "/tmp/key"},
						{Name: "VAULT_USE_SECRET_NAMES_AS_KEYS", Value: "true"},
						{Name: "VAULT_SECRET_VERSION", Value: "2"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/ca.pem"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "vault-tls", MountPath: "/etc/tls/ca.pem", SubPath: "ca.pem"},
						{Name: "secrets-consumer-env", MountPath: "/secrets-consumer"},
					},
				},
			},
		},
	}

	// subtests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := &mutatingWebhook{
				k8sClient: tt.fields.k8sClient,
			}
			// t.Logf("args: %+v", tt.args)
			got, err := mw.mutateContainers(tt.args.containers, tt.args.podSpec, tt.args.secretManagerConfig, tt.args.ns)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutatingWebhook.mutateContainers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.mutated {
				t.Errorf("mutatingWebhook.mutateContainers() = %v, want %v", got, tt.mutated)
			}
			if !cmp.Equal(tt.args.containers, tt.wantedContainers) {
				t.Errorf("mutatingWebhook.mutateContainers() = diff %v", cmp.Diff(tt.args.containers, tt.wantedContainers))
			}
		})
	}
}
