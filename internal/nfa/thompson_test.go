package nfa

import (
    "testing"
    "regex_project/internal/tokens"
)

func TestBuildFromPostfixSimple(t *testing.T) {
    toks := []tokens.Tok{
        {Kind: tokens.Lit, Val: "a"},
        {Kind: tokens.Lit, Val: "b"},
        {Kind: tokens.Op, Val: "."},
    }
    n := BuildFromPostfix(toks)
    if n.Start == nil || n.Accept == nil {
        t.Fatal("NFA Start/Accept should not be nil")
    }
    if len(n.States) < 2 {
        t.Errorf("expected more states, got %d", len(n.States))
    }
}
