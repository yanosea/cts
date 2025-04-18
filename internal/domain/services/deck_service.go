package services

import (
	"math/rand"
	"time"

	"github.com/yanosea/cts/internal/domain/entities"
)

// DeckService はデッキ関連の操作を提供するのじゃ
type DeckService struct{}

// NewDeckService はDeckServiceのインスタンスを生成するのじゃ
func NewDeckService() *DeckService {
	return &DeckService{}
}

// InitializeStarterDeck は初期デッキを作成するのじゃ
func (s *DeckService) InitializeStarterDeck() []entities.Card {
	deck := make([]entities.Card, 0, 12)

	// ストライクカードを5枚追加
	for i := 0; i < 5; i++ {
		deck = append(deck, entities.CreateStrikeCard())
	}

	// ディフェンドカードを4枚追加
	for i := 0; i < 4; i++ {
		deck = append(deck, entities.CreateDefendCard())
	}

	// バッシュカードを1枚追加
	deck = append(deck, entities.CreateBashCard())

	// ポンメルストライクカードを2枚追加
	for i := 0; i < 2; i++ {
		deck = append(deck, entities.CreatePommelStrikeCard())
	}

	return deck
}

// GetRandomCardReward はランダムな報酬カードを3枚生成するのじゃ
func (s *DeckService) GetRandomCardReward() []entities.Card {
	reward := make([]entities.Card, 3)

	// レア度の確率: コモン70%, アンコモン25%, レア5%
	for i := 0; i < 3; i++ {
		rarity := rand.Intn(100)
		if rarity < 70 {
			// コモンカード
			cardType := rand.Intn(2)
			if cardType == 0 {
				reward[i] = entities.CreateStrikeCard()
			} else {
				reward[i] = entities.CreatePommelStrikeCard()
			}
		} else if rarity < 95 {
			// アンコモンカード
			cardType := rand.Intn(2)
			if cardType == 0 {
				reward[i] = entities.CreateShockwaveCard()
			} else {
				reward[i] = entities.CreateInflameCard()
			}
		} else {
			// レアカード
			cardType := rand.Intn(2)
			if cardType == 0 {
				reward[i] = entities.CreateLimitBreakCard()
			} else {
				reward[i] = entities.CreateDemonFormCard()
			}
		}
	}

	return reward
}

// ShuffleDeck はデッキをシャッフルするのじゃ
func (s *DeckService) ShuffleDeck(deck []entities.Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}
