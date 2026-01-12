// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package iproute2

import (
	"context"
	"fmt"
	"net/netip"
	"os/exec"
)

// Run ip command
func runIP(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "ip", args...)
	cmd.Env = []string{}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s: %w", cmd.Args, err)
	}
	return nil
}

// Run iptables command
func runIPTables(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "iptables", args...)
	cmd.Env = []string{}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s: %w", cmd.Args, err)
	}
	return nil
}

// Run ip6tables command
func runIP6Tables(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, "ip6tables", args...)
	cmd.Env = []string{}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s: %w", cmd.Args, err)
	}
	return nil
}

func IPSrSetSourceAddress(ctx context.Context, address netip.Addr) error {
	return runIP(ctx, "sr", "tunsrc", "set", address.String())
}
