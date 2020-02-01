package proposing

import (
	"crypto/rand"
	"testing"

	. "github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/meta"
)

func TestProposerDistribution(t *testing.T) {
	const N = 1000

	var (
		balances       = []Gwei{32, 1, 16, 16, 16, 16, 16, 16}
		counts         = make([]int, len(balances))
		prop           = newFeature(balances)
		validatorCount = ValidatorIndex(len(balances))
	)

	for i := 0; i < N; i++ {
		var seed Root
		rand.Read(seed[:])

		idx := prop.computeProposerIndex(seqVI(0, validatorCount), seed)
		counts[int(idx)]++
	}

	t.Log("balances:", balances)
	t.Log("counts of election:", counts)
}

func seqVI(start, end ValidatorIndex) []ValidatorIndex {
	var out []ValidatorIndex
	for i := start; i < end; i++ {
		out = append(out, i)
	}
	return out
}

func newFeature(balances []Gwei) *ProposingFeature {
	type seedKey struct {
		e  Epoch
		dt BLSDomainType
	}

	const curEpoch = 42

	var (
		validatorCount = ValidatorIndex(len(balances))
		seeds          = make(map[seedKey]Root)
	)

	return &ProposingFeature{
		struct {
			meta.Versioning
			meta.BeaconCommittees
			meta.EffectiveBalances
			meta.ActiveIndices
			meta.CommitteeCount
			meta.EpochSeed
		}{
			Versioning: meta.VersioningFake{
				CurrentEpochFunc: func() Epoch {
					return curEpoch
				},
			},
			EffectiveBalances: meta.EffectiveBalancesFake{
				EffectiveBalanceFunc: func(index ValidatorIndex) Gwei {
					return balances[int(index)]
				},
			},
			EpochSeed: meta.EpochSeedFake{
				GetSeedFunc: func(epoch Epoch, domainType BLSDomainType) Root {
					sk := seedKey{e: epoch, dt: domainType}
					seed, ok := seeds[sk]
					if ok {
						return seed
					}

					rand.Read(seed[:])
					seeds[sk] = seed
					return seed
				},
			},
			ActiveIndices: meta.ActiveIndicesFake{
				GetActiveValidatorIndicesFunc: func(epoch Epoch) RegistryIndices {
					return RegistryIndices(seqVI(0, validatorCount))
				},
			},
		},
	}
}
