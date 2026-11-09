package domain

import (
	"fmt"
	"math"
)

type Link struct {
	From string `json:"fromStageId"`
	To   string `json:"toStageId"`
}
type Graph struct {
	Stages []Stage `json:"stages"`
	Links  []Link  `json:"links"`
}

func (g Graph) Validate() error {
	if len(g.Stages) > 100 || len(g.Links) > 1000 {
		return &InvalidPipelineError{Reason: "Graph exceeds 100 stages or 1000 connections"}
	}
	if err := ValidateStages(g.Stages); err != nil {
		return err
	}
	nodes := map[string]bool{}
	adj := map[string][]string{}
	edges := map[Link]bool{}
	for _, s := range g.Stages {
		if math.IsNaN(s.X) || math.IsNaN(s.Y) || math.IsInf(s.X, 0) || math.IsInf(s.Y, 0) || s.X < 0 || s.Y < 0 || s.X > 100000 || s.Y > 100000 {
			return &InvalidPipelineError{Reason: "Stage coordinates must be between 0 and 100000"}
		}
		if s.ID == "" || nodes[s.ID] {
			return &InvalidPipelineError{Reason: "Stage keys must be unique"}
		}
		nodes[s.ID] = true
	}
	for _, e := range g.Links {
		if !nodes[e.From] || !nodes[e.To] || e.From == e.To || edges[e] {
			return &InvalidPipelineError{Reason: "Invalid or duplicate connection"}
		}
		edges[e] = true
		adj[e.From] = append(adj[e.From], e.To)
	}
	if len(g.Stages) > 1 && len(g.Links) == 0 {
		return &InvalidPipelineError{Reason: "Connect the stages before saving"}
	}
	state := map[string]int{}
	var visit func(string) bool
	visit = func(id string) bool {
		if state[id] == 1 {
			return false
		}
		if state[id] == 2 {
			return true
		}
		state[id] = 1
		for _, to := range adj[id] {
			if !visit(to) {
				return false
			}
		}
		state[id] = 2
		return true
	}
	for id := range nodes {
		if !visit(id) {
			return &InvalidPipelineError{Reason: fmt.Sprintf("Connection at %s would create a cycle", id)}
		}
	}
	return nil
}
