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
		// Add policy variables here, use them instead
		cmd          *cobra.Command
		logger       logging.Logger
		outBuf       bytes.Buffer
		tempPackHome string
		configPath   string
		assert     = h.NewAssertionManager(t)
		defaultCfg   = config.Config{}
		neverCfg     = config.Config{
			Experimental: true,
			PullPolicy:   "never",
		}
		invalidCfg = config.Config{
			Experimental: true,
			PullPolicy:   "blah",
		}
		ifNotPresentCfg = config.Config{
			Experimental: true,
			PullPolicy:   "if-not-present",
		}
	)

	it.Before(func() {
		var err error
		logger = ilogging.NewLogWithWriters(&outBuf, &outBuf)
		tempPackHome, err = ioutil.TempDir("", "pack-home")
		h.AssertNil(t, err)
		configPath = filepath.Join(tempPackHome, "config.toml")

		cmd = commands.ConfigPullPolicy(logger, defaultCfg, configPath)
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
			it.Before(func() {
				cmd = commands.ConfigPullPolicy(logger, neverCfg, configPath)
			})

			it("lists never as pull policy", func() {
				cmd.SetArgs([]string{})
				h.AssertNil(t, cmd.Execute())
				output := outBuf.String()
				h.AssertEq(t, strings.TrimSpace(output), `Pull policy is never`)
			})
		})
		when("policy set to if-not-present in config", func() {
			it.Before(func() {
				cmd = commands.ConfigPullPolicy(logger, ifNotPresentCfg, configPath)
			})

			it("lists if-not-present as pull policy", func() {
				cmd.SetArgs([]string{})
				h.AssertNil(t, cmd.Execute())
				output := outBuf.String()
				h.AssertEq(t, strings.TrimSpace(output), `Pull policy is if-not-present`)
			})
		})
		when("invalid policy set", func() {
			it.Before(func() {
				cmd = commands.ConfigPullPolicy(logger, invalidCfg, configPath)
			})

			it("reports error", func() {
				cmd.SetArgs([]string{})
				err := cmd.Execute()
				h.AssertError(t, err, `invalid pull policy blah`)
			})
		})
		when("policy is specified", func() {
			it("should set the policy when policy is valid", func() {
				cmd.SetArgs([]string{"never"})
				assert.Succeeds(cmd.Execute())
				cfg, err := config.Read(configPath)
				assert.Nil(err)
				assert.Equal(cfg.PullPolicy, "never")
			})
			it("should fail when policy is invalid", func() {
				cmd.SetArgs([]string{"invalid"})

				err := cmd.Execute()
				h.AssertError(t, err, `invalid pull policy invalid`)
			})
		})
		when("run with --unset", func() {
			it.Before(func() {
				cmd = commands.ConfigPullPolicy(logger, neverCfg, configPath)
			})

			it("should reset to default pull policy", func() {
				cmd.SetArgs([]string{"--unset"})
				assert.Succeeds(cmd.Execute())

				cfg, err := config.Read(configPath)
				assert.Nil(err)
				assert.Equal(cfg.PullPolicy, "")
			})
		})
	})
}
