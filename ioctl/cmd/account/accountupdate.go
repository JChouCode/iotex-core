// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package account

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/iotexproject/go-pkgs/hash"
	"github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-core/ioctl/cmd/config"
	"github.com/iotexproject/iotex-core/ioctl/util"
	"github.com/iotexproject/iotex-core/pkg/log"
)

// accountUpdateCmd represents the account update command
var accountUpdateCmd = &cobra.Command{
	Use:   "update [ALIAS|ADDRESS]",
	Short: "Update password for IoTeX account",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		output, err := accountUpdate(args)
		if err == nil {
			fmt.Println(output)
		}
		return err
	},
}

func accountUpdate(args []string) (string, error) {
	account, err := util.GetAddress(args)
	if err != nil {
		return "", err
	}
	address, err := address.FromString(account)
	if err != nil {
		log.L().Error("failed to convert string into address", zap.Error(err))
		return "", err
	}
	// find the keystore and update
	ks := keystore.NewKeyStore(config.ReadConfig.Wallet,
		keystore.StandardScryptN, keystore.StandardScryptP)
	for _, v := range ks.Accounts() {
		if bytes.Equal(address.Bytes(), v.Address.Bytes()) {
			fmt.Printf("#%s: Enter current password\n", account)
			currentPassword, err := util.ReadSecretFromStdin()
			if err != nil {
				log.L().Error("failed to get current password", zap.Error(err))
				return "", err
			}
			_, err = ks.SignHashWithPassphrase(v, currentPassword, hash.ZeroHash256[:])
			if err != nil {
				return "", err
			}
			fmt.Printf("#%s: Enter new password\n", account)
			password, err := util.ReadSecretFromStdin()
			if err != nil {
				log.L().Error("failed to get password", zap.Error(err))
				return "", err
			}
			fmt.Printf("#%s: Enter new password again\n", account)
			passwordAgain, err := util.ReadSecretFromStdin()
			if err != nil {
				log.L().Error("failed to get password", zap.Error(err))
				return "", err
			}
			if password != passwordAgain {
				return "", ErrPasswdNotMatch
			}
			if err := ks.Update(v, currentPassword, password); err != nil {
				log.L().Error("failed to update keystore", zap.Error(err))
				return "", err
			}
			return fmt.Sprintf("Account #%s has been updated.", account), nil
		}
	}
	return "", fmt.Errorf("account #%s not found", account)
}
