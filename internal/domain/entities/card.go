package entities

// CardEffect はカードの効果を表す関数型じゃ
type CardEffect func(*Player, *Enemy)

// CardRarity はカードのレア度を表す型じゃ
type CardRarity int

// カードレア度の定義
const (
	Common CardRarity = iota
	Uncommon
	Rare
)

// CardType はカードの種類を表す型じゃ
type CardType int

// カード種類の定義
const (
	AttackCard CardType = iota
	SkillCard
	PowerCard
)

// Card はカードの基本構造を定義じゃ
type Card struct {
	Name        string
	Description string
	EnergyCost  int
	Rarity      CardRarity
	Type        CardType
	Effect      CardEffect // カードの効果を実装する関数じゃ
}

// CreateStrikeCard は基本的な攻撃カードを生成するのじゃ
func CreateStrikeCard() Card {
	return Card{
		Name:        "ストライク",
		Description: "6ダメージを与える",
		EnergyCost:  1,
		Rarity:      Common,
		Type:        AttackCard,
		Effect: func(p *Player, e *Enemy) {
			e.ApplyDamage(6)
		},
	}
}

// CreateDefendCard は基本的な防御カードを生成するのじゃ
func CreateDefendCard() Card {
	return Card{
		Name:        "ディフェンド",
		Description: "5ブロックを得る",
		EnergyCost:  1,
		Rarity:      Common,
		Type:        SkillCard,
		Effect: func(p *Player, e *Enemy) {
			p.AddBlock(5)
		},
	}
}

// CreateBashCard は基本的なバッシュカードを生成するのじゃ
func CreateBashCard() Card {
	return Card{
		Name:        "バッシュ",
		Description: "8ダメージを与え、2脆弱を付与する",
		EnergyCost:  2,
		Rarity:      Common,
		Type:        AttackCard,
		Effect: func(p *Player, e *Enemy) {
			e.ApplyDamage(8)
			e.ApplyVulnerable(2)
		},
	}
}

// CreatePommelStrikeCard はコモンの攻撃カードを生成するのじゃ
func CreatePommelStrikeCard() Card {
	return Card{
		Name:        "ポンメルストライク",
		Description: "9ダメージを与え、カードを1枚引く",
		EnergyCost:  1,
		Rarity:      Common,
		Type:        AttackCard,
		Effect: func(p *Player, e *Enemy) {
			e.ApplyDamage(9)
			p.DrawCount += 1 // カード引き処理はCombatServiceで実行
		},
	}
}

// CreateShockwaveCard はアンコモンの攻撃カードを生成するのじゃ
func CreateShockwaveCard() Card {
	return Card{
		Name:        "衝撃波",
		Description: "全ての敵に3脆弱と3弱体化を付与する",
		EnergyCost:  2,
		Rarity:      Uncommon,
		Type:        SkillCard,
		Effect: func(p *Player, e *Enemy) {
			e.ApplyVulnerable(3)
			e.ApplyWeak(3)
		},
	}
}

// CreateInflameCard はアンコモンのパワーカードを生成するのじゃ
func CreateInflameCard() Card {
	return Card{
		Name:        "発火",
		Description: "筋力を2得る",
		EnergyCost:  1,
		Rarity:      Uncommon,
		Type:        PowerCard,
		Effect: func(p *Player, e *Enemy) {
			p.AddStrength(2)
		},
	}
}

// CreateLimitBreakCard はレアのスキルカードを生成するのじゃ
func CreateLimitBreakCard() Card {
	return Card{
		Name:        "限界突破",
		Description: "筋力を2倍にする",
		EnergyCost:  3,
		Rarity:      Rare,
		Type:        SkillCard,
		Effect: func(p *Player, e *Enemy) {
			p.SetStrength(p.Strength * 2)
		},
	}
}

// CreateDemonFormCard はレアのパワーカードを生成するのじゃ
func CreateDemonFormCard() Card {
	return Card{
		Name:        "悪魔化",
		Description: "ターン開始時に筋力を3得る",
		EnergyCost:  3,
		Rarity:      Rare,
		Type:        PowerCard,
		Effect: func(p *Player, e *Enemy) {
			p.AddPower(&Power{
				Name:        "悪魔化",
				Description: "ターン開始時に筋力を3得る",
				OnTurnStart: func(p *Player, e *Enemy) {
					p.AddStrength(3)
				},
			})
		},
	}
}

// CreateAttackCard は旧関数の互換性のためにストライクカードを返す
func CreateAttackCard() Card {
	return CreateStrikeCard()
}
