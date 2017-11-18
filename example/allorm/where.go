// Autogenerated by github.com/posener/orm
package allorm

// WhereString adds a condition on String to the WHERE statement
func WhereString(op Op, val string) Where {
	return newWhere(op, "string", val)
}

// WhereStringIn adds an IN condition on String to the WHERE statement
func WhereStringIn(vals ...string) Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return newMulWhere(OpIn, "string", args...)
}

// WhereStringBetween adds a BETWEEN condition on String to the WHERE statement
func WhereStringBetween(low, high string) Where {
	return newMulWhere(OpBetween, "string", low, high)
}

// WhereInt adds a condition on Int to the WHERE statement
func WhereInt(op Op, val int) Where {
	return newWhere(op, "int", val)
}

// WhereIntIn adds an IN condition on Int to the WHERE statement
func WhereIntIn(vals ...int) Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return newMulWhere(OpIn, "int", args...)
}

// WhereIntBetween adds a BETWEEN condition on Int to the WHERE statement
func WhereIntBetween(low, high int) Where {
	return newMulWhere(OpBetween, "int", low, high)
}

// WhereBool adds a condition on Bool to the WHERE statement
func WhereBool(op Op, val bool) Where {
	return newWhere(op, "bool", val)
}

// WhereBoolIn adds an IN condition on Bool to the WHERE statement
func WhereBoolIn(vals ...bool) Where {
	args := make([]interface{}, len(vals))
	for i := range vals {
		args[i] = vals[i]
	}
	return newMulWhere(OpIn, "bool", args...)
}

// WhereBoolBetween adds a BETWEEN condition on Bool to the WHERE statement
func WhereBoolBetween(low, high bool) Where {
	return newMulWhere(OpBetween, "bool", low, high)
}
