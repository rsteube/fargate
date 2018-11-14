package cmd

import (
	"github.com/jpignata/fargate/elbv2"
	"github.com/jpignata/fargate/route53"
	"github.com/spf13/cobra"
)

type lbAliasOperation struct {
	lbOperation
	aliasDomain string
	lbName      string
	output      Output
	route53     route53.Client
}

func (o lbAliasOperation) execute() {
	loadBalancer, err := o.findLB(o.lbName)

	if err != nil {
		o.output.Fatal(err, "Could not alias load balancer")
		return
	}

	hostedZones, err := o.route53.ListHostedZones()

	if err != nil {
		o.output.Fatal(err, "Could not alias load balancer")
		return
	}

	if hostedZone, ok := hostedZones.FindSuperDomainOf(o.aliasDomain); ok {
		o.output.Debug("Creating alias record [API=route53 Action=CreateResourceRecordSet]")
		id, err := o.route53.CreateAlias(
			route53.CreateAliasInput{
				HostedZoneID:       hostedZone.ID,
				RecordType:         "A",
				Name:               o.aliasDomain,
				Target:             loadBalancer.DNSName,
				TargetHostedZoneID: loadBalancer.HostedZoneID,
			},
		)

		if err != nil {
			o.output.Fatal(err, "Could not alias load balancer")
			return
		}

		o.output.Debug("Created alias record [ChangeID=%s]", id)
		o.output.Info("Created alias record (%s -> %s)", o.aliasDomain, loadBalancer.DNSName)
	} else {
		o.output.Warn("Could not find hosted zone for %s", o.aliasDomain)
		o.output.Say("If you're hosting this domain elsewhere or in another AWS account, please manually create the alias record:", 1)
		o.output.Say("%s -> %s", 1, o.aliasDomain, loadBalancer.DNSName)
	}
}

var lbAliasCmd = &cobra.Command{
	Use:   "alias <load-balancer-name> <domain-name>",
	Args:  cobra.ExactArgs(2),
	Short: "Create a load balancer alias record",
	Long: `Create a load balancer alias record

Create an alias record to the load balancer for domains that are hosted within
Amazon Route 53 and within the same AWS account. If you're using another DNS
provider or host your domains in a different account, you will need to manually
create this record.  `,
	Run: func(cmd *cobra.Command, args []string) {
		lbAliasOperation{
			aliasDomain: args[1],
			lbName:      args[0],
			lbOperation: lbOperation{elbv2: elbv2.New(sess), output: output},
			output:      output,
			route53:     route53.New(sess),
		}.execute()
	},
}

func init() {
	lbAliasCmd.MarkZshCompPositionalArgumentCustom(1, "__fargate_completion_loadbalancer")
	lbAliasCmd.MarkZshCompPositionalArgumentCustom(2, "__fargate_completion_zone")
	lbCmd.AddCommand(lbAliasCmd)
}
