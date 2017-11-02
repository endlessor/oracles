package swcq

import (
	"encoding/json"
	"math/big"

	"bitbucket.org/sweetbridge/oracles/go-lib/liquidity"
	. "gopkg.in/check.v1"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var jsonLogSWCqueueDirectPledge = []byte(`{
  "anonymous": false,
  "inputs": [{
      "indexed": false, "name": "who", "type": "address"
    }, {
      "indexed": false, "name": "wad", "type": "uint128"
    }, {
      "indexed": false, "name": "currency", "type": "bytes3"
  }],
  "name": "LogSWCqueueDirectPledge",
  "type": "event"
}`)

type EventSuite struct {
	eventDirectPledge abi.Event
}

func (suite *EventSuite) SetUpSuite(c *C) {
	err := json.Unmarshal(jsonLogSWCqueueDirectPledge, &suite.eventDirectPledge)
	c.Assert(err, IsNil)
}

func (suite EventSuite) TestLogSWCqueueDirectPledge(c *C) {
	// event LogSWCqueueDirectPledge(address who, uint128 wad, bytes3 currency);
	var data = `{
      "address":"0xfdada2f6dfdd969b5f9297f07f3622e4dfe462d6",
      "topics":["0x3a7f5663aac61ec71b13259e6a35978a2fba9256cd41bfacd1840a8b1c6bdd43"],
      "data":"0x00000000000000000000000000ce0d46d924cc8437c806721496599fc3ffa2680000000000000000000000000000000000000000000000008ac7230489e800007573640000000000000000000000000000000000000000000000000000000000",
      "blockNumber":"0x13b45",
      "transactionHash":"0xb0b1d0d020b3962f7246cb2978a02d2755a69eb163e6e544a54cc573738f428f",
      "transactionIndex":"0x1",
      "blockHash":"0xe3483334b3ee3a97bef4d959b5798881e82ecbf724f51e43661942b8d67e2d61",
      "logIndex":"0x1",
      "removed":false}`

	var log = new(types.Log)
	err := log.UnmarshalJSON([]byte(data))
	c.Assert(err, IsNil)

	var expected = EventDirectPledge{
		Who:      common.HexToAddress("0x00Ce0d46d924CC8437c806721496599FC3FFA268"),
		Wad:      big.NewInt(0),
		Currency: liquidity.Currency{'u', 's', 'd'},
	}
	expected.Wad.SetString("10000000000000000000", 10)

	var o EventDirectPledge
	err = o.Unmarshal(*log)
	c.Assert(err, IsNil)

	c.Check(o.Who, Equals, expected.Who)
	c.Check(o.Wad.String(), Equals, expected.Wad.String())
	c.Check(o.Currency, Equals, expected.Currency)
}
