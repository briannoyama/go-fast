package threadless

type node[V any] struct {
	children [2]int
	v        V
}
