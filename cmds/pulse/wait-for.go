// Package pulse contains commands to interact with pulse..
package pulse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/taskcluster/taskcluster-cli/cmds/root"

	"github.com/spf13/cobra"

	"github.com/donovanhide/eventsource"
)

func init() {
	root.Command.AddCommand(&cobra.Command{
		Use:   "wait-for <taskId>",
		Short: "Wait for a task to be finished",
		RunE:  waitForTask,
	})
}

func waitForTask(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("%s expects argument <taskId>", cmd.Name())
	}

	type Binding struct {
		Exchange   string `json:"exchange"`
		RoutingKey string `json:"routingKeyPattern"`
	}

	type Bindings struct {
		Bindings []Binding `json:"bindings"`
	}

	bindings := []Binding{
		Binding{
			args[0],
			args[1]}}

	json_bindings, _ := json.Marshal(Bindings{Bindings: bindings})
	values := url.Values{"bindings": {string(json_bindings)}}
	eventsUrl := "https://taskcluster-events-staging.herokuapp.com/v1/connect/?" + values.Encode()

	stream, err := eventsource.Subscribe(eventsUrl, "")
	if err != nil {
		return fmt.Errorf("Error getting request: %v", err)
	}

	for {
		event, ok := <-stream.Events
		if ok == false {
			err := <-stream.Errors
			stream.Close()
			return fmt.Errorf("Error: %v", err)
		}
		fmt.Println(event.Event(), event.Data())
	}
	return nil
}