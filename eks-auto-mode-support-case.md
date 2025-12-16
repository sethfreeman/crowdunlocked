# AWS Support Case: EKS Auto Mode Not Provisioning Nodes

## Issue Summary
EKS Auto Mode is not provisioning compute nodes despite having pending pods with valid resource requests. The cluster has been configured with all required components including proper networking, IAM roles, and subnet tags, but no nodes are being created.

## Environment Details
- **AWS Account ID**: 179151668767
- **Region**: us-west-2
- **Cluster Name**: crowdunlocked-dev
- **EKS Version**: 1.34
- **Platform Version**: eks.9
- **Cluster Created**: 2025-12-08 22:37:30 (recreated with proper networking)

## Configuration Details

### EKS Auto Mode Configuration
```
ComputeConfig:
  enabled: true
  nodePools: ["general-purpose"]
  nodeRoleArn: arn:aws:iam::179151668767:role/crowdunlocked-dev-eks-node
```

### Authentication Mode
- API_AND_CONFIG_MAP (required for Auto Mode)

### VPC and Networking
- **VPC ID**: vpc-0c8e0a8f0e8e0a8f0 (with DNS support and hostnames enabled)
- **Private Subnets**: 3 subnets across us-west-2a, us-west-2b, us-west-2c
  - subnet-0a599f7e4ff3d01d5 (10.0.1.0/24, us-west-2a)
  - subnet-0a091e9f4920abd8e (10.0.2.0/24, us-west-2b)
  - subnet-006c5531e7cadb4cc (10.0.3.0/24, us-west-2c)
- **Public Subnets**: 3 subnets with Internet Gateway
- **NAT Gateways**: 3 (one per AZ for high availability)
- **Route Tables**: Properly configured (public → IGW, private → NAT)

### Subnet Tags (All Applied)
```
kubernetes.io/cluster/crowdunlocked-dev = "shared"
kubernetes.io/role/internal-elb = "1"
```

### IAM Role Configuration
**Node Role**: crowdunlocked-dev-eks-node

**Trust Policy** (updated for Auto Mode):
```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Service": "eks.amazonaws.com"
    },
    "Action": "sts:AssumeRole"
  }]
}
```

**Attached Policies**:
- AmazonEKSWorkerNodePolicy
- AmazonEKS_CNI_Policy
- AmazonEC2ContainerRegistryReadOnly

### Pending Workloads
4 Flux CD controller pods in flux-system namespace, all in Pending state:
- helm-controller (requests: 100m CPU, 64Mi memory)
- kustomize-controller (requests: 100m CPU, 64Mi memory)
- notification-controller (requests: 100m CPU, 64Mi memory)
- source-controller (requests: 100m CPU, 64Mi memory)

**Pod Events**: "no nodes available to schedule pods"

## Troubleshooting Steps Taken

1. ✅ Verified cluster status is ACTIVE
2. ✅ Verified compute_config is enabled with general-purpose node pool
3. ✅ Verified authentication_mode is API_AND_CONFIG_MAP
4. ✅ Added kubernetes.io/cluster/crowdunlocked-dev tags to all subnets
5. ✅ Updated node role trust policy from ec2.amazonaws.com to eks.amazonaws.com
6. ✅ Verified all required IAM policies are attached to node role
7. ✅ Verified VPC has proper networking (IGW, NAT gateways, route tables)
8. ✅ Verified subnets have proper tags for ELB integration
9. ✅ Enabled CloudWatch logging for cluster
10. ✅ Deleted and recreated pods to trigger node provisioning
11. ✅ Verified AWSServiceRoleForAmazonEKS service-linked role exists
12. ✅ Verified no EC2 instances or Auto Scaling groups exist (expected for Auto Mode)

## Expected Behavior
EKS Auto Mode should automatically provision compute nodes when pods are pending and cannot be scheduled due to lack of nodes.

## Actual Behavior
- Pods remain in Pending state indefinitely
- No compute nodes are provisioned
- No EC2 instances are launched
- No Auto Scaling groups are created
- No errors visible in CloudWatch logs
- `kubectl get nodes` returns "No resources found"

## Questions for AWS Support

1. Are there any additional IAM permissions required for EKS Auto Mode beyond the standard node policies?
2. Is there a specific service-linked role or additional configuration needed for Auto Mode compute provisioning?
3. Are there any known issues with EKS Auto Mode in us-west-2 on platform version eks.9?
4. Is there a way to view Auto Mode compute provisioning logs or status?
5. Does Auto Mode require any specific VPC or subnet configuration beyond the standard EKS requirements?

## Additional Context
- This cluster was initially created without proper networking (missing NAT gateways), then deleted and recreated with complete networking infrastructure
- The node role trust policy was initially set to ec2.amazonaws.com and later corrected to eks.amazonaws.com
- All infrastructure is managed via Terraform/OpenTofu and deployed through CI/CD
- We attended multiple re:Invent 2024 sessions about EKS Auto Mode and are eager to use this new feature

## Commands to Reproduce

```bash
# Check cluster status
aws eks describe-cluster --name crowdunlocked-dev --region us-west-2

# Check for nodes
kubectl get nodes

# Check pending pods
kubectl get pods -n flux-system

# Check pod events
kubectl describe pod -n flux-system <pod-name>

# Check for EC2 instances
aws ec2 describe-instances --region us-west-2 \
  --filters "Name=tag:eks:cluster-name,Values=crowdunlocked-dev"

# Check node role
aws iam get-role --role-name crowdunlocked-dev-eks-node
```

## Request
Please help identify why EKS Auto Mode is not provisioning compute nodes for this cluster despite having all required configuration in place.
