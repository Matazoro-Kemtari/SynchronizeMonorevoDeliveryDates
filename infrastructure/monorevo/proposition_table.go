package monorevo

// https://qiita.com/0829/items/c1e494bb128ade5f0872

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sclevine/agouti"
	"go.uber.org/zap"
)

// ものレボ案件一覧Repository
type PropositionTable struct {
	Propositions []monorevo.Proposition
	sugar        *zap.SugaredLogger
	comId        string
	userId       string
	userPass     string
	downloadDir  string
	tempDir      string
}

func NewPropositionTable(
	sugar *zap.SugaredLogger,
	comId,
	userId,
	userPass string,
) *PropositionTable {
	// 実行ディレクトリを取得する cronで実行時のカレントディレクトリ対策
	exeFile, _ := os.Executable()
	exePath := filepath.Dir(exeFile)
	return &PropositionTable{
		sugar:       sugar,
		comId:       comId,
		userId:      userId,
		userPass:    userPass,
		downloadDir: filepath.Join(exePath, "download"),
		tempDir:     filepath.Join(exePath, "temp"),
	}
}

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
	p.openCsvFile()

	// TODO: モデルに入れる

	return nil, nil
}

func (p *PropositionTable) PostRange([]monorevo.Proposition) error {
	panic(0)
	return nil
}

func (p *PropositionTable) getWebDriver() *agouti.WebDriver {
	_ = p.initializeDownloadDir()

	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("prefs", map[string]interface{}{
			"download.default_directory":   p.downloadDir,
			"download.prompt_for_download": false,
			"download.directory_upgrade":   true,
		}),
		agouti.ChromeOptions(
			"args", []string{
				// TODO: 開発中はコメントアウト
				// "--headless",
				// "--disable-gpu",
			},
		),
	)
	return driver
}

func (p *PropositionTable) loginToMonorevo(driver *agouti.WebDriver) (*agouti.Page, error) {
	page, err := driver.NewPage()
	if err != nil {
		p.sugar.Fatal("driver.NewPage", err)
	}

	// loginページを開く
	const MONOREVO_LOGIN_URL string = "https://app.monorevo.jp/base/auth/login.html"
	page.Navigate(MONOREVO_LOGIN_URL)

	time.Sleep(time.Second)

	// リダイレクトされるなら認証済み
	if url, _ := page.URL(); url != MONOREVO_LOGIN_URL {
		return page, nil
	}

	// ログインする
	page.FindByXPath(`//*[@id="inputCompany"]`).Fill(p.comId)
	page.FindByXPath(`//*[@id="inputLoginId"]`).Fill(p.userId)
	page.FindByXPath(`//*[@id="inputPassword"]`).Fill(p.userPass)
	page.FindByXPath(`//*[@id="app"]/div/div[3]/form/div/div[2]/div[5]/button`).Click()

	i := 0
	for i < 60 {
		check := page.FindByXPath(`/html/body/div[1]/div[2]/div[2]/div[1]/div[2]/div[1]/div/div[2]/div/div/div[2]/div[1]/div[1]/div[1]/div[1]/input`)
		if err := check.Click(); err == nil {
			break
		}
		time.Sleep(time.Second)
		i++
	}
	if i >= 60 {
		p.sugar.Fatal("ログインタイムアウト", i)
	}

	return page, nil
}

func (p *PropositionTable) movePropositionTablePage(page *agouti.Page) error {
	// メニューの案件一覧をクリックする
	const MONOREVO_PROPOSITION_TABLE = "https://app.monorevo.jp/smlot/order/list.html"
	err := page.Navigate(MONOREVO_PROPOSITION_TABLE)

	i := 0
	for i < 60 {
		btn := page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/form/table/tbody/tr[1]/td[1]/input`)
		if err := btn.Click(); err == nil {
			break
		}
		time.Sleep(time.Second)
		i++
	}
	return err
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

func (p *PropositionTable) initializeDownloadDir() error {
	if f, err := os.Stat(p.downloadDir); os.IsNotExist(err) || !f.IsDir() {
		p.sugar.Info("ダウンロードフォルダは存在しません", p.downloadDir)
	} else {
		p.sugar.Info("ダウンロードフォルダの削除を実行", p.downloadDir)
		if err := os.RemoveAll(p.downloadDir); err != nil {
			p.sugar.Fatal("ダウンロードフォルダの削除に失敗", err)
		}
	}
	return os.Mkdir(p.downloadDir, 0755)
}

func (p *PropositionTable) openCsvFile() error {
	// テンポラリフォルダの作成
	_ = p.initializeTempDir()

	// ファイル移動
	_ = p.moveDownloadToTemp()
}

func (p *PropositionTable) initializeTempDir() error {
	// WebDriverは何をダウンロードしたのかわからない
	// フォルダは今回ダウンロードしたもののみになる様にしておく必要がある
	if f, err := os.Stat(p.tempDir); os.IsNotExist(err) || !f.IsDir() {
		p.sugar.Info("テンポラリフォルダは存在しません", p.tempDir)
	} else {
		p.sugar.Info("テンポラリフォルダの削除を実行", p.tempDir)
		if err := os.RemoveAll(p.tempDir); err != nil {
			p.sugar.Fatal("テンポラリフォルダの削除に失敗", err)
		}
	}
	return os.Mkdir(p.tempDir, 0755)
}

func (p *PropositionTable) moveDownloadToTemp() error {
	files, err := ioutil.ReadDir(p.downloadDir)
	if err != nil {
		p.sugar.Fatal("ダウンロードフォルダのファイル一覧の取得失敗", err)
	}

	if len(files) == 0 {
		p.sugar.Fatal("ダウンロードフォルダにファイルが存在しない")
	}

	f := files[0]

	return os.Rename(
		filepath.Join(p.downloadDir, f.Name()),
		filepath.Join(p.tempDir, f.Name()),
	)
}
