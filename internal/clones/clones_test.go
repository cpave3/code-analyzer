package clones

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cameron/code-evaluator/internal/lang"
	lgo "github.com/cameron/code-evaluator/internal/lang/golang"
	lts "github.com/cameron/code-evaluator/internal/lang/typescript"
	"github.com/cameron/code-evaluator/internal/model"
)

type parser interface {
	lang.Adapter
	lang.SyntaxParser
}

func callables(t *testing.T, ad parser, files map[string]string, min int) []Callable {
	t.Helper()
	var out []Callable
	for p, src := range files {
		pf, err := ad.Parse([]byte(src))
		if err != nil || len(pf.ParseErrors) > 0 {
			t.Fatalf("%s: %v %v", p, err, pf.ParseErrors)
		}
		var syms []*model.Symbol
		for i := range pf.Symbols {
			pf.Symbols[i].File = p
			syms = append(syms, &pf.Symbols[i])
		}
		tree, _ := ad.ParseSyntax([]byte(src))
		out = append(out, Extract(p, ad.ID(), []byte(src), tree.RootNode(), syms, min)...)
		tree.Close()
	}
	return out
}

const goA = `package a

func Total(items []Item) int {
	sum := 0
	for _, it := range items {
		if it.Price > 100 {
			sum += it.Price * 2
		} else {
			sum += it.Price
		}
	}
	return sum
}

func Validate(in string) error {
	if in == "" {
		return errEmpty
	}
	parts := strings.Split(in, ",")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) > 64 {
			return fmt.Errorf("part %d too long: %q", i, p)
		}
		if strings.ContainsAny(p, "<>") {
			return fmt.Errorf("part %d has markup", i)
		}
	}
	return nil
}
`

// Renamed copy with one extra branch (near miss) and an unrelated function
// that embeds Validate's loop (fragment).
const goB = `package b

func Sum(rows []Row) int {
	acc := 0
	for _, r := range rows {
		if r.Cost > 50 {
			acc += r.Cost * 3
		} else if r.Cost < 0 {
			continue
		} else {
			acc += r.Cost
		}
	}
	return acc
}

func Import(raw string, sink func(string)) error {
	seen := map[string]bool{}
	fields := strings.Split(raw, ";")
	for j, f := range fields {
		f = strings.TrimSpace(f)
		if len(f) > 32 {
			return fmt.Errorf("field %d too long: %q", j, f)
		}
		if strings.ContainsAny(f, "{}") {
			return fmt.Errorf("field %d has braces", j)
		}
	}
	for _, f := range fields {
		if seen[f] {
			log.Printf("duplicate field %s", f)
			continue
		}
		seen[f] = true
		sink(strings.ToLower(f))
	}
	log.Printf("imported %d fields", len(seen))
	return nil
}
`

func TestGoNearMissAndFragment(t *testing.T) {
	cs := callables(t, lgo.New(), map[string]string{"a/a.go": goA, "b/b.go": goB}, 30)
	opt := DefaultOptions(30, 0.80)
	opt.FragmentTokens = 40
	near := NearMisses(cs, opt)
	if len(near) != 1 || near[0].A.Name != "Total" || near[0].B.Name != "Sum" || near[0].Similarity >= 1 || near[0].Similarity < 0.8 {
		t.Fatalf("near misses: %+v", near)
	}
	frags := Fragments(cs, opt, near)
	if len(frags) != 1 || frags[0].A.Name != "Validate" || frags[0].B.Name != "Import" || frags[0].Shared < 40 {
		t.Fatalf("fragments: %+v", frags)
	}
	if frags[0].A.StartLine <= 15 || frags[0].B.StartLine <= 17 {
		t.Errorf("fragment lines should point inside the functions: %+v", frags[0])
	}
}

func TestTypeScriptNearMiss(t *testing.T) {
	a := `export function price(items: Item[]): number {
  let total = 0;
  for (const it of items) {
    if (it.qty > 10) { total += it.unit * it.qty * 0.9; } else { total += it.unit * it.qty; }
  }
  return Math.round(total * 100) / 100;
}
`
	b := `export function cost(lines: Line[]): number {
  let sum = 0;
  for (const l of lines) {
    if (l.count > 5) { sum += l.price * l.count * 0.8; } else { sum += l.price * l.count; }
    if (l.taxed) { sum += 1; }
  }
  return Math.round(sum * 100) / 100;
}
`
	cs := callables(t, lts.New(), map[string]string{"src/a.ts": a, "src/b.ts": b}, 30)
	near := NearMisses(cs, DefaultOptions(30, 0.80))
	if len(near) != 1 || near[0].Language != "typescript" {
		t.Fatalf("near: %+v", near)
	}
}

func TestIdenticalStructureAndUnrelatedAreSkipped(t *testing.T) {
	copyA := goA
	cs := callables(t, lgo.New(), map[string]string{"a/a.go": goA, "c/c.go": copyA}, 30)
	opt := DefaultOptions(30, 0.80)
	if near := NearMisses(cs, opt); len(near) != 0 {
		t.Errorf("structural clones must be left to the structural groups: %+v", near)
	}
	if frags := Fragments(cs, opt, nil); len(frags) != 0 {
		t.Errorf("fragments inside structural clones: %+v", frags)
	}
}

func TestIntroduced(t *testing.T) {
	p := func(k string) Pair { return Pair{Key: k} }
	if got := Introduced([]Pair{p("x")}, []Pair{p("x"), p("y")}); len(got) != 1 || got[0].Key != "y" {
		t.Fatal(got)
	}
}

func TestPeriodicRepetitionIsNotAFragment(t *testing.T) {
	var b strings.Builder
	b.WriteString("package a\n\nfunc Table() []Row {\n\treturn []Row{\n")
	for i := 0; i < 12; i++ {
		b.WriteString("\t\t{Name: \"x\", Kind: KindA, Weight: 1, Note: \"n\"},\n")
	}
	b.WriteString("\t}\n}\n")
	cs := callables(t, lgo.New(), map[string]string{"a/a.go": b.String()}, 30)
	if frags := Fragments(cs, DefaultOptions(30, 0.85), nil); len(frags) != 0 {
		t.Fatalf("a repeated row inside one callable is not a copy: %+v", frags)
	}
}

// Near-miss partners tie on similarity whenever several copies differ from
// a callable by the same amount. Keeping the top MaxPairsPer must not depend
// on map order, or two runs (and a base and target) disagree.
func TestNearMissesDeterministicOnTies(t *testing.T) {
	body := func(op string) string {
		return `{
	sum := 0
	for _, it := range items {
		if it.Price > 100 {
			sum += it.Price ` + op + ` 2
		} else {
			sum += it.Price
		}
		if it.Qty > 3 {
			sum += it.Qty
		}
	}
	return sum
}
`
	}
	files := map[string]string{}
	for i, op := range []string{"*", "+", "-", "/", "%", "&", "|", "^"} {
		files[fmt.Sprintf("f%d.go", i)] = "package a\n\nfunc F(items []Item) int " + body(op)
	}
	cs := callables(t, lgo.New(), files, 20)
	opt := DefaultOptions(20, 0.85)
	keys := func() string {
		var ks []string
		for _, p := range NearMisses(cs, opt) {
			ks = append(ks, fmt.Sprintf("%s:%d-%d", p.Key, p.A.StartLine, p.B.StartLine))
		}
		return strings.Join(ks, "\n")
	}
	first := keys()
	if first == "" {
		t.Fatal("expected near-miss pairs")
	}
	for i := 0; i < 30; i++ {
		if got := keys(); got != first {
			t.Fatalf("run %d differs:\n%s\nvs\n%s", i, got, first)
		}
	}
}
