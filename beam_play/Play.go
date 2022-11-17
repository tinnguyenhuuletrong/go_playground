package beam_play

import (
	beam_play_grade "ttin.com/play2022/beam_play/grade"
	beam_play_wc "ttin.com/play2022/beam_play/wc"
)

func Play_WordCount() {
	beam_play_wc.StartWc()
}

func Play_StatsGrade() {
	beam_play_grade.StartGrade()
}
