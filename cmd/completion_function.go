package cmd

var bashCompletionFunction = `
function __fargate_completion_region {
  _values 'region' us-east-1 us-east-2 us-west-2 ap-southeast-1 ap-southeast-2 ap-northeast-1 eu-central-1 eu-west-1
}

function __fargate_completion_certificate {
  certificates=($(fargate certificate list | tail -n +2 | cut -f1 | sed "s_*_\\\*_g"))
  [ -z "$certificates" ] || _values 'certificate' $certificates
}

function __fargate_completion_log {
  while [ $# -gt 1 ]; do shift; done;
  logs=($(aws logs describe-log-groups | jq -r '.logGroups[].logGroupName | select(startswith("'$1'")) | sub("'$1'" ;"")'))
  [ -z "$logs" ] || _values 'log' $logs
}

function __fargate_completion_service {
  services=($(fargate service list | tail -n +2 | cut -f1))
  [ -z "$services" ] || _values 'service' $services
}

function __fargate_completion_task {
  tasks=($(fargate task list | tail -n +2 | cut -f1))
  [ -z "$tasks" ] || _values 'task' $tasks
}

function __fargate_completion_taskid {
  taskids=($(fargate task ps test | tail -n +2 | cut -f-2,4 | sed -e 's_:_\\:_' -e 's_\t_:_' -e 's_\t_(_' -e 's_$_)_'))
  [ -z "$taskids" ] || _describe 'taskid' taskids
 }

function __fargate_completion_loadbalancer {
  loadbalancers=($(fargate lb list | tail -n +2 | cut -f1))
  [ -z "$loadbalancers" ] || _values 'loadbalancer' $loadbalancers
}

function __fargate_completion_cpu {
  _values 'cpu' 256 512 1024 2048 4096
}

function __fargate_completion_memory {
  _values 'memory' 512 1024 2048 3072 4096 5120 6144 7168 8192 9216 10240 11264 12288 13312 14336 15360 16348
}

function __fargate_completion_port {
  _values 'port' 80 443 'http\:8080' 'https\:8443' 'tcp\:1935'
}

function __fargate_completion_zone {
  zones=($(aws route53 list-hosted-zones | jq -r '.HostedZones[].Name | rtrimstr(".")' | sed 's_^_._'))
  [ -z "$zones" ] || _values 'zone' $zones
}

function __fargate_completion_cluster {
  clusters=($(aws ecs list-clusters | jq -r '.clusterArns[] | sub(".*/"; "")'))
  [ -z "$clusters" ] || _values 'cluster' $clusters
}

function __fargate_completion_securitygroup {
  groups=($(aws ec2 describe-security-groups --query 'SecurityGroups[].{id: GroupId, name: GroupName}' | jq -r '.[] | "\(.id):\(.name)"'))
  [ -z "$groups" ] || _describe 'group' groups
}

function __fargate_completion_subnet {
  subnets=($(aws ec2 describe-subnets | jq -r '.Subnets[] | "\(.SubnetId):\(.AvailabilityZone)"'))
  [ -z "$subnets" ] || _describe 'subnet' subnets
}

function __fargate_completion_image {
  images=($(aws ecr describe-repositories | jq -r '.repositories[].repositoryUri'))
  [ -z "$images" ] || _values 'image' $images
}

function __fargate_completion_role {
  roles=($(aws iam list-roles | jq -r '.Roles[].RoleName'))
  [ -z "$roles" ] || _values 'role' $roles
}

function __fargate_completion_env {
  keys=($(fargate service env list $line | sed 's_=.*__'))
  [ -z "$keys" ] || _values 'key' $keys
}
`
