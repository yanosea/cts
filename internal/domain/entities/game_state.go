package entities

// GameState はゲームの現在の状態を表す型じゃ
type GameState int

const (
	StateMenu GameState = iota
	StateMap
	StateCombat
	StateReward
	StateRest
	StateShop
	StateEvent
	StateGameOver
)
