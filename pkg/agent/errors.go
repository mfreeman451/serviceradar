/*
 * Copyright 2025 Carver Automation Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agent

import "errors"

var (
	errInvalidPort          = errors.New("invalid port")
	errDetailsRequiredPorts = errors.New("details field is required for port checks")
	errDetailsRequiredGRPC  = errors.New("details field is required for gRPC checks")
	errDetailsRequiredSNMP  = errors.New("details field is required for SNMP checks")
	errInvalidDetailsFormat = errors.New("invalid details format: expected 'host:port'")
	errSNMPServiceUnhealthy = errors.New("SNMP service reported unhealthy")
	errInvalidDuration      = errors.New("invalid duration")
)
