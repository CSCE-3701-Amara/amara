package domain

import (
	"testing"
	"time"
)

func TestGraphValidation(t *testing.T) {
	stages := []Stage{}
	for _, id := range []string{"a", "b", "c", "d"} {
		stages = append(stages, Stage{ID: id, Name: id, Type: AssessmentStage, ModuleIDs: []string{"module"}, OpensAt: time.Now(), Deadline: time.Now().Add(time.Hour)})
	}
	g := Graph{Stages: stages, Links: []Link{{"a", "b"}, {"a", "c"}, {"b", "d"}, {"c", "d"}}}
	if err := g.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, e := range []Link{{"d", "a"}, {"a", "a"}, {"missing", "a"}, {"a", "b"}} {
		bad := g
		bad.Links = append(append([]Link{}, g.Links...), e)
		if bad.Validate() == nil {
			t.Errorf("accepted invalid edge %+v", e)
		}
	}
}

func TestGraphRejectsOutOfBoundsPosition(t *testing.T) {
	g := Graph{Stages: []Stage{{ID: "a", Name: "a", Type: AssessmentStage, ModuleIDs: []string{"module"}, OpensAt: time.Now(), Deadline: time.Now().Add(time.Hour), X: 100001}}}
	if g.Validate() == nil {
		t.Fatal("accepted oversized canvas")
	}
}
