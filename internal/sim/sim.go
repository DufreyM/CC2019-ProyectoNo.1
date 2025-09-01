package sim

import (
	"unicode/utf8"
	"regex_project/internal/dfa"
	"regex_project/internal/nfa"
)

// Simulación DFA
func RunDFA(d *dfa.DFA, w string) bool {
	cur := d.Start
	for len(w) > 0 {
		r, size := utf8.DecodeRuneInString(w)
		w = w[size:]
		sym := string(r)
		cur = d.Trans[cur][sym]
	}
	return d.Accepting[cur]
}

// e-closure
func epsClosure(set map[int]*nfa.State) map[int]*nfa.State {
	stack := make([]*nfa.State,0,len(set))
	for _, s := range set { stack = append(stack, s) }
	for len(stack)>0 {
		u := stack[len(stack)-1]; stack = stack[:len(stack)-1]
		for _, v := range u.Eps {
			if _,ok := set[v.ID]; !ok {
				set[v.ID] = v
				stack = append(stack, v)
			}
		}
	}
	return set
}

// move
func move(set map[int]*nfa.State, sym string) map[int]*nfa.State {
	out := map[int]*nfa.State{}
	for _, s := range set {
		if nexts, ok := s.Trans[sym]; ok {
			for _, v := range nexts {
				out[v.ID] = v
			}
		}
	}
	return out
}

// Simulación AFN en paralelo por cerradura-e
func RunNFA(n *nfa.NFA, w string) bool {
	cur := epsClosure(map[int]*nfa.State{ n.Start.ID: n.Start })
	for len(w)>0 {
		r, size := utf8.DecodeRuneInString(w)
		w = w[size:]
		sym := string(r)
		cur = epsClosure(move(cur, sym))
		if len(cur)==0 { return false }
	}
	// aceptación si algún estado aceptador está en el conjunto
	return cur[n.Accept.ID] != nil
}
