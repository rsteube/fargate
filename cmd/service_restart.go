package cmd

import (
	"github.com/jpignata/fargate/console"
	ECS "github.com/jpignata/fargate/ecs"
	"github.com/spf13/cobra"
)

type ServiceRestartOperation struct {
	ServiceName string
}

var serviceRestartCmd = &cobra.Command{
	Use:   "restart <service-name>",
	Short: "Restart service",
	Long: `Restart service

Creates a new set of tasks for the service and stops the previous tasks. This
is useful if your service needs to reload data cached from an external source,
for example.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		operation := &ServiceRestartOperation{
			ServiceName: args[0],
		}

		restartService(operation)
	},
}

func init() {
	serviceRestartCmd.MarkZshCompPositionalArgumentCustom(1, "__fargate_completion_service")
	serviceCmd.AddCommand(serviceRestartCmd)
}

func restartService(operation *ServiceRestartOperation) {
	ecs := ECS.New(sess, clusterName)

	ecs.RestartService(operation.ServiceName)
	console.Info("Restarted %s", operation.ServiceName)
}
