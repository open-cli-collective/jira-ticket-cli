package transitions

import (
	"fmt"

	"github.com/piekstra/jira-ticket-cli/api"
	"github.com/piekstra/jira-ticket-cli/internal/cmd/root"
	"github.com/spf13/cobra"
)

// Register registers the transitions commands
func Register(parent *cobra.Command, opts *root.Options) {
	cmd := &cobra.Command{
		Use:     "transitions",
		Aliases: []string{"transition", "tr"},
		Short:   "Manage issue transitions",
		Long:    "Commands for viewing and performing workflow transitions on issues.",
	}

	cmd.AddCommand(newListCmd(opts))
	cmd.AddCommand(newDoCmd(opts))

	parent.AddCommand(cmd)
}

func newListCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "list <issue-key>",
		Short: "List available transitions",
		Long:  "List the available workflow transitions for an issue.",
		Example: `  jira-ticket-cli transitions list MON-1234`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts, args[0])
		},
	}
}

func runList(opts *root.Options, issueKey string) error {
	v := opts.View()

	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	transitions, err := client.GetTransitions(issueKey)
	if err != nil {
		return err
	}

	if len(transitions) == 0 {
		v.Info("No transitions available for %s", issueKey)
		return nil
	}

	if opts.Output == "json" {
		return v.JSON(transitions)
	}

	headers := []string{"ID", "NAME", "TO STATUS"}
	var rows [][]string

	for _, t := range transitions {
		rows = append(rows, []string{t.ID, t.Name, t.To.Name})
	}

	return v.Table(headers, rows)
}

func newDoCmd(opts *root.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "do <issue-key> <transition>",
		Short: "Perform a transition",
		Long:  "Perform a workflow transition on an issue. The transition can be specified by name or ID.",
		Example: `  # Transition by name
  jira-ticket-cli transitions do MON-1234 "In Progress"

  # Transition by ID
  jira-ticket-cli transitions do MON-1234 21`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDo(opts, args[0], args[1])
		},
	}
}

func runDo(opts *root.Options, issueKey, transitionNameOrID string) error {
	v := opts.View()

	client, err := opts.APIClient()
	if err != nil {
		return err
	}

	// Get available transitions
	transitions, err := client.GetTransitions(issueKey)
	if err != nil {
		return err
	}

	// Find the transition
	var transitionID string

	// First try by ID
	for _, t := range transitions {
		if t.ID == transitionNameOrID {
			transitionID = t.ID
			break
		}
	}

	// Then try by name
	if transitionID == "" {
		if t := api.FindTransitionByName(transitions, transitionNameOrID); t != nil {
			transitionID = t.ID
		}
	}

	if transitionID == "" {
		v.Error("Transition '%s' not found", transitionNameOrID)
		v.Info("Available transitions:")
		for _, t := range transitions {
			v.Info("  %s: %s -> %s", t.ID, t.Name, t.To.Name)
		}
		return fmt.Errorf("transition not found: %s", transitionNameOrID)
	}

	if err := client.DoTransition(issueKey, transitionID); err != nil {
		return err
	}

	v.Success("Transitioned %s", issueKey)
	return nil
}
