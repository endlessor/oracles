// Copyright (c) 2017 Sweetbridge Stiftung (Sweetbridge Foundation)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"testing"

	"bitbucket.org/sweetbridge/oracles/go-lib/ethereum"
	. "gopkg.in/check.v1"
)

var flagIntegration = flag.Bool("integration", false, "Include integration tests")
var cf ethereum.ContractFactory

func Test(t *testing.T) { TestingT(t) }
func init() {
	setupFlags()
	cf = setupContracts()

	Suite(&PledgeS{})
}