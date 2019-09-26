package cmd

import (
	"fmt"

	"github.com/awslabs/fargatecli/console"
	ECS "github.com/awslabs/fargatecli/ecs"
	ELBV2 "github.com/awslabs/fargatecli/elbv2"
	"github.com/spf13/cobra"
	 zsh "github.com/rsteube/cobra-zsh-gen"
)

type ServiceDestroyOperation struct {
	ServiceName string
}

var serviceDestroyCmd = &cobra.Command{
	Use:   "destroy <service-name>",
	Short: "Destroy a service",
	Long: `Destroy service

In order to destroy a service, it must first be scaled to 0 running tasks.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceDestroyOperation{
			ServiceName: args[0],
		}

		destroyService(operation)
	},
}

func init() {
	zsh.Wrap(serviceDestroyCmd).MarkZshCompPositionalArgumentCustom(1, "__fargate_completion_service")
	serviceCmd.AddCommand(serviceDestroyCmd)
}

func destroyService(operation *ServiceDestroyOperation) {
	elbv2 := ELBV2.New(sess)
	ecs := ECS.New(sess, clusterName)
	service := ecs.DescribeService(operation.ServiceName)

	if service.DesiredCount > 0 {
		err := fmt.Errorf("%d tasks running, scale service to 0", service.DesiredCount)
		console.ErrorExit(err, "Cannot destroy service %s", operation.ServiceName)
	}

	if service.TargetGroupArn != "" {
		loadBalancerArn := elbv2.GetTargetGroupLoadBalancerArn(service.TargetGroupArn)
		loadBalancer := elbv2.DescribeLoadBalancerByARN(loadBalancerArn)
		listeners := elbv2.GetListeners(loadBalancerArn)

		for _, listener := range listeners {
			for _, rule := range elbv2.DescribeRules(listener.ARN) {
				if rule.TargetGroupARN == service.TargetGroupArn {
					if rule.IsDefault {
						defaultTargetGroupName := fmt.Sprintf(defaultTargetGroupFormat, loadBalancer.Name)
						defaultTargetGroupArn := elbv2.GetTargetGroupArn(defaultTargetGroupName)

						if defaultTargetGroupArn == "" {
							defaultTargetGroupArn, _ = elbv2.CreateTargetGroup(
								ELBV2.CreateTargetGroupParameters{
									Name:     defaultTargetGroupName,
									Port:     listeners[0].Port,
									Protocol: listeners[0].Protocol,
									VPCID:    loadBalancer.VPCID,
								},
							)
						}

						elbv2.ModifyListenerDefaultAction(listener.ARN, defaultTargetGroupArn)
					} else {
						elbv2.DeleteRule(rule.ARN)
					}
				}
			}
		}

		elbv2.DeleteTargetGroupByArn(service.TargetGroupArn)
	}

	ecs.DestroyService(operation.ServiceName)
	console.Info("Destroyed service %s", operation.ServiceName)
}
