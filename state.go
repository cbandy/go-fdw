package fdw

import "sync"

var initialized struct {
	handler Handler

	execs execStates
	plans planStates
}

type execStates struct{ sync.Map }

func (es *execStates) Load(key interface{}) *execState {
	i, _ := es.Map.Load(key)
	return i.(*execState)
}

type planStates struct{ sync.Map }

func (ps *planStates) Load(key interface{}) *planState {
	i, _ := ps.Map.Load(key)
	return i.(*planState)
}
