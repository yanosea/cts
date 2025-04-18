package entities

// PowerEffect はパワーの効果を表す関数型じゃ
type PowerEffect func(*Player, *Enemy)

// Power はゲーム中のパワー効果を定義するのじゃ
type Power struct {
	Name          string
	Description   string
	Duration      int // -1は永続的なパワーを意味するのじゃ
	OnTurnStart   PowerEffect
	OnTurnEnd     PowerEffect
	OnCardPlayed  PowerEffect
	OnDamageTaken PowerEffect
	OnDamageGiven PowerEffect
}
