package proposition

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"encoding/csv"
	"fmt"
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
		p.sugar.Error("ものレボにログインできなかった", err)
		return nil, fmt.Errorf("ものレボにログインできなかった error: %v", err)
	}

	// 案件一覧一覧画面に移動する
	if err := p.movePropositionTablePage(page); err != nil {
		p.sugar.Error("案件一覧一覧画面に移動できなかった", err)
		return nil, fmt.Errorf("案件一覧一覧画面に移動できなかった error: %v", err)
	}

	// ダウンロードする
	if err := p.downloadPropositionTable(page); err != nil {
		p.sugar.Error("案件一覧をダウンロードできなかった", err)
		return nil, fmt.Errorf("案件一覧をダウンロードできなかった error: %v", err)
	}

	// テンポラリフォルダの作成
	if err := p.initializeWorkDir(); err != nil {
		p.sugar.Error("作業フォルダの作成で失敗しました", err)
		return nil, fmt.Errorf("作業フォルダの作成で失敗しました error: %v", err)
	}

	// ファイル移動
	f, err := p.moveDownloadToWork()
	if err != nil {
		p.sugar.Error("ファイル移動で失敗しました", err)
		return nil, fmt.Errorf("ファイル移動で失敗しました error: %v", err)
	}

	// csvを開く
	csv, err := p.openCsvFile(f)
	if err != nil {
		p.sugar.Error("csvファイルを開く処理で失敗しました", err)
		return nil, fmt.Errorf("csvファイルを開く処理で失敗しました error: %v", err)
	}

	return csv, nil
}

func (p *PropositionTable) downloadPropositionTable(page *agouti.Page) error {
	// ダウンロードボタンを押す
	page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/div[1]/div[2]/div/div/button`).Click()
	page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/div[1]/div[2]/div/div/div/div[1]`).Click()

	// データ準備まで待つ
	// csvダウンロードポップアップ
	popup := page.FindByXPath(`/html/body/div[3]/div[2]/div`)
	for i := 0; i < 600; i++ {
		if v, _ := popup.Visible(); v {
			break
		}
		time.Sleep(time.Millisecond * 100)

		if i >= 600 {
			p.sugar.Error("ダウンロードタイムアウト", i)
			return fmt.Errorf("ダウンロードタイムアウト count: %v", i)
		}
	}
	time.Sleep(time.Second)

	// 実行ボタン押下
	for i := 0; ; i++ {
		spn := page.FindByXPath(`/html/body/div[3]/div[2]/div/div[3]/button[2]/span[2]/span/span`)
		btn := page.FindByXPath(`/html/body/div[3]/div[2]/div/div[3]/button[2]`)
		var err error
		p.sugar.Debugf("ダウンロード実行ボタン押下 counter: %d", i)
		if err = spn.Click(); err != nil {
			err = nil
			// spanをクリックしても反応しないことがあるから保険的
			err = btn.Check()
		}

		if i > 3 {
			p.sugar.Error("ダウンロード実行ボタン押下で失敗しました", err)
			return fmt.Errorf("ダウンロード実行ボタン押下で失敗しました error: %v", err)
		} else {
			time.Sleep(time.Millisecond * 100)
		}

		if err != nil {
			continue
		} else {
			break
		}
	}
	return nil
}

func (p *PropositionTable) openCsvFile(f string) ([]monorevo.Proposition, error) {
	// csvをパースする
	csv, err := p.deserializeCsv(f)
	if err != nil {
		p.sugar.Error("csvのパースに失敗しました", err)
		return nil, fmt.Errorf("csvのパースに失敗しました error: %v", err)
	}
	return csv, nil
}

func (p *PropositionTable) initializeWorkDir() error {
	if f, err := os.Stat(p.workDir); os.IsNotExist(err) || !f.IsDir() {
		p.sugar.Info("作業フォルダは存在しないため、削除しません", p.workDir)
	} else {
		p.sugar.Info("作業フォルダの削除を実行", p.workDir)
		if err := os.RemoveAll(p.workDir); err != nil {
			p.sugar.Error("作業フォルダの削除に失敗しました", err)
			return fmt.Errorf("作業フォルダの削除に失敗しました error: %v", err)
		}
	}
	return os.Mkdir(p.workDir, 0755)
}

func (p *PropositionTable) moveDownloadToWork() (string, error) {
	files, err := ioutil.ReadDir(p.downloadDir)
	if err != nil {
		p.sugar.Error("ダウンロードフォルダのファイル一覧の取得失敗", err)
		return "", fmt.Errorf("ダウンロードフォルダの削除に失敗しました error: %v", err)
	}

	if len(files) == 0 {
		p.sugar.Error("ダウンロードフォルダにファイルが存在しない", p.downloadDir)
		return "", fmt.Errorf("ダウンロードフォルダ(%v)にファイルが存在しない error: %v", p.downloadDir, err)
	}

	// 始めの1つをダウンロードしたファイルと推定
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
		p.sugar.Error("csvファイルが開けませんでした", err)
		return nil, fmt.Errorf("csvファイルが開けませんでした error: %v", err)
	}
	defer file.Close()

	// csv.NewReaderを使ってcsvを読み込む
	r := csv.NewReader(file)
	var fropositions []monorevo.Proposition
	for i := 1; ; i++ {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if i == 1 {
			// 1行目(ヘッダ)は無視する
			continue
		}

		// 納期を日付形式に変換
		d, err := time.Parse("2006/01/02", row[8])
		if err != nil {
			// 日付変換に失敗したレコードはパスする
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
