package dfa

import (
    "testing"
    "regex_project/internal/tokens"
    "regex_project/internal/parser"
    "regex_project/internal/nfa"
    "regex_project/internal/regex"
)

func TestFromNFA(t *testing.T) {
    expr := regex.Preprocess("ab")
    toks := tokens.Tokenize(expr)
    post := parser.ShuntingYardTokens(toks)
    n := nfa.BuildFromPostfix(post)
    d := FromNFA(n)

    if d.NumStates == 0 {
        t.Fatalf("expected some DFA states, got 0")
    }
}
