package meta

import (
	. "github.com/protolambda/zrnt/eth2/beacon/components/registry"
	. "github.com/protolambda/zrnt/eth2/core"
)

type ExitMeta interface {
	InitiateValidatorExit(index ValidatorIndex)
}

type BalanceMeta interface {
	GetBalance(index ValidatorIndex) Gwei
	IncreaseBalance(index ValidatorIndex, v Gwei)
	DecreaseBalance(index ValidatorIndex, v Gwei)
}

type RegistrySizeMeta interface {
	IsValidIndex(index ValidatorIndex) bool
	ValidatorCount() uint64
}

type PubkeyMeta interface {
	Pubkey(index ValidatorIndex) BLSPubkey
	ValidatorIndex(pubkey BLSPubkey) (index ValidatorIndex, exists bool)
}

type EffectiveBalanceMeta interface {
	EffectiveBalance(index ValidatorIndex) Gwei
	GetTotalEffectiveBalanceOf(indices []ValidatorIndex) (sum Gwei)
}

type FinalityMeta interface {
	Finalized() Checkpoint
	CurrentJustified() Checkpoint
	PreviousJustified() Checkpoint
}

type AttesterStatusMeta interface {
	GetAttesterStatus(index ValidatorIndex) AttesterStatus
}

type CompactValidatorMeta interface {
	PubkeyMeta
	EffectiveBalanceMeta
	IsSlashed(i ValidatorIndex) bool
}

type StakingMeta interface {
	EffectiveBalanceMeta
	GetTotalActiveEffectiveBalance(epoch Epoch) Gwei
}

type SlashingMeta interface {

}
type ValidatorMeta interface {
	Validator(index ValidatorIndex) *Validator
}

type VersioningMeta interface {
	Slot() Slot
	Epoch() Epoch
	PreviousEpoch() Epoch
	GetDomain(dom BLSDomainType, messageEpoch Epoch) BLSDomain
}

type Eth1Meta interface {
	DepositIndex() DepositIndex
	DepositCount() DepositIndex
	DepositRoot() Root
}

type DepositMeta interface {
	BalanceMeta
	AddNewValidator(pubkey BLSPubkey, withdrawalCreds Root, balance Gwei)
	IncrementDepositIndex()
}

type HeaderMeta interface {
	// Signing root of latest_block_header
	GetLatestBlockRoot() Root
}

type UpdateHeaderMeta interface {
	StoreHeaderData(slot Slot, parent Root, body Root)
	UpdateLatestBlockRoot(stateRoot Root) Root
}

type HistoryUpdateMeta interface {
	SetRecentRoots(slot Slot, blockRoot Root, stateRoot Root)
	UpdateHistoricalRoots()
}

type ProposingMeta interface {
	GetBeaconProposerIndex() ValidatorIndex
}

type ActivationExitMeta interface {
	GetChurnLimit(epoch Epoch) uint64
	ExitQueueEnd(epoch Epoch) Epoch
}

type FullRegistryMeta interface {
	ValidatorMeta
	GetActiveValidatorIndices() RegistryIndices
}

type ShardDeltaMeta interface {
	GetShardDelta(epoch Epoch) Shard
}

type ActiveIndexRootMeta interface {
	GetActiveIndexRoot(epoch Epoch) Root
	// TODO
	// indices := state.Validators.GetActiveValidatorIndices(epoch)
	// ssz.HashTreeRoot(indices, RegistryIndicesSSZ)
}

type CrosslinkTimingMeta interface {
	GetStartShard(epoch Epoch) Shard
	GetCommitteeCount(epoch Epoch) uint64
}

type CrosslinkCommitteeMeta interface {
	CrosslinkTimingMeta
	GetCrosslinkCommittee(epoch Epoch, shard Shard) []ValidatorIndex
}

type CrosslinkMeta interface {
	GetCurrentCrosslink(shard Shard) *Crosslink
	GetPreviousCrosslink(shard Shard) *Crosslink
}

type WinningCrosslinkMeta interface {
	GetWinningCrosslinkAndAttesters(epoch Epoch, shard Shard) (*Crosslink, ValidatorSet)
}

type RandomnessMeta interface {
	GetRandomMix(epoch Epoch) Root
	// epoch + EPOCHS_PER_HISTORICAL_VECTOR - MIN_SEED_LOOKAHEAD) // Avoid underflow
}
