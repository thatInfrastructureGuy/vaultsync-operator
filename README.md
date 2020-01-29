# VaultSync Operator
Periodically syncs secrets from various Vaults to Kubernetes Secrets. 

### Description
This project aims to simplify secret management. The idea is:
1. Store secrets in any of the industry standard vaults such as `Azure KeyVault`, `AWS Secrets Manager`, `GCP Secrets Manager` or `Hashicorp Vault`. 
2. These vaults are your _source of truth_.
3. Whenever secrets change in Vaults your application gets updated automatically with the newer values.

### Currently Supported Providers:
- [x] AWS Secrets Manager
- [x] Azure KeyVault Secrets
- [ ] Google Secrets Manager  : Waiting for Stable APIs
- [ ] Hashicorp Vault  : Non-requirement for initial release.

### Currently Supported Consumers:
- [x] Kubernetes Secrets

### Quick Start

1. Deploy the Operator

```
kubectl apply -f deploy/namespace.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/secret.yaml
kubectl apply -f deploy/crds/operator.thatinfrastructureguy.com_vaultsyncers_crd.yaml
kubectl apply -f deploy/operator.yaml
```

2. Set Your Cloud Credentials

> **Note**: Make sure your credentials have proper authorization to access azure keyvault / aws secrets manager.


_Azure:_
```
kubectl -n vaultsync create secret generic azure-credentials \
--from-literal AZURE_TENANT_ID=xxxxxxxxxxxxxx \
--from-literal AZURE_CLIENT_ID=xxxxxxxxxxxxxx \
--from-literal AZURE_CLIENT_SECRET=xxxxxxxxxxxxxx \
--dry-run -o yaml | kubectl -n vaultsync apply -f -
```

_AWS:_
```
kubectl -n vaultsync create secret generic aws-credentials \
--from-literal AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxxxxx \
--from-literal AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxx \
--from-literal AWS_DEFAULT_REGION=xxxxxxxxxxxxxxxxxxx \
--from-literal AWS_REGION=xxxxxxxxxxxxxxxxxxx \
--dry-run -o yaml | kubectl -n vaultsync apply -f -
```

3. Create the Custom Resource

_Azure:_ 
```
apiVersion: operator.thatinfrastructureguy.com/v1alpha1
kind: VaultSyncer
metadata:
  name: azure-vaultsyncer
  namespace: vaultsync
spec:
  provider: "azure"
  providerCredsSecret: "azure-credentials"
  vaultName: "myKeyVault"
  deploymentList: ""
```

_AWS:_ 
```
apiVersion: operator.thatinfrastructureguy.com/v1alpha1
kind: VaultSyncer
metadata:
  name: aws-vaultsyncer
  namespace: vaultsync
spec:
  provider: "aws"
  providerCredsSecret: "aws-credentials"
  vaultName: "mysecretsmanager"
  deploymentList: ""
```

### Custom Resource Values

Parameter | Description | Default
--- | --- | ---
`provider` | Cloud Provider currently supported `azure` and `aws` | `nil`
`providerCredsSecret` | Secret in `vaultsync` namespace where authn/authz credentials are stored. By default it points to `provider-credentials`. Create an empty secret if you are authorizing via IAM policies and do not need credentials.  | `provider-credentials`
`vaultName` | Azure KeyVault / AWS Secrets Manager name from where secrets will be pulled. | `nil`
`consumer` | This is defaulted to `kubernetes` secrets. In future other consumers may be supported. Eg: `jenkins`, `VM` | `kubernetes`
`secretName` | The name of the secret to be created/updated whenever the secrets are pulled from vault. If empty, name of secret is kept the same as name of the vault. | `vaultName`
`secretNamespace` | Namespace where the secret should be created. If empty, secret is created in `default` namespace. | `default` 
`deploymentList` | Comma seperate names of deployments which should be redeployed once the secret is updated. This is done in order for deployments to capture the newly updated secrets. Helpful when kubernetes secrets are mounted as environment variables or volumes. | `nil`
`statefulsetList` | Comma seperate names of statefulsets which should be redeployed once the secret is updated. This is done in order for statefulsets to capture the newly updated secrets. Helpful when kubernetes secrets are mounted as environment variables or volumes. | `nil` 
`refreshRate` | Determines how often check for updated secrets is done. Defaults to 60 seconds cycle. | `60` 
`convertHyphensToUnderscores` | Azure Keyvault does not support `_` in the key name. However, environment variables usually contain `_` eg: `AZURE_CLIENT_ID`. If set to `true`, you can store keys in Azure KeyVault as as `AZURE-CLIENT-ID` and they will be inserted in kubernetes secret as `AZURE_CLIENT_ID`.  | `false`
