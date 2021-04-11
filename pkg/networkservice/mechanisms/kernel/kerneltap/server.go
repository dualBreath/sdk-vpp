// Copyright (c) 2020-2021 Cisco and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build linux

package kerneltap

import (
	"context"

	"git.fd.io/govpp.git/api"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"

	"github.com/networkservicemesh/sdk/pkg/networkservice/utils/metadata"
)

type kernelTapServer struct {
	vppConn api.Connection
}

// NewServer - return a new Server chain element implementing the kernel mechanism with vpp using tapv2
func NewServer(vppConn api.Connection) networkservice.NetworkServiceServer {
	return &kernelTapServer{
		vppConn: vppConn,
	}
}

func (k *kernelTapServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	conn, err := next.Server(ctx).Request(ctx, request)
	if err != nil {
		return nil, err
	}
	if err := create(ctx, conn, k.vppConn, metadata.IsClient(k)); err != nil {
		_, _ = k.Close(ctx, conn)
		return nil, err
	}
	return conn, nil
}

func (k *kernelTapServer) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	if err := del(ctx, conn, k.vppConn, metadata.IsClient(k)); err != nil {
		return nil, err
	}
	return next.Server(ctx).Close(ctx, conn)
}
