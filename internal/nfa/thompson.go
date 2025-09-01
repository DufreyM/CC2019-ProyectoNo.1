package nfa

import (
	"regex_project/internal/tokens"
	"sort"
)

type State struct {
	ID   int
	Eps  []*State
	Trans map[string][]*State // etiqueta "" es epsilon
	Accept bool
}

type NFA struct {
	Start *State
	Accept *State
	States []*State
}

var sid int
func newState() *State {
	st := &State{ID: sid, Trans: make(map[string][]*State)}
	sid++
	return st
}

// literalSymbol convierte un literal tokenizado (p.ej., "\*") en su símbolo de 1 char: "*".
func literalSymbol(val string) string {
	runes := []rune(val)
	if len(runes) >= 2 && runes[0] == '\\' {
		return string(runes[1])
	}
	return val
}

type frag struct{ start, end *State }

// Thompson desde postfix de tokens
func BuildFromPostfix(post []tokens.Tok) *NFA {
	sid = 0
	var st []frag

	push := func(f frag){ st = append(st, f) }
	pop := func() frag { f:=st[len(st)-1]; st = st[:len(st)-1]; return f }
	pop2 := func() (frag, frag) { b:=pop(); a:=pop(); return a,b }

	for _, t := range post {
		switch t.Kind {
		case tokens.Lit:
			if t.Val == "ε" {
				// fragmento epsilon
				s := newState()
				e := newState(); e.Accept = true
				s.Eps = append(s.Eps, e)
				push(frag{s, e})
			} else {
				sym := literalSymbol(t.Val)
				s := newState()
				e := newState(); e.Accept = true
				s.Trans[sym] = append(s.Trans[sym], e)
				push(frag{s, e})
			}
		case tokens.Op:
			switch t.Val {
			case ".":
				a, b := pop2()
				// quitar aceptación intermedia
				a.end.Accept = false
				a.end.Eps = append(a.end.Eps, b.start)
				push(frag{a.start, b.end})
			case "|":
				a, b := pop2()
				s := newState()
				e := newState(); e.Accept = true
				a.end.Accept = false; b.end.Accept = false
				s.Eps = append(s.Eps, a.start, b.start)
				a.end.Eps = append(a.end.Eps, e)
				b.end.Eps = append(b.end.Eps, e)
				push(frag{s, e})
			case "*":
				a := pop()
				s := newState()
				e := newState(); e.Accept = true
				a.end.Accept = false
				s.Eps = append(s.Eps, a.start, e)
				a.end.Eps = append(a.end.Eps, a.start, e)
				push(frag{s, e})
			}
		}
	}
	if len(st) != 1 { panic("NFA mal formado") }
	fr := st[0]

	// recolectar estados por DFS
	visited := map[int]*State{}
	var order []*State
	var dfs func(*State)
	dfs = func(u *State){
		if _,ok := visited[u.ID]; ok { return }
		visited[u.ID] = u
		order = append(order, u)
		for _,v := range u.Eps { dfs(v) }
		for _, arr := range u.Trans {
			for _, v := range arr { dfs(v) }
		}
	}
	dfs(fr.start)

	// ordenar por ID para consistencia
	sort.Slice(order, func(i,j int) bool { return order[i].ID < order[j].ID })

	nfa := &NFA{Start: fr.start, Accept: fr.end, States: order}
	return nfa
}
