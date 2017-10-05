package main

import (
	"fmt"

	contracts "bitbucket.org/sweetbridge/oracles/go-contracts"
	"bitbucket.org/sweetbridge/oracles/go-lib/ethereum"
	"bitbucket.org/sweetbridge/oracles/go-lib/utils"
	"github.com/robert-zaremba/errstack"
)

func distributeSWC(records []Record) {
	_, cf := flags.MustEthFactory()
	swcC, addr, err := cf.GetSWC()
	utils.Assert(err, "Can't instantiate SWT contract")
	logger.Debug("Contract address", "swc", addr.Hex())

	checkOK(checkSWCbalance(records, swcC))
	// TODO - whitelists are failing
	// checkOK(checkWhitelist(records, cf)
	checkOK(transferSWC(records, swcC, cf))
}

func transferSWC(records []Record, swcC *contracts.SweetToken, cf ethereum.ContractFactory) errstack.E {
	if *flags.dryRun {
		logger.Debug("Dry run. Stopping execution.")
		return nil
	}

	txo := cf.Txo()
	for _, r := range records {
		wei := ethereum.ToWei(r.Amount)
		logger.Debug("Transfering", "amount.wei", wei, "dest", r.Address.Hex(),
			"nonce", txo.Nonce)
		tx, err := swcC.Transfer(txo, r.Address, wei)
		if err != nil {
			logger.Error("Can't transfer TOKEN", err)
			break
		} else {
			ethereum.LogTx("Transferred", tx)
			ethereum.IncTxoNonce(txo, tx)
			logger.Debug(">>>> nonce after", "txo", txo.Nonce, "tx", tx.Nonce())
		}
	}
	return nil
}

// TODO -- tests are still failing
func checkWhitelist(records []Record, cf ethereum.ContractFactory) errstack.E {
	swcLogic, _, err := cf.GetSWClogic()
	utils.Assert(err, "Can't instantiate SWTlogic contract")

	var listName = [32]byte{}
	copy(listName[:], []byte("whitelist"))
	logger.Debug("Constructing whitelist name", "bytes", listName)
	txo := cf.Txo()
	if ok, err := swcLogic.ListExists(nil, listName); err != nil {
		return errstack.WrapAsInf(err, "Can't check if whitelist exists")
	} else if !ok {
		tx, err := swcLogic.AddWhiteList(txo, listName)
		if err != nil {
			return errstack.WrapAsInf(err, "Can't create SWC whitelist")
		}
		ethereum.LogTx("SWC whitelist created", tx)
	}

	for _, r := range records {
		if ok, err := swcLogic.WhiteLists(nil, r.Address, listName); err != nil {
			return errstack.WrapAsInfF(err, "Can't check if user %s is in the whitelist",
				r.Address)
		} else if ok {
			continue
		}
		tx, err := swcLogic.AddToWhiteList(txo, listName, r.Address)
		if err != nil {
			return errstack.WrapAsInfF(err, "Can't add user %s to the whitelist", r.Address)
		}
		ethereum.LogTx(fmt.Sprintf("User %s added to the whitelist", r.Address), tx)
	}
	return nil

}

func checkSWCbalance(records []Record, token *contracts.SweetToken) errstack.E {
	var total uint64
	for _, r := range records {
		total += r.Amount
	}
	totalWei := ethereum.ToWei(total)
	k := ethereum.MustReadKeySimple(*flags.PkFile)
	logger.Debug("SWC distribution account holder", "address", k.Address.Hex())
	balance, err := token.BalanceOf(nil, k.Address)
	if err != nil {
		return errstack.WrapAsInf(err, "Can't check SWC balance")
	}
	if balance.Cmp(totalWei) < 0 {
		bInt := ethereum.WeiToInt(balance)
		return errstack.NewReqF("Not enough funds in the source account = %v, min_expected=%v",
			bInt, total)
	}
	logger.Debug("Distribution account balance", "swc.wei", balance.String(),
		"required", totalWei.String())
	return nil
}
