package tcell_screen

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/yanosea/cts/internal/interface/ui"
)

// ScreenAdapter はtcellライブラリを使用して画面表示を実装するのじゃ
type ScreenAdapter struct {
	screen tcell.Screen
}

// styleContainer はスタイル情報を持つ構造体じゃ
type styleContainer struct {
	style tcell.Style
}

// NewScreenAdapter はScreenAdapterのインスタンスを生成するのじゃ
func NewScreenAdapter() (*ScreenAdapter, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("スクリーンの初期化に失敗じゃ: %v", err)
	}

	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("スクリーンの初期化に失敗じゃ: %v", err)
	}

	// 基本的なスタイルを設定するのじゃ
	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	screen.Clear()

	return &ScreenAdapter{screen: screen}, nil
}

// Clear は画面をクリアするのじゃ
func (a *ScreenAdapter) Clear() {
	a.screen.Clear()
}

// Show は画面を更新するのじゃ
func (a *ScreenAdapter) Show() {
	a.screen.Show()
}

// DrawText はテキストを描画するのじゃ
func (a *ScreenAdapter) DrawText(x, y int, style ui.Style, text string) {
	// スタイル設定
	var tcellStyle tcell.Style

	if style == nil {
		tcellStyle = tcell.StyleDefault
	} else if styleContainer, ok := style.(*ui.StyleTypeContainer); ok && styleContainer.Type == ui.SelectedStyleType {
		// 選択中のスタイルはグレー背景、黒文字にする
		tcellStyle = tcell.StyleDefault.Background(tcell.ColorLightGray).Foreground(tcell.ColorBlack)
	} else {
		// その他のスタイルはデフォルトのまま
		tcellStyle = tcell.StyleDefault
	}

	// runewidthを使って正確なテキスト幅を計算するのじゃ
	width := runewidth.StringWidth(text)

	// 選択スタイルの場合は、必ず先に背景を描画するのじゃ
	if styleContainer, ok := style.(*ui.StyleTypeContainer); ok && styleContainer.Type == ui.SelectedStyleType {
		// テキストの全幅に対して背景色を描画するのじゃ
		for i := 0; i < width; i++ {
			// スペースを背景色付きで描画することで、背景を埋めるのじゃ
			a.screen.SetContent(x+i, y, ' ', nil, tcellStyle)
		}
	}

	// 文字を描画するためのx位置を追跡するのじゃ
	pos := 0

	// その上に文字を描画するのじゃ
	for _, r := range text {
		// 各文字の表示幅を取得するのじゃ
		w := runewidth.RuneWidth(r)

		// 文字を描画するのじゃ
		a.screen.SetContent(x+pos, y, r, nil, tcellStyle)

		// 次の文字の位置を更新するのじゃ
		pos += w
	}
}

// GetSize は画面サイズを取得するのじゃ
func (a *ScreenAdapter) GetSize() (width int, height int) {
	return a.screen.Size()
}

// PollEvent はイベントを取得するのじゃ
func (a *ScreenAdapter) PollEvent() ui.EventPort {
	return &EventAdapter{event: a.screen.PollEvent()}
}

// Sleep は指定されたミリ秒だけスリープするのじゃ
func (a *ScreenAdapter) Sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// Cleanup は画面を終了するのじゃ
func (a *ScreenAdapter) Cleanup() {
	a.screen.Fini()
}
