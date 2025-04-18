package entities

// Player はプレイヤーの状態を保持する構造体じゃ
type Player struct {
	Health      int
	MaxHealth   int
	Gold        int
	Block       int
	Deck        []Card
	Hand        []Card
	DrawPile    []Card
	DiscardPile []Card
	Energy      int
	MaxEnergy   int
	Strength    int
	Dexterity   int
	Vulnerable  int
	Weak        int
	DrawCount   int // ターン終了時に追加でドローするカード枚数
	Powers      []*Power
}

// NewPlayer はプレイヤーの新しいインスタンスを生成するのじゃ
func NewPlayer() *Player {
	return &Player{
		Health:      80,
		MaxHealth:   80,
		Gold:        0,
		Block:       0,
		Deck:        []Card{},
		Hand:        []Card{},
		DrawPile:    []Card{},
		DiscardPile: []Card{},
		Energy:      3,
		MaxEnergy:   3,
		Strength:    0,
		Dexterity:   0,
		Vulnerable:  0,
		Weak:        0,
		DrawCount:   0,
		Powers:      []*Power{},
	}
}

// ApplyDamage はプレイヤーにダメージを与えるのじゃ
func (p *Player) ApplyDamage(damage int) {
	if p.Block >= damage {
		p.Block -= damage
	} else {
		dmgAfterBlock := damage - p.Block
		p.Block = 0
		p.Health -= dmgAfterBlock
	}
}

// AddBlock はプレイヤーのブロック値を増加させるのじゃ
func (p *Player) AddBlock(amount int) {
	p.Block += amount
}

// ResetEnergy はプレイヤーのエナジーを最大値に戻すのじゃ
func (p *Player) ResetEnergy() {
	p.Energy = p.MaxEnergy
}

// IsDefeated はプレイヤーが敗北したかどうかを判定するのじゃ
func (p *Player) IsDefeated() bool {
	return p.Health <= 0
}

// AddStrength は筋力を増減させるのじゃ
func (p *Player) AddStrength(amount int) {
	p.Strength += amount
}

// SetStrength は筋力を設定するのじゃ
func (p *Player) SetStrength(amount int) {
	p.Strength = amount
}

// AddDexterity は敏捷性を増減させるのじゃ
func (p *Player) AddDexterity(amount int) {
	p.Dexterity += amount
}

// ApplyVulnerable は脆弱を付与するのじゃ
func (p *Player) ApplyVulnerable(amount int) {
	p.Vulnerable += amount
}

// ApplyWeak は弱体化を付与するのじゃ
func (p *Player) ApplyWeak(amount int) {
	p.Weak += amount
}

// AddPower はパワーを追加するのじゃ
func (p *Player) AddPower(power *Power) {
	p.Powers = append(p.Powers, power)
}

// ExecuteStartTurnPowers はターン開始時のパワー効果を実行するのじゃ
func (p *Player) ExecuteStartTurnPowers(enemy *Enemy) {
	for _, power := range p.Powers {
		if power.OnTurnStart != nil {
			power.OnTurnStart(p, enemy)
		}

		// 効果時間のあるパワーはカウントダウンするのじゃ
		if power.Duration > 0 {
			power.Duration--
		}
	}

	// 持続時間が0のパワーを削除するのじゃ
	newPowers := []*Power{}
	for _, power := range p.Powers {
		if power.Duration != 0 {
			newPowers = append(newPowers, power)
		}
	}
	p.Powers = newPowers
}

// ExecuteEndTurnPowers はターン終了時のパワー効果を実行するのじゃ
func (p *Player) ExecuteEndTurnPowers(enemy *Enemy) {
	for _, power := range p.Powers {
		if power.OnTurnEnd != nil {
			power.OnTurnEnd(p, enemy)
		}
	}
}
