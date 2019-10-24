package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_GenerateAST(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		result  string
		wantErr bool
	}{
		{name: "broken", src: " ~~~ ", wantErr: true},
		{name: "empty source", src: "", result: "nil"},
		{name: "plain test", src: "test", result: `
[]ast.Stmt {
	&ast.ExprStmt {
		X: &ast.Ident {
			Name: "test",
		},
	},
}`},
		{name: "advanced", src: `
		span, ctx := tracing.StartSpanWithParentByContext(c.App.Context, "api4:bot:createBot")
		c.App.Context = ctx
		defer span.Finish()		
`,
			result: `
[]ast.Stmt {
	&ast.AssignStmt {
		Lhs: []ast.Expr {
			&ast.Ident {
				Name: "span",
			},
			&ast.Ident {
				Name: "ctx",
			},
		},
		Tok: :=,
		Rhs: []ast.Expr {
			&ast.CallExpr {
				Fun: &ast.SelectorExpr {
					X: &ast.Ident {
						Name: "tracing",
					},
					Sel: &ast.Ident {
						Name: "StartSpanWithParentByContext",
					},
				},
				Lparen: 79,
				Args: []ast.Expr {
					&ast.SelectorExpr {
						X: &ast.SelectorExpr {
							X: &ast.Ident {
								Name: "c",
							},
							Sel: &ast.Ident {
								Name: "App",
							},
						},
						Sel: &ast.Ident {
							Name: "Context",
						},
					},
					&ast.BasicLit {
						ValuePos: 95,
						Kind: STRING,
						Value: "\"api4:bot:createBot\"",
					},
				},
				Ellipsis: 0,
			},
		},
	},
	&ast.AssignStmt {
		Lhs: []ast.Expr {
			&ast.SelectorExpr {
				X: &ast.SelectorExpr {
					X: &ast.Ident {
						Name: "c",
					},
					Sel: &ast.Ident {
						Name: "App",
					},
				},
				Sel: &ast.Ident {
					Name: "Context",
				},
			},
		},
		Tok: =,
		Rhs: []ast.Expr {
			&ast.Ident {
				Name: "ctx",
			},
		},
	},
	&ast.DeferStmt {
		Defer: 141,
		Call: &ast.CallExpr {
			Fun: &ast.SelectorExpr {
				X: &ast.Ident {
					Name: "span",
				},
				Sel: &ast.Ident {
					Name: "Finish",
				},
			},
			Lparen: 158,
			Ellipsis: 0,
		},
	},
}			
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := generateAST(wrapInPackage(tt.src), writer); (err != nil) != tt.wantErr {
				t.Errorf("generateAST() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); strings.TrimSpace(gotWriter) != strings.TrimSpace(tt.result) {
				t.Errorf("generateAST() = {%v}, want {%v}", strings.TrimSpace(gotWriter), strings.TrimSpace(tt.result))
			}
		})
	}
}
