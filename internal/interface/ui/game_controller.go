package ui

import (
	"fmt"

	"github.com/yanosea/cts/internal/usecase"
)

// GameController はゲームの入出力を制御するのじゃ
type GameController struct {
	screen         ScreenPort
	gameInteractor *usecase.GameInteractor
	// カーソル位置を保存する変数を追加
	cursorPosition int
	// カーソル行の最大値（選択肢の数など）
	cursorMaxPosition int
}

// NewGameController はGameControllerのインスタンスを生成するのじゃ
func NewGameController(screen ScreenPort, gameInteractor *usecase.GameInteractor) *GameController {
	return &GameController{
		screen:         screen,
		gameInteractor: gameInteractor,
		cursorPosition: 0,
		cursorMaxPosition: 0,
	}
}

// StartGame はゲームを開始するのじゃ
func (c *GameController) StartGame() {
	// 最初の戦闘を開始するのじゃ
	c.gameInteractor.StartNewCombat()

	// イベント処理を別のゴルーチンで実行するのじゃ
	go func() {
		for !c.gameInteractor.IsDone() {
			c.handleEvents()
		}
	}()

	// ゲームループじゃ
	for !c.gameInteractor.IsDone() {
		c.draw()
		c.screen.Sleep(16) // 約60FPSじゃ
	}
}

// handleEvents はイベントを処理する関数じゃ
func (c *GameController) handleEvents() {
	event := c.screen.PollEvent()

	// ESCキーまたはCtrl+Cでゲーム終了
	if event.IsExit() {
		c.gameInteractor.SetDone(true)
		return
	}

	// リサイズイベントの処理
	if event.IsResize() {
		c.screen.Show()
		return
	}

	// カーソル移動の処理
	c.handleCursorMovement(event)

	switch c.gameInteractor.State {
	case 1: // StateMap
		// マップ選択画面ではカーソルでノードを選択するのじゃ
		if event.IsEnter() || event.IsSpace() {
			if c.cursorPosition >= 0 && c.cursorPosition < len(c.gameInteractor.GameMap.CurrentNode.Connections) {
				c.gameInteractor.SelectMapNode(c.gameInteractor.GameMap.CurrentNode.Connections[c.cursorPosition])
				c.cursorPosition = 0 // カーソルをリセット
			}
		}

	case 2: // StateCombat
		// カーソルでカードを選択し、ENTERまたはSPACEでカード使用
		if event.IsEnter() || event.IsSpace() {
			if c.cursorPosition >= 0 && c.cursorPosition < len(c.gameInteractor.Player.Hand) {
				c.gameInteractor.UseCard(c.cursorPosition)
			}
		}

		// eキーまたはスペースでターン終了するのじゃ
		if event.IsEndTurn() {
			c.gameInteractor.EndTurn()
		}

	case 3: // StateReward
		// カーソルでカード報酬を選択するのじゃ
		if event.IsEnter() || event.IsSpace() {
			if c.cursorPosition >= 0 && c.cursorPosition < len(c.gameInteractor.CardRewards) {
				c.gameInteractor.SelectCardReward(c.cursorPosition)
				c.cursorPosition = 0 // カーソルをリセット
			}
		}

		// sキーで報酬をスキップするのじゃ
		if event.IsSKey() {
			c.gameInteractor.SkipCardReward()
			c.cursorPosition = 0 // カーソルをリセット
		}

	case 4: // StateRest
		// カーソル位置で選択肢を決定
		if event.IsEnter() || event.IsSpace() {
			if c.cursorPosition == 0 {
				c.gameInteractor.RestHeal()
			} else if c.cursorPosition == 1 {
				c.gameInteractor.RestUpgrade()
			}
			c.cursorPosition = 0 // カーソルをリセット
		}

	case 5: // StateShop
		// ショップ機能はまだ未実装
		// 何かキーを押すとマップに戻る
		if event.IsAnyKey() {
			c.gameInteractor.ReturnToMap()
			c.cursorPosition = 0 // カーソルをリセット
		}

	case 6: // StateEvent
		// イベント機能はまだ未実装
		// 何かキーを押すとマップに戻る
		if event.IsAnyKey() {
			c.gameInteractor.ReturnToMap()
			c.cursorPosition = 0 // カーソルをリセット
		}

	case 7: // StateGameOver
		// 何かキーを押すとゲームを終了するのじゃ
		if event.IsAnyKey() {
			c.gameInteractor.SetDone(true)
		}
	}
}

// カーソル移動を処理する関数じゃ
func (c *GameController) handleCursorMovement(event EventPort) {
	// カーソル操作 - 上下移動
	if event.IsUp() && c.cursorPosition > 0 {
		c.cursorPosition--
	} else if event.IsDown() && c.cursorPosition < c.cursorMaxPosition-1 {
		c.cursorPosition++
	}
}

// draw は画面を描画する関数じゃ
func (c *GameController) draw() {
	c.screen.Clear()

	// 画面サイズを確認するのじゃ
	width, height := c.screen.GetSize()

	// 最小サイズを確認するのじゃ
	if width < 80 || height < 24 {
		c.drawSizeWarning(width, height)
	} else {
		// タイトルを表示するのじゃ
		c.screen.DrawText(1, 1, DefaultStyle(), "Slay the CLI")

		switch c.gameInteractor.State {
		case 1: // StateMap
			c.drawMapScreen(width, height)
		case 2: // StateCombat
			c.drawCombatScreen(width, height)
		case 3: // StateReward
			c.drawRewardScreen(width, height)
		case 4: // StateRest
			c.drawRestScreen(width, height)
		case 5: // StateShop
			c.drawShopScreen(width, height)
		case 6: // StateEvent
			c.drawEventScreen(width, height)
		case 7: // StateGameOver
			c.drawGameOverScreen(width, height)
		}
	}

	c.screen.Show()
}

// 画面サイズが小さすぎる場合の警告を表示するのじゃ
func (c *GameController) drawSizeWarning(width, height int) {
	centerX := width / 2
	centerY := height / 2

	warningText := "画面サイズが小さすぎる！"
	sizeText := fmt.Sprintf("現在: %dx%d 最小: 80x24", width, height)

	c.screen.DrawText(max(0, centerX-len(warningText)/2), centerY-1, DefaultStyle(), warningText)
	c.screen.DrawText(max(0, centerX-len(sizeText)/2), centerY+1, DefaultStyle(), sizeText)
}

// 戦闘画面を描画する関数じゃ
func (c *GameController) drawCombatScreen(width, height int) {
	centerX := width / 2

	// プレイヤー情報を表示するのじゃ（右側に配置）
	healthInfo := fmt.Sprintf("体力: %d/%d", c.gameInteractor.Player.Health, c.gameInteractor.Player.MaxHealth)
	blockInfo := fmt.Sprintf("ブロック: %d", c.gameInteractor.Player.Block)
	goldInfo := fmt.Sprintf("ゴールド: %d", c.gameInteractor.Player.Gold)
	energyInfo := fmt.Sprintf("エナジー: %d/%d", c.gameInteractor.Player.Energy, c.gameInteractor.Player.MaxEnergy)

	// プレイヤー情報を画面右側に表示するのじゃ
	statusX := width - 25
	c.screen.DrawText(statusX, height-6, DefaultStyle(), healthInfo)
	c.screen.DrawText(statusX, height-5, DefaultStyle(), blockInfo)
	c.screen.DrawText(statusX, height-4, DefaultStyle(), goldInfo)
	c.screen.DrawText(statusX, height-3, DefaultStyle(), energyInfo)

	// 敵の情報を表示するのじゃ
	enemyInfo := fmt.Sprintf("%s (%d/%d)", c.gameInteractor.Enemy.Name, c.gameInteractor.Enemy.Health, c.gameInteractor.Enemy.MaxHealth)
	enemyBlockInfo := fmt.Sprintf("ブロック: %d", c.gameInteractor.Enemy.Block)
	enemyIntention := fmt.Sprintf("意図: %s %d", c.gameInteractor.Enemy.Intention, c.gameInteractor.Enemy.Damage)

	c.screen.DrawText(centerX-len(enemyInfo)/2, 3, DefaultStyle(), enemyInfo)
	c.screen.DrawText(centerX-len(enemyBlockInfo)/2, 4, DefaultStyle(), enemyBlockInfo)
	c.screen.DrawText(centerX-len(enemyIntention)/2, 5, DefaultStyle(), enemyIntention)

	// 手札を表示するのじゃ（左側に配置）
	c.screen.DrawText(1, height-12, DefaultStyle(), "手札:")

	// カーソルの最大位置を設定（手札の枚数）
	c.cursorMaxPosition = min(len(c.gameInteractor.Player.Hand), 5) // 最大表示枚数を超えないように

	maxCardsToShow := 5 // 最大表示枚数
	for i, card := range c.gameInteractor.Player.Hand {
		if i >= maxCardsToShow {
			break // 最大表示枚数を超えたら表示しないのじゃ
		}
		
		// カードの情報を作成
		cardInfo := fmt.Sprintf("%s (%dエナジー) - %s", card.Name, card.EnergyCost, card.Description)
		
		// カーソル位置に応じてスタイルを変更（背景色のみで選択表示）
		if i == c.cursorPosition {
			c.screen.DrawText(1, height-11+i, SelectedStyle(), cardInfo)
		} else {
			c.screen.DrawText(1, height-11+i, DefaultStyle(), cardInfo)
		}
	}

	// 手札が多い場合は省略を表示するのじゃ
	if len(c.gameInteractor.Player.Hand) > maxCardsToShow {
		c.screen.DrawText(1, height-11+maxCardsToShow, DefaultStyle(), fmt.Sprintf("...他 %d 枚", len(c.gameInteractor.Player.Hand)-maxCardsToShow))
	}

	// 山札と捨て札の情報を表示するのじゃ
	deckInfo := fmt.Sprintf("山札: %d枚 捨て札: %d枚", len(c.gameInteractor.Player.DrawPile), len(c.gameInteractor.Player.DiscardPile))
	c.screen.DrawText(width-len(deckInfo)-1, height-1, DefaultStyle(), deckInfo)

	// 操作説明を表示するのじゃ
	c.screen.DrawText(1, height-1, DefaultStyle(), "操作: i/,:選択 ;//:決定 e/;:ターン終了 q:終了")
}

// マップ画面を描画する関数じゃ
func (c *GameController) drawMapScreen(width, height int) {
	centerX := width / 2

	// タイトルを表示
	mapTitle := "ダンジョンマップ"
	c.screen.DrawText(centerX-len(mapTitle)/2, 3, DefaultStyle(), mapTitle)

	// 現在のフロア情報を表示
	floorInfo := fmt.Sprintf("現在のフロア: %d", c.gameInteractor.GameMap.CurrentNode.Position.Floor)
	c.screen.DrawText(centerX-len(floorInfo)/2, 5, DefaultStyle(), floorInfo)

	// 接続されたノードを表示
	c.screen.DrawText(centerX-10, 7, DefaultStyle(), "選択可能なノード:")

	// カーソルの最大位置を設定（接続ノードの数）
	c.cursorMaxPosition = len(c.gameInteractor.GameMap.CurrentNode.Connections)

	for i, node := range c.gameInteractor.GameMap.CurrentNode.Connections {
		nodeInfo := fmt.Sprintf("%s (フロア %d)", node.GetNodeTypeString(), node.Position.Floor)
		
		// カーソル位置に応じてスタイルを変更（背景色のみで選択表示）
		if i == c.cursorPosition {
			c.screen.DrawText(centerX-len(nodeInfo)/2, 9+i, SelectedStyle(), nodeInfo)
		} else {
			c.screen.DrawText(centerX-len(nodeInfo)/2, 9+i, DefaultStyle(), nodeInfo)
		}
	}

	// プレイヤー情報を表示
	playerInfo := fmt.Sprintf("体力: %d/%d  ゴールド: %d", c.gameInteractor.Player.Health, c.gameInteractor.Player.MaxHealth, c.gameInteractor.Player.Gold)
	c.screen.DrawText(centerX-len(playerInfo)/2, height-5, DefaultStyle(), playerInfo)

	// 操作説明
	c.screen.DrawText(centerX-20, height-3, DefaultStyle(), "操作: i/,:選択 ;//:決定 q:終了")
}

// 休憩場所画面を描画する関数じゃ
func (c *GameController) drawRestScreen(width, height int) {
	centerX := width / 2

	// タイトルを表示
	restTitle := "休憩所"
	c.screen.DrawText(centerX-len(restTitle)/2, 3, DefaultStyle(), restTitle)

	// 選択肢を表示
	healOption := fmt.Sprintf("回復 (体力の30%%回復)")
	upgradeOption := "カードアップグレード (未実装)"

	// カーソルの最大位置を設定（選択肢の数）
	c.cursorMaxPosition = 2

	// カーソル位置に応じてスタイルを変更（背景色のみで選択表示）
	if c.cursorPosition == 0 {
		c.screen.DrawText(centerX-len(healOption)/2, height/2-1, SelectedStyle(), healOption)
	} else {
		c.screen.DrawText(centerX-len(healOption)/2, height/2-1, DefaultStyle(), healOption)
	}

	if c.cursorPosition == 1 {
		c.screen.DrawText(centerX-len(upgradeOption)/2, height/2+1, SelectedStyle(), upgradeOption)
	} else {
		c.screen.DrawText(centerX-len(upgradeOption)/2, height/2+1, DefaultStyle(), upgradeOption)
	}

	// プレイヤー情報を表示
	playerInfo := fmt.Sprintf("体力: %d/%d", c.gameInteractor.Player.Health, c.gameInteractor.Player.MaxHealth)
	c.screen.DrawText(centerX-len(playerInfo)/2, height-5, DefaultStyle(), playerInfo)

	// 操作説明
	c.screen.DrawText(centerX-20, height-3, DefaultStyle(), "操作: i/,:選択 ;//:決定 q:終了")
}

// ショップ画面を描画する関数じゃ
func (c *GameController) drawShopScreen(width, height int) {
	centerX := width / 2

	// タイトルを表示
	shopTitle := "ショップ"
	c.screen.DrawText(centerX-len(shopTitle)/2, 3, DefaultStyle(), shopTitle)

	// まだ実装されていないことを表示
	notImplementedText := "ショップ機能はまだ実装されていません"
	continueText := "何かキーを押してマップに戻る..."

	c.screen.DrawText(centerX-len(notImplementedText)/2, height/2-1, DefaultStyle(), notImplementedText)
	c.screen.DrawText(centerX-len(continueText)/2, height/2+1, DefaultStyle(), continueText)

	// プレイヤーの所持金を表示
	goldInfo := fmt.Sprintf("所持金: %dゴールド", c.gameInteractor.Player.Gold)
	c.screen.DrawText(centerX-len(goldInfo)/2, height-5, DefaultStyle(), goldInfo)
}

// イベント画面を描画する関数じゃ
func (c *GameController) drawEventScreen(width, height int) {
	centerX := width / 2

	// タイトルを表示
	eventTitle := "イベント"
	c.screen.DrawText(centerX-len(eventTitle)/2, 3, DefaultStyle(), eventTitle)

	// まだ実装されていないことを表示
	notImplementedText := "イベント機能はまだ実装されていません"
	continueText := "何かキーを押してマップに戻る..."

	c.screen.DrawText(centerX-len(notImplementedText)/2, height/2-1, DefaultStyle(), notImplementedText)
	c.screen.DrawText(centerX-len(continueText)/2, height/2+1, DefaultStyle(), continueText)
}

// 報酬画面を描画する関数じゃ
func (c *GameController) drawRewardScreen(width, height int) {
	centerX := width / 2

	// 勝利メッセージを表示するのじゃ
	victoryText := "敵を倒した！"
	goldText := fmt.Sprintf("報酬: %dゴールド", c.gameInteractor.Player.Gold)

	c.screen.DrawText(centerX-len(victoryText)/2, height/4, DefaultStyle(), victoryText)
	c.screen.DrawText(centerX-len(goldText)/2, height/4+1, DefaultStyle(), goldText)

	// カード報酬を表示するのじゃ
	c.screen.DrawText(centerX-5, height/4+3, DefaultStyle(), "カード報酬:")

	// カーソルの最大位置を設定（カード報酬の数）
	c.cursorMaxPosition = len(c.gameInteractor.CardRewards)

	for i, card := range c.gameInteractor.CardRewards {
		cardInfo := fmt.Sprintf("%s (%dエナジー) - %s", card.Name, card.EnergyCost, card.Description)
		
		// カーソル位置に応じてスタイルを変更（背景色のみで選択表示）
		if i == c.cursorPosition {
			c.screen.DrawText(centerX-len(cardInfo)/2, height/4+5+i, SelectedStyle(), cardInfo)
		} else {
			c.screen.DrawText(centerX-len(cardInfo)/2, height/4+5+i, DefaultStyle(), cardInfo)
		}
	}

	// 操作説明
	skipText := "s: スキップ"
	c.screen.DrawText(centerX-len(skipText)/2, height/4+10, DefaultStyle(), skipText)
	
	// 操作説明を追加
	vimText := "i/,:選択 ;//:決定 q:終了"
	c.screen.DrawText(centerX-len(vimText)/2, height/4+12, DefaultStyle(), vimText)
}

// ゲームオーバー画面を描画する関数じゃ
func (c *GameController) drawGameOverScreen(width, height int) {
	centerX := width / 2

	// ゲームオーバーメッセージを表示するのじゃ
	gameOverText := "ゲームオーバー"
	scoreText := fmt.Sprintf("獲得したゴールド: %d", c.gameInteractor.Player.Gold)
	exitText := "何かキーを押して終了..."

	c.screen.DrawText(centerX-len(gameOverText)/2, height/2-1, DefaultStyle(), gameOverText)
	c.screen.DrawText(centerX-len(scoreText)/2, height/2, DefaultStyle(), scoreText)
	c.screen.DrawText(centerX-len(exitText)/2, height/2+2, DefaultStyle(), exitText)
}

// 2つの整数の最大値を返す関数じゃ
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 2つの整数の最小値を返す関数じゃ
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
