package main

import (
	"github.com/innovia/secrets-consumer-webhook/registry"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

type secretManagerConfig struct {
	aws
	gcp
	vault
	explicitSecrets bool // only get secrets that match the prefix `secret:`
}

// MutatingWebhook holds k8s client interface
type mutatingWebhook struct {
	k8sClient kubernetes.Interface
	registry  registry.ImageRegistry
	logger    log.FieldLogger
}
