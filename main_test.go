package main

import (
	"testing"

	cmp "github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
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
		smCfg.vault.config.vaultCACert = "vault-ca.pem"
		smCfg.vault.config.tokenPath = "/tmp/key"
		smCfg.vault.config.backend = "kubernetes"
		smCfg.vault.config.useSecretNamesAsKeys = true
		smCfg.vault.config.kubernetesBackend = "/alt/kubernetes/path"
		smCfg.vault.config.version = "5"
	case "vault-multi":
		smCfg.vault.config.enabled = true
		smCfg.vault.config.addr = "https://vault:8200"
		smCfg.vault.config.role = "x-role"
		smCfg.vault.config.tlsSecretName = "vault-tls"
		smCfg.vault.config.vaultCACert = "vault-ca.pem"
		smCfg.vault.config.backend = "kubernetes"
		smCfg.vault.config.secretConfigs = []string{
			`{"path": "/some/secret/path-1", "version": "3", "use-secret-names-as-keys":  true}`,
			`{"path": "/some/secret/path-2"}`,
			`{"path": "/some/secret/path-3", "use-secret-names-as-keys":  true}`,
		}

	case "vault-gcp":
		smCfg.vault.config.enabled = true
		smCfg.vault.config.addr = "https://vault:8200"
		smCfg.vault.config.path = "/secret/data/top-secret"
		smCfg.vault.config.role = "x-role"
		smCfg.vault.config.tlsSecretName = "vault-tls"
		smCfg.vault.config.vaultCACert = "vault-ca.pem"
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
		smCfg.vault.config.vaultCACert = "vault-ca.pem"
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
		},
		{
			name: "Will mutate container for AWS",
			fields: fields{
				k8sClient: fake.NewSimpleClientset(),
			},
			args: args{
				containers: []corev1.Container{
					{
						Name:    "AWSContainer",
						Image:   "some-image-aws",
						Command: []string{"/bin/bash"},
						Args:    []string{"-c", "echo 'ACCESS_KEY: $AWS_ACCESS_KEY'"},
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
					Name:    "AWSContainer",
					Image:   "some-image-aws",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args: []string{
						"aws",
						"--region=us-west-2",
						"--secret-name=test-aws-secret",
						"--role-arn=arn:aws:iam::user:role/secretmanger",
						"--previous-version=true",
						"--",
						"/bin/bash",
						"-c",
						"echo 'ACCESS_KEY: $AWS_ACCESS_KEY'",
					},
					Env: []corev1.EnvVar{
						{Name: "SOME_VARIABLE", Value: "non-of-your-business"},
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
						Name:    "MyGCPContainer",
						Image:   "some-gcp-image",
						Command: []string{"/bin/bash"},
						Args:    []string{"-c", "echo 'API_KEY: $API_KEY'"},
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
					Name:    "MyGCPContainer",
					Image:   "some-gcp-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args: []string{
						"gcp",
						"--project-id=project-x",
						"--secret-name=gcp-test-secret",
						"--secret-version=5",
						"--google-application-credentials=/var/run/secret/cloud.google.com/service-account.json",
						"--",
						"/bin/bash",
						"-c",
						"echo 'API_KEY: $API_KEY'",
					},
					Env: []corev1.EnvVar{
						{Name: "HOST", Value: "127.0.0.1"},
					},
					VolumeMounts: []corev1.VolumeMount{
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
						Command: []string{"/bin/bash"},
						Args:    []string{"-c", "echo 'API_KEY: $API_KEY'"},
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
					Args: []string{
						"vault",
						"--role=x-role",
						"--kubernetes-backend=/alt/kubernetes/path",
						"--token-path=/tmp/key",
						"--path=/secret/data/top-secret",
						"--names-as-keys",
						"--version=5",
						"--",
						"/bin/bash",
						"-c",
						"echo 'API_KEY: $API_KEY'",
					},
					Env: []corev1.EnvVar{
						{Name: "HOST", Value: "127.0.0.1"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/vault-ca.pem"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "vault-tls", MountPath: "/etc/tls", SubPath: "vault-ca.pem"},
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
						Command: []string{"/bin/bash"},
						Args:    []string{"-c", "echo 'API_KEY: $API_KEY'"},
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
					Args: []string{
						"vault",
						"--role=x-role",
						"--backend=gcp",
						"--google-application-credentials=/var/run/secret/cloud.google.com/service-account.json",
						"--token-path=/tmp/key",
						"--path=/secret/data/top-secret",
						"--names-as-keys",
						"--",
						"/bin/bash",
						"-c",
						"echo 'API_KEY: $API_KEY'",
					},
					Env: []corev1.EnvVar{
						{Name: "HOST", Value: "127.0.0.1"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/vault-ca.pem"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "google-cloud-key", MountPath: "/var/run/secret/cloud.google.com"},
						{Name: "vault-tls", MountPath: "/etc/tls", SubPath: "vault-ca.pem"},
						{Name: "secrets-consumer-env", MountPath: "/secrets-consumer"},
					},
				},
			},
		}, {
			name: "Will mutate container for Vault with a multiple secret-configs",
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
				secretManagerConfig: getSecretManagerConfig("vault-multi"),
			},
			mutated: true,
			wantErr: false,
			wantedContainers: []corev1.Container{
				{
					Name:    "MyContainer",
					Image:   "some-image",
					Command: []string{"/secrets-consumer/secrets-consumer-env"},
					Args: []string{
						"vault",
						"--role=x-role",
						`--secret-config={"path": "/some/secret/path-1", "version": "3", "use-secret-names-as-keys":  true}`,
						`--secret-config={"path": "/some/secret/path-2"}`,
						`--secret-config={"path": "/some/secret/path-3", "use-secret-names-as-keys":  true}`,
						"--",
						"/app",
					},
					Env: []corev1.EnvVar{
						{Name: "API_KEY", Value: "vault:API_KEY"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/vault-ca.pem"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "vault-tls", MountPath: "/etc/tls", SubPath: "vault-ca.pem"},
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
					Args: []string{
						"vault",
						"--role=x-role",
						"--token-path=/tmp/key",
						"--path=/secret/data/top-secret",
						"--names-as-keys",
						"--version=2",
						"--",
						"/app",
					},
					Env: []corev1.EnvVar{
						{Name: "API_KEY", Value: "vault:API_KEY"},
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_CACERT", Value: "/etc/tls/vault-ca.pem"},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "vault-tls", MountPath: "/etc/tls", SubPath: "vault-ca.pem"},
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
