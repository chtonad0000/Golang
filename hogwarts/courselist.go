//go:build !solution

package hogwarts

var (
	graph   = map[string][]string{}
	used    = map[string]bool{}
	ans     []string
	inStack = map[string]bool{}
)

func dfs(v string) {
	used[v] = true
	inStack[v] = true
	for _, to := range graph[v] {
		if !used[to] {
			dfs(to)
		} else if inStack[to] {
			panic(to)
		}
	}

	inStack[v] = false
	ans = append(ans, v)
}

func topologicalSort() {
	used = make(map[string]bool)
	inStack = make(map[string]bool)
	ans = []string{}

	for vertex := range graph {
		if !used[vertex] {
			dfs(vertex)
		}
	}
}

func GetCourseList(prereqs map[string][]string) []string {
	graph = prereqs
	topologicalSort()
	return ans
}
