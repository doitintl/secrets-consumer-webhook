package main

import (
	corev1 "k8s.io/api/core/v1"
)

type aws struct {
	config struct {
		enabled         bool
		region          string
		secretName      string
		previousVersion string
		roleARN         string
	}
}

func (aws *aws) mutateContainer(container corev1.Container) corev1.Container {
	envVars := aws.setEnvVars()
	container.Env = append(container.Env, envVars...)
	return container
}

func (aws *aws) setEnvVars() []corev1.EnvVar {
	var envVars []corev1.EnvVar
	envVars = append(envVars, []corev1.EnvVar{
		{
			Name:  "SECRET_MANAGER",
			Value: "aws",
		},
		{
			Name:  "SECRET_NAME",
			Value: aws.config.secretName,
		}, {
			Name:  "REGION",
			Value: aws.config.region,
		}, {
			Name:  "ROLE_ARN",
			Value: aws.config.roleARN,
		}, {
			Name:  "PREVIOUS_VERSION",
			Value: aws.config.previousVersion,
		},
	}...)

	return envVars
}
