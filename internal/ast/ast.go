package ast

import (
	"fmt"
	"log"
	"os"
	"regex_project/internal/tokens"
)

type Node struct {
	Value string
	Left, Right *Node
}

func BuildFromPostfix(post []tokens.Tok) *Node {
	var st []*Node
	push := func(n *Node){ st = append(st, n) }
	pop := func() *Node { n:=st[len(st)-1]; st=st[:len(st)-1]; return n }

	for _, t := range post {
		switch t.Kind {
		case tokens.Lit:
			push(&Node{Value: t.Val})
		case tokens.Op:
			if t.Val == "*" {
				if len(st) < 1 { log.Fatal("'*' sin operando") }
				a := pop()
				push(&Node{Value: "*", Left: a})
			} else { // . |
				if len(st) < 2 { log.Fatalf("Operador %s sin operandos", t.Val) }
				b := pop(); a := pop()
				push(&Node{Value: t.Val, Left: a, Right: b})
			}
		default:
			log.Fatal("Token inesperado en postfix")
		}
	}
	if len(st) != 1 { log.Fatalf("AST mal formado, quedan %d nodos", len(st)) }
	return st[0]
}

var nodeID int

func ExportDOT(root *Node, filename string) {
	file, err := os.Create(filename)
	if err != nil { log.Fatalf("Error creando archivo DOT: %v", err) }
	defer file.Close()

	file.WriteString("digraph AST {\n")
	nodeID = 0
	writeNode(file, root)
	file.WriteString("}\n")
}

func writeNode(file *os.File, node *Node) int {
	if node == nil { return -1 }
	id := nodeID
	fmt.Fprintf(file, "  node%d [label=\"%s\"];\n", id, node.Value)
	nodeID++
	if node.Left != nil {
		leftID := writeNode(file, node.Left)
		fmt.Fprintf(file, "  node%d -> node%d;\n", id, leftID)
	}
	if node.Right != nil {
		rightID := writeNode(file, node.Right)
		fmt.Fprintf(file, "  node%d -> node%d;\n", id, rightID)
	}
	return id
}
