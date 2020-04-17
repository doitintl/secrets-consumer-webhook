# Kubernetes Secrets Consumer Webhook

The following webhook is completely based on [Banazi Cloud Vault Secrets Webhook](https://github.com/banzaicloud/bank-vaults/tree/master/cmd/vault-secrets-webhook)
however, it does not support configmap or secret mutations, as well as consulTemplates like the original does.

This version was rewrite to allow AWS, GCP and VAULT secrets manager, as well as treat vault secret paths as wildcard or a directory containing multiple secrets where each secret name is the key, and a single value for it.

Another variation is allowing to get explicit secrets vs all secrets from path.

This Mutation webhook will mutate a Pod based on annotations and automatically inject secrets from various secrets managers like [AWS Secret Manager](https://aws.amazon.com/secrets-manager/), [GCP Secret Manager](https://cloud.google.com/secret-manager) or [Hashicorp Vault](https://www.vaultproject.io/) using its companion tool [secrets-consumer-env](https://github.com/doitintl/secrets-consumer-env)

Please note, this is a single secret manager setup, this tool doesn't support fetching secrets from multiple secrets managers nor it should.

## How this mutation webhook works

This mutation webhook watch for events where a new Pod is requested via the API, if the object is a Pod and it has specific annotations, the webhook will read these annotations and convert them to env vars or volumes needed for the `secrets-consumer-env` tool to run.

It will create an init container with `secrets-consumer-env` image in it, as well as an in-memory shared volume that will also be mounted for your container.

The init container will copy the binary of `secrets-consumer-env` into the shared volume.

The webhook will also change your command to be prefixed by the command `secrets-consumer-env`

## Setting up Vault Kubernetes Backend Authentication

Vault can authenticate to kubernetes using a kubernetes service account

It does this by another service account called `vault-reviewer` with an auth-delegator permissions, which allows it to pass another service account token for authentication to the kubernetes master.

Once the authentication to kubernetes is successful,
Vault returns a client token that can be used to login to Vault,
Vault will check a mapping between a vault role, service account, namespace and the policy to allows/deny the access.

This `vault-reviewer` service account token will be configured inside the vault using vault CLI.
let’s create the service account for that vault-reviewer

Create the service account for that vault-reviewer:

**Please note**; if you have set up Vault on any other namespace, make sure to update this file accordingly.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: vault-reviewer
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: role-tokenreview-binding
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: vault-reviewer
  namespace: default
```

```bash
kubectl apply -f vault-reviewer.yaml
```

Enable the Kubernetes auth backend:

Make sure you are logged in to vault using the root token

```bash
$ vault login
  Token (will be hidden):
```

```bash
$ vault auth enable kubernetes
Success! Enabled kubernetes auth method at: kubernetes/
```

Configure Vault with the vault-reviewer token and Kubernetes CA:

**Note**: if you setup vault on any other namespace set the -n <namespace> flag after each kubectl command

```bash
  VAULT_SA_TOKEN_NAME=$(kubectl get sa vault-reviewer -o jsonpath="{.secrets[*]['name']}")

  SA_JWT_TOKEN=$(kubectl get secret "$VAULT_SA_TOKEN_NAME" -o jsonpath="{.data.token}" | base64 --decode; echo)

  SA_CA_CRT=$(kubectl get secret "$VAULT_SA_TOKEN_NAME" -o jsonpath="{.data['ca\.crt']}" | base64 --decode; echo)

  $ vault write auth/kubernetes/config token_reviewer_jwt="$SA_JWT_TOKEN" kubernetes_host=https://kubernetes.default kubernetes_ca_cert="$SA_CA_CRT"

  Success! Data written to: auth/kubernetes/config
```

### example of role mapping to service account and namespace in vault

```bash
vault write auth/kubernetes/role/tester \
  bound_service_account_names=tester \
  bound_service_account_namespaces=default \
  policies=test_policy \
  ttl=1h
```

## Setting up Vault GCP Backend Authentication

### Enable the Google Cloud auth method

```bash
vault auth enable gcp
```

Configure the auth method credentials:

```bash
 vault write auth/gcp/config credentials=@/path/to/credentials.json
```

If you are using instance credentials or want to specify credentials via an environment variable, you can skip this step.

### Create a named role

#### For an iam-type role

```bash
  vault write auth/gcp/role/my-iam-role \
  type="iam" \
  policies="dev,prod" \
  bound_service_accounts="my-service@my-project.iam.gserviceaccount.com"
```

#### For a gce-type role

```bash
vault write auth/gcp/role/my-gce-role \
  type="gce" \
  policies="dev,prod" \
  bound_projects="my-project1,my-project2" \
  bound_zones="us-east1-b" \
  bound_labels="foo:bar,zip:zap" \
  bound_service_accounts="my-service@my-project.iam.gserviceaccount.com"
```

### Required GCP Permissions

#### Vault Server Permissions

For iam-type Vault roles, Vault can be given the following roles:

```console
roles/iam.serviceAccountKeyAdmin
```

For gce-type Vault roles, Vault can be given the following roles:

```console
roles/compute.viewer
```

If you instead wish to create a custom role with only the exact GCP permissions required, use the following list of permissions:

```console
iam.serviceAccounts.get
iam.serviceAccountKeys.get
compute.instances.get
compute.instanceGroups.list
```

#### Permissions For Authenticating Against Vault

Note that the previously mentioned permissions are given to the Vault servers. The IAM service account or GCE instance that is authenticating against Vault must have the following role:

```console
roles/iam.serviceAccountTokenCreator
```

## Installation

Before you install this chart you must create a namespace for it, this is due to the order in which the resources in the charts are applied (Helm collects all of the resources in a given Chart and it's dependencies, groups them by resource type, and then installs them in a predefined order (see [here](https://github.com/helm/helm/blob/release-2.10/pkg/tiller/kind_sorter.go#L29) - Helm 2.10).

**Note: this namespace is excluded from the mutations**
The `MutatingWebhookConfiguration` gets created before the actual backend Pod which serves as the webhook itself, Kubernetes would like to mutate that pod as well, but it is not ready to mutate yet (infinite recursion in logic).

The namespace must have a label of `name` with the namespace name as it's value.

set the target namespace name or skip for the default name: vswh

```bash
export WEBHOOK_NS=`<namespace>`
```

```bash
WEBHOOK_NS=${WEBHOOK_NS:-secrets-consumer-wh}
kubectl create namespace "${WEBHOOK_NS}"
kubectl label ns "${WEBHOOK_NS}" name="${WEBHOOK_NS}"
```

Use the helm chart to install the webhook:

```bash
helm upgrade --namespace ${WEBHOOK_NS} --install secrets-consumer-webhook helm-chart --wait
```

**NOTE**: `--wait` is necessary because of Helm timing issues, please see [this issue](https://github.com/banzaicloud/banzai-charts/issues/888).

### About GKE Private Clusters

When Google configure the control plane for private clusters, they automatically configure VPC peering between your Kubernetes cluster’s network in a separate Google managed project.

The auto-generated rules **only** open ports 10250 and 443 between masters and nodes. This means that in order to use the webhook component with a GKE private cluster, you must configure an additional firewall rule to allow your masters CIDR to access your webhook pod using the port 8443.

You can read more information on how to add firewall rules for the GKE control plane nodes in the [GKE docs](https://cloud.google.com/kubernetes-engine/docs/how-to/private-clusters#add_firewall_rules).

### Auto detecting container entrypoint or command

The webhook will attempt to query the metadata for the container image if no explicit command is given for `secrets-consumer-env` to work properly, If your container is on a private repo, you can set your docker repo credentials via the `imagePullSecrets` attribute of the container.

You can also specify a default secret being used by the webhook for cases where a pod has no imagePullSecrets specified. To make this work you have to set the environment variables `DEFAULT_IMAGE_PULL_SECRET` and `DEFAULT_IMAGE_PULL_SECRET_NAMESPACE` when deploying the secrets-consumer-webhook. Have a look at the values.yaml of the vault-secrets-webhook helm chart to see how this is done.

**NOTE:** If you EC2 nodes are having ECR instance role added the webhook can request an ECR access token through that role automatically, instead of an explicit imagePullSecret

## explicit vs non-explicit (get all) secrets

You have the option to select which secrets you want to expose to your process, or get all secrets

To explicitly select secrets from the secret manager, add an env var to your pod using the following convention:

```yaml
env:
- name:  <variable name to export>
  value: vault:<vault key name from secret>
```

### Annotations

#### AWS secret manager

| Name| Description | Required | Default|
| :--- |:---|:---:|:---|
|"aws.secret.manager/enabled"| enable the AWS secret manager | - | false |
|"aws.secret.manager/region" | AWS secret manager region | No | us-east-1 |
|"aws.secret.manager/role-arn" | AWS IAM Role to access the secret | No | - |
|"aws.secret.manager/secret-name" | secret name | Yes | - |
|"aws.secret.manager/previous-version" | if the secret is rotated, set to "true" | No | - |

#### GCP secret manager

| Name| Description | Required | Default|
| :--- |:---|:---:|:---|
|"gcp.secret.manager/enabled"| enable the GCP secret manager | - | false |
|"gcp.secret.manager/project-id" | GCP Project ID | Yes | - |
|"gcp.secret.manager/gcp-service-account-key-secret-name" | GCP IAM service account secret name (file name **must be** `service-account.json`) | No | Google Default Application Credentials |
|"gcp.secret.manager/secret-name" | secret name | Yes | - |
|"gcp.secret.manager/secret-version" | specify the secret version as string | No | Latest |

#### Vault secret manager

| Name| Description | Required | Default|
| :--- |:---|:---:|:---|
|"vault.security/enabled"| enable the Vault secret manager | - | false |
|"vault.security/vault-addr" | Vault cluster service address | Yes | - |
|"vault.security/vault-path" | Vault secret path  | Yes | - |
|"vault.security/vault-secret-version" | Vault secret version (if using v2 secret engine)  | Yes | - |
|"vault.security/vault-use-secret-names-as-keys" | treat secret path ending with `/` as directory where secret name is the key and a single value in each  | No | - |
|"vault.security/vault-role" | Vault role to access the secret path  | Yes | - |
|"vault.security/vault-tls-secret-name" | Vault TLS secret name  | No | Latest |
|"vault.security/k8s-token-path" | alternate kubernetes service account token path  | No | `/var/run/secrets/kubernetes.io/serviceaccount/token` |

Vault can be used with 2 backend authentications (GCP / Kubernetes)

##### Kubernetes backend authentication

Default authentication method

| Name| Description | Required | Default|
| :--- |:---|:---:|:---|
|"vault.security/k8s-token-path" | alternate kubernetes service account token path  | No | `/var/run/secrets/kubernetes.io/serviceaccount/token` |

##### GCP Backend authentication

Use GCP service account to authenticate to Vault

| Name| Description | Required | Default|
| :--- |:---|:---:|:---|
|"vault.security/gcp-service-account-key-secret-name" | GCP IAM service account secret name (file name **must be** `service-account.json`) to login with gcp  | No | Latest |
|"vault.security/vault-tls-secret-name" | Vault TLS secret name  | No | Latest |
