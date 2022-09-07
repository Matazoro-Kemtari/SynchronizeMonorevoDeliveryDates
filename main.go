package main

import (
	"SynchronizeMonorevoDeliveryDates/infrastructure/appsetting"
	"SynchronizeMonorevoDeliveryDates/infrastructure/reportsetting"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// バージョン埋め込む
// INFO: https://qiita.com/irotoris/items/4aae9ad483bf08915688
var version string
var revision string

func main() {
	// コマンドライン引数を取得
	isVersion := flag.Bool("version", false, "バージョンを表示する")
	flag.Parse()

	// バージョン表示
	if *isVersion {
		fmt.Printf("version: %s-%s\n", version, revision)
		os.Exit(0)
	}

	// 実行ディレクトリを取得する cronで実行時のカレントディレクトリ対策
	exeFile, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(exeFile)
	if err := os.Chdir(exePath); err != nil {
		panic(err)
	}

	// ログファイルの設定
	// logFile := filepath.Join(exePath, "app_log.json")
	logFile := "app_log.json"

	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)
	// https://qiita.com/emonuh/items/28dbee9bf2fe51d28153#config%E7%B7%A8
	enc := zapcore.NewJSONEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			MessageKey:     "Message",
			StacktraceKey:  "Stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	)
	sink := zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   "./" + logFile, // ファイル名
			MaxSize:    1,              // ローテーションするファイルサイズ(megabytes)
			MaxAge:     7,              // 古いログを保持する日数
			MaxBackups: 3,              // 保持する古いログの最大ファイル数
			LocalTime:  false,          // バックアップファイルの時刻フォーマットをサーバローカル時間指定
			Compress:   false,          // ローテーションされたファイルのgzip圧縮
		},
	)
	logger := zap.New(zapcore.NewCore(enc, sink, level))
	defer logger.Sync()
	sugar := logger.Sugar()

	// .envファイルから環境変数を読み込む
	err_read := godotenv.Load()
	if err_read != nil {
		sugar.Fatal(err_read)
	}

	appConfigFile := "appsettings.json"
	ap, err := appsetting.NewLoadableSetting(sugar).Load(appConfigFile)
	if err != nil {
		sugar.Fatalf("%vの読み込みに失敗しました error: %v", appConfigFile, err)
	}
	reportConfigFile := "reportsettings.json"
	rp, err := reportsetting.NewLoadableSetting(sugar).Load(reportConfigFile)
	if err != nil {
		sugar.Fatalf("%vの読み込みに失敗しました error: %v", reportConfigFile, err)
	}
	synchronize := InitializeSynchronize(sugar, ap, rp)
	if err := synchronize.Synchronize(); err != nil {
		sugar.Fatal(err)
	}
}
