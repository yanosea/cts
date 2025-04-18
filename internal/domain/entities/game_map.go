package entities

import (
	"math/rand"
)

// GameMap はゲームの全体マップを表す構造体じゃ
type GameMap struct {
	Nodes       [][]*MapNode // フロアごとのノードの配列じゃ
	CurrentNode *MapNode
}

// NewGameMap はゲームマップの新しいインスタンスを生成するのじゃ
func NewGameMap(floorCount int, nodesPerFloor int) *GameMap {
	gameMap := &GameMap{
		Nodes: make([][]*MapNode, floorCount),
	}

	// 各フロアにノードを配置するのじゃ
	for floor := 0; floor < floorCount; floor++ {
		gameMap.Nodes[floor] = make([]*MapNode, nodesPerFloor)

		for x := 0; x < nodesPerFloor; x++ {
			// ノードタイプを決定するのじゃ
			var nodeType NodeType
			if floor == floorCount-1 {
				nodeType = NodeBoss
			} else if floor == 0 {
				nodeType = NodeEnemy // 最初のフロアは必ず敵
			} else {
				// ランダムにノードタイプを設定するのじゃ
				r := rand.Intn(100)
				switch {
				case r < 60:
					nodeType = NodeEnemy
				case r < 70:
					nodeType = NodeElite
				case r < 85:
					nodeType = NodeRest
				case r < 95:
					nodeType = NodeShop
				default:
					nodeType = NodeEvent
				}
			}

			// ノードを作成するのじゃ
			gameMap.Nodes[floor][x] = NewMapNode(nodeType, floor, x)
		}
	}

	// ノード間の接続を設定するのじゃ
	for floor := 0; floor < floorCount-1; floor++ {
		for x := 0; x < nodesPerFloor; x++ {
			currentNode := gameMap.Nodes[floor][x]

			// 次のフロアの接続可能なノードをランダムに2-3個選ぶのじゃ
			connectCount := rand.Intn(2) + 2 // 2か3
			for i := 0; i < connectCount; i++ {
				targetX := rand.Intn(nodesPerFloor)
				currentNode.AddConnection(gameMap.Nodes[floor+1][targetX])
			}
		}
	}

	// スタート地点を設定するのじゃ
	gameMap.CurrentNode = gameMap.Nodes[0][rand.Intn(nodesPerFloor)]
	gameMap.CurrentNode.Visited = true

	return gameMap
}

// MoveToNode はマップ上の指定されたノードに移動するのじゃ
func (m *GameMap) MoveToNode(node *MapNode) bool {
	// 現在のノードから接続されていれば移動可能じゃ
	for _, connection := range m.CurrentNode.Connections {
		if connection == node {
			m.CurrentNode = node
			node.Visited = true
			return true
		}
	}
	return false
}
