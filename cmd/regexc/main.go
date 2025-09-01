package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"regex_project/internal/regex"
	"regex_project/internal/tokens"
	"regex_project/internal/parser"
	"regex_project/internal/ast"
	"regex_project/internal/nfa"
	"regex_project/internal/dfa"
	"regex_project/internal/graphviz"
	"regex_project/internal/sim"
)

func main() {
	var (
		fileIn string
		oneR string
		w string
		outdir string
		eps string
	)
	flag.StringVar(&fileIn, "f", "", "Archivo con ERs (una por línea)")
	flag.StringVar(&oneR, "r", "", "Una ER por parámetro")
	flag.StringVar(&w, "w", "", "Cadena a simular (opcional)")
	flag.StringVar(&outdir, "out", "out", "Directorio de salida")
	flag.StringVar(&eps, "eps", "ε", "Símbolo para epsilon")
	flag.Parse()

	if fileIn == "" && oneR == "" {
		fmt.Println("Uso: -r \"(a|b)*abb\" [-w cadena] [-out dir] [-eps ε]  |  -f input.txt [-w cadena] [-out dir]")
		return
	}

	if err := os.MkdirAll(outdir, 0o755); err != nil {
		log.Fatal(err)
	}

	var exprs []string
	if oneR != "" {
		exprs = []string{oneR}
	} else {
		f, err := os.Open(fileIn)
		if err != nil { log.Fatal(err) }
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" { continue }
			exprs = append(exprs, line)
		}
		if err := sc.Err(); err != nil { log.Fatal(err) }
	}

	for i, raw := range exprs {
		k := i+1
		fmt.Printf("ER #%d: %q\n", k, raw)

		// Preprocesado
		pre := regex.Preprocess(strings.ReplaceAll(raw, "\\e", eps))

		// Tokenización y Shunting Yard
		toks := tokens.Tokenize(pre)
		post := parser.ShuntingYardTokens(toks)

		// AST (opcional, export)
		astRoot := ast.BuildFromPostfix(post)
		graphviz.ExportAST(astRoot, filepath.Join(outdir, fmt.Sprintf("ast_%d.dot", k)))

		// NFA y DFA
		n := nfa.BuildFromPostfix(post)
		graphviz.ExportNFA(n, filepath.Join(outdir, fmt.Sprintf("nfa_%d.dot", k)))

		d := dfa.FromNFA(n)
		graphviz.ExportDFA(d, filepath.Join(outdir, fmt.Sprintf("dfa_raw_%d.dot", k)))

		md := dfa.Minimize(d)
		graphviz.ExportDFA(md, filepath.Join(outdir, fmt.Sprintf("dfa_min_%d.dot", k)))

		// Postfix string legible
		var sb strings.Builder
		for _, t := range post {
			if t.Kind == tokens.Lit || t.Kind == tokens.Op { sb.WriteString(t.Val) } else { sb.WriteRune('?') }
		}

		// Simulaciones
		resNFA, resDFA, resMin := "n/a", "n/a", "n/a"
		if w != "" {
			if sim.RunNFA(n, w) { resNFA = "sí" } else { resNFA = "no" }
			if sim.RunDFA(d, w) { resDFA = "sí" } else { resDFA = "no" }
			if sim.RunDFA(md, w) { resMin = "sí" } else { resMin = "no" }
		}

		// Reporte
		rep := filepath.Join(outdir, fmt.Sprintf("report_%d.txt", k))
		f, _ := os.Create(rep)
		fmt.Fprintf(f, "Expresión: %s\nPreprocesada: %s\nPostfix: %s\n", raw, pre, sb.String())
		if w != "" {
			fmt.Fprintf(f, "Simulación con w=%q -> AFN: %s, AFD: %s, AFD_min: %s\n", w, resNFA, resDFA, resMin)
		}
		f.Close()

		fmt.Printf("  -> DOTs y reporte en %s (postfix=%s)\n", outdir, sb.String())
	}
}
