package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type vault struct {
	config struct {
		enabled                        bool
		addr                           string
		tlsSecretName                  string
		vaultCACert                    string
		path                           string
		role                           string
		tokenPath                      string
		authPath                       string
		backend                        string
		kubernetesBackend              string
		useSecretNamesAsKeys           bool
		gcpServiceAccountKeySecretName string
		version                        string
		secretConfigs                  []string
	}
}

func (vault *vault) mutateContainer(container corev1.Container) corev1.Container {
	envVars := vault.setEnvVars()
	container.Env = append(container.Env, envVars...)

	// Mount google service account key if given
	if vault.config.gcpServiceAccountKeySecretName != "" {
		container.VolumeMounts = append(container.VolumeMounts, []corev1.VolumeMount{
			{
				Name:      VolumeMountGoogleCloudKeyName,
				MountPath: VolumeMountGoogleCloudKeyPath,
			},
		}...)
	}

	if vault.config.tlsSecretName != "" {
		mountPath := fmt.Sprintf("%s/%s", VaultTLSMountPath, vault.config.vaultCACert)
		volumeName := VaultTLSVolumeName

		container.Env = append(container.Env, []corev1.EnvVar{
			{
				Name:  "VAULT_CACERT",
				Value: mountPath,
			},
		}...)
		container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
			Name:      volumeName,
			MountPath: VaultTLSMountPath,
			SubPath:   vault.config.vaultCACert,
		})
	} else {
		container.Env = append(container.Env, []corev1.EnvVar{
			{
				Name:  "VAULT_SKIP_VERIFY",
				Value: "true",
			},
		}...)
	}

	container = vault.setArgs(container)
	return container
}

func (vault *vault) setArgs(c corev1.Container) corev1.Container {
	args := []string{"vault"}
	args = append(args, fmt.Sprintf("--role=%s", vault.config.role))

	if vault.config.backend == "gcp" {
		args = append(args, "--backend=gcp")
		if vault.config.gcpServiceAccountKeySecretName != "" {
			args = append(args, fmt.Sprintf("--google-application-credentials=%s/%s", VolumeMountGoogleCloudKeyPath, GCPServiceAccountCredentialsFileName))
		}
	}

	if vault.config.kubernetesBackend != "" {
		args = append(args, fmt.Sprintf("--kubernetes-backend=%s", vault.config.kubernetesBackend))
	}

	if vault.config.tokenPath != "" {
		args = append(args, fmt.Sprintf("--token-path=%s", vault.config.tokenPath))
	}

	for _, s := range vault.config.secretConfigs {
		args = append(args, fmt.Sprintf("--secret-config=%s", s))
	}

	if vault.config.path != "" {
		args = append(args, fmt.Sprintf("--path=%s", vault.config.path))
	}

	if vault.config.useSecretNamesAsKeys {
		args = append(args, "--names-as-keys")
	}

	if vault.config.version != "" {
		args = append(args, fmt.Sprintf("--version=%s", vault.config.version))
	}

	args = append(args, "--")
	// args = append(args, fmt.Sprintf("%s", strings.Join(c.Args, " ")))
	args = append(args, c.Args...)
	c.Args = args
	return c
}

func (vault *vault) setEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar
	envVars = append(envVars, []corev1.EnvVar{
		{
			Name:  "VAULT_ADDR",
			Value: vault.config.addr,
		},
	}...)

	return envVars
}
