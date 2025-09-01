package ast

import (
    "testing"
    "regex_project/internal/tokens"
)

func TestBuildFromPostfix(t *testing.T) {
    // Expresión: ab.
    toks := []tokens.Tok{
        {Kind: tokens.Lit, Val: "a"},
        {Kind: tokens.Lit, Val: "b"},
        {Kind: tokens.Op, Val: "."},
    }
    root := BuildFromPostfix(toks)

    if root.Value != "." {
        t.Errorf("expected root '.', got %q", root.Value)
    }
    if root.Left.Value != "a" || root.Right.Value != "b" {
        t.Errorf("AST children mismatch: got %q and %q", root.Left.Value, root.Right.Value)
    }
}
