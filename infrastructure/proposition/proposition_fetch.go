package proposition

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sclevine/agouti"
)

// ものレボから案件一覧を取得する
func (p *PropositionTable) FetchAll() ([]monorevo.Proposition, error) {
	// webdriverを初期化する
	driver := p.getWebDriver()
	defer driver.Stop()
	driver.Start()

	// ログインする
	page, err := p.loginToMonorevo(driver)
	if err != nil {
		p.sugar.Fatal("ものレボにログインできなかった", err)
	}

	// 案件一覧一覧画面に移動する
	if err := p.movePropositionTablePage(page); err != nil {
		p.sugar.Fatal("案件一覧一覧画面に移動できなかった", err)
	}

	// ダウンロードする
	p.downloadPropositionTable(page)

	// csvを開く
	csv, err := p.openCsvFile()
	if err != nil {
		p.sugar.Fatal("csvファイルを開く処理で失敗しました", err)
	}

	return csv, nil
}

func (p *PropositionTable) downloadPropositionTable(page *agouti.Page) {
	// ダウンロードボタンを押す
	page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/div[1]/div[2]/div/div/button`).Click()
	page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/div[1]/div[2]/div/div/div/div[1]`).Click()

	// データ準備まで待つ
	i := 0
	for i < 60 {
		selector := page.FindByXPath(`/html/body/div[3]/div[2]/div`)
		if v, _ := selector.Visible(); v {
			break
		}
		time.Sleep(time.Second)
		i++
	}

	if i >= 60 {
		p.sugar.Fatal("ダウンロードタイムアウト", i)
	}

	// 実行ボタン押下
	page.FindByXPath(`/html/body/div[3]/div[2]/div/div[3]/button[2]`).Click()
	time.Sleep(time.Second)
}

func (p *PropositionTable) openCsvFile() ([]monorevo.Proposition, error) {
	// テンポラリフォルダの作成
	if err := p.initializeWorkDir(); err != nil {
		p.sugar.Fatal("テンポラリフォルダの作成で失敗しました", err)
	}

	// ファイル移動
	n, err := p.moveDownloadToWork()
	if err != nil {
		p.sugar.Fatal("ファイル移動で失敗しました", err)
	}

	// csvをパースする
	csv, err := p.deserializeCsv(n)
	if err != nil {
		p.sugar.Fatal("csvのパースに失敗しました", err)
	}
	return csv, nil
}

func (p *PropositionTable) initializeWorkDir() error {
	if f, err := os.Stat(p.workDir); os.IsNotExist(err) || !f.IsDir() {
		p.sugar.Info("テンポラリフォルダは存在しません", p.workDir)
	} else {
		p.sugar.Info("テンポラリフォルダの削除を実行", p.workDir)
		if err := os.RemoveAll(p.workDir); err != nil {
			p.sugar.Fatal("テンポラリフォルダの削除に失敗", err)
		}
	}
	return os.Mkdir(p.workDir, 0755)
}

func (p *PropositionTable) moveDownloadToWork() (string, error) {
	files, err := ioutil.ReadDir(p.downloadDir)
	if err != nil {
		p.sugar.Fatal("ダウンロードフォルダのファイル一覧の取得失敗", err)
	}

	if len(files) == 0 {
		p.sugar.Fatal("ダウンロードフォルダにファイルが存在しない")
	}

	f := files[0]

	return f.Name(),
		os.Rename(
			filepath.Join(p.downloadDir, f.Name()),
			filepath.Join(p.workDir, f.Name()),
		)
}

func (p *PropositionTable) deserializeCsv(name string) ([]monorevo.Proposition, error) {
	file, err := os.Open(filepath.Join(p.workDir, name))
	if err != nil {
		p.sugar.Fatal("ファイルが開けませんでした", err)
	}
	defer file.Close()

	// csv.NewReaderを使ってcsvを読み込む
	r := csv.NewReader(file)
	var fropositions []monorevo.Proposition
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		// 納期を日付形式に変換
		d, err := time.Parse("2006/01/02", row[8])
		if err != nil {
			p.sugar.Error("csvファイルの納期の日付形式変換でエラー発生", row[8], err)
			continue
		}

		// ものレボ案件配列に蓄積
		fropositions = append(
			fropositions,
			*monorevo.NewProposition(
				row[0],
				row[1],
				d,
			),
		)
	}
	return fropositions, nil
}