package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	pubcfg "github.com/buildpacks/pack/config"

	"github.com/buildpacks/pack/internal/config"
	"github.com/buildpacks/pack/logging"
)

func ConfigPullPolicy(logger logging.Logger, cfg config.Config, cfgPath string) *cobra.Command {
	var unset bool
	cmd := &cobra.Command{
		Use:     "pull-policy <if-not-present | always | never>",
		Args:    cobra.MaximumNArgs(1),
		Short:   "List, set, or unset global pull policy",
		Aliases: []string{"pull-policy"},
		RunE: logError(logger, func(cmd *cobra.Command, args []string) error {
			var pullPolicy pubcfg.PullPolicy
			var err error
			if len(args) == 0 {
				pullPolicy, err = pubcfg.ParsePullPolicy(cfg.PullPolicy)
				if err != nil {
					return err
				}
			} else {
				newPullPolicy := args[0]

				if newPullPolicy == cfg.PullPolicy {
					logger.Infof("Pull policy is already set to %s", newPullPolicy)
					return nil
				}

				pullPolicy, err = pubcfg.ParsePullPolicy(newPullPolicy)
				if err != nil {
					return err
				}

				cfg.PullPolicy = newPullPolicy
				if err := config.Write(cfg, cfgPath); err != nil {
					return errors.Wrap(err, "writing to config")
				}
			}
			if unset {
				cfg.PullPolicy = ""
				if err := config.Write(cfg, cfgPath); err != nil {
					return errors.Wrap(err, "writing to config")
				}
				logger.Infof("Unsetting configured pull policy")
				pullPolicy, err = pubcfg.ParsePullPolicy(cfg.PullPolicy)
				if err != nil {
					return err
				}
			}

			logger.Infof("Pull policy is %s", pullPolicy.String())
			return nil
		}),
	}
	cmd.Flags().BoolVarP(&unset, "unset", "u", false, "Unset pull policy")

	AddHelpFlag(cmd, "pull-policy")
	return cmd
}
