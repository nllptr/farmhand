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

    # Check for skaffold
    if command -v skaffold &> /dev/null
    then
      echo "✅ Skaffold"
    else
      echo "❌ Skaffold not found. Make sure it is installed. https://skaffold.dev/docs/install/"
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
    ;;

  *)
    echo "You need to enter a valid sub command. Valid sub commands are:"
    echo "pre\nsetup"
    ;;
esac
