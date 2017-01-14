package recommend

import (
	"math"
	"sort"
)

var Critics = make(map[string]map[string]float64)

//------------------ p.11 ------------------
// person1とperson2の距離を基にした類似性スコアを返す
func SimDistance(prefs map[string]map[string]float64, person1, person2 string) float64 {
	si := make(map[string]bool)

	for k := range prefs[person1] {
		if _, ok := prefs[person2][k]; ok {
			si[k] = true
		}
	}

	if len(si) == 0 {
		return 0
	}

	var sum_of_squares float64
	for k := range si {
		sum_of_squares += math.Pow(prefs[person1][k]-prefs[person2][k], 2.0)
	}

	return 1 / (1 + sum_of_squares)
}

func SimPearson(prefs map[string]map[string]float64, p1, p2 string) float64 {
	si := make(map[string]bool)

	for k := range prefs[p1] {
		if _, ok := prefs[p2][k]; ok {
			si[k] = true
		}
	}

	n := float64(len(si))

	if n == 0 {
		return 0
	}

	var sum1, sum2 float64
	var sum1Sq, sum2Sq float64
	var pSum float64
	for k := range si {
		// すべての嗜好の合計
		sum1 += prefs[p1][k]
		sum2 += prefs[p2][k]

		// 平方の合計
		sum1Sq += math.Pow(prefs[p1][k], 2.0)
		sum2Sq += math.Pow(prefs[p2][k], 2.0)

		// 積の合計
		pSum += prefs[p1][k] * prefs[p2][k]
	}

	num := pSum - (sum1 * sum2 / n)
	den := math.Sqrt((sum1Sq - math.Pow(sum1, 2.0)/n) * (sum2Sq - math.Pow(sum2, 2.0)/n))
	if den == 0 {
		return 0
	}

	return num / den

}

//------------------------------------------

//------------------ p.15 ------------------
func TopMatches(prefs map[string]map[string]float64, person string, n int, similarity func(map[string]map[string]float64, string, string) float64) ScorePairList {
	scores := make(ScorePairList, len(prefs)-1)
	var idx int
	for k := range prefs {
		if k == person {
			continue
		}

		scores[idx] = ScorePair{
			Score: similarity(prefs, person, k),
			Key:   k,
		}
		idx++
	}
	sort.Sort(sort.Reverse(scores))

	return scores[:n]
}

//------------------------------------------
