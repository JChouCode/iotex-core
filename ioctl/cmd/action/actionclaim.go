// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"github.com/spf13/cobra"

	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-core/ioctl/util"
)

// actionClaimCmd represents the action claim command
var actionClaimCmd = &cobra.Command{
	Use:   "claim AMOUNT_IOTX [DATA] [-s SIGNER] [-l GAS_LIMIT] [-p GASPRICE]",
	Short: "Claim rewards from rewarding fund",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		amount, err := util.StringToRau(args[0], util.IotxDecimalNum)
		if err != nil {
			return err
		}
		payload := make([]byte, 0)
		if len(args) == 2 {
			payload = []byte(args[1])
		}
		sender, err := signer()
		if err != nil {
			return err
		}
		gasLimit := gasLimitFlag.Value().(uint64)
		if gasLimit == 0 {
			gasLimit = action.ClaimFromRewardingFundBaseGas +
				action.ClaimFromRewardingFundGasPerByte*uint64(len(payload))
		}
		gasPriceRau, err := gasPriceInRau()
		nonce, err := nonce(sender)
		if err != nil {
			return err
		}
		act := (&action.ClaimFromRewardingFundBuilder{}).SetAmount(amount).SetData(payload).Build()

		return sendAction((&action.EnvelopeBuilder{}).SetNonce(nonce).
			SetGasPrice(gasPriceRau).
			SetGasLimit(gasLimit).
			SetAction(&act).Build(),
			sender,
		)
	},
}

func init() {
	registerWriteCommand(actionClaimCmd)
}
