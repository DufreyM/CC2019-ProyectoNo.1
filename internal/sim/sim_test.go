package sim

import (
	"regex_project/internal/dfa"
	"regex_project/internal/nfa"
	"regex_project/internal/parser"
	"regex_project/internal/regex"
	"regex_project/internal/tokens"
	"testing"
)

func TestRunSimulations(t *testing.T) {
	expr := regex.Preprocess("a(b|c)*")
	toks := tokens.Tokenize(expr)
	post := parser.ShuntingYardTokens(toks)

    n := nfa.BuildFromPostfix(post)
    d := dfa.FromNFA(n)

    cases := []struct {
        w    string
        want bool
    }{
        {"a", true},
        {"ab", true},
        {"accc", true},
        {"b", false},
    }

    for _, c := range cases {
        if got := RunNFA(n, c.w); got != c.want {
            t.Errorf("RunNFA(%q) = %v, want %v", c.w, got, c.want)
        }
        if got := RunDFA(d, c.w); got != c.want {
            t.Errorf("RunDFA(%q) = %v, want %v", c.w, got, c.want)
        }
    }
}
