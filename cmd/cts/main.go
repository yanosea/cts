package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/yanosea/cts/internal/infrastructure/tcell_screen"
	"github.com/yanosea/cts/internal/interface/ui"
	"github.com/yanosea/cts/internal/usecase"
)

func main() {
	// 乱数の種を設定するのじゃ
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// スクリーンアダプタを初期化するのじゃ
	screenAdapter, err := tcell_screen.NewScreenAdapter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	defer screenAdapter.Cleanup()

	// ゲームのインタラクタを初期化するのじゃ
	gameInteractor := usecase.NewGameInteractor()

	// ゲームコントローラを初期化するのじゃ
	gameController := ui.NewGameController(screenAdapter, gameInteractor)

	// ゲームを開始するのじゃ
	gameController.StartGame()
}
