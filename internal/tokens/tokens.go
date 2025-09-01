package tokens

type TokKind int

const (
	Lit TokKind = iota // literal (incluye ε y escapes como \*)
	Op                 // operadores: . | *
	LPar
	RPar
)

type Tok struct {
	Kind TokKind
	Val  string // para Lit u operador
}

// Un operador válido
func IsOpRune(r rune) bool { return r == '.' || r == '|' || r == '*' }

// Tokenize: une \x como un único literal
func Tokenize(s string) []Tok {
	var toks []Tok
	r := []rune(s)
	for i := 0; i < len(r); i++ {
		c := r[i]
		if c == '\\' && i+1 < len(r) {
			toks = append(toks, Tok{Kind: Lit, Val: string([]rune{'\\', r[i+1]})})
			i++
			continue
		}
		switch c {
		case '(':
			toks = append(toks, Tok{Kind: LPar})
		case ')':
			toks = append(toks, Tok{Kind: RPar})
		default:
			if IsOpRune(c) {
				toks = append(toks, Tok{Kind: Op, Val: string(c)})
			} else {
				toks = append(toks, Tok{Kind: Lit, Val: string(c)})
			}
		}
	}
	return toks
}
