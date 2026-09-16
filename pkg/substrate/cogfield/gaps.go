package cogfield

import (
	"fmt"
	"sort"
)

type GapType string

const (
	GapStructural GapType = "structural"
	GapSemantic   GapType = "semantic"
	GapMissing    GapType = "missing"
)

// SemanticGap represents a gap in the knowledge graph.
// It is used by the curiosity loop (REM phase) to generate hypotheses.
type SemanticGap struct {
	GapType           GapType `json:"gap_type"`
	SubjectID         string  `json:"subject_id"`
	SubjectLabel      string  `json:"subject_label"`
	Relation          string  `json:"relation,omitempty"`
	ObjectLabel       string  `json:"object_label,omitempty"`
	Severity          float64 `json:"severity"`
	Explanation       string  `json:"explanation"`
	SuggestedQuestion string  `json:"suggested_question"`
}

// DetectGaps scans the graph for anomalies and unresolved structures.
func DetectGaps(g *Graph) []SemanticGap {
	var gaps []SemanticGap

	nodeMap := make(map[string]Node, len(g.Nodes))
	for _, n := range g.Nodes {
		nodeMap[n.ID] = n
	}

	incomingEdges := make(map[string][]Edge)
	outgoingEdges := make(map[string][]Edge)
	for _, e := range g.Edges {
		outgoingEdges[e.Source] = append(outgoingEdges[e.Source], e)
		incomingEdges[e.Target] = append(incomingEdges[e.Target], e)
	}

	// 1. Structural Gaps
	// Node exists but is poorly understood (low strength, no outgoing edges, but has incoming edges)
	for _, n := range g.Nodes {
		inEdges := incomingEdges[n.ID]
		outEdges := outgoingEdges[n.ID]

		// Using 4.5 as threshold for "low strength" on a 0-10 scale
		if n.Strength < 4.5 && len(outEdges) == 0 && len(inEdges) > 0 {
			baseSeverity := (10.0 - n.Strength) / 10.0
			centralityBonus := 1.0 + float64(len(inEdges))*0.1
			severity := baseSeverity * centralityBonus
			if severity > 1.0 {
				severity = 1.0
			}

			rel := "?"
			srcLabel := "?"
			if len(inEdges) > 0 {
				rel = inEdges[0].Relation
				if srcNode, ok := nodeMap[inEdges[0].Source]; ok {
					srcLabel = srcNode.Label
				}
			}

			gaps = append(gaps, SemanticGap{
				GapType:           GapStructural,
				SubjectID:         n.ID,
				SubjectLabel:      n.Label,
				Relation:          rel,
				ObjectLabel:       "",
				Severity:          severity,
				Explanation:       fmt.Sprintf("%q is referenced (%d incoming) but strength=%.1f and no outgoing edges. We know it exists but cannot reason through it.", n.Label, len(inEdges), n.Strength),
				SuggestedQuestion: fmt.Sprintf("What is %q? How does it relate to %q?", n.Label, srcLabel),
			})
		}
	}

	// 2. Semantic Gaps
	// Edge exists with high weight but no intermediate nodes explain WHY the relation holds.
	// We'll approximate this by checking for direct edges of certain relations (causes, similar)
	// that lack a common neighbor path (A -> X -> B).
	for _, e := range g.Edges {
		if e.Weight < 5.0 {
			continue // low-confidence edges are structural issues
		}
		if e.Relation != "causes" && e.Relation != "similar" {
			continue // only check specific relation types for now
		}

		src, hasSrc := nodeMap[e.Source]
		tgt, hasTgt := nodeMap[e.Target]
		if !hasSrc || !hasTgt {
			continue
		}

		hasExplanation := hasExplanatoryPath(e.Source, e.Target, e.Relation, outgoingEdges)
		if !hasExplanation {
			severity := e.Weight / 10.0
			if e.Relation == "causes" {
				severity *= 1.4
			}
			if severity > 1.0 {
				severity = 1.0
			}

			gaps = append(gaps, SemanticGap{
				GapType:           GapSemantic,
				SubjectID:         src.ID,
				SubjectLabel:      src.Label,
				Relation:          e.Relation,
				ObjectLabel:       tgt.Label,
				Severity:          severity,
				Explanation:       fmt.Sprintf("We assert %q %s %q with weight %.1f, but no intermediate nodes explain WHY this holds.", src.Label, e.Relation, tgt.Label, e.Weight),
				SuggestedQuestion: fmt.Sprintf("Why does %q %s %q? What is the mechanism or intermediate step?", src.Label, e.Relation, tgt.Label),
			})
		}
	}

	// 3. Missing Gaps
	// Edge points to a Target ID that is NOT in the Nodes list.
	for _, e := range g.Edges {
		if _, ok := nodeMap[e.Target]; !ok {
			src, hasSrc := nodeMap[e.Source]
			srcLabel := "?"
			if hasSrc {
				srcLabel = src.Label
			}
			gaps = append(gaps, SemanticGap{
				GapType:           GapMissing,
				SubjectID:         "",
				SubjectLabel:      e.Target, // We use target ID as label since we don't have it
				Relation:          e.Relation,
				ObjectLabel:       srcLabel,
				Severity:          0.5,
				Explanation:       fmt.Sprintf("%q appears as a target for %q but has no corresponding node.", e.Target, srcLabel),
				SuggestedQuestion: fmt.Sprintf("What is %q?", e.Target),
			})
		}
	}

	// Deduplicate by (GapType, Subject, Object) and rank by severity
	seen := make(map[string]bool)
	var deduped []SemanticGap

	sort.Slice(gaps, func(i, j int) bool {
		return gaps[i].Severity > gaps[j].Severity
	})

	for _, g := range gaps {
		key := fmt.Sprintf("%s:%s:%s", g.GapType, g.SubjectLabel, g.ObjectLabel)
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, g)
		}
	}

	if len(deduped) > 20 {
		return deduped[:20]
	}
	return deduped
}

// hasExplanatoryPath attempts to find a 2-hop path (A -> X -> B) that could explain A -> B.
func hasExplanatoryPath(src, tgt, relation string, outgoingEdges map[string][]Edge) bool {
	// Simple BFS up to depth 2
	srcOut := outgoingEdges[src]
	for _, e1 := range srcOut {
		if e1.Target == tgt {
			continue // skip the direct edge itself
		}
		// intermediate node
		mid := e1.Target
		midOut := outgoingEdges[mid]
		for _, e2 := range midOut {
			if e2.Target == tgt {
				// Found a path A -> mid -> B
				// In a full implementation, we'd check if e1.Relation and e2.Relation
				// logically compose to explain 'relation'.
				return true
			}
		}
	}
	return false
}
