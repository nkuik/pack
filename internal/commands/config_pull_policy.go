package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/buildpacks/pack/internal/config"
	"github.com/buildpacks/pack/internal/style"
	"github.com/buildpacks/pack/logging"
)

// When no flag for policy passed, use default policy
// Default policy should be:
//  1. "always" if not set in config
//  2. value in config if set in config
//  3. Set if the pull-policy with policy command is run
// When flag is passed, should take precedence to default value

func SetPullPolicy(logger logging.Logger, cfg config.Config, cfgPath string) *cobra.Command {
	var unset bool
	cmd := &cobra.Command{
		Use:     "pull-policy <if-not-present | always | never>",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Set or unset global pull policy",
		Aliases: []string{"pull-policy"},
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				pullPolicy, err := config.ParsePullPolicy(cfg.PullPolicy)
				if err != nil {
					return errors.Wrapf(err, "parsing pull policy %s", cfg.PullPolicy)
				}
				logger.Infof("Pull policy is %s", pullPolicy.String())
			} else {
				newPullPolicy := args[0]

				if newPullPolicy == cfg.PullPolicy {
					logger.Infof("Pull policy is already set to %s", style.Symbol(newPullPolicy))
					return nil
				}
				pullPolicy, err := config.ParsePullPolicy(newPullPolicy)
				if err != nil {
					return errors.Wrapf(err, "parsing pull policy %s", newPullPolicy)
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
