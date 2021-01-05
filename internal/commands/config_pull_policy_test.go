package commands_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heroku/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"github.com/spf13/cobra"

	"github.com/buildpacks/pack/internal/commands"
	"github.com/buildpacks/pack/internal/config"
	ilogging "github.com/buildpacks/pack/internal/logging"
	"github.com/buildpacks/pack/logging"
	h "github.com/buildpacks/pack/testhelpers"
)

func TestConfigPullPolicy(t *testing.T) {
	color.Disable(true)
	defer color.Disable(false)
	spec.Run(t, "ConfigPullPolicyCommand", testConfigPullPolicyCommand, spec.Random(), spec.Report(report.Terminal{}))
}

func testConfigPullPolicyCommand(t *testing.T, when spec.G, it spec.S) {
	var (
		cmd          *cobra.Command
		logger       logging.Logger
		outBuf       bytes.Buffer
		tempPackHome string
		configPath   string
		assert       = h.NewAssertionManager(t)
		cfg          = config.Config{}
	)

	it.Before(func() {
		var err error
		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)
		tempPackHome, err = ioutil.TempDir("", "pack-home")
		h.AssertNil(t, err)
		configPath = filepath.Join(tempPackHome, "config.toml")

		cmd = commands.ConfigPullPolicy(logger, cfg, configPath)
		cmd.SetOut(logging.GetWriterForLevel(logger, logging.InfoLevel))
	})

	it.After(func() {
		h.AssertNil(t, os.RemoveAll(tempPackHome))
	})

	when("#ConfigPullPolicy", func() {
		when("no policy is specified", func() {
			it("lists default pull policy", func() {
				cmd.SetArgs([]string{})

				h.AssertNil(t, cmd.Execute())

				output := outBuf.String()
				h.AssertEq(t, strings.TrimSpace(output), `Pull policy is always`)
			})
		})

		when("policy set to never in config", func() {
			it("lists never as pull policy", func() {
				cfg.PullPolicy = "never"
				cmd = commands.ConfigPullPolicy(logger, cfg, configPath)
				cmd.SetArgs([]string{})

				h.AssertNil(t, cmd.Execute())

				output := outBuf.String()
				h.AssertEq(t, strings.TrimSpace(output), `Pull policy is never`)
			})
		})

		when("policy set to if-not-present in config", func() {
			it("lists if-not-present as pull policy", func() {
				cfg.PullPolicy = "if-not-present"
				cmd = commands.ConfigPullPolicy(logger, cfg, configPath)
				cmd.SetArgs([]string{})

				h.AssertNil(t, cmd.Execute())

				output := outBuf.String()
				h.AssertEq(t, strings.TrimSpace(output), `Pull policy is if-not-present`)
			})
		})

		when("policy provided is the same as configured pull policy", func() {
			it("provides a helpful message", func() {
				cfg.PullPolicy = "if-not-present"
				cmd = commands.ConfigPullPolicy(logger, cfg, configPath)
				cmd.SetArgs([]string{"if-not-present"})

				h.AssertNil(t, cmd.Execute())

				output := outBuf.String()
				h.AssertEq(t, strings.TrimSpace(output), `Pull policy is already set to if-not-present`)
			})
			it("it does not change the configured policy", func() {
				cmd = commands.ConfigPullPolicy(logger, cfg, configPath)
				cmd.SetArgs([]string{"never"})

				cmd = commands.ConfigPullPolicy(logger, cfg, configPath)
				cmd.SetArgs([]string{"never"})
				assert.Succeeds(cmd.Execute())

				readCfg, err := config.Read(configPath)
				assert.Nil(err)
				assert.Equal(readCfg.PullPolicy, "never")
			})
		})

		when("invalid policy is specified", func() {
			it("reports error", func() {
				cfg.PullPolicy = "unknown-policy"
				cmd = commands.ConfigPullPolicy(logger, cfg, configPath)
				cmd.SetArgs([]string{})

				err := cmd.Execute()
				h.AssertError(t, err, `invalid pull policy unknown-policy`)
			})
		})

		when("valid policy is specified", func() {
			it("sets the policy in config", func() {
				cmd.SetArgs([]string{"never"})
				assert.Succeeds(cmd.Execute())

				readCfg, err := config.Read(configPath)
				assert.Nil(err)
				assert.Equal(readCfg.PullPolicy, "never")
			})
			it("should fail invalid policy given", func() {
				cmd.SetArgs([]string{"unknown-policy"})

				err := cmd.Execute()
				h.AssertError(t, err, `invalid pull policy unknown-policy`)
			})
		})

		when("--unset", func() {
			it("removes set policy and resets to default pull policy", func() {
				cmd.SetArgs([]string{"never"})
				cmd = commands.ConfigPullPolicy(logger, cfg, configPath)

				cmd.SetArgs([]string{"--unset"})
				assert.Succeeds(cmd.Execute())

				cfg, err := config.Read(configPath)
				assert.Nil(err)
				assert.Equal(cfg.PullPolicy, "")
			})
		})
	})
}
