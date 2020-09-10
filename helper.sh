#!/bin/sh

subcommand=$1

case $subcommand in
  pre)
    echo "Checking prerequisites..."

    # Check for docker
    if command -v docker &> /dev/null
    then
      echo "✅ Docker"
    else
      echo "❌ Docker not found. Make sure it is installed. https://www.docker.com/get-started"
    fi

    # Check for kind
    if command -v kind &> /dev/null
    then
      echo "✅ Kind"
    else
      echo "❌ Kind not found. Make sure it is installed. https://kind.sigs.k8s.io/docs/user/quick-start/"
    fi
    
    # Check for kubectl
    if command -v kubectl &> /dev/null
    then
      echo "✅ Kubectl"
    else
      echo "❌ Kubectl not found. Make sure it is installed. https://kubernetes.io/docs/tasks/tools/install-kubectl/"
    fi
    ;;
  setup)
    echo "Setting up environment...\n"
    echo "Creating kind cluster..."
    kind create cluster --config=./tools/kind-config.yaml
    echo "✅ Cluster creation complete\n"

    echo "Installing ingress controller..."
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.35.0/deploy/static/provider/kind/deploy.yaml
    kubectl wait --namespace ingress-nginx \
    --for=condition=ready pod \
    --selector=app.kubernetes.io/component=controller \
    --timeout=90s
    echo "✅ Ingress controller installed\n"

    echo "Creating secrets..."
    clientId=$(awk '/AUTH_CLIENT_ID/{split($0, arr, "="); print arr[2]}' .env)
    echo $clientId
    clientSecret=$(awk '/AUTH_CLIENT_SECRET/{split($0, arr, "="); print arr[2]}' .env)
    echo $clientSecret
    kubectl create secret generic auth --from-literal=clientID=$clientId --from-literal=clientSecret=$clientSecret
    echo "✅ Secrets created\n"

    echo "Building docker images..."
    docker-compose build
    echo "✅ Docker images built\n"
    
    echo "Loading docker images into cluster..."
    kind load docker-image auth:dev
    kind load docker-image settings:dev
    echo "✅ Docker images loaded\n"

    echo "Patching validation webhook..."
    # The ingress admission validator needs to be updated with the correct service path. I suspect this has to do
    # with k8s 1.16 deprecations
    kubectl patch validatingwebhookconfigurations.admissionregistration.k8s.io ingress-nginx-admission -p '{"webhooks": [{"name": "validate.nginx.ingress.kubernetes.io", "clientConfig": {"service": {"path": "/networking.k8s.io/v1beta1/ingresses"}}}]}'
    echo "✅ Validating webhook patched\n"
    ;;
  build)
    echo "Building docker images..."
    docker-compose build
    echo "✅ Docker images built\n"
    ;;
  load)
    echo "Loading docker images into cluster..."
    kind load docker-image auth:dev
    kind load docker-image settings:dev
    echo "✅ Docker images loaded\n"
    ;;
  secrets)
    echo "Creating secrets..."
    clientId=$(awk '/AUTH_CLIENT_ID/{split($0, arr, "="); print arr[2]}' .env)
    echo $clientId
    clientSecret=$(awk '/AUTH_CLIENT_SECRET/{split($0, arr, "="); print arr[2]}' .env)
    echo $clientSecret
    kubectl create secret generic auth --from-literal=clientID=$clientId --from-literal=clientSecret=$clientSecret
    echo "✅ Secrets created\n"
    ;;
  apply)
    echo "Applying k8s specifications..."
    kubectl apply -f ./k8s/deployments
    kubectl apply -f ./k8s/services
    kubectl apply -f ./k8s/ingress.yaml
    echo "✅ K8s specifications applied\n"
    ;;
  *)
    echo "You need to enter a valid sub command. Valid sub commands are:"
    echo "pre\nsetup\nbuild\nload\nsecrets\napply"
    ;;
esac
