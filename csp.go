package csp

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

func searchWithFringe(f Fringe, p *BinaryCSP) string {
	startState := p.getStartState()
	closedSet := make(map[*CSPSearchNode]bool)
	f.push(startState)
	expanded := 0

	for {
		if f.isEmpty() {
			return "Failed!"
		}

		nextNode := f.pop()
		if p.isGoalState(nextNode) {
			return fmt.Sprintf("Expanded %d nodes and found: %s", expanded, nextNode.String())
		}

		if _, ok := closedSet[nextNode]; !ok {
			// node not in closed set
			closedSet[nextNode] = true
			for _, child := range expand(nextNode, p) {
				f.push(child)
			}
			expanded++
		}

	}
}

type Fringe interface {
	push(n *CSPSearchNode)
	pop() *CSPSearchNode
	isEmpty() bool
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

type MostFreeFringe struct {
	nodes   []*CSPSearchNode
	problem *BinaryCSP
}

func (f *MostFreeFringe) push(n *CSPSearchNode) {
	f.nodes = append(f.nodes, n)
}

func (f *MostFreeFringe) pop() (n *CSPSearchNode) {
	return n
}

func (f *MostFreeFringe) isEmpty() bool {
	return len(f.nodes) == 0
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
	options     map[string]bool
}

type CSPSearchNode struct {
	assignments    map[CSPVar]CSPVal
	forwardDomains map[CSPVar][]CSPVal
}

func (n *CSPSearchNode) copy() (c *CSPSearchNode) {
	c = &CSPSearchNode{
		assignments:    make(map[CSPVar]CSPVal),
		forwardDomains: make(map[CSPVar][]CSPVal),
	}
	for k, v := range n.assignments {
		c.assign(k, v)
	}
	for k, v := range n.forwardDomains {
		c.forwardDomains[k] = v
	}
	return c
}

func (n *CSPSearchNode) assign(v CSPVar, val CSPVal) {
	n.assignments[v] = val
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

func (p *BinaryCSP) isLegal(s *CSPSearchNode) bool {
	for pair, f := range p.constraints {
		if s.assignments[pair.a] != "" && s.assignments[pair.b] != "" &&
			!f(s.assignments[pair.a], s.assignments[pair.b]) {
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
				s.assign(v, possibleValue)
				// FAIL ON VIOLATION
				if (p.options["checkViolation"] && p.isLegal(s)) || !p.options["checkViolation"] {
					successors = append(successors, s)
				}
			}
			// ORDERING
			if p.options["ordering"] {
				break
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
