// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package config

import (
	"net"
	"slices"
	"strings"
)

type Bsid struct {
	BsidPrefix   *string  `yaml:"bsid-prefix,omitempty"`
	SegmentsList []string `yaml:"segments-list"`
}

func (a *Bsid) ToIPRoute2() string {
	return strings.Join(a.SegmentsList[:], ",")
}

func (a *Bsid) ReverseSegmentsList() []net.IP {
	res := []net.IP{}
	for _, v := range slices.Backward(a.SegmentsList) {
		ip := net.ParseIP(v)
		res = append(res, ip)
	}
	return res
}
