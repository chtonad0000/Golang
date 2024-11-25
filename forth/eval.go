//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type Evaluator struct {
	Stack          []int
	BaseCommands   map[string]func() (*Evaluator, error)
	CustomCommands map[string][]string
}

func (e *Evaluator) Plus() (*Evaluator, error) {
	if len(e.Stack) < 2 {
		return e, errors.New("evaluator stack has less than 2 elements")
	} else {
		a := e.Stack[len(e.Stack)-1]
		b := e.Stack[len(e.Stack)-2]
		e.Stack = e.Stack[:len(e.Stack)-2]
		e.Stack = append(e.Stack, a+b)
	}
	return e, nil
}

func (e *Evaluator) Subtract() (*Evaluator, error) {
	if len(e.Stack) < 2 {
		return e, errors.New("evaluator stack has less than 2 elements")
	} else {
		a := e.Stack[len(e.Stack)-1]
		b := e.Stack[len(e.Stack)-2]
		e.Stack = e.Stack[:len(e.Stack)-2]
		e.Stack = append(e.Stack, b-a)
	}
	return e, nil
}

func (e *Evaluator) Multiply() (*Evaluator, error) {
	if len(e.Stack) < 2 {
		return e, errors.New("evaluator stack has less than 2 elements")
	} else {
		a := e.Stack[len(e.Stack)-1]
		b := e.Stack[len(e.Stack)-2]
		e.Stack = e.Stack[:len(e.Stack)-2]
		e.Stack = append(e.Stack, a*b)
	}
	return e, nil
}

func (e *Evaluator) Divide() (*Evaluator, error) {
	if len(e.Stack) < 2 {
		return e, errors.New("evaluator stack has less than 2 elements")
	} else {
		a := e.Stack[len(e.Stack)-1]
		b := e.Stack[len(e.Stack)-2]
		e.Stack = e.Stack[:len(e.Stack)-2]
		if a == 0 {
			return e, errors.New("division by zero")
		}
		e.Stack = append(e.Stack, b/a)
	}
	return e, nil
}

func (e *Evaluator) Dup() (*Evaluator, error) {
	if len(e.Stack) < 1 {
		return e, errors.New("evaluator stack has less than 1 element")
	} else {
		a := e.Stack[len(e.Stack)-1]
		e.Stack = append(e.Stack, a)
	}
	return e, nil
}

func (e *Evaluator) Over() (*Evaluator, error) {
	if len(e.Stack) < 2 {
		return e, errors.New("evaluator stack has less than 2 elements")
	} else {
		b := e.Stack[len(e.Stack)-2]
		e.Stack = append(e.Stack, b)
	}
	return e, nil
}

func (e *Evaluator) Drop() (*Evaluator, error) {
	if len(e.Stack) < 1 {
		return e, errors.New("evaluator stack has less than 1 element")
	} else {
		e.Stack = e.Stack[:len(e.Stack)-1]
	}
	return e, nil
}

func (e *Evaluator) Swap() (*Evaluator, error) {
	if len(e.Stack) < 2 {
		return e, errors.New("evaluator stack has less than 2 elements")
	} else {
		a := e.Stack[len(e.Stack)-1]
		b := e.Stack[len(e.Stack)-2]
		e.Stack = e.Stack[:len(e.Stack)-2]
		e.Stack = append(e.Stack, a)
		e.Stack = append(e.Stack, b)
	}
	return e, nil
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	e := &Evaluator{
		Stack:          []int{},
		BaseCommands:   make(map[string]func() (*Evaluator, error)),
		CustomCommands: make(map[string][]string),
	}
	e.BaseCommands["+"] = e.Plus
	e.BaseCommands["-"] = e.Subtract
	e.BaseCommands["*"] = e.Multiply
	e.BaseCommands["/"] = e.Divide
	e.BaseCommands["dup"] = e.Dup
	e.BaseCommands["over"] = e.Over
	e.BaseCommands["drop"] = e.Drop
	e.BaseCommands["swap"] = e.Swap

	return e
}

func (e *Evaluator) AddCustom(command string) error {
	words := strings.Split(command[2:len(command)-2], " ")
	if len(words) <= 1 {
		return errors.New("invalid definition")
	}
	_, err := strconv.Atoi(words[0])
	if err == nil {
		return errors.New("trying redefine number")
	}
	words[0] = strings.ToLower(words[0])
	old := e.CustomCommands[words[0]]
	e.CustomCommands[words[0]] = []string{}
	for i := 1; i < len(words); i++ {
		a, err := strconv.Atoi(words[i])
		if err == nil {
			e.BaseCommands[words[i]] = func() (*Evaluator, error) {
				e.Stack = append(e.Stack, a)
				return e, nil
			}
		}
		words[i] = strings.ToLower(words[i])
		if words[i] == words[0] {
			e.CustomCommands[words[0]] = append(e.CustomCommands[words[0]], old...)
		} else if _, exists := e.CustomCommands[words[i]]; exists {
			e.CustomCommands[words[0]] = append(e.CustomCommands[words[0]], e.CustomCommands[words[i]]...)
		} else {
			e.CustomCommands[words[0]] = append(e.CustomCommands[words[0]], words[i])
		}
	}
	return nil
}

func (e *Evaluator) SolveCustom(command string) error {
	for _, word := range e.CustomCommands[command] {
		if opFunc, exists := e.BaseCommands[word]; exists {
			_, er := opFunc()
			if er != nil {
				return er
			}
		} else {
			return errors.New("unknown command")
		}
	}

	return nil
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	if len(row) >= 5 && row[:2] == ": " && row[len(row)-2:] == " ;" {
		err := e.AddCustom(row)
		if err != nil {
			return nil, err
		}
	} else {
		operations := strings.Split(row, " ")
		for _, operation := range operations {
			i, err := strconv.Atoi(operation)
			if err == nil {
				e.Stack = append(e.Stack, i)
			} else {
				operation = strings.ToLower(operation)
				if _, exists := e.CustomCommands[operation]; exists {
					err := e.SolveCustom(operation)
					if err != nil {
						return nil, err
					}
				} else if opFunc, exists := e.BaseCommands[operation]; exists {
					_, er := opFunc()
					if er != nil {
						return nil, er
					}
				} else {
					return nil, errors.New("unknown command")
				}
			}
		}
	}
	return e.Stack, nil
}
