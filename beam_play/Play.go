package beam_play

import (
	combine "ttin.com/play2022/beam_play/combine"
	beam_play_grade "ttin.com/play2022/beam_play/grade"
	groupby "ttin.com/play2022/beam_play/groupby"
	beam_play_wc "ttin.com/play2022/beam_play/wc"
)

func Play_WordCount() {
	beam_play_wc.StartWc()
}

func Play_StatsGrade() {
	beam_play_grade.StartGrade()
}

func Play_GroupBy() {
	groupby.StartGroupBy()
}

func Play_CombineSum() {
	combine.StartCombineSum()
}
