// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package tasks_api

import (
	"context"
)

// Task to be run
type TaskUnit interface {
	Name() string
	Run(ctx context.Context) error
}
