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

    # Check for kustomize
    if command -v kustomize &> /dev/null
    then
      echo "✅ Kustomize"
    else
      echo "❌ Kustomize not found. Make sure it is installed. https://kubernetes-sigs.github.io/kustomize/installation/"
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

    echo "Patching validation webhook..."
    # The ingress admission validator needs to be updated with the correct service path. I suspect this has to do
    # with k8s 1.16 deprecations
    kubectl patch validatingwebhookconfigurations.admissionregistration.k8s.io ingress-nginx-admission -p '{"webhooks": [{"name": "validate.nginx.ingress.kubernetes.io", "clientConfig": {"service": {"path": "/networking.k8s.io/v1beta1/ingresses"}}}]}'
    echo "✅ Validating webhook patched\n"

    # echo "Creating secrets..."
    # clientId=$(awk '/AUTH_CLIENT_ID/{split($0, arr, "="); print arr[2]}' .env)
    # echo $clientId
    # clientSecret=$(awk '/AUTH_CLIENT_SECRET/{split($0, arr, "="); print arr[2]}' .env)
    # echo $clientSecret
    # kubectl create secret generic auth --from-literal=clientID=$clientId --from-literal=clientSecret=$clientSecret
    # echo "✅ Secrets created\n"

    echo "Building docker images..."
    docker-compose build
    echo "✅ Docker images built\n"
    
    echo "Loading docker images into cluster..."
    kind load docker-image auth:dev
    kind load docker-image settings:dev
    echo "✅ Docker images loaded\n"

    echo "🎉 Everything is done, but you probably need to wait a minute or two, before the admission controller patch takes effect. 🎉"
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
    if kubectl get namespaces farmhand-dev &> /dev/null
    then
      echo "Namespace 'farmhand-dev' already exists"
    else
      kubectl create namespace farmhand-dev
    fi
    kubectl apply -k ./k8s/dev
    echo "✅ K8s specifications applied\n"
    ;;
  delete)
    echo "Deleting k8s resources..."
    if kubectl get namespaces farmhand-dev &> /dev/null
    then
      kubectl delete all --all -n farmhand-dev
      echo "✅ K8s resources deleted\n"
    else
      echo "K8s namespace 'farmhand-dev' did not exist"
    fi
    ;;
  *)
    echo "You need to enter a valid sub command. Valid sub commands are:"
    echo "pre\nsetup\nbuild\nload\nsecrets\napply"
    ;;
esac
