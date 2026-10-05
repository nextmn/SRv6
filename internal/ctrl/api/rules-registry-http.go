// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package ctrl_api

import (
	"net/http"
)

type RulesRegistryHTTP interface {
	GetRule(w http.ResponseWriter, r *http.Request)
	GetRules(w http.ResponseWriter, r *http.Request)
	DeleteRule(w http.ResponseWriter, r *http.Request)
	EnableRule(w http.ResponseWriter, r *http.Request)
	DisableRule(w http.ResponseWriter, r *http.Request)
	SwitchRule(w http.ResponseWriter, r *http.Request)
	PostRule(w http.ResponseWriter, r *http.Request)
	UpdateAction(w http.ResponseWriter, r *http.Request)
}
