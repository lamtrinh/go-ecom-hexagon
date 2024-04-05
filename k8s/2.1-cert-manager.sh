helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \ --create-namespace \
  --version v1.14.4 \
  --set installCRDs=true
