#!/bin/bash

set -e

# Configuration
NAMESPACE="preview"
SERVICE_ACCOUNT_NAME="github-preview-deployer"
CLUSTER_NAME=$(kubectl config current-context)
SECRET_NAME=""
GITHUB_REPO="edkadigital/startmeup"

# Create namespace if it doesn't exist
kubectl get namespace $NAMESPACE >/dev/null 2>&1 || kubectl create namespace $NAMESPACE

echo "Creating service account in namespace $NAMESPACE..."
kubectl create serviceaccount $SERVICE_ACCOUNT_NAME -n $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

# Create Role with full admin permissions in preview namespace
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: $NAMESPACE
  name: preview-admin
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
EOF

# Bind the role to the service account
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: github-preview-admin-binding
  namespace: $NAMESPACE
subjects:
- kind: ServiceAccount
  name: $SERVICE_ACCOUNT_NAME
  namespace: $NAMESPACE
roleRef:
  kind: Role
  name: preview-admin
  apiGroup: rbac.authorization.k8s.io
EOF

# Create an additional ClusterRole for namespace operations (needed for --create-namespace flag in Helm)
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: preview-namespace-admin
rules:
- apiGroups: [""]
  resources: ["namespaces"]
  verbs: ["get", "list", "watch", "create", "update", "patch"]
- apiGroups: ["apiextensions.k8s.io"]
  resources: ["customresourcedefinitions"]
  verbs: ["get", "list"]
EOF

# Create ClusterRoleBinding
cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: github-preview-namespace-binding
subjects:
- kind: ServiceAccount
  name: $SERVICE_ACCOUNT_NAME
  namespace: $NAMESPACE
roleRef:
  kind: ClusterRole
  name: preview-namespace-admin
  apiGroup: rbac.authorization.k8s.io
EOF

# Get the service account token secret name
SECRET_NAME=$(kubectl get serviceaccount $SERVICE_ACCOUNT_NAME -n $NAMESPACE -o jsonpath='{.secrets[0].name}')

if [ -z "$SECRET_NAME" ]; then
  # For Kubernetes 1.24+, create a token manually
  echo "Creating token for service account..."
  SECRET_NAME="${SERVICE_ACCOUNT_NAME}-token"
  kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: ${SECRET_NAME}
  namespace: ${NAMESPACE}
  annotations:
    kubernetes.io/service-account.name: ${SERVICE_ACCOUNT_NAME}
type: kubernetes.io/service-account-token
EOF
fi

# Wait for the secret to be properly initialized
echo "Waiting for secret to be initialized..."
sleep 3

# Get service account details
TOKEN=$(kubectl get secret $SECRET_NAME -n $NAMESPACE -o jsonpath='{.data.token}' | base64 --decode)
CA_CERT=$(kubectl get secret $SECRET_NAME -n $NAMESPACE -o jsonpath='{.data.ca\.crt}' | base64 --decode)
SERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
CA_CERT_FILE="${TEMP_DIR}/ca.crt"
KUBE_CONFIG_FILE="${TEMP_DIR}/preview-kubeconfig.yaml"

# Write the CA certificate to a file
echo "$CA_CERT" > "$CA_CERT_FILE"

# Generate the kubeconfig
cat > "$KUBE_CONFIG_FILE" <<EOF
apiVersion: v1
kind: Config
preferences: {}

clusters:
- cluster:
    certificate-authority-data: $(cat "$CA_CERT_FILE" | base64 | tr -d '\n')
    server: $SERVER
  name: $CLUSTER_NAME

contexts:
- context:
    cluster: $CLUSTER_NAME
    namespace: $NAMESPACE
    user: $SERVICE_ACCOUNT_NAME
  name: $SERVICE_ACCOUNT_NAME@$CLUSTER_NAME

current-context: $SERVICE_ACCOUNT_NAME@$CLUSTER_NAME

users:
- name: $SERVICE_ACCOUNT_NAME
  user:
    token: $TOKEN
EOF

echo "======================================================"
echo "Kubeconfig file created at: $KUBE_CONFIG_FILE"
echo "Content of the kubeconfig (to copy to GitHub Secret):"
echo "======================================================"
cat "$KUBE_CONFIG_FILE"
echo "======================================================"
echo "Instructions:"
echo "1. Copy the above kubeconfig content"
echo "2. Go to your GitHub repository: https://github.com/$GITHUB_REPO"
echo "3. Navigate to Settings > Secrets > Actions"
echo "4. Create a new repository secret named 'KUBE_CONFIG'"
echo "5. Paste the content and save"

echo "Removing temporary directory..."
rm -rf "$TEMP_DIR"

echo "Setup complete!" 