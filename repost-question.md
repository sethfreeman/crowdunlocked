# EKS Auto Mode Not Provisioning Nodes Despite Pending Pods

## Summary
I have an EKS 1.34 cluster with Auto Mode enabled, proper networking, correct IAM roles, and pending pods, but Auto Mode is not provisioning any compute nodes.

## Environment
- **Region**: us-west-2
- **EKS Version**: 1.34
- **Platform Version**: eks.9
- **Cluster Created**: December 8, 2024

## Configuration

### Auto Mode Config
```
compute_config {
  enabled                      = true
  node_pools                   = ["general-purpose"]
  node_role_arn                = arn:aws:iam::ACCOUNT:role/cluster-eks-node
}

access_config {
  authentication_mode = "API_AND_CONFIG_MAP"
}
```

### Networking
- VPC with DNS support enabled
- 3 private subnets across 3 AZs (10.0.1.0/24, 10.0.2.0/24, 10.0.3.0/24)
- 3 public subnets with Internet Gateway
- 3 NAT Gateways (one per AZ)
- Route tables properly configured

### Subnet Tags
All subnets tagged with:
```
kubernetes.io/cluster/CLUSTER-NAME = "shared"
kubernetes.io/role/internal-elb = "1"  (private subnets)
kubernetes.io/role/elb = "1"  (public subnets)
```

### IAM Node Role
Trust policy:
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

Attached policies:
- AmazonEKSWorkerNodePolicy
- AmazonEKS_CNI_Policy
- AmazonEC2ContainerRegistryReadOnly

## Problem
I have 4 Flux CD pods pending in the flux-system namespace with reasonable resource requests (100m CPU, 64Mi memory each). The pods show "no nodes available to schedule pods" but Auto Mode is not provisioning any nodes.

```bash
$ kubectl get nodes
No resources found

$ kubectl get pods -n flux-system
NAME                                       READY   STATUS    RESTARTS   AGE
helm-controller-68578f8447-n9msx           0/1     Pending   0          30m
kustomize-controller-7ddfbb5875-c4h5m      0/1     Pending   0          30m
notification-controller-6d766f87cf-8xzhd   0/1     Pending   0          30m
source-controller-6679d8bdb-xknq9          0/1     Pending   0          30m
```

## What I've Verified
✅ Cluster status is ACTIVE  
✅ compute_config shows enabled=true  
✅ Node role has correct trust policy (eks.amazonaws.com)  
✅ All required IAM policies attached  
✅ Subnets have cluster-specific tags  
✅ VPC networking is complete (IGW, NAT, routes)  
✅ CloudWatch logging enabled  
✅ AWSServiceRoleForAmazonEKS exists  
✅ No EC2 instances or ASGs (expected for Auto Mode)  

## Questions
1. Are there additional IAM permissions needed for Auto Mode beyond the standard node policies?
2. Is there a way to view Auto Mode provisioning logs or status?
3. Does Auto Mode require any specific configuration I'm missing?
4. Are there known issues with Auto Mode on EKS 1.34 platform version eks.9?

I attended multiple re:Invent sessions on Auto Mode and am excited to use this feature, but can't get it working despite following all documentation.

Any help would be greatly appreciated!

---
**Tags**: eks, eks-auto-mode, kubernetes, compute, node-provisioning
