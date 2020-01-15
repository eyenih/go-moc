package moc

type Iterator interface {
	Next() (interface{}, error)
	Done() bool
}

type StateMachine interface {
	Transition(interface{}) error
}

func Execute(it Iterator, sm StateMachine) error {
	for !it.Done() {
		if i, err := it.Next(); err != nil {
			return err
		} else {
			if err := sm.Transition(i); err != nil {
				return err
			}
		}
	}

	return nil
}
