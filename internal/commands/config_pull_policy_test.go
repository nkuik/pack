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
		defaultCfg   = config.Config{
			Experimental: true,
		}
		// neverCfg = config.Config{
		// 	Experimental: true,
		// 	PullPolicy:   "never",
		// }
		// ifNotPresentCfg = config.Config{
		// 	Experimental: true,
		// 	PullPolicy:   "if-not-present",
		// }
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

	// when("-h", func() {
	// 	it("prints available commands", func() {
	// 		cmd.SetArgs([]string{"-h"})
	// 		h.AssertNil(t, cmd.Execute())
	// 		output := outBuf.String()
	// 		h.AssertContains(t, output, "Usage:")
	// 		for _, command := range []string{"add", "remove", "list"} {
	// 			h.AssertContains(t, output, command)
	// 		}
	// 	})
	// })

	when("no arguments", func() {
		it("lists default pull policy", func() {
			cmd.SetArgs([]string{})
			h.AssertNil(t, cmd.Execute())
			output := outBuf.String()
			h.AssertEq(t, strings.TrimSpace(output), `Pull policy is always`)
		})
	})
}
