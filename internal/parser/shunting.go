package parser

import "regex_project/internal/tokens"

var prec = map[string]int{"*":3, ".":2, "|":1}

func leftAssoc(op string) bool { return op != "*" }

// Shunting Yard sobre tokens
func ShuntingYardTokens(toks []tokens.Tok) []tokens.Tok {
	var out []tokens.Tok
	var st []tokens.Tok
	for _, t := range toks {
		switch t.Kind {
		case tokens.Lit:
			out = append(out, t)
		case tokens.Op:
			for len(st) > 0 && st[len(st)-1].Kind == tokens.Op {
				top := st[len(st)-1].Val
				if (leftAssoc(t.Val) && prec[t.Val] <= prec[top]) || (!leftAssoc(t.Val) && prec[t.Val] < prec[top]) {
					out = append(out, st[len(st)-1])
					st = st[:len(st)-1]
				} else { break }
			}
			st = append(st, t)
		case tokens.LPar:
			st = append(st, t)
		case tokens.RPar:
			for len(st) > 0 && st[len(st)-1].Kind != tokens.LPar {
				out = append(out, st[len(st)-1])
				st = st[:len(st)-1]
			}
			if len(st) > 0 && st[len(st)-1].Kind == tokens.LPar {
				st = st[:len(st)-1]
			}
		}
	}
	for i := len(st)-1; i >= 0; i-- {
		out = append(out, st[i])
	}
	return out
}
