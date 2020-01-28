# VaultSync Operator

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

```
### AZURE KEYVAULT
kubectl -n vaultsync create secret generic azure-credentials \
--from-literal AZURE_TENANT_ID=xxxxxxxxxxxxxx \
--from-literal AZURE_CLIENT_ID=xxxxxxxxxxxxxx \
--from-literal AZURE_CLIENT_SECRET=xxxxxxxxxxxxxx \
--dry-run -o yaml | kubectl -n vaultsync apply -f -

### AWS SECRETS MANAGER
kubectl -n vaultsync create secret generic aws-credentials \
--from-literal AWS_ACCESS_KEY_ID=xxxxxxxxxxxxxxxxxxx \
--from-literal AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxx \
--from-literal AWS_DEFAULT_REGION=xxxxxxxxxxxxxxxxxxx \
--from-literal AWS_REGION=xxxxxxxxxxxxxxxxxxx \
--dry-run -o yaml | kubectl -n vaultsync apply -f -
```

3. Create the Custom Resource

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
```

### Custom Resource Values

Parameter | Description | Default
--- | --- | ---
`provider` |  | `nil`
`providerCredsSecret` |  | `provider-credentials`
`vaultName` |  | `nil`
`consumer` |  | `kubernetes`
`secretName` |  | `vaultName`
`secretNamespace` |  | `default` 
`deploymentList` |  | `nil`
`statefulsetList` |  | `nil` 
`refreshRate` |  | `60` 
`convertHyphensToUnderscores` |  | `false`
