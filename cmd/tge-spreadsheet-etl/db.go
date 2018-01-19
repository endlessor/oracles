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
	"strings"
	"time"

	"bitbucket.org/sweetbridge/oracles/go-lib/liquidity"
	"github.com/robert-zaremba/errstack"
	"github.com/robert-zaremba/go-pgt"
)

// DirectBuy represents the swc direct buy for the next distribution.
type DirectBuy struct {
	ID        pgt.UUID           `sql:"direct_buy_id,pk"`
	UserID    pgt.UUID           `sql:"individual_id"`
	Email     string             `sql:"email"`
	TrancheID uint64             `sql:"tranche_id,notnull"`
	AmountOut float64            `sql:"amount_out,notnull"`
	AmountIn  float64            `sql:"amount_in,notnull"`
	Currency  liquidity.Currency `sql:"currency_id"`
	UsdRate   float64            `sql:"usd_rate,notnull"`
	SenderID  string             `sql:"sender_id"`

	TransactionDate time.Time `sql:"transaction_date,notnull"`
	CreatedAt       time.Time `sql:"created_at,notnull"`
	UpdatedAt       time.Time `sql:"updated_at,notnull"`
}

// missing user is not reported as an error!
func findUserByEmail(r Record) (pgt.UUID, errstack.E) {
	var user pgt.UUID
	_, err := db.QueryOne(&user, "SELECT id FROM individual WHERE email_address = ?", r.Email)
	if err != nil {
		if strings.Index(err.Error(), "no rows in") < 0 {
			return user, errstack.WrapAsInf(err, "can't query DB")
		}
		logger.Error("Can't find user", "email", r.Email, "name", r.FullName)
	}
	return user, nil
}

func createDBRecords(records []Record) ([]DirectBuy, errstack.E) {
	now := time.Now()
	var dbs []DirectBuy
	for _, r := range records {
		user, err := findUserByEmail(r)
		if err != nil {
			return dbs, err
		}
		d := DirectBuy{
			pgt.RandomUUID(), user, r.Email, r.TrancheID, r.AmountSWC,
			r.AmountIn, r.Currency, r.UsdRate, r.SenderID, r.Timestamp, now, now}
		dbs = append(dbs, d)
	}
	return dbs, nil
}

func insertRecords(records []Record) errstack.E {
	dbs, err := createDBRecords(records)
	if err != nil {
		return err
	}
	logger.Info("All direct_buys created. Inserting into DB")
	return errstack.WrapAsInf(db.Insert(&dbs), "DB direct_buy insert")
	// for i := range dbs {
	// 	err = errstack.WrapAsInf(db.Insert(&dbs[i]), "DB direct_buy insert")
	// 	if err != nil {
	// 		fmt.Println(dbs[i])
	// 		return err
	// 	}
	// }
	// return nil
}
