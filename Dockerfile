# Copyright Louis Royer and the NextMN contributors. All rights reserved.
# Use of this source code is governed by a MIT-style license that can be
# found in the LICENSE file.
# SPDX-License-Identifier: MIT

FROM golang:1.26.6 AS builder
WORKDIR /src
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -o /usr/local/bin/srv6

FROM alpine:3.24.1
RUN apk add --no-cache iptables iproute2
# TODO: create those files using a new subcommand
COPY etc/iproute2/rt_tables.d/nextmn.conf /etc/iproute2/rt_tables.d/nextmn.conf
COPY etc/iproute2/rt_protos.d/nextmn.conf /etc/iproute2/rt_protos.d/nextmn.conf
# TODO: integrate the following as a configuration option
COPY --chmod=+x <<EOF /usr/local/bin/remove-default-routes.sh
#!/usr/bin/env sh
ip -6 r del default
ip -4 r del default
EOF

COPY --from=builder /usr/local/bin/srv6 /usr/local/bin/srv6
ENTRYPOINT ["srv6"]
CMD ["--help"]
HEALTHCHECK --interval=1m --timeout=1s --retries=3 --start-period=5s --start-interval=100ms \
CMD ["srv6", "healthcheck"]
