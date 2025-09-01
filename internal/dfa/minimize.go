package dfa

// Minimización por "table-filling" (marca de pares).
// Requiere DFA totalizado (este FromNFA ya lo totaliza).

func Minimize(d *DFA) *DFA {
	n := d.NumStates
	// mark[i][j] where i<j
	mark := make([][]bool, n)
	for i := range mark {
		mark[i] = make([]bool, n)
	}

	// Inicial: marcar pares donde uno es aceptador y el otro no
	for i := 0; i < n; i++ {
		for j := i+1; j < n; j++ {
			if d.Accepting[i] != d.Accepting[j] {
				mark[i][j] = true
			}
		}
	}

	changed := true
	for changed {
		changed = false
		for i := 0; i < n; i++ {
			for j := i+1; j < n; j++ {
				if mark[i][j] { continue }
				// si para algún símbolo las imágenes están marcadas, marcar (i,j)
				dist := false
				for _, a := range d.Alphabet {
					pi := d.Trans[i][a]
					pj := d.Trans[j][a]
					ii, jj := pi, pj
					if ii > jj { ii, jj = jj, ii }
					if mark[ii][jj] {
						dist = true
						break
					}
				}
				if dist {
					mark[i][j] = true
					changed = true
				}
			}
		}
	}

	// Construir clases de equivalencia
	class := make([]int, n)
	for i := 0; i < n; i++ { class[i] = -1 }
	numC := 0
	for i := 0; i < n; i++ {
		if class[i] != -1 { continue }
		class[i] = numC
		for j := i+1; j < n; j++ {
			if mark[i][j] { continue }
			class[j] = numC
		}
		numC++
	}

	// Transiciones del DFA mínimo
	minTrans := map[int]map[string]int{}
	minAccept := map[int]bool{}
	for s := 0; s < n; s++ {
		cs := class[s]
		if minTrans[cs] == nil {
			minTrans[cs] = map[string]int{}
			minAccept[cs] = d.Accepting[s]
		} else if d.Accepting[s] {
			minAccept[cs] = true
		}
		for _, a := range d.Alphabet {
			ns := d.Trans[s][a]
			minTrans[cs][a] = class[ns]
		}
	}
	// Calcular estado inicial en clase
	startC := class[d.Start]

	// Normalizar ids compactos 0..k-1
	remap := map[int]int{}
	next := 0
	for c := 0; c < numC; c++ {
		if _, ok := minTrans[c]; ok {
			remap[c] = next
			next++
		}
	}
	finalTrans := map[int]map[string]int{}
	finalAccept := map[int]bool{}
	for c, t := range minTrans {
		nc := remap[c]
		finalTrans[nc] = map[string]int{}
		finalAccept[nc] = minAccept[c]
		for a, dst := range t {
			finalTrans[nc][a] = remap[dst]
		}
	}
	return &DFA{
		Start: remap[startC],
		Accepting: finalAccept,
		Trans: finalTrans,
		Alphabet: d.Alphabet,
		NumStates: len(finalTrans),
	}
}
