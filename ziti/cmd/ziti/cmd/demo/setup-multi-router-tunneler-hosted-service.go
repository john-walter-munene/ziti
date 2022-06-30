/*
	Copyright NetFoundry Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package demo

import (
	_ "embed"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/api"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/common"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	tutorial2 "github.com/openziti/ziti/ziti/cmd/ziti/cmd/tutorial"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/tutorial/actions"
	"github.com/openziti/ziti/ziti/cmd/ziti/tutorial"
	"github.com/spf13/cobra"
	"time"
)

//go:embed setup-scripts/multi-router-tunneler-hosted.md
var multiRouterTunnelerHostedScriptSource []byte

type multiRouterTunnelerHosted struct {
	api.Options
	tutorial2.TutorialOptions
	interactive bool
}

func newMultiRouterTunnelerHostedCmd(p common.OptionsProvider) *cobra.Command {
	options := &multiRouterTunnelerHosted{
		Options: api.Options{
			CommonOptions: p(),
		},
	}

	cmd := &cobra.Command{
		Use:   "multi-router-tunneler-hosted",
		Short: "Walks you through hosting configuration for a simple echo service, hosted by router-embedded tunnelers and multiple applications",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.run()
			cmdhelper.CheckErr(err)
		},
		SuggestFor: []string{},
	}

	// allow interspersing positional args and flags
	cmd.Flags().SetInterspersed(true)
	cmd.Flags().StringVar(&options.ControllerUrl, "controller-url", "", "The Ziti controller URL to use")
	cmd.Flags().StringVarP(&options.Username, "username", "u", "", "The Ziti controller username to use")
	cmd.Flags().StringVarP(&options.Password, "password", "p", "", "The Ziti controller password to use")
	cmd.Flags().DurationVar(&options.NewlinePause, "newline-pause", time.Millisecond*10, "How long to pause between lines when scrolling")
	cmd.Flags().BoolVar(&options.interactive, "interactive", false, "Interactive mode, waiting for user input")
	options.AddCommonFlags(cmd)

	return cmd
}

func (self *multiRouterTunnelerHosted) run() error {
	t := tutorial.NewRunner()
	t.NewLinePause = self.NewlinePause
	t.AssumeDefault = !self.interactive

	t.RegisterActionHandler("ziti", &actions.ZitiRunnerAction{})
	t.RegisterActionHandler("ziti-login", &actions.ZitiEnsureLoggedIn{
		LoginParams: &self.TutorialOptions,
	})
	t.RegisterActionHandler("keep-session-alive", &actions.KeepSessionAliveAction{})
	t.RegisterActionHandler("ziti-create-config", &actions.ZitiCreateConfigAction{})
	t.RegisterActionHandler("ziti-for-each", &actions.ZitiForEach{})

	return t.Run(multiRouterTunnelerHostedScriptSource)
}
