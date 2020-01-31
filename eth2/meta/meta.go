package meta

import (
	"github.com/protolambda/zrnt/eth2/beacon/validator"
	. "github.com/protolambda/zrnt/eth2/core"
)

type Exits interface {
	InitiateValidatorExit(currentEpoch Epoch, index ValidatorIndex)
}

type Balance interface {
	GetBalance(index ValidatorIndex) Gwei
	IncreaseBalance(index ValidatorIndex, v Gwei)
	DecreaseBalance(index ValidatorIndex, v Gwei)
}

type BalanceDeltas interface {
	ApplyDeltas(deltas *Deltas)
}

type AttestationDeltas interface {
	AttestationDeltas() *Deltas
}

type RegistrySize interface {
	IsValidIndex(index ValidatorIndex) bool
	ValidatorCount() uint64
}

type Pubkeys interface {
	Pubkey(index ValidatorIndex) BLSPubkey
	ValidatorIndex(pubkey BLSPubkey) (index ValidatorIndex, exists bool)
}

type EffectiveBalances interface {
	EffectiveBalance(index ValidatorIndex) Gwei
	SumEffectiveBalanceOf(indices []ValidatorIndex) (sum Gwei)
}

type EffectiveBalancesFake struct {
	EffectiveBalanceFunc      func(index ValidatorIndex) Gwei
	SumEffectiveBalanceOfFunc func(indices []ValidatorIndex) (sum Gwei)
}

func (e EffectiveBalancesFake) EffectiveBalance(index ValidatorIndex) Gwei {
	return e.EffectiveBalanceFunc(index)
}

func (e EffectiveBalancesFake) SumEffectiveBalanceOf(indices []ValidatorIndex) Gwei {
	return e.SumEffectiveBalanceOfFunc(indices)
}

type EffectiveBalancesUpdate interface {
	UpdateEffectiveBalances()
}

type Finality interface {
	Finalized() Checkpoint
	CurrentJustified() Checkpoint
	PreviousJustified() Checkpoint
}

type Justification interface {
	Justify(checkpoint Checkpoint)
}

type EpochAttestations interface {
	RotateEpochAttestations()
}

type AttesterStatuses interface {
	GetAttesterStatuses() []AttesterStatus
}

type SlashedIndices interface {
	IsSlashed(i ValidatorIndex) bool
	FilterUnslashed(indices []ValidatorIndex) []ValidatorIndex
}

type CompactCommittees interface {
	Pubkeys
	EffectiveBalances
	SlashedIndices
	GetCompactCommitteesRoot(epoch Epoch) Root
}

type Staking interface {
	// Staked = Active effective balance
	GetTotalStake() Gwei
	GetAttestersStake(statuses []AttesterStatus, mask AttesterFlag) Gwei
}

type Slashing interface {
	GetIndicesToSlash(withdrawal Epoch) (out []ValidatorIndex)
}

type SlashingHistory interface {
	ResetSlashings(epoch Epoch)
}

type Slasher interface {
	SlashValidator(slashedIndex ValidatorIndex, whistleblowerIndex *ValidatorIndex)
}

type Validators interface {
	Validator(index ValidatorIndex) *validator.Validator
}

type Versioning interface {
	CurrentSlot() Slot
	CurrentEpoch() Epoch
	PreviousEpoch() Epoch
	CurrentVersion() Version
	GetDomain(dom BLSDomainType, messageEpoch Epoch) BLSDomain
}

type VersioningFake struct {
	CurrentSlotFunc    func() Slot
	CurrentEpochFunc   func() Epoch
	PreviousEpochFunc  func() Epoch
	CurrentVersionFunc func() Version
	GetDomainFunc      func(dom BLSDomainType, messageEpoch Epoch) BLSDomain
}

func (f VersioningFake) CurrentSlot() Slot {
	return f.CurrentSlotFunc()
}

func (f VersioningFake) CurrentEpoch() Epoch {
	return f.CurrentEpochFunc()
}

func (f VersioningFake) CurrentVersion() Version {
	return f.CurrentVersionFunc()
}

func (f VersioningFake) PreviousEpoch() Epoch {
	return f.PreviousEpochFunc()
}

func (f VersioningFake) GetDomain(dom BLSDomainType, messageEpoch Epoch) BLSDomain {
	return f.GetDomainFunc(dom, messageEpoch)
}

type Eth1Voting interface {
	ResetEth1Votes()
}

type Deposits interface {
	DepIndex() DepositIndex
	DepCount() DepositIndex
	DepRoot() Root
}

type Onboarding interface {
	AddNewValidator(pubkey BLSPubkey, withdrawalCreds Root, balance Gwei)
}

type Depositing interface {
	IncrementDepositIndex()
}

type LatestHeader interface {
	// Signing root of latest_block_header
	GetLatestBlockRoot() Root
}

type LatestHeaderUpdate interface {
	UpdateLatestBlockRoot(stateRoot Root) Root
}

type History interface {
	GetBlockRootAtSlot(slot Slot) Root
	GetBlockRoot(epoch Epoch) Root
}

type HistoryUpdate interface {
	SetRecentRoots(slot Slot, blockRoot Root, stateRoot Root)
	UpdateStateRoot(root Root)
	UpdateHistoricalRoots()
}

type EpochSeed interface {
	// Retrieve the seed for beacon proposer indices.
	GetSeed(epoch Epoch, domainType BLSDomainType) Root
}

type EpochSeedFake struct {
	GetSeedFunc func(epoch Epoch, domainType BLSDomainType) Root
}

func (f EpochSeedFake) GetSeed(epoch Epoch, domainType BLSDomainType) Root {
	return f.GetSeedFunc(epoch, domainType)
}

type Proposers interface {
	GetBeaconProposerIndex(slot Slot) ValidatorIndex
}

type ActivationExit interface {
	GetChurnLimit(epoch Epoch) uint64
	ExitQueueEnd(epoch Epoch) Epoch
}

type ActivationQeueue interface {
	ProcessActivationQueue(activationEpoch Epoch, currentEpoch Epoch)
}

type ActiveValidatorCount interface {
	GetActiveValidatorCount(epoch Epoch) uint64
}

type ValidatorEpochData interface {
	WithdrawableEpoch(index ValidatorIndex) Epoch
}

type ActiveIndices interface {
	IsActive(index ValidatorIndex, epoch Epoch) bool
	GetActiveValidatorIndices(epoch Epoch) RegistryIndices
	ComputeActiveIndexRoot(epoch Epoch) Root
}

type ActiveIndicesFake struct {
	IsActiveFunc                  func(index ValidatorIndex, epoch Epoch) bool
	GetActiveValidatorIndicesFunc func(epoch Epoch) RegistryIndices
	ComputeActiveIndexRootFunc    func(epoch Epoch) Root
}

func (f ActiveIndicesFake) IsActive(index ValidatorIndex, epoch Epoch) bool {
	return f.IsActiveFunc(index, epoch)
}

func (f ActiveIndicesFake) GetActiveValidatorIndices(epoch Epoch) RegistryIndices {
	return f.GetActiveValidatorIndicesFunc(epoch)
}

func (f ActiveIndicesFake) ComputeActiveIndexRoot(epoch Epoch) Root {
	return f.ComputeActiveIndexRootFunc(epoch)
}

type CommitteeCount interface {
	GetCommitteeCountAtSlot(slot Slot) uint64
}

type CommitteeCountFake struct {
	GetCommitteeCountAtSlotFunc func(slot Slot) uint64
}

func (f CommitteeCountFake) GetCommitteeCountAtSlot(slot Slot) uint64 {
	return f.GetCommitteeCountAtSlotFunc(slot)
}

type BeaconCommittees interface {
	GetBeaconCommittee(slot Slot, index CommitteeIndex) []ValidatorIndex
}

type BeaconCommitteesFake struct {
	GetBeaconCommitteeFunc func(slot Slot, index CommitteeIndex) []ValidatorIndex
}

func (f BeaconCommitteesFake) GetBeaconCommittee(slot Slot, index CommitteeIndex) []ValidatorIndex {
	return f.GetBeaconCommitteeFunc(slot, index)
}

type Randao interface {
	PrepareRandao(epoch Epoch)
}

type Randomness interface {
	GetRandomMix(epoch Epoch) Root
}
