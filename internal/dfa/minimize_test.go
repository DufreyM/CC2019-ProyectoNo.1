// internal/dfa/minimize_test.go
package dfa

import "testing"

func TestMinimizeSimple(t *testing.T) {
    d := &DFA{
        Start: 0,
        Accepting: map[int]bool{0: false, 1: true, 2: true},
        Trans: map[int]map[string]int{
            0: {"a": 1},
            1: {"a": 1},
            2: {"a": 1},
        },
        Alphabet: []string{"a"},
        NumStates: 3,
    }

    md := Minimize(d)
    if md.NumStates >= d.NumStates {
        t.Errorf("expected fewer states after minimization, got %d", md.NumStates)
    }
}
