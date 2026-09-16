package cogfield

import (
	"testing"
)

func TestDetectGaps(t *testing.T) {
	g := &Graph{
		Nodes: []Node{
			{ID: "A", Label: "Concept A", Strength: 8.0},
			{ID: "B", Label: "Concept B", Strength: 2.0}, // Structural gap candidate (low strength)
			{ID: "C", Label: "Concept C", Strength: 9.0},
			{ID: "D", Label: "Concept D", Strength: 7.0},
			{ID: "X", Label: "Concept X", Strength: 6.0}, // Intermediate
		},
		Edges: []Edge{
			// Structural gap for B: has incoming from A, no outgoing
			{Source: "A", Target: "B", Relation: "is_related_to", Weight: 6.0},
			
			// Semantic gap: C causes D directly with high weight, no intermediate
			{Source: "C", Target: "D", Relation: "causes", Weight: 8.0},
			
			// Missing gap: A points to E which doesn't exist
			{Source: "A", Target: "E", Relation: "references", Weight: 5.0},
			
			// Explained relation: A -> X -> D
			{Source: "A", Target: "D", Relation: "similar", Weight: 7.0}, // Direct
			{Source: "A", Target: "X", Relation: "part_of", Weight: 6.0}, // Path 1
			{Source: "X", Target: "D", Relation: "causes", Weight: 6.0},  // Path 2
		},
	}

	gaps := DetectGaps(g)

	var hasStructural, hasSemantic, hasMissing bool
	for _, gap := range gaps {
		switch gap.GapType {
		case GapStructural:
			if gap.SubjectID == "B" {
				hasStructural = true
			}
		case GapSemantic:
			if gap.SubjectID == "C" && gap.ObjectLabel == "Concept D" {
				hasSemantic = true
			}
			if gap.SubjectID == "A" && gap.ObjectLabel == "Concept D" {
				t.Errorf("A -> D should not be a semantic gap because it has an explanatory path A -> X -> D")
			}
		case GapMissing:
			if gap.SubjectLabel == "E" {
				hasMissing = true
			}
		}
	}

	if !hasStructural {
		t.Errorf("Expected to detect a structural gap for Concept B")
	}
	if !hasSemantic {
		t.Errorf("Expected to detect a semantic gap for Concept C causes Concept D")
	}
	if !hasMissing {
		t.Errorf("Expected to detect a missing gap for E")
	}
}
