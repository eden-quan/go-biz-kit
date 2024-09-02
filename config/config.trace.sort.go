package config

type traceArray []*traceObject

func (t traceArray) Len() int {
	return len(t)
}

func (t traceArray) Less(i, j int) bool {
	return t[i].path < t[j].path
}

func (t traceArray) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
