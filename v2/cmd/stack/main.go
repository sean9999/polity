package main

import (
	"fmt"
	"strconv"
	"strings"
)

type stack []int

// func (p *stack) push(n int) {
// 	s := *p
// 	s = append(s, n)
// }

// func (p *stack) pop() int {
// 	s := *p
// 	val := s[len(s)-1]
// 	s = s[:len(s)-1]
// 	return val
// }

func push(s stack, n int) stack {
	s = append(s, n)
	return s
}

func pop(s stack) (int, stack, bool) {
	if len(s) == 0 {
		return 0, nil, false
	}
	lastval := s[len(s)-1]
	s = s[:len(s)-1]
	return lastval, s, true
}

func EvaluateReversePolishNotation(expression string) int {
	// TODO: Initialize a slice to simulate a stack for holding integer values

	stack := make(stack, 0)

	// TODO: Split the expression into tokens using whitespace as the delimiter
	tokens := strings.Split(expression, " ")

	// TODO: Iterate over each token in the split expression

	for _, token := range tokens {
		if token == "+" || token == "-" {
			n2, newStack, _ := pop(stack)
			stack = newStack
			n1, newStack, _ := pop(stack)
			stack = newStack
			if token == "+" {
				stack = push(stack, n1+n2)
			} else {
				stack = push(stack, n1-n2)
			}
		} else {
			n, _ := strconv.Atoi(token)
			stack = push(stack, n)
		}
	}

	// TODO: If the token is an operator ('+' or '-'), pop the top two elements from the stack,
	// perform the corresponding operation, and push the result back onto the stack

	// TODO: If the token is an operand, parse it to an integer and push it onto the stack

	xx := stack[0]
	fmt.Println(xx)

	// TODO: Return the final result, which should be the only element left in the stack
	return stack[0] // Placeholder return statement
}

func main() {
	// The expression "1 2 + 4 -" is "(1 + 2) - 4"
	fmt.Println(EvaluateReversePolishNotation("1 2 + 4 -")) // Expected output: -1
}
