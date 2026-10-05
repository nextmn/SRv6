// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package tasks

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"time"

	"github.com/nextmn/json-api/healthcheck"
	"github.com/nextmn/logrus-formatter/httplog"
	app_api "github.com/nextmn/srv6/internal/app/api"
	"github.com/nextmn/srv6/internal/ctrl"
	ctrl_api "github.com/nextmn/srv6/internal/ctrl/api"

	"github.com/sirupsen/logrus"
)

// HttpServerTask starts an http server
type HttpServerTask struct {
	WithName
	WithState
	srv               *http.Server
	httpAddr          netip.AddrPort
	rulesRegistryHTTP ctrl_api.RulesRegistryHTTP
	setupRegistry     app_api.Registry
}

// Create a new HttpServerTask
func NewHttpServerTask(name string, httpAddr netip.AddrPort, setupRegistry app_api.Registry) *HttpServerTask {
	return &HttpServerTask{
		WithName:          NewName(name),
		WithState:         NewState(),
		srv:               nil,
		httpAddr:          httpAddr,
		rulesRegistryHTTP: nil,
		setupRegistry:     setupRegistry,
	}
}

// Init
func (t *HttpServerTask) RunInit(ctx context.Context) error {
	if t.setupRegistry == nil {
		return fmt.Errorf("Registry is nil")
	}
	db, ok := t.setupRegistry.DB()
	if !ok {
		return fmt.Errorf("DB is not in Registry")
	}
	rr := ctrl.NewRulesRegistry(db)
	t.rulesRegistryHTTP = rr
	r := http.NewServeMux()
	r.HandleFunc("GET /status", t.Status)
	r.HandleFunc("POST /rules", t.rulesRegistryHTTP.PostRule)
	r.HandleFunc("GET /rules/{uuid}", t.rulesRegistryHTTP.GetRule)
	r.HandleFunc("GET /rules", t.rulesRegistryHTTP.GetRules)
	r.HandleFunc("PATCH /rules/{uuid}/enable", t.rulesRegistryHTTP.EnableRule)
	r.HandleFunc("PATCH /rules/{uuid}/disable", t.rulesRegistryHTTP.DisableRule)
	r.HandleFunc("PATCH /rules/switch/{enable_uuid}/{disable_uuid}", t.rulesRegistryHTTP.SwitchRule)
	r.HandleFunc("DELETE /rules/{uuid}", t.rulesRegistryHTTP.DeleteRule)
	r.HandleFunc("PATCH /rules/{uuid}/update-action", t.rulesRegistryHTTP.UpdateAction)
	logger := httplog.NewRequestLoggerMiddleware(r)
	t.srv = &http.Server{
		Addr:    t.httpAddr.String(),
		Handler: logger,
	}

	l, err := net.Listen("tcp", t.srv.Addr)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		if err := t.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Error("HTTP Server error")
		}
	}(l)
	t.state = true
	return nil
}

func (t *HttpServerTask) Status(w http.ResponseWriter, req *http.Request) {
	status := healthcheck.Status{
		Ready: true,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	json.MarshalWrite(w, status)
}

// Exit
func (t *HttpServerTask) RunExit(ctx context.Context) error {
	t.state = false
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second) // context.Background() is already Done()
	defer cancel()
	if err := t.srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Info("HTTP Server Shutdown")
	}
	return nil
}
