package testutils

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func MustGetHclExpression(s string) hcl.Expression {
	expr, d := hclsyntax.ParseExpression([]byte("\""+s+"\""), "", hcl.InitialPos)
	if d.HasErrors() {
		panic("Expected no error for valid HCL expression")
	}
	return expr
}
