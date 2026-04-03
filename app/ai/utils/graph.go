package utils

// GraphFile is the adjacency-list based graph file the AI works with.
type GraphFile struct {
	SchemaVersion string      `json:"schema_version"`
	Tree          GraphTree   `json:"tree"`
	Nodes         []*GraphNode `json:"nodes"`
	Meta          GraphMeta   `json:"meta"`
}

// GraphTree stores graph-level identifiers.
type GraphTree struct {
	Name      string `json:"name"`
	TopNodeID string `json:"top_node_id"`
}

// GraphNode stores node content and adjacency relations.
type GraphNode struct {
	NodeID           string   `json:"node_id"`
	NodeType         string   `json:"node_type"`
	Label            string   `json:"label"`
	Description      string   `json:"description"`
	GateType         string   `json:"gate_type"`
	PointsTo         []string `json:"points_to"`
	PointedBy        []string `json:"pointed_by"`
	EvidenceChunkIDs []string `json:"evidence_chunk_ids"`
}

// GraphMeta stores versioning and provenance metadata for the tmp graph.
type GraphMeta struct {
	Version        string   `json:"version"`
	BasedOnVersion string   `json:"based_on_version"`
	GeneratedAt    string   `json:"generated_at"`
	SourceChunkIDs []string `json:"source_chunk_ids"`
}
