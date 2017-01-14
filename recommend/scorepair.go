package recommend

type ScorePair struct {
	Score float64
	Key   string
}

type ScorePairList []ScorePair

func (p ScorePairList) Len() int {
	return len(p)
}

func (p ScorePairList) Less(i, j int) bool {
	return p[i].Score < p[j].Score
}

func (p ScorePairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
