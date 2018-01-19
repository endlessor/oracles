// Copyright (c) 2017 Sweetbridge Inc.
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

	"bitbucket.org/sweetbridge/oracles/go-lib/log"
	"bitbucket.org/sweetbridge/oracles/go-lib/setup"
	"github.com/go-pg/pg"
	"github.com/robert-zaremba/log15/rollbar"
)

var (
	db     *pg.DB
	logger = log.Root()
)

type mainFlags struct {
	setup.PgFlags
	setup.RollbarFlags
	port *string
}

var flags = mainFlags{PgFlags: setup.NewPgFlags(),
	RollbarFlags: setup.NewRollbarFlags(),
	port:         flag.String("port", "8000", "The HTTP listening port")}

func init() {
	setup.FlagSimpleInit("tge-spreadsheet-etl", "confirmed_payments.csv", flags.Rollbar, flags.PgFlags)
	db = flags.MustConnect()
}

func main() {
	defer rollbar.WaitForRollbar(logger)

	records, err := readRecords(flag.Arg(0))
	checkOK(err)
	checkOK(insertRecords(records))
}

func checkOK(err error) {
	if err != nil {
		logger.Fatal("Bad request", err)
	}
}
