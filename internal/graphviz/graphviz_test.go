package graphviz

import (
    "os"
    "testing"
    "regex_project/internal/tokens"
    "regex_project/internal/parser"
    "regex_project/internal/nfa"
    "regex_project/internal/dfa"
    "regex_project/internal/regex"
)

func TestExportGraphviz(t *testing.T) {
    expr := regex.Preprocess("a|b")
    toks := tokens.Tokenize(expr)
    post := parser.ShuntingYardTokens(toks)
    n := nfa.BuildFromPostfix(post)
    d := dfa.FromNFA(n)

    // AST
    nfa.BuildFromPostfix(post) // reusamos
    if err := ExportNFA(n, "test_nfa.dot"); err != nil {
        t.Fatal(err)
    }
    if err := ExportDFA(d, "test_dfa.dot"); err != nil {
        t.Fatal(err)
    }

    os.Remove("test_nfa.dot")
    os.Remove("test_dfa.dot")
}
