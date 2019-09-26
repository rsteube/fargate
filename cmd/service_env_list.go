package cmd

import (
	"fmt"
	ECS "github.com/awslabs/fargatecli/ecs"
	"github.com/spf13/cobra"
	 zsh "github.com/rsteube/cobra-zsh-gen"
)

type ServiceEnvListOperation struct {
	ServiceName string
}

var serviceEnvListCmd = &cobra.Command{
	Use:   "list <service-name>",
	Short: "Show environment variables",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceEnvListOperation{
			ServiceName: args[0],
		}

		serviceEnvList(operation)
	},
}

func init() {
	zsh.Wrap(serviceEnvListCmd).MarkZshCompPositionalArgumentCustom(1, "__fargate_completion_service")
	serviceEnvCmd.AddCommand(serviceEnvListCmd)
}

func serviceEnvList(operation *ServiceEnvListOperation) {
	ecs := ECS.New(sess, clusterName)
	service := ecs.DescribeService(operation.ServiceName)
	envVars := ecs.GetEnvVarsFromTaskDefinition(service.TaskDefinitionArn)

	for _, envVar := range envVars {
		fmt.Printf("%s=%s\n", envVar.Key, envVar.Value)
	}
}
