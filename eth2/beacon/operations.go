package beacon

import (
	. "github.com/protolambda/zrnt/eth2/beacon/components/attestations"
	. "github.com/protolambda/zrnt/eth2/beacon/components/deposits"
	. "github.com/protolambda/zrnt/eth2/beacon/components/exits"
	. "github.com/protolambda/zrnt/eth2/beacon/components/slashings"
	. "github.com/protolambda/zrnt/eth2/beacon/components/transfers"
	. "github.com/protolambda/zrnt/eth2/core"
)

type ProposerSlashings []ProposerSlashing

func (_ *ProposerSlashings) Limit() uint32 {
	return MAX_PROPOSER_SLASHINGS
}

type AttesterSlashings []AttesterSlashing

func (_ *AttesterSlashings) Limit() uint32 {
	return MAX_ATTESTER_SLASHINGS
}

type Attestations []Attestation

func (_ *Attestations) Limit() uint32 {
	return MAX_ATTESTATIONS
}

type Deposits []Deposit

func (_ *Deposits) Limit() uint32 {
	return MAX_DEPOSITS
}

type Transfers []Transfer

func (_ *Transfers) Limit() uint32 {
	return MAX_TRANSFERS
}

type VoluntaryExits []VoluntaryExit

func (_ *VoluntaryExits) Limit() uint32 {
	return MAX_VOLUNTARY_EXITS
}
