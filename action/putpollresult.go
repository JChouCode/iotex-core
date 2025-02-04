// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package action

import (
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/pkg/errors"

	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	"github.com/iotexproject/iotex-core/pkg/version"
	"github.com/iotexproject/iotex-core/state"
)

// PutPollResult represents put the poll result from gravity chain.
type PutPollResult struct {
	AbstractAction

	height     uint64
	candidates state.CandidateList
}

// NewPutPollResult instantiates a putting poll result action struct.
func NewPutPollResult(
	nonce uint64,
	height uint64,
	candidates state.CandidateList,
) *PutPollResult {
	return &PutPollResult{
		AbstractAction: AbstractAction{
			version:  version.ProtocolVersion,
			nonce:    nonce,
			gasLimit: 0,
			gasPrice: big.NewInt(0),
		},
		height:     height,
		candidates: candidates,
	}
}

// LoadProto converts a proto message into put block action.
func (r *PutPollResult) LoadProto(putPollResultPb *iotextypes.PutPollResult) error {
	if putPollResultPb == nil {
		return errors.New("empty action proto to load")
	}
	if r == nil {
		return errors.New("nil action to load proto")
	}
	*r = PutPollResult{}

	r.height = putPollResultPb.Height

	return r.candidates.LoadProto(putPollResultPb.Candidates)
}

// Proto converts put poll result action into a proto message.
func (r *PutPollResult) Proto() *iotextypes.PutPollResult {
	return &iotextypes.PutPollResult{
		Height:     r.height,
		Candidates: r.candidates.Proto(),
	}
}

// Height returns put poll result height.
func (r *PutPollResult) Height() uint64 { return r.height }

// Candidates returns the list of candidates.
func (r *PutPollResult) Candidates() state.CandidateList { return r.candidates }

// ProducerPublicKey return producer public key.
func (r *PutPollResult) ProducerPublicKey() crypto.PublicKey { return r.SrcPubkey() }

// Serialize returns the byte representation of put poll result action.
func (r *PutPollResult) Serialize() []byte {
	return byteutil.Must(proto.Marshal(r.Proto()))
}

// IntrinsicGas returns the intrinsic gas of a put poll result action
func (r *PutPollResult) IntrinsicGas() (uint64, error) {
	return 0, nil
}

// Cost returns the total cost of a put poll result action
func (r *PutPollResult) Cost() (*big.Int, error) {
	return big.NewInt(0), nil
}
