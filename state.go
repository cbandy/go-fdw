package fdw

import "sync"

var initialized struct {
	handler Handler

	execs execStates
}

type execStates struct{ sync.Map }

func (es *execStates) Load(key interface{}) *execState {
	i, _ := es.Map.Load(key)
	return i.(*execState)
}
