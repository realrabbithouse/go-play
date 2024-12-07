package algs

/* 编辑距离: 动态规划 */
func editDistanceDP(s string, t string) int {
	/*
		状态 dp[i,j] 对应的问题描述: 将 s 的前 i 个元素替换为 t 的前 j 个元素所需的最少编辑距离
		dp[i,j] = min{dp[i,j-1], dp[i-1,j-1], dp[i-1,j]} + 1
		如果 s[i-1] = t[j-1], 那么 dp[i,j] = dp[i-1,j-1]
		初始条件: dp[0,0] = 0, dp[0,j] = j, dp[i,0] = i
	*/
	dp := make([][]int, len(s)+1)
	for i := 0; i <= len(s); i++ {
		dp[i] = make([]int, len(t)+1)
	}

	for i := 1; i <= len(s); i++ {
		dp[i][0] = i
	}
	for j := 1; j <= len(t); j++ {
		dp[0][j] = j
	}

	for i := 1; i <= len(s); i++ {
		for j := 1; j <= len(t); j++ {
			if s[i-1] == t[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min(dp[i-1][j], dp[i-1][j-1], dp[i][j-1]) + 1
			}
		}
	}

	return dp[len(s)][len(t)]
}
