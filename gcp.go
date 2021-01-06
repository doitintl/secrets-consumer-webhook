package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type gcp struct {
	config struct {
		enabled                     bool
		projectID                   string
		secretName                  string
		secretVersion               string
		serviceAccountKeySecretName string
	}
}

func (gcp *gcp) mutateContainer(container corev1.Container) corev1.Container {
	envVars := gcp.setEnvVars()
	container.Env = append(container.Env, envVars...)
	// Mount google service account key if given
	if gcp.config.serviceAccountKeySecretName != "" {
		container.VolumeMounts = append(container.VolumeMounts, []corev1.VolumeMount{
			{
				Name:      VolumeMountGoogleCloudKeyName,
				MountPath: VolumeMountGoogleCloudKeyPath,
			},
		}...)
	}
	return container
}

func (gcp *gcp) setEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar
	envVars = append(envVars, []corev1.EnvVar{
		{
			Name:  "SECRET_MANAGER",
			Value: "gcp",
		},
		{
			Name:  "SECRET_NAME",
			Value: gcp.config.secretName,
		}, {
			Name:  "PROJECT_ID",
			Value: gcp.config.projectID,
		}, {
			Name:  "SECRET_VERSION",
			Value: gcp.config.secretVersion,
		},
	}...)

	if gcp.config.secretName != "" {
		envVars = append(envVars, []corev1.EnvVar{
			{
				Name:  "GOOGLE_APPLICATION_CREDENTIALS",
				Value: fmt.Sprintf("%s/%s", VolumeMountGoogleCloudKeyPath, GCPServiceAccountCredentialsFileName),
			},
		}...)
	}

	return envVars
}
