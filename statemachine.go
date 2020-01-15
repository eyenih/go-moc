package moc

type Iterator interface {
	Next() interface{}
	Done() bool
}

type StateMachine interface {
	Transition(interface{}) error
}

func Execute(it Iterator, sm StateMachine) error {
	for !it.Done() {
		if err := sm.Transition(it.Next()); err != nil {
			return err
		}
	}

	return nil
}
