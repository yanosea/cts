package tcell_screen

import (
	"github.com/gdamore/tcell/v2"
	// "github.com/yanosea/cts/internal/interface/ui"
)

// EventAdapter はtcellイベントをEventPortインターフェースに変換するのじゃ
type EventAdapter struct {
	event tcell.Event
}

// IsExit はイベントが終了要求かどうかを判定するのじゃ
func (e *EventAdapter) IsExit() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// ESCキー、Ctrl+C、または'q'キーで終了
		return ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC ||
			(ev.Key() == tcell.KeyRune && ev.Rune() == 'q')
	}
	return false
}

// IsEndTurn はイベントがターン終了かどうかを判定するのじゃ
func (e *EventAdapter) IsEndTurn() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// 'e'キーまたはEnterキーでターン終了
		return (ev.Key() == tcell.KeyRune && ev.Rune() == 'e') || ev.Key() == tcell.KeyEnter || 
			(ev.Key() == tcell.KeyRune && ev.Rune() == ';') // セミコロンを追加
	}
	return false
}

// IsAnyKey はイベントが何かしらのキー入力かどうかを判定するのじゃ
func (e *EventAdapter) IsAnyKey() bool {
	_, ok := e.event.(*tcell.EventKey)
	return ok
}

// IsResize はイベントがリサイズイベントかどうかを判定するのじゃ
func (e *EventAdapter) IsResize() bool {
	_, ok := e.event.(*tcell.EventResize)
	return ok
}

// GetCardIndex はイベントからカードインデックスを取得するのじゃ
// 既存の数字キー入力もサポートしておくのじゃ
func (e *EventAdapter) GetCardIndex() int {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyRune {
			switch ev.Rune() {
			case '1':
				return 0
			case '2':
				return 1
			case '3':
				return 2
			case '4':
				return 3
			case '5':
				return 4
			}
		}
	}
	return -1
}

// IsSKey はSキーが押されたかを判定するのじゃ
func (e *EventAdapter) IsSKey() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		return ev.Key() == tcell.KeyRune && (ev.Rune() == 's' || ev.Rune() == 'S')
	}
	return false
}

// IsKey1 は1キーが押されたかを判定するのじゃ
func (e *EventAdapter) IsKey1() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		return ev.Key() == tcell.KeyRune && ev.Rune() == '1'
	}
	return false
}

// IsKey2 は2キーが押されたかを判定するのじゃ
func (e *EventAdapter) IsKey2() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		return ev.Key() == tcell.KeyRune && ev.Rune() == '2'
	}
	return false
}

// Vim風の操作に必要なイベント判定を追加するのじゃ

// IsUp は上方向キーが押されたかを判定するのじゃ
func (e *EventAdapter) IsUp() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// 上矢印キー、'k'キー、または'i'キーで上
		return ev.Key() == tcell.KeyUp || 
			(ev.Key() == tcell.KeyRune && (ev.Rune() == 'k' || ev.Rune() == 'K' || ev.Rune() == 'i' || ev.Rune() == 'I'))
	}
	return false
}

// IsDown は下方向キーが押されたかを判定するのじゃ
func (e *EventAdapter) IsDown() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// 下矢印キー、'j'キー、または','(コンマ)キーで下
		return ev.Key() == tcell.KeyDown || 
			(ev.Key() == tcell.KeyRune && (ev.Rune() == 'j' || ev.Rune() == 'J' || ev.Rune() == ',' || ev.Rune() == '.'))
	}
	return false
}

// IsLeft は左方向キーが押されたかを判定するのじゃ
func (e *EventAdapter) IsLeft() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// 左矢印キー、'h'キー、または'j'キーで左
		return ev.Key() == tcell.KeyLeft || 
			(ev.Key() == tcell.KeyRune && (ev.Rune() == 'h' || ev.Rune() == 'H' || ev.Rune() == 'j' || ev.Rune() == 'J'))
	}
	return false
}

// IsRight は右方向キーが押されたかを判定するのじゃ
func (e *EventAdapter) IsRight() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// 右矢印キー、'l'キー、または'l'キーで右
		return ev.Key() == tcell.KeyRight || 
			(ev.Key() == tcell.KeyRune && (ev.Rune() == 'l' || ev.Rune() == 'L'))
	}
	return false
}

// IsEnter はEnterキーが押されたかを判定するのじゃ
func (e *EventAdapter) IsEnter() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// Enterキーまたは';'(セミコロン)キーで決定
		return ev.Key() == tcell.KeyEnter || (ev.Key() == tcell.KeyRune && ev.Rune() == ';')
	}
	return false
}

// IsSpace はスペースキーが押されたかを判定するのじゃ
func (e *EventAdapter) IsSpace() bool {
	switch ev := e.event.(type) {
	case *tcell.EventKey:
		// スペースキーまたは'/'(スラッシュ)キーで決定
		return (ev.Key() == tcell.KeyRune && (ev.Rune() == ' ' || ev.Rune() == '/'))
	}
	return false
}
