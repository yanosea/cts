package usecase

import (
	"math/rand"

	"github.com/yanosea/cts/internal/domain/entities"
	"github.com/yanosea/cts/internal/domain/services"
)

// GameInteractor はゲーム全体のユースケースを実装するのじゃ
type GameInteractor struct {
	Player        *entities.Player
	Enemy         *entities.Enemy
	GameMap       *entities.GameMap
	CardRewards   []entities.Card
	State         entities.GameState
	DeckService   *services.DeckService
	CombatService *services.CombatService
	Done          bool
}

// NewGameInteractor はGameInteractorのインスタンスを生成するのじゃ
func NewGameInteractor() *GameInteractor {
	deckService := services.NewDeckService()
	combatService := services.NewCombatService(deckService)

	player := entities.NewPlayer()
	player.Deck = deckService.InitializeStarterDeck()

	// ゲームマップを生成するのじゃ
	gameMap := entities.NewGameMap(15, 4) // 15フロア、各フロア4ノード

	return &GameInteractor{
		Player:        player,
		Enemy:         nil,
		GameMap:       gameMap,
		CardRewards:   []entities.Card{},
		State:         entities.StateMap, // マップ画面から開始
		DeckService:   deckService,
		CombatService: combatService,
		Done:          false,
	}
}

// StartNewCombat は新しい戦闘を開始するのじゃ
func (i *GameInteractor) StartNewCombat() {
	// デッキを山札にセットするのじゃ
	i.Player.DrawPile = i.Player.Deck
	i.Player.DiscardPile = []entities.Card{}
	i.Player.Hand = []entities.Card{}

	// デッキをシャッフルするのじゃ
	i.DeckService.ShuffleDeck(i.Player.DrawPile)

	// 現在のマップノードタイプに基づいて敵を生成するのじゃ
	if i.GameMap.CurrentNode.Type == entities.NodeEnemy {
		// ランダムに通常の敵を選択するのじゃ
		if rand.Intn(2) == 0 {
			i.Enemy = entities.NewSlimeEnemy()
		} else {
			i.Enemy = entities.NewJawWormEnemy()
		}
	} else if i.GameMap.CurrentNode.Type == entities.NodeElite {
		// とりあえずアゴムシで代用
		i.Enemy = entities.NewJawWormEnemy()
		// 強化しておくのじゃ
		i.Enemy.MaxHealth += 20
		i.Enemy.Health += 20
		i.Enemy.AddStrength(2)
	} else if i.GameMap.CurrentNode.Type == entities.NodeBoss {
		// とりあえず強化版のアゴムシで代用
		i.Enemy = entities.NewJawWormEnemy()
		// 大幅強化するのじゃ
		i.Enemy.MaxHealth *= 3
		i.Enemy.Health = i.Enemy.MaxHealth
		i.Enemy.AddStrength(5)
		i.Enemy.Name = "超アゴムシ"
	}

	// バフ、デバフをリセットするのじゃ
	i.Player.Vulnerable = 0
	i.Player.Weak = 0

	// プレイヤーのエナジーをリセットするのじゃ
	i.Player.ResetEnergy()

	// 初期手札を引くのじゃ
	i.CombatService.DrawCards(i.Player, 5)

	// パワー効果を実行するのじゃ
	i.Player.ExecuteStartTurnPowers(i.Enemy)

	// 戦闘状態にセットするのじゃ
	i.State = entities.StateCombat
}

// SelectMapNode はマップ上のノードを選択するのじゃ
func (i *GameInteractor) SelectMapNode(node *entities.MapNode) bool {
	if i.GameMap.MoveToNode(node) {
		// ノードの種類に応じた状態に移行するのじゃ
		switch node.Type {
		case entities.NodeEnemy, entities.NodeElite, entities.NodeBoss:
			i.StartNewCombat()
		case entities.NodeRest:
			i.State = entities.StateRest
		case entities.NodeShop:
			i.State = entities.StateShop
		case entities.NodeEvent, entities.NodeTreasure:
			i.State = entities.StateEvent
		}
		return true
	}
	return false
}

// RestHeal は休憩所で回復するのじゃ
func (i *GameInteractor) RestHeal() {
	healAmount := i.Player.MaxHealth / 3
	i.Player.Health = min(i.Player.MaxHealth, i.Player.Health+healAmount)
	i.ReturnToMap()
}

// RestUpgrade は休憩所でカードをアップグレードするのじゃ（未実装）
func (i *GameInteractor) RestUpgrade() {
	// カードアップグレード機能はまだ未実装
	i.ReturnToMap()
}

// ReturnToMap はマップ画面に戻るのじゃ
func (i *GameInteractor) ReturnToMap() {
	i.State = entities.StateMap
}

// UseCard はカードを使用するのじゃ
func (i *GameInteractor) UseCard(cardIndex int) bool {
	success := i.CombatService.UseCard(i.Player, i.Enemy, cardIndex)

	// カードの追加ドロー効果を処理するのじゃ
	if success && i.Player.DrawCount > 0 {
		i.CombatService.DrawCards(i.Player, i.Player.DrawCount)
		i.Player.DrawCount = 0
	}

	// 敵の体力が0以下なら報酬画面へ移るのじゃ
	if success && i.Enemy.IsDefeated() {
		i.State = entities.StateReward

		// 敵の種類によって報酬を変えるのじゃ
		if i.GameMap.CurrentNode.Type == entities.NodeEnemy {
			i.Player.Gold += 10
		} else if i.GameMap.CurrentNode.Type == entities.NodeElite {
			i.Player.Gold += 25
		} else if i.GameMap.CurrentNode.Type == entities.NodeBoss {
			i.Player.Gold += 50
		}

		// カード報酬を生成するのじゃ
		i.CardRewards = i.DeckService.GetRandomCardReward()
	}

	return success
}

// EndTurn はターンを終了するのじゃ
func (i *GameInteractor) EndTurn() {
	// 手札を捨て札に移すのじゃ
	i.Player.DiscardPile = append(i.Player.DiscardPile, i.Player.Hand...)
	i.Player.Hand = []entities.Card{}

	// パワー効果を実行するのじゃ
	i.Player.ExecuteEndTurnPowers(i.Enemy)

	// バフ/デバフの効果時間を減少させるのじゃ
	if i.Player.Vulnerable > 0 {
		i.Player.Vulnerable--
	}
	if i.Player.Weak > 0 {
		i.Player.Weak--
	}
	if i.Enemy.Vulnerable > 0 {
		i.Enemy.Vulnerable--
	}
	if i.Enemy.Weak > 0 {
		i.Enemy.Weak--
	}

	// 敵のアクションを実行するのじゃ
	i.Enemy.PerformAction(i.Player)

	// プレイヤーの体力が0以下ならゲームオーバー
	if i.Player.IsDefeated() {
		i.State = entities.StateGameOver
	} else {
		// 新しいターンの準備をするのじゃ
		i.Player.ResetEnergy()
		i.CombatService.DrawCards(i.Player, 5)
		i.Player.ExecuteStartTurnPowers(i.Enemy)
	}
}

// SelectCardReward は報酬からカードを選択するのじゃ
func (i *GameInteractor) SelectCardReward(cardIndex int) bool {
	if cardIndex >= 0 && cardIndex < len(i.CardRewards) {
		// 選択したカードをデッキに追加するのじゃ
		i.Player.Deck = append(i.Player.Deck, i.CardRewards[cardIndex])
		i.CardRewards = []entities.Card{}
		i.ReturnToMap()
		return true
	}
	return false
}

// SkipCardReward はカード報酬をスキップするのじゃ
func (i *GameInteractor) SkipCardReward() {
	i.CardRewards = []entities.Card{}
	i.ReturnToMap()
}

// SetDone はゲーム終了フラグを設定するのじゃ
func (i *GameInteractor) SetDone(done bool) {
	i.Done = done
}

// IsDone はゲームが終了したかどうかを返すのじゃ
func (i *GameInteractor) IsDone() bool {
	return i.Done
}

// ヘルパー関数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
