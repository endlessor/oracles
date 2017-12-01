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

	"bitbucket.org/sweetbridge/oracles/go-contracts"
	"bitbucket.org/sweetbridge/oracles/go-lib/ethereum"
	"bitbucket.org/sweetbridge/oracles/go-lib/log"
	"bitbucket.org/sweetbridge/oracles/go-lib/middleware"
	"bitbucket.org/sweetbridge/oracles/go-lib/setup"
	"bitbucket.org/sweetbridge/oracles/go-lib/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/log15/rollbar"
)

type mainFlags struct {
	setup.BaseOracleFlags
	port      *string
	apiKey    *string
	bcyNet    *string
	txTimeout *int
}

func (f mainFlags) Check() error {
	if len(*f.apiKey) != 32 {
		return errstack.NewReq("`-bcy-key` must be a 32 bytes string")
	}
	if len(*f.bcyNet) < 4 {
		return errstack.NewReq("`-bcy-net` must be at least 4 characters long")
	}

	return f.BaseOracleFlags.Check()
}

var (
	flags                   mainFlags
	logger                  = log.Root()
	client                  *ethclient.Client
	cf                      ethereum.ContractFactory
	forwarderFactory        *contracts.ForwarderFactory
	forwarderFactoryAddress common.Address
)

const serviceName = "payment-forwarder"

func setupFlags() {
	flags = mainFlags{
		BaseOracleFlags: setup.NewBaseOracleFlags(),
		port:            flag.String("port", "8000", "The HTTP listening port"),
		apiKey:          flag.String("bcy-key", "", "BlockCypher API Token [required]"),
		bcyNet:          flag.String("bcy-net", "", "BlockCypher network (main or test) [required]"),
		txTimeout:       flag.Int("tx-timeout", 600, "how many seconds should the daemon wait for the transaction receipt?"),
	}
	setup.FlagSimpleInit(serviceName, *flags.Rollbar, flags)
}

func setupContracts() {
	var err error
	client, cf = flags.MustEthFactory()
	forwarderFactory, forwarderFactoryAddress, err = cf.GetForwarderFactory()
	utils.Assert(err, "Can't instantiate ForwardFactory contract")
}

func main() {
	setupFlags()
	defer rollbar.WaitForRollbar(logger)
	initBcyAPI()
	setupContracts()

	handler, r := middleware.StdRouter(serviceName)
	r.Post("/btc", handleBtcCreate)
	r.Post("/eth", handleEthCreate)
	setup.HTTPServer(serviceName, *flags.port, handler)
}
