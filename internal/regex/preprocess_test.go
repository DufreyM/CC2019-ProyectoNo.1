package regex

import "testing"

func TestPreprocess(t *testing.T) {
    cases := []struct {
        in, want string
    }{
        {"a(b|c)*", "a.(b|c)*"},
        {"[ab]c", "(a|b).c"},
        {"a+", "(a).(a)*"},
        {"a?", "(a|ε)"},
    }

    for _, c := range cases {
        got := Preprocess(c.in)
        if got != c.want {
            t.Errorf("Preprocess(%q) = %q, want %q", c.in, got, c.want)
        }
    }
}
