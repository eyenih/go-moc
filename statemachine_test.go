package moc

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testByteSliceIterator struct {
	current int
	s       []byte
}

func (it *testByteSliceIterator) Next() interface{} {
	defer func() { it.current++ }()
	return it.s[it.current]
}

func (it *testByteSliceIterator) Done() bool {
	return it.current >= len(it.s)
}

type testSM struct {
	inputs      []byte
	transitions int

	err error
}

func (fsm *testSM) Transition(i interface{}) error {
	if fsm.err != nil {
		return fsm.err
	}

	fsm.inputs = append(fsm.inputs, i.(byte))
	fsm.transitions++

	return nil
}

func TestExecute(t *testing.T) {
	gofakeit.Seed(0)
	t.Run("expected", func(t *testing.T) {
		content := gofakeit.FirstName()
		it := &testByteSliceIterator{s: []byte(content)}
		fsm := &testSM{}
		err := Execute(it, fsm)
		require.NoError(t, err)

		assert.Equal(t, len(content), it.current)
		assert.Equal(t, content, string(fsm.inputs))
		assert.Equal(t, len(content), fsm.transitions)
	})

	t.Run("fsm error caught", func(t *testing.T) {
		it := &testByteSliceIterator{s: []byte(gofakeit.Letter())}
		fsm := &testSM{err: errors.New(gofakeit.FirstName())}
		err := Execute(it, fsm)
		assert.Equal(t, fsm.err, err)
	})
}

func BenchmarkExecute(b *testing.B) {
	it := &testByteSliceIterator{s: []byte("test")}
	fsm := &testSM{}
	var err error

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		err = Execute(it, fsm)

		if err != nil {
			b.Error(err)
		}
	}

	// only for insturcting the compiler to really execute
	resultErr := err

	if resultErr != err {
		panic(err)
	}
}
