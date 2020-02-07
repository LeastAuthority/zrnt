package proposing

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	. "github.com/protolambda/zrnt/eth2/core"
	"github.com/protolambda/zrnt/eth2/meta"

	"github.com/montanaflynn/stats"
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
		times = 32
		N     = 64 * Ki
	)

	var tcs []distTestCase

	// small set sizes
	for i := uint(3); i < 6; i++ {
		validatorCount := 1 << i
		tcs = append(tcs,
			distTestCase{
				Name:           fmt.Sprintf("%d-0-%d", validatorCount, N),
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: validatorCount,
				BigOffset:      0,
			},
			distTestCase{
				Name:           fmt.Sprintf("%d-1-%d", validatorCount, N),
				N:              N,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: validatorCount,
				BigOffset:      1,
			})
	}

	// set of 128 validators with different offsets of large validator
	for i := 0; i < 2; i++ {
		tcs = append(tcs, distTestCase{
			Name:           fmt.Sprintf("4096-%d-%d", i, N),
			N:              N,
			BigBalance:     32 * Gi,
			SmallBalance:   16*Gi + 1*Ki,
			ValidatorCount: 128,
			BigOffset:      i,
		})
	}

	// larger validator sets with variable number of elections
	for i := uint(3); i < 13; i++ {
		validatorCount := 1 << i
		electionCount := 1 << (i + 10)
		tcs = append(tcs,
			distTestCase{
				Name:           fmt.Sprintf("%d-0-%d", validatorCount, electionCount),
				N:              electionCount,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: validatorCount,
				BigOffset:      0,
			},
			distTestCase{
				Name:           fmt.Sprintf("%d-1-%d", validatorCount, electionCount),
				N:              electionCount,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: validatorCount,
				BigOffset:      1,
			},
		)
	}

	// 4Ki validator set with different number of elections
	for i := uint(3); i < 9; i++ {
		electionCount := 4096 << i
		tcs = append(tcs,
			distTestCase{
				Name:           fmt.Sprintf("4096-0-%d", electionCount),
				N:              electionCount,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 4096,
				BigOffset:      0,
			},
			distTestCase{
				Name:           fmt.Sprintf("4096-1-%d", electionCount),
				N:              electionCount,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: 4096,
				BigOffset:      1,
			},
		)
	}

	// 16Ki validators with 128Ki elections
	func() {
		const (
			nValidators = 16 * Ki
			nElections  = nValidators << 3
		)

		tcs = append(tcs,
			distTestCase{
				Name:           fmt.Sprintf("%d-0-%d", nValidators, nElections),
				N:              nElections,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: nValidators,
				BigOffset:      0,
			},
			distTestCase{
				Name:           fmt.Sprintf("%d-1-%d", nValidators, nElections),
				N:              nElections,
				BigBalance:     32 * Gi,
				SmallBalance:   16 * Gi,
				ValidatorCount: nValidators,
				BigOffset:      1,
			},
		)
	}()

	var results = make(chan distTestResult, times*len(tcs))

	fmt.Fprintln(os.Stderr, "writing test data to results.csv")

	f, err := os.OpenFile("results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// write table headers if file was opened for the first time
	if fi, err := f.Stat(); err == nil && fi.Size() == 0 {
		// what a mouthful
		_, err = fmt.Fprintln(f, "validatorCount,electionsCount,bigValidatorIndex,"+
			"bigValidatorElectionsAvg,bigValidatorElectionsStdev,bigValidatorElectionsRelStdev,"+
			"smallValidatorElectionsAvg,smallValidatorElectionsStdev,smallValidatorElectionsRelStdev,"+
			"ratioBigAvgToSmallAvg")
	}
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		type average struct {
			BigCount      float32
			SmallCountAvg float32
		}

		var (
			avg            average
			avgs           = make(map[string]average)
			bigCounts      = make(map[string][]float64)
			smallCountAvgs = make(map[string][]float64)
			doneMap        = make(map[string]int)
			res            distTestResult
		)

		for i := 0; i < times*len(tcs); i++ {
			res = <-results

			avg = avgs[res.Name]
			avg.BigCount += float32(res.BigCount()) / times
			avg.SmallCountAvg += float32(res.SmallCountAvg()) / times
			avgs[res.Name] = avg

			bigCounts[res.Name] = append(bigCounts[res.Name], float64(res.BigCount()))
			smallCountAvgs[res.Name] = append(smallCountAvgs[res.Name], float64(res.SmallCountAvg()))

			doneMap[res.Name] += 1

			if doneMap[res.Name] == times {
				bigStdev, err := stats.StandardDeviation(bigCounts[res.Name])
				if err != nil {
					panic(err)
				}

				smallStdev, err := stats.StandardDeviation(smallCountAvgs[res.Name])
				if err != nil {
					panic(err)
				}

				fmt.Fprintf(f,
					"%d,%d,%d,"+
						"%f,%f,%f,"+
						"%f,%f,%f,"+
						"%f\n",
					res.ValidatorCount, res.N, res.BigOffset,
					avg.BigCount, bigStdev, float32(bigStdev)/avg.BigCount,
					avg.SmallCountAvg, smallStdev, float32(smallStdev)/avg.SmallCountAvg,
					avg.BigCount/avg.SmallCountAvg)
			}
		}
	}()

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

func (tr distTestResult) SmallCountAvg() float32 {
	return float32(sum(tr.Counts)-tr.Counts[tr.BigOffset]) / float32(tr.ValidatorCount-1)
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
		/*
			l := len(fmt.Sprintf("%d", tc.N))
			formatN := fmt.Sprintf("N=%%%dd\n", l)
			formatI := fmt.Sprintf("i=%%%dd\n", l)

			fmt.Printf(formatN, tc.N)
		*/
		for i := 0; i < tc.N; i++ {
			/*
				if i&0xffff == 0 {
					fmt.Printf(formatI, i)
				}
			*/

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

		/*
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
		*/
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
