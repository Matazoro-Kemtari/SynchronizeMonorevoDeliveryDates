package proposition

// https://qiita.com/0829/items/c1e494bb128ade5f0872

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sclevine/agouti"
	"go.uber.org/zap"
)

type MonorevoUserConfig struct {
	comId    string
	userId   string
	userPass string
}

func NewMonorevoUserConfig() *MonorevoUserConfig {
	return &MonorevoUserConfig{
		comId:    os.Getenv("MONOREVO_COMPANY_ID"),
		userId:   os.Getenv("MONOREVO_USER_ID"),
		userPass: os.Getenv("MONOREVO_USER_PASSWORD"),
	}
}

func TestMonorevoUserConfigCreate(comId, userId, userPass string) *MonorevoUserConfig {
	return &MonorevoUserConfig{
		comId:    comId,
		userId:   userId,
		userPass: userPass,
	}
}

// ものレボ案件一覧Repository
type PropositionTable struct {
	sugar       *zap.SugaredLogger
	comId       string
	userId      string
	userPass    string
	downloadDir string
	workDir     string
}

func NewPropositionTable(
	sugar *zap.SugaredLogger,
	cnf *MonorevoUserConfig,
) *PropositionTable {
	// 実行ディレクトリを取得する cronで実行時のカレントディレクトリ対策
	exeFile, _ := os.Executable()
	exePath := filepath.Dir(exeFile)
	return &PropositionTable{
		sugar:       sugar,
		comId:       cnf.comId,
		userId:      cnf.userId,
		userPass:    cnf.userPass,
		downloadDir: filepath.Join(exePath, "download"),
		workDir:     filepath.Join(exePath, "work"),
	}
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
				"--no-sandbox",
				"--disable-dev-shm-usage", // /dev/shmパーティションの使用を禁止し、パーティションが小さすぎることによる、クラッシュを回避する。 dockerなどのVM環境下では、設定したほうがクラッシュする確率が減る。
				// "window-size=500,400",                  // 画面を小さく
				"--blink-settings=imagesEnabled=false", // 画像を読み込まない
				"lang=ja",
				"--disable-desktop-notifications",
				"--ignore-certificate-errors", // sslまわりのエラーを許容する
				"--disable-extensions",        // Extensionを利用しない
				// "--user-agent=Mozilla/5.0 (iPhone; CPU iPhone OS 10_2 like Mac OS X) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0 Mobile/14C92 Safari/602.1')", // UAの設定。ここではiPhoneに偽装している
			},
		),
	)
	return driver
}

func (p *PropositionTable) loginToMonorevo(driver *agouti.WebDriver) (*agouti.Page, error) {
	page, err := driver.NewPage()
	if err != nil {
		p.sugar.Error("driver.NewPage", err)
		return nil, fmt.Errorf("driver.NewPage error: %v", err)
	}

	// loginページを開く
	p.sugar.Info("loginページを開く")
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

	check := page.FindByXPath(`/html/body/div[1]/div[2]/div[2]/div[1]/div[2]/div[1]/div/div[2]/div/div/div[2]/div[1]/div[1]/div[1]/div[1]/input`)
	for i := 0; ; i++ {
		if err := check.Click(); err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)

		if i >= 60 {
			p.sugar.Error("ログインタイムアウト", i)
			return nil, fmt.Errorf("ログインタイムアウト count: %v", i)
		}
	}
	p.sugar.Infof("ログイン成功: 会社ID(%v) ログインID(%v)", p.comId, p.userId)

	return page, nil
}

func (p *PropositionTable) movePropositionTablePage(page *agouti.Page) error {
	// メニューの案件一覧をクリックする
	const MONOREVO_PROPOSITION_TABLE = "https://app.monorevo.jp/smlot/order/list.html"
	err := page.Navigate(MONOREVO_PROPOSITION_TABLE)

	btn := page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/form/table/tbody/tr[1]/td[1]/input`)
	for i := 0; ; i++ {
		if err := btn.Click(); err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)

		if i >= 60 {
			p.sugar.Error("案件一覧に移動タイムアウト", i)
			return fmt.Errorf("案件一覧に移動タイムアウト count: %v", i)
		}
	}
	p.sugar.Info("案件一覧ページに移動成功")

	return err
}

func (p *PropositionTable) initializeDownloadDir() error {
	// WebDriverは何をダウンロードしたのかわからない
	// フォルダは今回ダウンロードしたもののみになる様にしておく必要がある
	if f, err := os.Stat(p.downloadDir); os.IsNotExist(err) || !f.IsDir() {
		p.sugar.Info("ダウンロードフォルダは存在しないため、削除しません", p.downloadDir)
	} else {
		p.sugar.Info("ダウンロードフォルダの削除を実行", p.downloadDir)
		if err := os.RemoveAll(p.downloadDir); err != nil {
			p.sugar.Error("ダウンロードフォルダの削除に失敗", err)
			return fmt.Errorf("ダウンロードフォルダの削除に失敗 error: %v", err)
		}
	}
	return os.Mkdir(p.downloadDir, 0755)
}
