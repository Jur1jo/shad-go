//go:build !solution

package hogwarts

func GetCourseList(prereqs map[string][]string) []string {
	graphDegs := make(map[string]int)
	inverseGraph := make(map[string][]string)
	for key, courses := range prereqs {
		if _, ok := prereqs[key]; !ok {
			graphDegs[key] = 0
		}
		for _, s := range courses {
			graphDegs[key] += 1
			if _, ok := graphDegs[s]; !ok {
				graphDegs[s] = 0
			}
			if _, ok := inverseGraph[s]; !ok {
				inverseGraph[s] = make([]string, 0)
			}
			inverseGraph[s] = append(inverseGraph[s], key)
		}
	}
	stackActualCourses := make([]string, 0)
	for key, deg := range graphDegs {
		if deg == 0 {
			stackActualCourses = append(stackActualCourses, key)
		}
	}
	ans := make([]string, 0)
	for len(stackActualCourses) > 0 {
		currentCourse := stackActualCourses[len(stackActualCourses)-1]
		ans = append(ans, currentCourse)
		stackActualCourses = stackActualCourses[:len(stackActualCourses)-1]
		for _, s := range inverseGraph[currentCourse] {
			graphDegs[s] -= 1
			if graphDegs[s] == 0 {
				stackActualCourses = append(stackActualCourses, s)
			}
		}
	}
	if len(ans) != len(graphDegs) {
		panic("Cycle dependence")
	}
	return ans
}
