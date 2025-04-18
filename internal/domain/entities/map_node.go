package entities

// NodeType はマップノードの種類を表す型じゃ
type NodeType int

// マップノードの種類を定義するのじゃ
const (
	NodeEnemy NodeType = iota
	NodeElite
	NodeBoss
	NodeRest
	NodeShop
	NodeTreasure
	NodeEvent
)

// MapNode はゲームマップの1つのノードを表す構造体じゃ
type MapNode struct {
	Type     NodeType
	Position struct {
		Floor int
		X     int
	}
	Connections []*MapNode // このノードから進めるノードのリストじゃ
	Visited     bool
}

// NewMapNode はMapNodeの新しいインスタンスを生成するのじゃ
func NewMapNode(nodeType NodeType, floor, x int) *MapNode {
	return &MapNode{
		Type: nodeType,
		Position: struct {
			Floor int
			X     int
		}{
			Floor: floor,
			X:     x,
		},
		Connections: []*MapNode{},
		Visited:     false,
	}
}

// AddConnection はこのノードから進めるノードを追加するのじゃ
func (n *MapNode) AddConnection(node *MapNode) {
	n.Connections = append(n.Connections, node)
}

// GetNodeTypeString はノードタイプを文字列で返すのじゃ
func (n *MapNode) GetNodeTypeString() string {
	switch n.Type {
	case NodeEnemy:
		return "敵"
	case NodeElite:
		return "エリート"
	case NodeBoss:
		return "ボス"
	case NodeRest:
		return "休憩"
	case NodeShop:
		return "店"
	case NodeTreasure:
		return "宝箱"
	case NodeEvent:
		return "イベント"
	default:
		return "不明"
	}
}
