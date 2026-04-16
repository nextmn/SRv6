// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package tasks

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// HookSingle
type SingleHook struct {
	command *string
	name    string
}

// Creates a new SingleHook
func NewSingleHook(name string, cmd *string) SingleHook {
	return SingleHook{
		name:    name,
		command: cmd,
	}
}

func (h SingleHook) Name() string {
	return h.name
}

// Runs the command of the SingleHook
func (h SingleHook) Run(ctx context.Context) error {
	if h.command == nil {
		// nothing to do
		return nil
	}
	cmd := exec.CommandContext(ctx, *h.command)
	// XXX: We are executing arbitrary (used defined) command here!!
	// We need to preserve environment variables because this is currently
	// the indended way to pass arguments (for example an ip address) to the command.
	// TODO: update the "hook" mechanism to use a reworked version of docker-setup instead
	env := cmd.Environ()
	cmd.Env = []string{}
	for _, e := range env {
		if !strings.HasPrefix("POSTGRES", e) { // filter some env variables that should not be used by routing scripts
			cmd.Env = append(cmd.Env, e)
		}
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s: %w", cmd.Args[0], err)
	}
	return nil
}
