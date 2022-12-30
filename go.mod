module github.com/innovia/secrets-consumer-webhook

go 1.13

require (
	github.com/aws/aws-sdk-go v1.33.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.7.3-0.20190327010347-be7ac8be2ae0
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/google/go-cmp v0.4.0
	github.com/heroku/docker-registry-client v0.0.0-20190909225348-afc9e1acc3d5
	github.com/opencontainers/image-spec v1.0.1
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.4.1
	github.com/sirupsen/logrus v1.4.2
	github.com/slok/kubewebhook v0.8.0
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	k8s.io/api v0.17.4-beta.0
	k8s.io/apimachinery v0.17.4-beta.0
	k8s.io/client-go v0.17.4-beta.0
	sigs.k8s.io/controller-runtime v0.5.0
)
