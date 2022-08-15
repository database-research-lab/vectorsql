// Copyright 2019 The OctoSQL Authors.
// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package planners

import (
	"reflect"
	"strconv"
	"strings"

	"base/errors"
	"parsers/sqlparser"
)

// 解析AST上的表达式为逻辑执行计划 
// alias: 
// expr: 表达式所对应的AST节点 
func parseExpression(aliases map[string]IPlan, expr sqlparser.Expr) (IPlan, error) {

	// 根据不同的表达式的类型会生成不同的执行计划 
	switch expr := expr.(type) {

	// 列名表达式 
	case *sqlparser.ColName:
		name := expr.Name.String()
		if aliases != nil {
			if p, ok := aliases[name]; ok {
				return p, nil
			}
		}
		return NewVariablePlan(name), nil

	// 字面值常量的形式 
	case *sqlparser.SQLVal:
		var err error
		var val interface{}
		
		// 根据常量的不同类型，转换为对应的值 
		switch expr.Type {
		case sqlparser.IntVal:
			var i int64
			i, err = strconv.ParseInt(string(expr.Val), 10, 64)
			val = int(i)
		case sqlparser.FloatVal:
			val, err = strconv.ParseFloat(string(expr.Val), 64)
		case sqlparser.StrVal:
			val = string(expr.Val)
		default:
			err = errors.Errorf("Constant value type unsupported")
		}
		if err != nil {
			return nil, err
		}
		// 然后返回一个ConstantPlan，执行的时候就原样返回自己持有的值 
		return NewConstantPlan(val), nil

	// 函数表达式，这样就可以加一些自定义的函数之类的 
	case *sqlparser.FuncExpr:
		// 函数名统一转换为大写 
		funcName := strings.ToUpper(expr.Name.String())
		switch len(expr.Exprs) {
		case 1:
			expr, err := parseFunctionArgument(aliases, expr.Exprs[0].(*sqlparser.AliasedExpr))
			if err != nil {
				return nil, err
			}
			return NewUnaryExpressionPlan(funcName, expr), nil
		case 2:
			left, err := parseFunctionArgument(aliases, expr.Exprs[0].(*sqlparser.AliasedExpr))
			if err != nil {
				return nil, err
			}
			right, err := parseFunctionArgument(aliases, expr.Exprs[1].(*sqlparser.AliasedExpr))
			if err != nil {
				return nil, err
			}
			return NewBinaryExpressionPlan(funcName, left, right), nil
		default:
			args := make([]IPlan, len(expr.Exprs))
			for i, expr := range expr.Exprs {
				aliased, ok := expr.(*sqlparser.AliasedExpr)
				if !ok {
					return nil, errors.Errorf("Unsupported argument %v of type %v", expr, reflect.TypeOf(expr))
				}
				arg, err := parseFunctionArgument(aliases, aliased)
				if err != nil {
					return nil, err
				}
				args[i] = arg
			}
			return NewFunctionExpressionPlan(funcName, args...), nil
		}

	case *sqlparser.BinaryExpr:
		left, err := parseExpression(aliases, expr.Left)
		if err != nil {
			return nil, err
		}
		right, err := parseExpression(aliases, expr.Right)
		if err != nil {
			return nil, err
		}
		return NewBinaryExpressionPlan(expr.Operator, left, right), nil

	case *sqlparser.ComparisonExpr:
		left, err := parseExpression(aliases, expr.Left)
		if err != nil {
			return nil, err
		}
		right, err := parseExpression(aliases, expr.Right)
		if err != nil {
			return nil, err
		}
		return NewBinaryExpressionPlan(expr.Operator, left, right), nil

	case *sqlparser.OrExpr:
		left, err := parseExpression(aliases, expr.Left)
		if err != nil {
			return nil, err
		}
		right, err := parseExpression(aliases, expr.Right)
		if err != nil {
			return nil, err
		}
		return NewBinaryExpressionPlan("OR", left, right), nil

	case *sqlparser.AndExpr:
		left, err := parseExpression(aliases, expr.Left)
		if err != nil {
			return nil, err
		}
		right, err := parseExpression(aliases, expr.Right)
		if err != nil {
			return nil, err
		}
		return NewBinaryExpressionPlan("AND", left, right), nil
		
	case *sqlparser.ParenExpr:
		return parseExpression(aliases, expr.Expr)
	}

	return nil, errors.Errorf("Unsupported expression %+v %+v", expr, reflect.TypeOf(expr))
}

func parseFunctionArgument(aliases map[string]IPlan, expr *sqlparser.AliasedExpr) (IPlan, error) {
	subExpr, err := parseExpression(aliases, expr.Expr)
	if err != nil {
		return nil, errors.Wrapf(err, "Couldn't parse argument")
	}
	return subExpr, nil
}

// 带别名的查询方式 
func parseAliasedTableExpression(expr *sqlparser.AliasedTableExpr) (IPlan, error) {
	switch subExpr := expr.Expr.(type) {
	case sqlparser.TableName:
		return NewScanPlan(subExpr.Name.String(), subExpr.Qualifier.String()), nil
	default:
		return nil, errors.Errorf("Unsupported aliased table expression:%+v", expr.Expr)
	}
}

func parseTableValuedFunction(aliases map[string]IPlan, expr *sqlparser.TableValuedFunction) (IPlan, error) {
	var arguments []IPlan
	name := expr.Name.String()

	for i := range expr.Args {
		argument, err := parseTableValuedFunctionArgument(aliases, expr.Args[i])
		if err != nil {
			return nil, errors.Wrapf(err, "Couldn't parse table valued function argument \"%v\" with index %v", expr.Args[i].Name.String(), i)
		}
		arguments = append(arguments, argument)
	}
	return NewTableValuedFunctionPlan(name, NewMapPlan(arguments...)), nil
}

func parseTableValuedFunctionArgument(aliases map[string]IPlan, expr *sqlparser.TableValuedFunctionArgument) (IPlan, error) {
	name := expr.Name.String()

	switch expr := expr.Value.(type) {
	case *sqlparser.ExprTableValuedFunctionArgumentValue:
		parsed, err := parseExpression(aliases, expr.Expr)
		if err != nil {
			return nil, errors.Wrapf(err, "Couldn't parse table valued function argument expression \"%v\"", expr.Expr)
		}
		return NewTableValuedFunctionExpressionPlan(name, parsed), nil
	default:
		return nil, errors.Errorf("Invalid table valued function argument: %v", expr)
	}
}

// 解析from查询 
func parseFrom(expr sqlparser.TableExpr) (IPlan, error) {
	// from可能会有多种形式 
	switch expr := expr.(type) {
	// 带别名的方式查询 
	case *sqlparser.AliasedTableExpr:
		return parseAliasedTableExpression(expr)
	// 这个是啥类型呢 
	case *sqlparser.ParenTableExpr:
		return parseFrom(expr.Exprs[0])
	// 这个又是啥类型呢 
	case *sqlparser.TableValuedFunction:
		return parseTableValuedFunction(nil, expr)
	default:
		return nil, errors.Errorf("Unsupported table expression:%+v", expr)
	}
}

// 解析查询的字段都有哪些，解析的时候会考虑变量表 
func parseFields(aliased map[string]IPlan, sel sqlparser.SelectExprs) (*MapPlan, error) {

	fields := NewMapPlan()

	// 如果不是*的话，则对每个字段做转换 
	if _, ok := sel[0].(*sqlparser.StarExpr); !ok {
		for i, expr := range sel {
			
			aliasedExpression, ok := expr.(*sqlparser.AliasedExpr)
			if !ok {
				return nil, errors.Errorf("Expected aliased expression in select on index:%v, got:%+v %+v", i, expr, reflect.TypeOf(expr))
			}

			child, err := parseExpression(aliased, aliasedExpression.Expr)
			if err != nil {
				return nil, err
			}

			if aliasedExpression.As.String() != "" {
				child = NewAliasedExpressionPlan(aliasedExpression.As.String(), child)
			}

			fields.Add(child)
		}
	}

	// 如果是*选择所有字段的的话，则直接返回，认为是一个字段就行了，后面ProjectionPlan能够识别得出来并转换  
	return fields, nil
}

// 把别名生成一个map是啥意思，是为了让别的地方能够引用得到这些别名吗 
func parseAliases(fields *MapPlan) (map[string]IPlan, error) {
	aliases := make(map[string]IPlan)
	if err := fields.Walk(func(plan IPlan) (bool, error) {
		switch p := plan.(type) {
		case *AliasedExpressionPlan:
			aliases[p.As] = p.Expr
		}
		return true, nil
	}); err != nil {
		return nil, err
	}
	return aliases, nil
}

func parseWhere(aliases map[string]IPlan, expr sqlparser.Expr) (IPlan, error) {
	return parseExpression(aliases, expr)
}

func parseGroupBy(aliases map[string]IPlan, groupby sqlparser.GroupBy) (*MapPlan, error) {
	all := NewMapPlan()
	for i, g := range groupby {
		expr, err := parseExpression(aliases, g)
		if err != nil {
			return nil, errors.Errorf("couldn't parse group by expression with index %v", i)
		}
		all.Add(expr)
	}
	return all, nil
}

func parseOrderBy(orderBy sqlparser.OrderBy) ([]Order, error) {
	orders := make([]Order, len(orderBy))

	for i, field := range orderBy {
		expr, err := parseExpression(nil, field.Expr)
		if err != nil {
			return nil, errors.Errorf("couldn't parse order by expression with index %v", i)
		}
		orders[i].Expression = expr
		orders[i].Direction = field.Direction
	}
	return orders, nil
}

func parseLimit(limit *sqlparser.Limit) (IPlan, error) {
	if limit.Offset == nil {
		limit.Offset = sqlparser.NewIntVal([]byte("0"))
	}
	offsetPlan, err := parseExpression(nil, limit.Offset)
	if err != nil {
		return nil, errors.Wrapf(err, "Couldn't parse offset")
	}
	rowcountPlan, err := parseExpression(nil, limit.Rowcount)
	if err != nil {
		return nil, errors.Wrapf(err, "Couldn't parse limit")
	}
	return NewLimitPlan(offsetPlan, rowcountPlan), nil
}
