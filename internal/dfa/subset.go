package dfa

import (
	"regex_project/internal/nfa"
	"sort"
	"strconv"
	"strings"
)

type DFA struct {
	Start int
	Accepting map[int]bool
	Trans map[int]map[string]int
	Alphabet []string
	NumStates int
	NFAsets map[int][]int // opcional: set de NFA states por DFA
}

func epsilonClosure(set map[int]*nfa.State) map[int]*nfa.State {
	stack := make([]*nfa.State,0,len(set))
	for _, s := range set { stack = append(stack, s) }
	for len(stack)>0 {
		u := stack[len(stack)-1]; stack = stack[:len(stack)-1]
		for _,v := range u.Eps {
			if _,ok := set[v.ID]; !ok {
				set[v.ID] = v
				stack = append(stack, v)
			}
		}
	}
	return set
}

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

func keyOf(set map[int]*nfa.State) string {
	ids := make([]int,0,len(set))
	for id := range set { ids = append(ids, id) }
	sort.Ints(ids)
	sb := strings.Builder{}
	for i, id := range ids {
		if i>0 { sb.WriteByte(',') }
		sb.WriteString(strconv.Itoa(id))
	}
	return sb.String()
}

// gather alphabet (símbolos de 1 char) desde NFA
func alphabetFromNFA(n *nfa.NFA) []string {
	mp := map[string]bool{}
	for _, s := range n.States {
		for sym := range s.Trans {
			if sym == "" { continue }
			mp[sym] = true
		}
	}
	al := make([]string,0,len(mp))
	for k := range mp { al = append(al, k) }
	sort.Strings(al)
	return al
}

// Subset construction
func FromNFA(n *nfa.NFA) *DFA {
	alpha := alphabetFromNFA(n)
	// e-closure({start})
	startSet := epsilonClosure(map[int]*nfa.State{n.Start.ID:n.Start})
	idxOf := map[string]int{ keyOf(startSet):0 }
	queue := []map[int]*nfa.State{ startSet }
	dfatr := map[int]map[string]int{ 0: {} }
	accept := map[int]bool{ 0: startSet[n.Accept.ID] != nil }
	sets := map[int][]int{ 0: {} }
	for id := range startSet {
		sets[0] = append(sets[0], id)
	}
	sort.Ints(sets[0])

	for qi := 0; qi < len(queue); qi++ {
		S := queue[qi]
		u := idxOf[keyOf(S)]
		if dfatr[u] == nil { dfatr[u] = map[string]int{} }
		for _, sym := range alpha {
			M := move(S, sym)
			if len(M)==0 {
				// dead move handled later
				continue
			}
			T := epsilonClosure(M)
			k := keyOf(T)
			v, ok := idxOf[k]
			if !ok {
				v = len(queue)
				idxOf[k] = v
				queue = append(queue, T)
				dfatr[v] = map[string]int{}
				accept[v] = T[n.Accept.ID] != nil
				sets[v] = []int{}
				for id := range T { sets[v] = append(sets[v], id) }
				sort.Ints(sets[v])
			}
			dfatr[u][sym] = v
		}
	}

	// totalizar DFA con estado pozo
	dead := -1
	for u := 0; u < len(queue); u++ {
		if dfatr[u] == nil { dfatr[u] = map[string]int{} }
		for _, sym := range alpha {
			if _,ok := dfatr[u][sym]; !ok {
				if dead == -1 {
					dead = len(queue)
					dfatr[dead] = map[string]int{}
					accept[dead] = false
					sets[dead] = []int{}
				}
				dfatr[u][sym] = dead
			}
		}
	}
	if dead != -1 {
		for _, sym := range alpha {
			dfatr[dead][sym] = dead
		}
	}

	d := &DFA{
		Start: 0,
		Accepting: accept,
		Trans: dfatr,
		Alphabet: alpha,
		NumStates: len(dfatr),
		NFAsets: sets,
	}
	return d
}
