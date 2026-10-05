// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package ctrl

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"github.com/nextmn/srv6/internal/database"

	"github.com/nextmn/json-api/jsonapi"
	"github.com/nextmn/json-api/jsonapi/n4tosrv6"

	"github.com/gofrs/uuid/v5"
	"github.com/sirupsen/logrus"
)

// A RulesRegistry contains rules for an headend
type RulesRegistry struct {
	db *database.Database
}

func NewRulesRegistry(db *database.Database) *RulesRegistry {
	return &RulesRegistry{
		db: db,
	}
}

func (rr *RulesRegistry) GetRule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	id := req.PathValue("uuid")
	iduuid, err := uuid.FromString(id)
	if err != nil {
		logrus.WithError(err).Error("Bad UUID")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "bad uuid", Error: err})
		return
	}
	rule, err := rr.db.GetRule(req.Context(), iduuid)
	if err != nil {
		logrus.WithError(err).Error("Could not get rule from database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not get rule from database", Error: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.MarshalWrite(w, rule)
}

func (rr *RulesRegistry) GetRules(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	rules, err := rr.db.GetRules(req.Context())
	if err != nil {
		logrus.WithError(err).Error("Could not get all rules from database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not get all rules from database", Error: err})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.MarshalWrite(w, rules)
}

func (rr *RulesRegistry) DeleteRule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	id := req.PathValue("uuid")
	iduuid, err := uuid.FromString(id)
	if err != nil {
		logrus.WithError(err).Error("Bad UUID")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "bad uuid", Error: err})
		return
	}
	err = rr.db.DeleteRule(req.Context(), iduuid)
	if err != nil {
		logrus.WithError(err).Error("Could not delete rule in the database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not delete rule in the database", Error: err})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rr *RulesRegistry) EnableRule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	id := req.PathValue("uuid")
	iduuid, err := uuid.FromString(id)
	if err != nil {
		logrus.WithError(err).Error("Bad UUID")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "bad uuid", Error: err})
		return
	}
	err = rr.db.EnableRule(req.Context(), iduuid)
	if err != nil {
		logrus.WithError(err).Error("Could not enable rule in the database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not enable rule in the database", Error: err})
		return
		//TODO: check if rule not found
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rr *RulesRegistry) DisableRule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	id := req.PathValue("uuid")
	iduuid, err := uuid.FromString(id)
	if err != nil {
		logrus.WithError(err).Error("Bad UUID")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "bad uuid", Error: err})
		return
	}
	err = rr.db.DisableRule(req.Context(), iduuid)
	if err != nil {
		logrus.WithError(err).Error("Could not disable rule in the database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not disable rule in the database", Error: err})
		return
		//TODO: check if rule not found
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rr *RulesRegistry) SwitchRule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	idEnable := req.PathValue("enable_uuid")
	idDisable := req.PathValue("disable_uuid")
	iduuidEnable, err := uuid.FromString(idEnable)
	if err != nil {
		logrus.WithError(err).Error("Bad UUID")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "bad uuid", Error: err})
		return
	}
	iduuidDisable, err := uuid.FromString(idDisable)
	if err != nil {
		logrus.WithError(err).Error("Bad UUID")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "bad uuid", Error: err})
		return
	}
	err = rr.db.SwitchRule(req.Context(), iduuidEnable, iduuidDisable)
	if err != nil {
		logrus.WithError(err).Error("Could not Switch rule in the database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not switch rule in the database", Error: err})
		return
		//TODO: check if rule not found
	}
	w.WriteHeader(http.StatusNoContent)
}

// Post a new rule
func (rr *RulesRegistry) PostRule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	var rule n4tosrv6.Rule
	if err := json.UnmarshalRead(req.Body, &rule); err != nil {
		logrus.WithError(err).Error("could not deserialize")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not deserialize", Error: err})
		return
	}
	id, err := rr.db.InsertRule(req.Context(), rule)
	if err != nil {
		logrus.WithError(err).Error("Could not insert rule in the database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "failed to insert rule", Error: err})
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/rules/%s", id))
	w.WriteHeader(http.StatusCreated)
	json.MarshalWrite(w, rule)
}

// Update action of a rule
func (rr *RulesRegistry) UpdateAction(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	id_rule := req.PathValue("uuid")
	iduuid_rule, err := uuid.FromString(id_rule)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "bad uuid", Error: err})
		return
	}
	var action n4tosrv6.Action
	if err := json.UnmarshalRead(req.Body, &action); err != nil {
		logrus.WithError(err).Error("could not deserialize")
		w.WriteHeader(http.StatusBadRequest)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not deserialize", Error: err})
		return
	}
	err = rr.db.UpdateAction(req.Context(), iduuid_rule, action)
	if err != nil {
		logrus.WithError(err).Error("Could not update Action for this rule in the database")
		w.WriteHeader(http.StatusInternalServerError)
		json.MarshalWrite(w, jsonapi.MessageWithError{Message: "could not update Action for this rule in the database", Error: err})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
