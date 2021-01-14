module github.com/innovia/secrets-consumer-webhook

go 1.15

require (
	github.com/aws/aws-sdk-go v1.36.26
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/google/go-cmp v0.5.2
	github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5
	github.com/opencontainers/image-spec v1.0.1
	github.com/patrickmn/go-cache v1.0.0
	github.com/prometheus/client_golang v1.9.0
	github.com/sirupsen/logrus v1.7.0
	github.com/slok/kubewebhook v0.11.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	sigs.k8s.io/controller-runtime v0.7.0
)
