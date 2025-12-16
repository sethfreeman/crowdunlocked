# EKS Auto Mode Migration Guide

## Current Status: Using Managed Node Groups

We are currently using **EKS Managed Node Groups** instead of **EKS Auto Mode** due to technical issues encountered during initial setup. This document explains why we made this decision and provides a migration path for when Auto Mode is ready.

## Why We're Not Using EKS Auto Mode (December 2024)

### Issues Encountered

1. **NodeClass Creation Failure**
   - Auto Mode failed to create a functional NodeClass
   - Instance profile creation failed with error: `InstanceProfileCreationFailed`
   - NodeClass status showed `ValidationSucceeded=False`

2. **Missing Service-Linked Role**
   - Required `AWSServiceRoleForEC2Spot` was missing
   - Auto Mode couldn't create this role automatically

3. **Node Provisioning Failure**
   - Despite pending pods with valid resource requests, no nodes were provisioned
   - NodePool showed `NodeClassNotFound` status
   - No EC2 instances or Auto Scaling groups created

4. **Undocumented Requirements**
   - Additional EC2 permissions needed beyond standard policies
   - Instance profile management unclear in documentation
   - Subnet tagging requirements not well documented

### Troubleshooting Attempts

We exhaustively tried the following solutions:

✅ **Networking Configuration**
- Complete VPC with NAT gateways, Internet Gateway, route tables
- Proper subnet tags: `kubernetes.io/cluster/CLUSTER-NAME = "shared"`
- Multi-AZ deployment across 3 availability zones

✅ **IAM Configuration**
- Correct trust policy: `eks.amazonaws.com` (not `ec2.amazonaws.com`)
- All standard policies: AmazonEKSWorkerNodePolicy, AmazonEKS_CNI_Policy, AmazonEC2ContainerRegistryReadOnly
- Additional EC2 permissions for Auto Mode

✅ **Service-Linked Roles**
- Created `AWSServiceRoleForEC2Spot` manually
- Verified `AWSServiceRoleForAmazonEKS` exists

✅ **Cluster Configuration**
- Authentication mode: `API_AND_CONFIG_MAP` (required for Auto Mode)
- EKS version 1.34 with platform version eks.9
- CloudWatch logging enabled

✅ **Pod Configuration**
- Valid resource requests (100m CPU, 64Mi memory)
- No restrictive node selectors or taints
- Multiple test deployments created

### Root Cause Analysis

The core issue appears to be that **EKS Auto Mode is a brand new feature** (announced at AWS re:Invent 2024) with:
- Incomplete documentation
- Potential bugs in the initial release
- Missing error reporting and debugging capabilities
- Unclear instance profile management

## Current Architecture: Managed Node Groups

### Configuration

```hcl
resource "aws_eks_node_group" "main" {
  cluster_name    = aws_eks_cluster.main.name
  node_group_name = "crowdunlocked-dev-nodes"
  node_role_arn   = aws_iam_role.eks_node.arn
  subnet_ids      = aws_subnet.private[*].id
  
  capacity_type  = "ON_DEMAND"
  instance_types = ["t3.medium"]
  
  scaling_config {
    desired_size = 2
    max_size     = 4
    min_size     = 1
  }
}
```

### Benefits of Current Approach

✅ **Reliability**: Mature, well-tested technology
✅ **Documentation**: Comprehensive AWS documentation and community resources
✅ **Predictability**: Known behavior and troubleshooting procedures
✅ **Support**: Full AWS support coverage
✅ **Control**: Explicit configuration of instance types and scaling

### Limitations vs Auto Mode

❌ **Manual Scaling**: Requires explicit scaling configuration
❌ **Instance Selection**: Manual instance type selection
❌ **Operational Overhead**: More configuration to manage
❌ **Cost Optimization**: Less automatic right-sizing

## Migration Path to EKS Auto Mode

### Prerequisites for Migration

Before attempting migration to Auto Mode, ensure:

1. **AWS Documentation Updates**
   - Complete Auto Mode documentation available
   - Clear troubleshooting guides published
   - Known issues and limitations documented

2. **Feature Maturity**
   - Auto Mode has been in production for 6+ months
   - Community adoption and success stories
   - AWS support confirms stability

3. **Testing Environment**
   - Test Auto Mode in a separate cluster first
   - Validate node provisioning works reliably
   - Confirm all workloads schedule correctly

### Migration Steps

#### Phase 1: Preparation

1. **Create Test Branch**
   ```bash
   git checkout -b feature/test-auto-mode
   ```

2. **Update Terraform Configuration**
   ```hcl
   # Remove managed node group
   # resource "aws_eks_node_group" "main" { ... }
   
   # Add Auto Mode configuration
   resource "aws_eks_cluster" "main" {
     # ... existing config ...
     
     compute_config {
       enabled                      = true
       node_pools                   = ["general-purpose"]
       node_role_arn                = aws_iam_role.eks_node.arn
     }
   }
   ```

3. **Update IAM Role Trust Policy**
   ```hcl
   # Change from ec2.amazonaws.com to eks.amazonaws.com
   assume_role_policy = jsonencode({
     Version = "2012-10-17"
     Statement = [{
       Action = "sts:AssumeRole"
       Effect = "Allow"
       Principal = {
         Service = "eks.amazonaws.com"  # Changed from ec2.amazonaws.com
       }
     }]
   })
   ```

4. **Add Auto Mode Permissions**
   ```hcl
   resource "aws_iam_role_policy" "eks_auto_mode_ec2" {
     name = "eks-auto-mode-ec2-permissions"
     role = aws_iam_role.eks_node.id
     
     policy = jsonencode({
       Version = "2012-10-17"
       Statement = [{
         Effect = "Allow"
         Action = [
           "ec2:RunInstances",
           "ec2:CreateFleet",
           "ec2:CreateLaunchTemplate",
           # ... additional permissions
         ]
         Resource = "*"
       }]
     })
   }
   ```

#### Phase 2: Testing

1. **Deploy to Test Environment**
   ```bash
   cd infra/terraform/dev
   tofu plan -var-file=dev.tfvars
   tofu apply -var-file=dev.tfvars
   ```

2. **Verify Node Provisioning**
   ```bash
   # Check cluster status
   kubectl get nodes
   
   # Check NodeClass and NodePool
   kubectl get nodeclass
   kubectl get nodepool
   kubectl describe nodepool general-purpose
   
   # Deploy test workload
   kubectl create deployment test-nginx --image=nginx
   kubectl scale deployment test-nginx --replicas=5
   
   # Verify nodes scale up
   kubectl get nodes -w
   ```

3. **Validate Workloads**
   ```bash
   # Test Flux controllers
   kubectl get pods -n flux-system
   
   # Test application deployments
   kubectl get pods -A
   
   # Check resource utilization
   kubectl top nodes
   kubectl top pods -A
   ```

#### Phase 3: Production Migration

1. **Create Migration PR**
   - Include all terraform changes
   - Document expected downtime (cluster recreation required)
   - Get team approval for maintenance window

2. **Execute Migration**
   ```bash
   # This will recreate the cluster
   tofu apply -var-file=dev.tfvars
   
   # Reconfigure kubectl
   aws eks update-kubeconfig --name crowdunlocked-dev --region us-west-2
   
   # Bootstrap Flux
   flux bootstrap github --context=crowdunlocked-dev \
     --owner=sethfreeman --repository=crowdunlocked \
     --branch=develop --path=flux/clusters/dev --personal
   ```

3. **Post-Migration Validation**
   - Verify all applications are running
   - Test auto-scaling behavior
   - Monitor for 48 hours
   - Update documentation

### Rollback Plan

If Auto Mode fails after migration:

1. **Immediate Rollback**
   ```bash
   git checkout feature/managed-node-groups
   cd infra/terraform/dev
   tofu apply -var-file=dev.tfvars
   ```

2. **Re-bootstrap Applications**
   ```bash
   ./scripts/setup-flux.sh
   ```

## Monitoring Auto Mode Readiness

### AWS Resources to Watch

1. **AWS Documentation**
   - https://docs.aws.amazon.com/eks/latest/userguide/auto-mode.html
   - AWS EKS User Guide updates

2. **AWS re:Post**
   - Monitor Auto Mode questions and AWS responses
   - Look for success stories and resolved issues

3. **AWS Support**
   - Check if Auto Mode moves out of "preview" status
   - Review AWS support policy coverage

### Success Criteria for Migration

Before migrating to Auto Mode, we need:

✅ **Stable Node Provisioning**: Consistent node creation for pending pods
✅ **Reliable NodeClass Management**: No instance profile creation failures
✅ **Complete Documentation**: Clear setup and troubleshooting guides
✅ **Community Validation**: Multiple success stories from other users
✅ **AWS Support Confirmation**: Full support coverage for Auto Mode issues

## Conclusion

**Current Decision**: Use managed node groups for reliability and immediate productivity.

**Future Path**: Migrate to Auto Mode when it reaches production maturity (estimated Q2 2025).

**Benefits of Waiting**:
- Avoid production issues with bleeding-edge features
- Let AWS resolve initial bugs and documentation gaps
- Benefit from community knowledge and best practices
- Ensure full AWS support coverage

This approach follows our **CI/CD First** principle: "Do not push features if CI/CD is broken" - we prioritize working infrastructure over cutting-edge features.

---

**Last Updated**: December 2024  
**Next Review**: March 2025  
**Status**: Managed Node Groups (Stable)