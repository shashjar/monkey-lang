package evaluator

import (
	"monkey/object"
)

var builtins = map[string]*object.BuiltIn{
	"puts":   object.GetBuiltInByName("puts"),
	"len":    object.GetBuiltInByName("len"),
	"first":  object.GetBuiltInByName("first"),
	"last":   object.GetBuiltInByName("last"),
	"rest":   object.GetBuiltInByName("rest"),
	"append": object.GetBuiltInByName("append"),
	"join":   object.GetBuiltInByName("join"),
	"split":  object.GetBuiltInByName("split"),
	"sum":    object.GetBuiltInByName("sum"),
}
