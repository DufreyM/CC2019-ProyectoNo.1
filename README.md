# Proyecto 1 – Expresiones regulares → AFN → AFD → AFD mínimo

Pipeline completo en Go:
1) Preprocesado de la ER (clases de caracteres, +, ?, \ε, y concatenación explícita `.`).
2) Tokenización con **escapes como un solo literal** (p.ej., `\*`).
3) Shunting Yard sobre *tokens* → **postfix**.
4) Construcción de **AST** (opcional, se exporta a DOT).
5) Construcción de **AFN** (Thompson) desde postfix.
6) **AFD** por construcción de subconjuntos.
7) **Minimización** de AFD (table-filling; DFA totalizado con estado pozo).
8) **Simulación** de AFN y AFD/AFD mínimo con la cadena `w`.
9) Export a **Graphviz DOT** para AST, AFN, AFD y AFD mínimo.

> El símbolo de ε por defecto es `ε` (configurable con `-eps`).  
> **Todos los símbolos del alfabeto son de longitud 1**; un literal escapado (p.ej., `\*`) equivale al símbolo `*`.

## Uso rápido

```bash
cd regex_project
go run ./cmd/regexc -r "(a|b)*abb" -w "babb" -out out
# para el archivo 
go run ./cmd/regexc -f input.txt -w "babb" -out out
```

Salidas en `out/` para cada ER:  
- `ast_k.dot`, `nfa_k.dot`, `dfa_raw_k.dot`, `dfa_min_k.dot`  
- `report_k.txt` con postfix y resultados de simulación

**Nota**: Para PNG, con Graphviz instalado:
```bash
dot -Tpng out/dfa_min_1.dot -o out/dfa_min_1.png
```

## Estructura
```
cmd/regexc/main.go
internal/regex/preprocess.go
internal/tokens/tokens.go
internal/parser/shunting.go
internal/ast/ast.go
internal/nfa/thompson.go
internal/dfa/subset.go
internal/dfa/minimize.go
internal/graphviz/graphviz.go
internal/sim/sim.go
```

## Pruebas rápidas
```
go run ./cmd/regexc -r "[ab]*abb" -w "babb" -out out
go run ./cmd/regexc -r "(abc)+" -w "abcabc" -out out
go run ./cmd/regexc -r "0?" -w "" -out out
go run ./cmd/regexc -r "\*a*" -w "*aa" -out out
```
