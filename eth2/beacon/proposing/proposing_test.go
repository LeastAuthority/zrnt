package proposing

import (
	"crypto/rand"
	"fmt"
	"os"
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
	const (
		times = 48
		N     = Mi
	)

	var (
		tcs = []distTestCase{
			{
				Name:           "blll",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 4,
			},
			{
				Name:           "lbll",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 4,
				BigOffset:      1,
			},
			{
				Name:           "blll..8",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 8,
			},
			{
				Name:           "lbll..8",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 8,
				BigOffset:      1,
			},
			{
				Name:           "blll..16",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 16,
			},
			{
				Name:           "lbll..16",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 16,
				BigOffset:      1,
			},
			{
				Name:           "blll..32",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 32,
			},
			{
				Name:           "lbll..32",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 32,
				BigOffset:      1,
			},
			{
				Name:           "blll..64",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 64,
			},
			{
				Name:           "lbll..64",
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 64,
				BigOffset:      1,
			},
		}
		results = make(chan distTestResult, times*len(tcs))
	)

	f, err := os.Create("results")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("manage", func(t *testing.T) {
		t.Parallel()
		defer f.Close()

		var (
			ratioAvgs  = make(map[string]float32)
			doneMap    = make(map[string]int)
			res        distTestResult
			resultList []distTestResult
		)

		for i := 0; i < times*len(tcs); i++ {
			res = <-results

			ratioAvgs[res.Name] += float32(res.BigCount()) / float32(res.SmallCountAvg()) / times
			doneMap[res.Name] += 1
			resultList = append(resultList, res)

			if doneMap[res.Name] == times {
				fmt.Fprintf(f, "%s,%d,%d,%f\n", res.Name, res.BigCount(), res.SmallCountAvg(), ratioAvgs[res.Name])
			}
		}

		t.Log(ratioAvgs)
	})

	for _, tc := range tcs {

		tc := tc
		t.Run(tc.String(), func(t *testing.T) {
			for i := 0; i < times; i++ {
				t.Run(fmt.Sprint(i), tc.makeTest(i, results))
			}
		})
	}

}

type distTestResult struct {
	distTestCase

	Round  int
	Counts []int
}

func (tr distTestResult) BigCount() int {
	return tr.Counts[tr.BigOffset]
}

func (tr distTestResult) SmallCountAvg() int {
	return (sum(tr.Counts) - tr.Counts[tr.BigOffset]) / (tr.ValidatorCount - 1)
}

type distTestCase struct {
	Name                     string
	N                        int
	BigBalance, SmallBalance Gwei
	ValidatorCount           int
	BigOffset                int
}

func (tc distTestCase) String() string {
	return tc.Name
}

func (tc distTestCase) makeTest(round int, results chan<- distTestResult) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		var (
			balances       = make([]Gwei, tc.ValidatorCount)
			counts         = make([]int, len(balances))
			validatorCount = ValidatorIndex(len(balances))
			prop           = newFeature(balances)
		)

		// populate balances
		for i := range balances {
			balances[i] = tc.SmallBalance
		}
		balances[tc.BigOffset] = tc.BigBalance

		// prepare logging format strings
		l := len(fmt.Sprintf("%d", tc.N))
		formatN := fmt.Sprintf("N=%%%dd\n", l)
		formatI := fmt.Sprintf("i=%%%dd\n", l)

		fmt.Printf(formatN, tc.N)
		for i := 0; i < tc.N; i++ {
			if i&0xffff == 0 {
				fmt.Printf(formatI, i)
			}

			var seed Root
			rand.Read(seed[:])

			idx := prop.computeProposerIndex(seqVI(0, validatorCount), seed)
			counts[int(idx)]++
		}

		results <- distTestResult{
			distTestCase: tc,
			Round:        round,
			Counts:       counts,
		}

		var (
			smallAvg      = (sum(counts) - counts[tc.BigOffset]) / (len(balances) - 1)
			countFactor   = float32(counts[tc.BigOffset]) / float32(smallAvg)
			depositFactor = float32(balances[1]) / float32(balances[0])
		)

		t.Log("count of big validator:", counts[tc.BigOffset])
		t.Log("average of small validators:", smallAvg)
		t.Log("large validator count is ahead by factor:", countFactor)
		t.Log("large validator deposit is ahead by factor:", depositFactor)
		t.Log("averate iterations between election of single little validator:", tc.N/smallAvg)
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
