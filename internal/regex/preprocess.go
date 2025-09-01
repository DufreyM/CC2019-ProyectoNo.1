package regex

import (
	"strings"
	"unicode"
)

// Expandir escapes simples comunes visibles \n, \t, \{, \}, \\
// No "desescapa" operadores como \* o \| porque esos son literales.
func ExpandEscapes(input string) string {
	r := strings.NewReplacer(
		`\\`, `\`,
		`\n`, "\n",
		`\t`, "\t",
		`\{`, "{",
		`\}`, "}",
	)
	return r.Replace(input)
}

// [abc] -> (a|b|c). Conserva escapes dentro del rango como literales.
func ExpandCharClasses(s string) string {
	var out []rune
	r := []rune(s)
	i := 0
	for i < len(r) {
		if r[i] == '\\' && i+1 < len(r) {
			out = append(out, r[i], r[i+1])
			i += 2
			continue
		}
		if r[i] == '[' {
			j := i + 1
			var items []string
			esc := false
			for j < len(r) {
				if !esc && r[j] == '\\' && j+1 < len(r) {
					esc = true
					j++
					continue
				}
				if !esc && r[j] == ']' {
					break
				}
				if esc {
					items = append(items, string([]rune{'\\', r[j]}))
					esc = false
				} else {
					items = append(items, string(r[j]))
				}
				j++
			}
			if j >= len(r) || r[j] != ']' || len(items) == 0 {
				// no se cerró, dejar literal '['
				out = append(out, r[i])
				i++
				continue
			}
			out = append(out, '(')
			for k, it := range items {
				out = append(out, []rune(it)...)
				if k != len(items)-1 {
					out = append(out, '|')
				}
			}
			out = append(out, ')')
			i = j + 1
			continue
		}
		out = append(out, r[i])
		i++
	}
	return string(out)
}

// Maneja + y ? expandiéndolos a primitivas: G+ -> (G).(G)* ; G? -> (G|ε)
func HandleExtensions(expr string) string {
	var out []rune
	runes := []rune(expr)
	escaped := false

	extractGroup := func() []rune {
		if len(out) == 0 {
			return nil
		}
		if out[len(out)-1] == ')' {
			cnt := 0
			for j := len(out) - 1; j >= 0; j-- {
				if out[j] == ')' {
					cnt++
				} else if out[j] == '(' {
					cnt--
				}
				if cnt == 0 {
					g := append([]rune(nil), out[j:]...)
					out = out[:j]
					return g
				}
			}
			return nil
		}
		g := []rune{out[len(out)-1]}
		out = out[:len(out)-1]
		return g
	}

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		if escaped {
			out = append(out, '\\', c)
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}

		if (c == '+' || c == '?') {
			group := extractGroup()
			if group == nil {
				continue
			}
			switch c {
			case '+':
				out = append(out, '(')
				out = append(out, group...)
				out = append(out, ')', '.')
				out = append(out, '(')
				out = append(out, group...)
				out = append(out, ')', '*')
			case '?':
				out = append(out, '(')
				out = append(out, group...)
				out = append(out, '|', 'ε', ')')
			}
		} else {
			out = append(out, c)
		}
	}
	return string(out)
}

// Inserta '.' de concatenación donde aplica.
func InsertConcatOperators(regex string) string {
	var result []rune
	runes := []rune(regex)
	escaped := false

	isLiteral := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsDigit(r) || r == 'ε'
	}

	for i := 0; i < len(runes); i++ {
		c := runes[i]

		if escaped {
			result = append(result, '\\', c)
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		result = append(result, c)

		if i+1 < len(runes) {
			next := runes[i+1]
			if (isLiteral(c) || c == '*' || c == ')' ) &&
				(isLiteral(next) || next == '(' || next == '\\') {
				result = append(result, '.')
			}
		}
	}
	return string(result)
}

// Preprocess aplica el orden sugerido.
func Preprocess(raw string) string {
	s := ExpandEscapes(raw)
	s = ExpandCharClasses(s)
	s = HandleExtensions(s)
	s = InsertConcatOperators(s)
	return s
}
