package main

import (
	"encoding/json"
	"fmt"
)

func expand(n *CSPSearchNode, p *BinaryCSP) (children []*CSPSearchNode) {
	for _, s := range p.getSuccessors(n) {
		children = append(children, s)
	}
	return children
}

func searchWithFringe(f *StackFringe, p *BinaryCSP) string {
	startState := p.getStartState()
	closedSet := make(map[*CSPSearchNode]bool)
	f.push(startState)

	for {
		if f.isEmpty() {
			return "Failed!"
		}

		nextNode := f.pop()
		if p.isGoalState(nextNode) {
			return nextNode.String()
		}

		if _, ok := closedSet[nextNode]; !ok {
			// node not in closed set
			closedSet[nextNode] = true
			for _, child := range expand(nextNode, p) {
				f.push(child)
			}
		}

	}
}

type StackFringe []*CSPSearchNode

func (f *StackFringe) push(n *CSPSearchNode) {
	*f = append((*f), n)
}

func (f *StackFringe) pop() *CSPSearchNode {
	n := (*f)[len(*f)-1]
	(*f) = (*f)[:len(*f)-1]
	return n
}

func (f *StackFringe) isEmpty() bool {
	return len(*f) == 0
}

type CSPVar string
type CSPVal string
type CSPPair struct {
	a CSPVar
	b CSPVar
}
type BinaryCSP struct {
	vars        []CSPVar
	domain      []CSPVal
	assignments map[CSPVar]CSPVal
	constraints map[CSPPair]func(CSPVal, CSPVal) bool
}

type CSPSearchNode struct {
	assignments map[CSPVar]CSPVal
}

func (n *CSPSearchNode) copy() (c *CSPSearchNode) {
	c = &CSPSearchNode{
		assignments: make(map[CSPVar]CSPVal),
	}
	for k, v := range n.assignments {
		c.assignments[k] = v
	}
	return c
}

func (p *BinaryCSP) getStartState() *CSPSearchNode {
	assignments := make(map[CSPVar]CSPVal)
	for _, v := range p.vars {
		assignments[v] = ""
	}
	return &CSPSearchNode{
		assignments: assignments,
	}
}

func (p *BinaryCSP) isGoalState(s *CSPSearchNode) bool {
	for pair, f := range p.constraints {
		if !f(s.assignments[pair.a], s.assignments[pair.b]) {
			return false
		}
	}
	for _, a := range s.assignments {
		if a == "" {
			return false
		}
	}
	return true
}

func (p *BinaryCSP) getSuccessors(n *CSPSearchNode) (successors []*CSPSearchNode) {
	for v, a := range n.assignments {
		if a == "" {
			for _, possibleValue := range p.domain {
				s := n.copy()
				s.assignments[v] = possibleValue
				successors = append(successors, s)
			}
		}
	}
	return successors
}

func (n *CSPSearchNode) expand(p *BinaryCSP) (expanded []*CSPSearchNode) {
	successors := p.getSuccessors(n)
	for _, s := range successors {
		expanded = append(expanded, n.step(s))
	}
	return expanded
}

func (n *CSPSearchNode) step(next *CSPSearchNode) *CSPSearchNode {
	return n
}

func (n *CSPSearchNode) makePath() []*CSPSearchNode {
	return []*CSPSearchNode{n}
}

func (n *CSPSearchNode) String() string {
	json, _ := json.Marshal(n.assignments)
	return string(json)
}

func main() {

	var neq = func(a CSPVal, b CSPVal) bool {
		return a != b
	}

	sampleProblem := &BinaryCSP{
		vars:   []CSPVar{"A", "B", "C", "D"},
		domain: []CSPVal{"Red", "Blue", "Green", "Yellow"},
		constraints: map[CSPPair]func(CSPVal, CSPVal) bool{
			CSPPair{
				a: "A",
				b: "B",
			}: neq,
			CSPPair{
				a: "B",
				b: "C",
			}: neq,
			CSPPair{
				a: "C",
				b: "D",
			}: neq,
			CSPPair{
				a: "D",
				b: "A",
			}: neq,
			CSPPair{
				a: "D",
				b: "B",
			}: neq,
			CSPPair{
				a: "A",
				b: "C",
			}: neq,
		},
	}

	fringe := &StackFringe{}
	fmt.Println(searchWithFringe(fringe, sampleProblem))
}
