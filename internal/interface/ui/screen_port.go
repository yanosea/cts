package ui

// ScreenPort は画面表示のインターフェースを定義するのじゃ
type ScreenPort interface {
	Clear()
	Show()
	DrawText(x, y int, style Style, text string)
	GetSize() (width int, height int)
	PollEvent() EventPort
	Sleep(ms int)
	Cleanup()
}

// EventPort はイベントのインターフェースを定義するのじゃ
type EventPort interface {
	IsExit() bool
	IsEndTurn() bool
	IsAnyKey() bool
	GetCardIndex() int
	IsResize() bool
	IsSKey() bool
	IsKey1() bool
	IsKey2() bool
	// Vim風の操作に必要なイベント判定を追加するのじゃ
	IsUp() bool
	IsDown() bool
	IsLeft() bool
	IsRight() bool
	IsEnter() bool
	IsSpace() bool
}

// Style は表示スタイルを表すのじゃ
type Style interface {
	// スタイル関連のメソッドが必要になれば追加するのじゃ
}

// スタイル定義用の定数を追加するのじゃ
type StyleType int

const (
	DefaultStyleType StyleType = iota
	SelectedStyleType
)

// StyleTypeContainer はスタイルの種類を保持する構造体じゃ
type StyleTypeContainer struct {
	Type StyleType
}

// DefaultStyle はデフォルトスタイルを返すのじゃ
func DefaultStyle() Style {
	return &StyleTypeContainer{Type: DefaultStyleType}
}

// SelectedStyle は選択中のスタイルを返すのじゃ
func SelectedStyle() Style {
	return &StyleTypeContainer{Type: SelectedStyleType}
}
