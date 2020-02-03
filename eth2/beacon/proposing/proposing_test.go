package proposing

import (
	"crypto/rand"
	"fmt"
	"testing"

	. "github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/meta"
)

func eth2gwei(eths []uint64) []Gwei {
	var out = make([]Gwei, len(eths))

	for i := range eths {
		out[i] = Gwei(eths[i]) * 1000 * 1000 * 1000
	}

	return out
}

const (
	K = 1000
	M = K * K
	G = M * K

	Ki = 1024
	Mi = Ki * Ki
	Gi = Mi * Ki
)

func TestProposerDistribution(t *testing.T) {
	const times = 32
	var ratios = make(chan float32, times)

	for i := 0; i < times; i++ {
		for _, tc := range []distTestCase{
			{name: fmt.Sprintf("N:1Mi-bBlnc:32Gi-lBlnc:32-i:%d", i), n: 1 * Mi, bigBlnc: 32 * Gi, lilBlnc: 16 * Gi},
		} {
			t.Run(tc.String(), tc.makeTest(ratios))
		}
	}

	var ratioAvg float32
	for i := 0; i < times; i++ {
		ratioAvg += <-ratios
	}
	ratioAvg = ratioAvg / times

	t.Logf("The 32 ETH validator gets elected %f times as often as any of the 16 ETH validators (expected: 32/16=2).", ratioAvg)
}

type distTestCase struct {
	name             string
	n                int
	bigBlnc, lilBlnc Gwei
}

func (tc distTestCase) String() string {
	return tc.name
}

func (tc distTestCase) makeTest(ratios chan<- float32) func(*testing.T) {
	return func(t *testing.T) {
		//t.Parallel()

		var (
			N              = tc.n
			lilBlnc        = tc.lilBlnc
			bigBlnc        = tc.bigBlnc
			balances       = []Gwei{bigBlnc, lilBlnc, lilBlnc, lilBlnc, lilBlnc, lilBlnc, lilBlnc, lilBlnc}
			counts         = make([]int, len(balances))
			prop           = newFeature(balances)
			validatorCount = ValidatorIndex(len(balances))
		)

		l := len(fmt.Sprintf("%d", N))
		formatN := fmt.Sprintf("N=%%%dd\n", l)
		formatI := fmt.Sprintf("i=%%%dd\n", l)
		fmt.Printf(formatN, N)
		for i := 0; i < N; i++ {
			if i&0xffff == 0 {
				fmt.Printf(formatI, i)
			}

			var seed Root
			rand.Read(seed[:])

			idx := prop.computeProposerIndex(seqVI(0, validatorCount), seed)
			counts[int(idx)]++
		}

		t.Log("balances:", balances)
		t.Log("counts of election:", counts)

		var (
			smallAvg      = sum(counts[1:]) / 7
			countFactor   = float32(counts[0]) / float32(smallAvg)
			depositFactor = float32(balances[0]) / float32(balances[1])
		)

		ratios <- countFactor

		t.Log("average of small validators:", smallAvg)
		t.Log("large validator count is ahead by factor:", countFactor)
		t.Log("large validator deposit is ahead by factor:", depositFactor)
		t.Log("averate iterations between election of single little validator:", N/smallAvg)
	}
}

func sum(is []int) int {
	var sum int
	for _, i := range is {
		sum += i
	}
	return sum
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
