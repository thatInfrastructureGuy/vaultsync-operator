# vaultsync-operator

##### Steps
kubectl apply -f deploy/namespace.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/crds/...
kubectl apply -f deploy/operator.yaml

kubectl -n vaultsync create secret generic my-provider-creds \
--from-literal \
