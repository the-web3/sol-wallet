package node

import "math/big"

/*

  "transaction": {
	"message": {
		"accountKeys": [
			{
				"pubkey": "4k7tCcpCq67P1p34MVYe5k6vtNAFbQqVo1xYXtTDCQPD",
				"signer": true,
				"source": "transaction",
				"writable": true
			},
			{
				"pubkey": "CVD2nHU2doMu2Ycqq2V532iQdKnGR7WzzFxHx9NVAerU",
				"signer": false,
				"source": "transaction",
				"writable": true
			},
			{
				"pubkey": "G5segjMzCWnu3zayK3bQAHc7DzthE2fjL1HFhU4hZ37R",
				"signer": false,
				"source": "transaction",
				"writable": true
			},
			{
				"pubkey": "AQqQHiLZi6JyZ2z2BUJxpmFgchyknGR8sfs49c32iZem",
				"signer": false,
				"source": "transaction",
				"writable": true
			},
			{
				"pubkey": "9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin",
				"signer": false,
				"source": "transaction",
				"writable": false
			},
			{
				"pubkey": "11111111111111111111111111111111",
				"signer": false,
				"source": "transaction",
				"writable": false
			}
		],
		"instructions": [
			{
				"accounts": [
					"CVD2nHU2doMu2Ycqq2V532iQdKnGR7WzzFxHx9NVAerU",
					"G5segjMzCWnu3zayK3bQAHc7DzthE2fjL1HFhU4hZ37R",
					"AQqQHiLZi6JyZ2z2BUJxpmFgchyknGR8sfs49c32iZem",
					"G5segjMzCWnu3zayK3bQAHc7DzthE2fjL1HFhU4hZ37R",
					"G5segjMzCWnu3zayK3bQAHc7DzthE2fjL1HFhU4hZ37R"
				],
				"data": "12VeXEVRR",
				"programId": "9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin",
				"stackHeight": null
			},
			{
				"parsed": {
					"info": {
						"destination": "4k7tCcpCq67P1p34MVYe5k6vtNAFbQqVo1xYXtTDCQPD",
						"lamports": 7035,
						"source": "4k7tCcpCq67P1p34MVYe5k6vtNAFbQqVo1xYXtTDCQPD"
					},
					"type": "transfer"
				},
				"program": "system",
				"programId": "11111111111111111111111111111111",
				"stackHeight": null
			}
		],
		"recentBlockhash": "GS7APHdgBjCFuHsEAo2DxvRcCWWnMfAzYvP4iKVv9vDu"
	},
	"signatures": [
		"G6wz1rFZaGRbVUa9qPumYvmhNA3cxYXD8BCgZfztLfaJAAFP3rhQ74uEEza2wSSADBtiLHM5hoFD2jcAnaaYfiT"
	]
},
*/

type Transaction struct {
	Destination string   `json:"destination"`
	Source      string   `json:"source"`
	Lamports    *big.Int `json:"lamports"`
	Type        string   `json:"type"`
	Fee         *big.Int `json:"fee"`
}
