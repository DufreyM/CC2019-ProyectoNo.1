package parser

import (
	"regex_project/internal/regex"
	"regex_project/internal/tokens"
	"testing"
)

func TestShuntingYardTokens(t *testing.T) {
    expr := regex.Preprocess("a|(b*c)")
    toks := tokens.Tokenize(expr)
    out := ShuntingYardTokens(toks)

    t.Log("postfix:", out)

    // Lo que devuelve tu implementación
    want := []string{"a", "b", "*", "c", ".", "|"}

    if len(out) != len(want) {
        t.Fatalf("expected %d tokens, got %d", len(want), len(out))
    }

    for i := range want {
        if out[i].Val != want[i] {
            t.Errorf("pos %d: got %q, want %q", i, out[i].Val, want[i])
        }
    }
}
