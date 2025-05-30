// Copyright 2023 Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package tasks

import (
	"context"

	"github.com/nextmn/srv6/internal/constants"
	"github.com/nextmn/srv6/internal/iproute2"
)

// TaskBlackhole
type TaskBlackhole struct {
	WithName
	WithState
	table iproute2.Table
}

// Create a new TaskBlackhole
func NewTaskBlackhole(name string, table_name string) *TaskBlackhole {
	return &TaskBlackhole{
		WithName:  NewName(name),
		WithState: NewState(),
		table:     iproute2.NewTable(table_name, constants.RT_PROTO_NEXTMN),
	}
}

// Create blackhole
func (t *TaskBlackhole) RunInit(ctx context.Context) error {
	if err := t.table.AddDefaultRoutesBlackhole(); err != nil {
		return err
	}
	t.state = true
	return nil
}

// Delete blackhole
func (t *TaskBlackhole) RunExit() error {
	if err := t.table.DelDefaultRoutesBlackhole(); err != nil {
		return err
	}
	t.state = false
	return nil
}
