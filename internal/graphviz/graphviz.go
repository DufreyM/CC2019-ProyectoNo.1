package graphviz

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"regex_project/internal/ast"
	"regex_project/internal/dfa"
	"regex_project/internal/nfa"
)

func escapeLabel(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}

// AST
func ExportAST(root *ast.Node, filename string) {
	ast.ExportDOT(root, filename)
}

// NFA
func ExportNFA(n *nfa.NFA, filename string) error {
	f, err := os.Create(filename)
	if err != nil { return err }
	defer f.Close()

	fmt.Fprintln(f, "digraph NFA {")
	fmt.Fprintln(f, `  rankdir=LR;`)
	for _, s := range n.States {
		shape := "circle"
		if s.Accept { shape = "doublecircle" }
		fmt.Fprintf(f, `  q%d [shape=%s];`+"\n", s.ID, shape)
	}
	// marcador de inicio
	fmt.Fprintln(f, `  start [shape=plaintext,label=""];`)
	fmt.Fprintf(f, "  start -> q%d;\n", n.Start.ID)
	for _, s := range n.States {
		for _, v := range s.Eps {
			fmt.Fprintf(f, `  q%d -> q%d [label="ε"];`+"\n", s.ID, v.ID)
		}
		for sym, arr := range s.Trans {
			for _, v := range arr {
				fmt.Fprintf(f, `  q%d -> q%d [label="%s"];`+"\n", s.ID, v.ID, escapeLabel(sym))
			}
		}
	}
	fmt.Fprintln(f, "}")
	return nil
}

// DFA
func ExportDFA(d *dfa.DFA, filename string) error {
	f, err := os.Create(filename)
	if err != nil { return err }
	defer f.Close()

	fmt.Fprintln(f, "digraph DFA {")
	fmt.Fprintln(f, `  rankdir=LR;`)

	ids := make([]int,0,len(d.Trans))
	for s := range d.Trans { ids = append(ids, s) }
	sort.Ints(ids)

	for _, s := range ids {
		shape := "circle"
		if d.Accepting[s] { shape = "doublecircle" }
		fmt.Fprintf(f, `  D%d [shape=%s];`+"\n", s, shape)
	}
	fmt.Fprintln(f, `  start [shape=plaintext,label=""];`)
	fmt.Fprintf(f, "  start -> D%d;\n", d.Start)

	for _, s := range ids {
		for _, a := range d.Alphabet {
			dst := d.Trans[s][a]
			fmt.Fprintf(f, `  D%d -> D%d [label="%s"];`+"\n", s, dst, escapeLabel(a))
		}
	}
	fmt.Fprintln(f, "}")
	return nil
}
