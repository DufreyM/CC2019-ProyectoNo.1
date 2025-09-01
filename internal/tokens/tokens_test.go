package tokens

import "testing"

func TestTokenize(t *testing.T) {
    toks := Tokenize("a|(b*c)")
    kinds := []TokKind{Lit, Op, LPar, Lit, Op, Lit, RPar}

    if len(toks) != len(kinds) {
        t.Fatalf("expected %d tokens, got %d", len(kinds), len(toks))
    }

    for i, k := range kinds {
        if toks[i].Kind != k {
            t.Errorf("at %d: got kind %v, want %v", i, toks[i].Kind, k)
        }
    }
}
