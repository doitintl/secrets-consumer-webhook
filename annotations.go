package main

const (

	// AnnotationAWSSecretManagerEnabled if enabled it will use AWS secret manager
	AnnotationAWSSecretManagerEnabled = "aws.secret.manager/enabled"

	// AnnotaionAWSSecretManagerRegion the region for which the secret manager is set
	AnnotaionAWSSecretManagerRegion = "aws.secret.manager/region"

	// AnnotaionAWSSecretManagerRoleARN if specified it will assume the role for fetching the secret
	AnnotaionAWSSecretManagerRoleARN = "aws.secret.manager/role-arn"

	// AnnotaionAWSSecretManagerSecretName aws secret manager secret name to fetch
	AnnotaionAWSSecretManagerSecretName = "aws.secret.manager/secret-name"

	// AnnotaionAWSSecretManagerPreviousVersion when used will retrive the previous version for the secret
	// note that AWS only supports single previous vresion
	AnnotaionAWSSecretManagerPreviousVersion = "aws.secret.manager/previous-version"

	// AnnotaionGCPSecretManagerEnabled if enabled use GCP secret manager
	AnnotaionGCPSecretManagerEnabled = "gcp.secret.manager/enabled"

	// AnnotaionGCPSecretManagerProjectID the gcp project id to use for the secret manager
	AnnotaionGCPSecretManagerProjectID = "gcp.secret.manager/project-id"

	// AnnotaionGCPSecretManagerSecretName the name of the GCP secret
	AnnotaionGCPSecretManagerSecretName = "gcp.secret.manager/secret-name"

	// AnnotaionGCPSecretManagerSecretVersion the version number for the secret
	AnnotaionGCPSecretManagerSecretVersion = "gcp.secret.manager/secret-version"

	// AnnotaionGCPSecretManagerGCPServiceAccountKeySecretName is the secret name where the GCP service account credentials
	// are stored and has teh permissions to access the secret
	AnnotaionGCPSecretManagerGCPServiceAccountKeySecretName = "gcp.secret.manager/gcp-service-account-key-secret-name"

	// AnnotaionVaultEnabled if enabled use vault as the secret manager
	AnnotaionVaultEnabled = "vault.secret.manager/enabled"

	// AnnotationVaultService vault address in the http/https format including the port number
	// for example https://vault.vault.svc:8200
	AnnotationVaultService = "vault.secret.manager/service"

	// AnnotaionVaultAuthPath specifies the mount path to be used for the Kubernetes auto-auth method.
	AnnotaionVaultAuthPath = "vault.secret.manager/auth-path"

	// AnnotaionVaultSecretPath the secret path in vault - will auto detect if kv2 is used and auto-append `data` to it
	AnnotaionVaultSecretPath = "vault.secret.manager/path"

	// AnnotationVaultRole specifies the role to be used for the Kubernetes auto-auth method.
	AnnotationVaultRole = "vault.secret.manager/role"

	// AnnotaionVaultGCPServiceAccountKeySecretName The secret name that holds the GCP service account credentials
	AnnotaionVaultGCPServiceAccountKeySecretName = "vault.secret.manager/gcp-service-account-key-secret-name"

	// AnnotationVaultTLSSecret is the name of the Kubernetes secret containing
	// client TLS certificates and keys.
	AnnotationVaultTLSSecret = "vault.secret.manager/tls-secret"

	// AnnotationVaultCACert is the filename of the CA certificate used to verify Vault's
	// CA certificate.
	AnnotationVaultCACert = "vault.secret.manager/ca-cert"

	// AnnotaionVaultK8sTokenPath override the token that will be used for vault authentication
	AnnotaionVaultK8sTokenPath = "vault.secret.manager/k8s-token-path"

	// AnnotaionVaultUseSecretNamesAsKeys is used with a path that has a tree under it,
	// will be using the secret names as the keys and ignore the real key in the secret itself
	AnnotaionVaultUseSecretNamesAsKeys = "vault.secret.manager/use-secret-names-as-keys"

	// AnnotaionVaultSecretVersion get the specified secret version, default to latest version
	AnnotaionVaultSecretVersion = "vault.secret.manager/secret-version"

	// AnnotationVaultMultiSecretPrefix allow multi secret by order
	// vault.secret.manager/secret-config-1: '{"Path": "secrets/v2/plain/secrets/path/app", "Version": "2", "use-secret-names-as-keys": "true"}'
	AnnotationVaultMultiSecretPrefix = "vault.secret.manager/secret-config-"
)
