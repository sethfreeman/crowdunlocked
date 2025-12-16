#!/bin/bash
set -e

export AWS_PAGER=""

echo "Waiting for EKS cluster to be ready..."
while true; do
  STATUS=$(AWS_PROFILE=crowdunlocked-dev aws eks describe-cluster \
    --name crowdunlocked-dev \
    --region us-west-2 \
    --query 'cluster.status' \
    --output text 2>/dev/null || echo "NOT_FOUND")
  
  if [ "$STATUS" = "ACTIVE" ]; then
    echo "✓ Cluster is ACTIVE"
    break
  elif [ "$STATUS" = "NOT_FOUND" ]; then
    echo "  Cluster not found yet, waiting..."
  else
    echo "  Cluster status: $STATUS"
  fi
  sleep 30
done

echo ""
echo "Updating kubeconfig..."
AWS_PROFILE=crowdunlocked-dev aws eks update-kubeconfig \
  --name crowdunlocked-dev \
  --region us-west-2 \
  --alias crowdunlocked-dev

echo ""
echo "Creating EKS access entry for OrganizationAccountAccessRole..."
AWS_PROFILE=crowdunlocked-dev aws eks create-access-entry \
  --cluster-name crowdunlocked-dev \
  --principal-arn arn:aws:iam::179151668767:role/OrganizationAccountAccessRole \
  --type STANDARD \
  --region us-west-2 2>/dev/null || echo "Access entry already exists"

echo ""
echo "Associating cluster admin policy..."
AWS_PROFILE=crowdunlocked-dev aws eks associate-access-policy-to-principal \
  --cluster-name crowdunlocked-dev \
  --principal-arn arn:aws:iam::179151668767:role/OrganizationAccountAccessRole \
  --policy-arn arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy \
  --access-scope type=cluster \
  --region us-west-2 2>/dev/null || echo "Policy already associated"

echo ""
echo "Checking for nodes..."
kubectl get nodes --context crowdunlocked-dev || echo "No nodes yet (this is expected with EKS Auto Mode until pods are scheduled)"

echo ""
echo "Bootstrapping Flux..."
flux bootstrap github \
  --context=crowdunlocked-dev \
  --owner=sethfreeman \
  --repository=crowdunlocked \
  --branch=develop \
  --path=flux/clusters/dev \
  --personal

echo ""
echo "✓ Flux bootstrap complete!"
echo ""
echo "Checking Flux status..."
flux get all --context crowdunlocked-dev

echo ""
echo "Watching for pods..."
kubectl get pods -A --context crowdunlocked-dev

echo ""
echo "Setup complete! Flux will now sync your applications from the git repository."
echo "Monitor with: kubectl get pods -A --context crowdunlocked-dev -w"
