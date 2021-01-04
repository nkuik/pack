package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	pubcfg "github.com/buildpacks/pack/config"

	"github.com/buildpacks/pack/internal/config"
	"github.com/buildpacks/pack/internal/style"
	"github.com/buildpacks/pack/logging"
)

func ConfigPullPolicy(logger logging.Logger, cfg config.Config, cfgPath string) *cobra.Command {
	var unset bool
	cmd := &cobra.Command{
		Use:     "pull-policy <if-not-present | always | never>",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Set or unset global pull policy",
		Aliases: []string{"pull-policy"},
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				pullPolicy, err := pubcfg.ParsePullPolicy(cfg.PullPolicy)
				if err != nil {
					return err
				}
				logger.Infof("Pull policy is %s", pullPolicy.String())
			} else {
				newPullPolicy := args[0]

				if newPullPolicy == cfg.PullPolicy {
					logger.Infof("Pull policy is already set to %s", style.Symbol(newPullPolicy))
					return nil
				}
				pullPolicy, err := pubcfg.ParsePullPolicy(newPullPolicy)
				if err != nil {
					return err
				}
				cfg.PullPolicy = newPullPolicy
				if err = config.Write(cfg, cfgPath); err != nil {
					return errors.Wrap(err, "writing to config")
				}
				logger.Infof("New pull policy is %s", pullPolicy.String())
			}
			if unset {
				cfg.PullPolicy = ""
				if err := config.Write(cfg, cfgPath); err != nil {
					return errors.Wrap(err, "writing to config")
				}
				logger.Infof("Resetting pull policy to always")
			}
			return nil
		}),
	}
	cmd.Flags().BoolVarP(&unset, "unset", "u", false, "Unset pull policy")

	AddHelpFlag(cmd, "pull-policy")
	return cmd
}
