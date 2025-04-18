package services

import (
	"github.com/yanosea/cts/internal/domain/entities"
)

// CombatService は戦闘関連のロジックを提供するのじゃ
type CombatService struct {
	DeckService *DeckService
}

// NewCombatService はCombatServiceのインスタンスを生成するのじゃ
func NewCombatService(deckService *DeckService) *CombatService {
	return &CombatService{
		DeckService: deckService,
	}
}

// UseCard はカードを使用するのじゃ
func (s *CombatService) UseCard(player *entities.Player, enemy *entities.Enemy, cardIndex int) bool {
	if cardIndex < 0 || cardIndex >= len(player.Hand) {
		return false
	}

	card := player.Hand[cardIndex]
	if player.Energy < card.EnergyCost {
		return false
	}

	// カードの効果を実行するのじゃ
	card.Effect(player, enemy)
	player.Energy -= card.EnergyCost

	// 使用したカードを捨て札に移すのじゃ
	player.DiscardPile = append(player.DiscardPile, card)
	player.Hand = append(player.Hand[:cardIndex], player.Hand[cardIndex+1:]...)

	return true
}

// DrawCards はカードを引くのじゃ
func (s *CombatService) DrawCards(player *entities.Player, count int) {
	for i := 0; i < count; i++ {
		// ドローパイルが空なら、捨て札をシャッフルしてドローパイルにするのじゃ
		if len(player.DrawPile) == 0 {
			player.DrawPile = player.DiscardPile
			player.DiscardPile = []entities.Card{}
			s.DeckService.ShuffleDeck(player.DrawPile)
		}

		// ドローパイルが空でなければ、1枚引くのじゃ
		if len(player.DrawPile) > 0 {
			player.Hand = append(player.Hand, player.DrawPile[0])
			player.DrawPile = player.DrawPile[1:]
		}
	}
}
