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

package directbuy

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-pg/pg"
	"github.com/robert-zaremba/errstack"
)

// CreateDirectBuyReport creates aggregated SWC distribution report of DirectBuys from given
// tranche id. `title` is used for the report field. TrancheID is automatically appended to
// the title using "_" as the separator
func CreateDirectBuyReport(title, trancheID string, db *pg.DB) ([]*ReportRecord, errstack.E) {
	ds, err := GetPendingDirectBuys(trancheID, db)
	if err != nil {
		return nil, err
	}
	logger.Debug("Aggregating pending direct buy", "tranche_id", trancheID, "num_records", len(ds))
	title += "_" + trancheID
	var byUser = map[string]*ReportRecord{}
	for _, d := range ds {
		if d.Status != StatusPending {
			logger.Warn("Obtained directbuy with not pending status",
				"id", d.ID, "status", d.Status)
			continue
		}
		addr, err := FindDistributionAccount(d.UserID, db)
		if err != nil {
			if !err.IsReq() {
				return nil, err
			}
			logger.Warn("DirectBuy with unmatched individual (ID, distribution account)", err,
				"direct_buy_id", d.ID, "user_email", d.Email, "individual_id", d.UserID,
				"distribution_account", addr.String())
			continue
		}
		if sr, ok := byUser[d.UserID.String()]; !ok {
			byUser[d.UserID.String()] = &ReportRecord{
				title,
				addr,
				d.AmountOut,
				fmt.Sprintf("email: %s; individual_id: %v", d.Email, d.UserID),
				false}
		} else {
			sr.Amount += d.AmountOut
		}
		// if *flags.setPendingStatus {
		// 	d.Status = directbuy.StatusPending
		// 	if err := directbuy.UpdateStatus(d.ID, d.Status, db); err != nil {
		// 		return err
		// 	}
		// }
	}
	var records = make([]*ReportRecord, len(byUser))
	var i = 0
	for _, v := range byUser {
		records[i] = v
		i++
	}
	return records, nil
}

// ReportRecord represent SWC distribution record
// aligned to swc-distributor
type ReportRecord struct {
	List    string
	Address common.Address
	Amount  float64
	Comment string
	Done    bool
}
