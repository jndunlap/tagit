package processor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

func Run(root string, dry bool, logf func(string, ...any)) error {
	return filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() {
			if skip(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !goFile(p) {
			return nil
		}
		return process(p, dry, logf)
	})
}

func skip(n string) bool {
	switch n {
	case "vendor", "testdata", "tmp":
		return true
	}
	return strings.HasPrefix(n, ".") && n != "."
}

func goFile(p string) bool { return strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go") }

func process(p string, dry bool, logf func(string, ...any)) error {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, p, nil, parser.ParseComments)
	if err != nil {
		logf("parse-fail %s %v", p, err)
		return nil
	}

	if !update(f, logf) {
		return nil
	}

	if dry {
		logf("dry %s", p)
		return nil
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fs, f); err != nil {
		return err
	}
	logf("write %s", p)
	return os.WriteFile(p, buf.Bytes(), 0644)
}

func update(file *ast.File, logf func(string, ...any)) bool {
	changed := false
	for _, d := range file.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, s := range gd.Specs {
			ts, ok := s.(*ast.TypeSpec)
			if !ok {
				continue
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if tagged(ts, gd) && fixStruct(st, logf) {
				changed = true
			}
		}
	}
	return changed
}

func tagged(ts *ast.TypeSpec, gd *ast.GenDecl) bool {
	return has(ts.Doc) || has(gd.Doc)
}

func has(cg *ast.CommentGroup) bool {
	if cg == nil {
		return false
	}
	for _, c := range cg.List {
		if strings.Contains(strings.ToLower(c.Text), "tagit") {
			return true
		}
	}
	return false
}

func fixStruct(st *ast.StructType, logf func(string, ...any)) bool {
	c := false
	for _, f := range st.Fields.List {
		if fixField(f, logf) {
			c = true
		}
	}
	return c
}

func fixField(f *ast.Field, logf func(string, ...any)) bool {
	if len(f.Names) == 0 {
		return false
	}
	name := f.Names[0].Name

	orig := ""
	raw := ""
	if f.Tag != nil {
		orig = f.Tag.Value
		raw = strings.Trim(orig, "`")
	}
	tag := reflect.StructTag(raw)
	if tag.Get("json") != "" && tag.Get("db") != "" {
		return false
	}

	snake := camel(name)
	var add []string
	if tag.Get("json") == "" {
		add = append(add, fmt.Sprintf(`json:"%s"`, snake))
	}
	if tag.Get("db") == "" {
		add = append(add, fmt.Sprintf(`db:"%s"`, snake))
	}
	if len(add) == 0 {
		return false
	}

	newTag := "`" + strings.TrimSpace(strings.Join(append([]string{raw}, add...), " ")) + "`"
	if newTag == orig {
		return false
	}

	f.Tag = &ast.BasicLit{Kind: token.STRING, Value: newTag}
	logf("field %s -> %s", name, newTag)
	return true
}

var r1 = regexp.MustCompile("(.)([A-Z][a-z]+)")
var r2 = regexp.MustCompile("([a-z0-9])([A-Z])")

func camel(s string) string {
	s = r1.ReplaceAllString(s, "${1}_${2}")
	s = r2.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(s)
}
