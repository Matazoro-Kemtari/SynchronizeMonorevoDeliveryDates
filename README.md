# SynchronizeMono-revoDeliveryDates
ものレボの納期を同期させる

## 技術資料
<img src="images/sudoDiagram.png" alt="sudo図"/>

## 開発環境インストール
vscodeは入っている前提で説明する。

こちらを参考に、拡張機能をインストールして、設定を行う。
[VSCodeでGo言語の開発環境を構築する](https://qiita.com/melty_go/items/c977ba594efcffc8b567)

さらに読み進めてプロジェクトのセットアップを行う。
もしくは、gitからダウンロードする。

launch.jsonを記述する。
```json
{
    // IntelliSense を使用して利用可能な属性を学べます。
    // 既存の属性の説明をホバーして表示します。
    // 詳細情報は次を確認してください: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Workspace",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}"
        },
        {
            "name": "Launch File",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${fileDirname}"
        }

    ]
}
```

### 必要ライブラリ

DIツール
[Goでwireを使って依存性注入（DI）する](https://rinoguchi.net/2022/06/go_wire_id.html)
[GoのDIツールwireで知っておくと良いこと](https://christina04.hatenablog.com/entry/google-wire)
[GoのプロジェクトのDIをWireを使ってシンプルに](https://qiita.com/momotaro98/items/0b75a37048833dd6d324)
```
$ go install github.com/google/wire/cmd/wire@latest
```
モック
[Goでメソッドを簡単にモック化する【gomock】](https://qiita.com/gold-kou/items/81562f9142323b364a60)
```
$ go get github.com/golang/mock/gomock
$ go install github.com/golang/mock/mockgen
```
環境変数
[【Go】.envファイルをGolangでも使用するためのライブラリ「godotenv」](https://qiita.com/sola-msr/items/fb7d6889d7bd7a6705d0)
```
$ go get -u github.com/joho/godotenv
```
ロギング
[golangの高速な構造化ログライブラリ「zap」の使い方](https://qiita.com/emonuh/items/28dbee9bf2fe51d28153)
```
$ go get -u go.uber.org/zap
```
