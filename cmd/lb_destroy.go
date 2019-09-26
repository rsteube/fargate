package cmd

import (
	"fmt"

	"github.com/awslabs/fargatecli/console"
	ELBV2 "github.com/awslabs/fargatecli/elbv2"
	"github.com/spf13/cobra"
	 zsh "github.com/rsteube/cobra-zsh-gen"
)

type LoadBalancerDestroyOperation struct {
	LoadBalancerName string
}

var loadBalancerDestroyCmd = &cobra.Command{
	Use:   "destroy <load-balancer-name>",
	Short: "Destroy load balancer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		operation := &LoadBalancerDestroyOperation{
			LoadBalancerName: args[0],
		}

		destroyLoadBalancer(operation)
	},
}

func init() {
	zsh.Wrap(loadBalancerDestroyCmd).MarkZshCompPositionalArgumentCustom(1, "__fargate_completion_loadbalancer")
	lbCmd.AddCommand(loadBalancerDestroyCmd)
}

func destroyLoadBalancer(operation *LoadBalancerDestroyOperation) {
	elbv2 := ELBV2.New(sess)

	elbv2.DeleteLoadBalancer(operation.LoadBalancerName)
	elbv2.DeleteTargetGroup(fmt.Sprintf(defaultTargetGroupFormat, operation.LoadBalancerName))
	console.Info("Destroyed load balancer %s", operation.LoadBalancerName)
}
