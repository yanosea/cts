package entities

// Enemy は敵の状態を保持する構造体じゃ
type Enemy struct {
	Name       string
	Health     int
	MaxHealth  int
	Block      int
	Intention  string
	Damage     int
	Strength   int
	Vulnerable int
	Weak       int
	NextAction func(*Enemy, *Player)   // 敵の次の行動を定義する関数じゃ
	Patterns   []func(*Enemy, *Player) // 敵の行動パターンのリストじゃ
	PatternIdx int                     // 現在の行動パターンのインデックスじゃ
}

// NewSlimeEnemy はスライム敵のインスタンスを生成するのじゃ
func NewSlimeEnemy() *Enemy {
	// 攻撃パターン
	attackAction := func(e *Enemy, p *Player) {
		e.Intention = "攻撃"
		e.Damage = 5 + e.Strength
		p.ApplyDamage(e.Damage)
	}

	// 防御パターン
	defendAction := func(e *Enemy, p *Player) {
		e.Intention = "防御"
		e.Damage = 0
		e.AddBlock(5)
	}

	enemy := &Enemy{
		Name:       "スライム",
		Health:     20,
		MaxHealth:  20,
		Block:      0,
		Intention:  "攻撃",
		Damage:     5,
		Strength:   0,
		Vulnerable: 0,
		Weak:       0,
		Patterns:   []func(*Enemy, *Player){attackAction, defendAction},
		PatternIdx: 0,
	}

	// 初期行動をセット
	enemy.NextAction = enemy.Patterns[0]

	return enemy
}

// NewJawWormEnemy はアゴムシ敵のインスタンスを生成するのじゃ
func NewJawWormEnemy() *Enemy {
	// 通常攻撃パターン
	attackAction := func(e *Enemy, p *Player) {
		e.Intention = "攻撃"
		e.Damage = 11 + e.Strength
		p.ApplyDamage(e.Damage)
	}

	// 防御パターン
	defendAction := func(e *Enemy, p *Player) {
		e.Intention = "防御"
		e.Damage = 0
		e.AddBlock(6)
	}

	// 強化パターン
	buffAction := func(e *Enemy, p *Player) {
		e.Intention = "強化"
		e.Damage = 0
		e.Strength += 3
		e.AddBlock(6)
	}

	enemy := &Enemy{
		Name:       "アゴムシ",
		Health:     40,
		MaxHealth:  40,
		Block:      0,
		Intention:  "強化",
		Damage:     0,
		Strength:   0,
		Vulnerable: 0,
		Weak:       0,
		Patterns:   []func(*Enemy, *Player){buffAction, attackAction, defendAction},
		PatternIdx: 0,
	}

	// 初期行動をセット
	enemy.NextAction = enemy.Patterns[0]

	return enemy
}

// ApplyDamage は敵にダメージを与えるのじゃ
func (e *Enemy) ApplyDamage(damage int) {
	if e.Block >= damage {
		e.Block -= damage
	} else {
		dmgAfterBlock := damage - e.Block
		e.Block = 0
		e.Health -= dmgAfterBlock
	}
}

// AddBlock は敵のブロック値を増加させるのじゃ
func (e *Enemy) AddBlock(amount int) {
	e.Block += amount
}

// IsDefeated は敵が倒されたかどうかを判定するのじゃ
func (e *Enemy) IsDefeated() bool {
	return e.Health <= 0
}

// PerformAction は敵のアクションをプレイヤーに対して実行するのじゃ
func (e *Enemy) PerformAction(player *Player) {
	e.NextAction(e, player)

	// 次の行動パターンを選択するのじゃ
	e.PatternIdx = (e.PatternIdx + 1) % len(e.Patterns)
	e.NextAction = e.Patterns[e.PatternIdx]
}

// ApplyVulnerable は脆弱を付与するのじゃ
func (e *Enemy) ApplyVulnerable(amount int) {
	e.Vulnerable += amount
}

// ApplyWeak は弱体化を付与するのじゃ
func (e *Enemy) ApplyWeak(amount int) {
	e.Weak += amount
}

// AddStrength は筋力を増加させるのじゃ
func (e *Enemy) AddStrength(amount int) {
	e.Strength += amount
}
