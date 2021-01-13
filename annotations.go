package main

const (

	// AnnotationAWSSecretManagerEnabled if enabled it will use AWS secret manager
	AnnotationAWSSecretManagerEnabled = "aws.secret.manager/enabled"

	// AnnotationAWSSecretManagerRegion the region for which the secret manager is set
	AnnotationAWSSecretManagerRegion = "aws.secret.manager/region"

	// AnnotationAWSSecretManagerRoleARN if specified it will assume the role for fetching the secret
	AnnotationAWSSecretManagerRoleARN = "aws.secret.manager/role-arn"

	// AnnotationAWSSecretManagerSecretName aws secret manager secret name to fetch
	AnnotationAWSSecretManagerSecretName = "aws.secret.manager/secret-name"

	// AnnotationAWSSecretManagerPreviousVersion when used will retrive the previous version for the secret
	// note that AWS only supports single previous vresion
	AnnotationAWSSecretManagerPreviousVersion = "aws.secret.manager/previous-version"

	// AnnotationGCPSecretManagerEnabled if enabled use GCP secret manager
	AnnotationGCPSecretManagerEnabled = "gcp.secret.manager/enabled"

	// AnnotationGCPSecretManagerProjectID the gcp project id to use for the secret manager
	AnnotationGCPSecretManagerProjectID = "gcp.secret.manager/project-id"

	// AnnotationGCPSecretManagerSecretName the name of the GCP secret
	AnnotationGCPSecretManagerSecretName = "gcp.secret.manager/secret-name"

	// AnnotationGCPSecretManagerSecretVersion the version number for the secret
	AnnotationGCPSecretManagerSecretVersion = "gcp.secret.manager/secret-version"

	// AnnotationGCPSecretManagerGCPServiceAccountKeySecretName is the secret name where the GCP service account credentials
	// are stored and has teh permissions to access the secret
	AnnotationGCPSecretManagerGCPServiceAccountKeySecretName = "gcp.secret.manager/gcp-service-account-key-secret-name"

	// AnnotationVaultEnabled if enabled use vault as the secret manager
	AnnotationVaultEnabled = "vault.secret.manager/enabled"

	// AnnotationVaultService vault address in the http/https format including the port number
	// for example https://vault.vault.svc:8200
	AnnotationVaultService = "vault.secret.manager/service"

	// AnnotationVaultAuthPath specifies the mount path to be used for the Kubernetes auto-auth method.
	AnnotationVaultAuthPath = "vault.secret.manager/auth-path"

	// AnnotationVaultSecretPath the secret path in vault - will auto detect if kv2 is used and auto-append `data` to it
	AnnotationVaultSecretPath = "vault.secret.manager/path"

	// AnnotationVaultRole specifies the role to be used for the Kubernetes auto-auth method.
	AnnotationVaultRole = "vault.secret.manager/role"

	// AnnotationVaultGCPServiceAccountKeySecretName The secret name that holds the GCP service account credentials
	AnnotationVaultGCPServiceAccountKeySecretName = "vault.secret.manager/gcp-service-account-key-secret-name"

	// AnnotationVaultTLSSecret is the name of the Kubernetes secret containing
	// client TLS certificates and keys.
	AnnotationVaultTLSSecret = "vault.secret.manager/tls-secret"

	// AnnotationVaultCACert is the filename of the CA certificate used to verify Vault's
	// CA certificate.
	AnnotationVaultCACert = "vault.secret.manager/ca-cert"

	// AnnotationVaultK8sTokenPath override the token that will be used for vault authentication
	AnnotationVaultK8sTokenPath = "vault.secret.manager/k8s-token-path"

	// AnnotationVaultUseSecretNamesAsKeys is used with a path that has a tree under it,
	// will be using the secret names as the keys and ignore the real key in the secret itself
	AnnotationVaultUseSecretNamesAsKeys = "vault.secret.manager/use-secret-names-as-keys"

	// AnnotationVaultSecretVersion get the specified secret version, default to latest version
	AnnotationVaultSecretVersion = "vault.secret.manager/secret-version"

	// AnnotationVaultMultiSecretPrefix allow multi secret by order
	// vault.secret.manager/secret-config-1: '{"Path": "secrets/v2/plain/secrets/path/app", "Version": "2", "use-secret-names-as-keys": "true"}'
	AnnotationVaultMultiSecretPrefix = "vault.secret.manager/secret-config-"
)
