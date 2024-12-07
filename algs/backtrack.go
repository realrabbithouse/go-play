package algs

//
// In order to apply eightQueenBT to a specific class of problems, one must provide the testdata P for
// the particular instance of the problem that is to be solved, and six procedural parameters, root,
// reject, accept, first, next, and output. These procedures should take the instance testdata P as a
// parameter and should do the following:
//
// root(P): return the partial candidate at the root of the search tree.
// reject(P,c): return true only if the partial candidate c is not worth completing.
// accept(P,c): return true if c is a solution of P, and false otherwise.
// first(P,c): generate the first extension of candidate c.
// next(P,s): generate the next alternative extension of a candidate, after the extension s.
// output(P,c): use the solution c of P, as appropriate to the application.
//
// The eightQueenBT algorithm reduces the problem to the call backtrack(root(P)), where backtrack is
// the following recursive procedure:
//
// procedure backtrack(c) is
//    if reject(P, c) then return
//    if accept(P, c) then output(P, c)
//    s ← first(P, c)
//    while s ≠ NULL do
//        backtrack(s)
//        s ← next(P, s)
//
// 1. Decision Problem – In this, we search for a feasible solution.
// 2. Optimization Problem – In this, we search for the best solution.
// 3. Enumeration Problem – In this, we find all feasible solutions.
//

/*
回溯算法框架
func backtrack(state *State, choices []Choice, res *[]State) {
	// 判断是否为解
	if isSolution(state) {
		// 记录解
		recordSolution(state, res)
		// 不再继续搜索
		return
	}
	// 遍历所有选择
	for _, choice := range choices {
		// 剪枝：判断选择是否合法
		if isValid(state, choice) {
			// 尝试：做出选择，更新状态
			makeChoice(state, choice)
			backtrack(state, choices, res)
			// 回退：撤销选择，恢复到之前的状态
			undoChoice(state, choice)
		}
	}
}
*/

// nQueens solves the N-Queens problem for a given board size n.
// It returns all possible solutions where n queens can be placed on an n x n chessboard
// such that no two queens threaten each other.
//
// Parameters:
//
//	n - The size of the chessboard (n x n) and the number of queens to place.
//
// Returns:
//
//	A slice of slices of integers, where each inner slice represents a valid solution.
//	Each integer in the inner slice represents the column position of a queen in the corresponding row.
func nQueens(n int) [][]int {
	var (
		state []int
		res   [][]int
	)
	nQueensBT(state, n, &res)
	return res
}

func nQueensBT(state []int, n int, res *[][]int) {
	if len(state) == n {
		chess := make([]int, n)
		copy(chess, state)
		*res = append(*res, chess)
		return
	}

	for i := 0; i < n; i++ {
		if accept(state, i) {
			state = append(state, i)
			nQueensBT(state, n, res)
			state = state[:len(state)-1] // backtrack to try next column
		}
	}
}

func accept(state []int, k int) bool {
	N := len(state)
	for i := 0; i < N; i++ {
		if state[i] == k ||
			state[i]-k == N-i ||
			k-state[i] == N-i {
			return false
		}
	}
	return true
}
