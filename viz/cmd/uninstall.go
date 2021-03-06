package cmd

import (
	"context"
	"fmt"
	"os"

	pkgCmd "github.com/linkerd/linkerd2/pkg/cmd"
	"github.com/linkerd/linkerd2/pkg/k8s"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func newCmdUninstall() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Args:  cobra.NoArgs,
		Short: "Output Kubernetes resources to uninstall the linkerd-viz extension",
		Long: `Output Kubernetes resources to uninstall the linkerd-viz extension.

This command provides all Kubernetes namespace-scoped and cluster-scoped resources (e.g services, deployments, RBACs, etc.) necessary to uninstall the Linkerd-viz extension.`,
		Example: `linkerd viz uninstall | kubectl delete -f -`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := uninstallRunE(cmd.Context())
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return nil
		},
	}

	return cmd
}

func uninstallRunE(ctx context.Context) error {
	k8sAPI, err := k8s.NewAPI(kubeconfigPath, kubeContext, impersonate, impersonateGroup, 0)
	if err != nil {
		return err
	}

	extensionReq, err := labels.NewRequirement(k8s.LinkerdExtensionLabel, selection.In, []string{ExtensionName, LegacyExtensionName})
	if err != nil {
		return err
	}

	selector := labels.NewSelector().Add(*extensionReq)
	return pkgCmd.Uninstall(ctx, k8sAPI, selector.String())
}
