//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
)

type operation struct {
	f []func(e *Evaluator) error
}

func (e *Evaluator) eval(op string) func(e *Evaluator) error {
	if f, ok := e.op[strings.ToLower(op)]; ok {
		return f.eval
	} else {
		return func(e *Evaluator) error {
			num, err := strconv.Atoi(op)
			e.st = append(e.st, num)
			return err
		}
	}
}

func (op *operation) eval(e *Evaluator) error {
	for _, f := range op.f {
		err := f(e)
		if err != nil {
			return err
		}
	}
	return nil
}

type Evaluator struct {
	st []int
	op map[string]operation
}

func (e *Evaluator) getDelTwoSt() (int, int) {
	valFirst := e.st[len(e.st)-1]
	valSecond := e.st[len(e.st)-2]
	e.st = e.st[:len(e.st)-2]
	return valFirst, valSecond
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{op: map[string]operation{
		"+": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) < 2 {
						return errors.New("error")
					}
					valFirst, valSecond := e.getDelTwoSt()
					e.st = append(e.st, valSecond+valFirst)
					return nil
				},
			},
		},
		"-": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) < 2 {
						return errors.New("error")
					}
					valFirst, valSecond := e.getDelTwoSt()
					e.st = append(e.st, valSecond-valFirst)
					return nil
				},
			},
		},
		"*": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) < 2 {
						return errors.New("error")
					}
					valFirst, valSecond := e.getDelTwoSt()
					e.st = append(e.st, valSecond*valFirst)
					return nil
				},
			},
		},
		"/": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) < 2 {
						return errors.New("error")
					}
					valFirst, valSecond := e.getDelTwoSt()
					if valFirst == 0 {
						return errors.New("divide by zero")
					}
					e.st = append(e.st, valSecond/valFirst)
					return nil
				},
			},
		},
		"dup": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) == 0 {
						return errors.New("error")
					}
					e.st = append(e.st, e.st[len(e.st)-1])
					return nil
				},
			},
		},
		"over": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) < 2 {
						return errors.New("error")
					}
					val1, val2 := e.getDelTwoSt()
					e.st = append(e.st, val2, val1, val2)
					return nil
				},
			},
		},
		"drop": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) == 0 {
						return errors.New("error")
					}
					e.st = e.st[:len(e.st)-1]
					return nil
				},
			},
		},
		"swap": operation{
			f: []func(e *Evaluator) error{
				func(e *Evaluator) error {
					if len(e.st) < 2 {
						return errors.New("error")
					}
					val1, val2 := e.getDelTwoSt()
					e.st = append(e.st, val1, val2)
					return nil
				},
			},
		},
	}}
}

func (e *Evaluator) createNewOperation(word string, operations []string) {
	var op operation
	for _, curOp := range operations {
		op.f = append(op.f, e.eval(curOp))
	}
	e.op[strings.ToLower(word)] = op
}

func (e *Evaluator) initOperations(row string) error {
	initSectionWord := false
	initWord := false
	var newWord string
	posStartInitWord := -1
	operations := strings.Split(row, " ")
	for i, word := range operations {
		switch word {
		case ":":
			if initSectionWord {
				panic("MEOW1")
			}
			initSectionWord = true
		case ";":
			if !initSectionWord {
				panic("MEOW2")
			}
			if !initWord {
				panic("MEOW3")
			}
			if _, err := strconv.Atoi(newWord); err == nil {
				return errors.New("can't redefine number")
			}
			e.createNewOperation(newWord, operations[posStartInitWord+1:i])
			initSectionWord = false
		default:
			if initSectionWord {
				if !initWord {
					posStartInitWord = i
					initWord = true
					newWord = word
				}
			}
		}
	}
	return nil
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	if err := e.initOperations(row); err != nil {
		return e.st, err
	}
	initSectionWord := false
	operations := strings.Split(row, " ")
	for _, op := range operations {
		switch op {
		case ":":
			initSectionWord = true
		case ";":
			initSectionWord = false
		default:
			if !initSectionWord {
				err := e.eval(op)(e)
				if err != nil {
					return e.st, err
				}
			}
		}
	}
	return e.st, nil
}
