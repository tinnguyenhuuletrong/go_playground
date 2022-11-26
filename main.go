package main

import networktcp "ttin.com/play2022/network_tcp"

func main() {
	// network.Play()
	// beam_play.Play_StatsGrade()
	// beam_play.Play_WordCount()
	// beam_play.Play_GroupBy()
	// beam_play.Play_CombineSum()

	// grpc_play.Play_Grpc_Twirp()

	// faninout.Play_FanInOut()

	networktcp.CreateTCPServer("localhost:3000")
}
