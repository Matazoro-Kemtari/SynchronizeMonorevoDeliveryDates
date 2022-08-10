package proposition

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

func (p *PropositionTable) PostRange(postableProposition []monorevo.Proposition) error {
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

	// 案件検索をする
	for _, v := range postableProposition {
		if r, err := p.searchPropositionTable(page, v.WorkedNumber); err != nil {
			p.sugar.Fatal("案件検索ができなかった", err)
		} else if !r {
			p.sugar.Infof("作業No(%v)の該当がなかった", v.WorkedNumber)
			continue
		}

		// 納期を更新する
		p.updatedDeliveryDate(page, v.WorkedNumber, v.DeliveryDate)
	}

	return nil
}

type hasRecord bool

func (p *PropositionTable) searchPropositionTable(page *agouti.Page, workNum string) (hasRecord, error) {
	// 検索条件を開く
	openBtn := page.FindByXPath(`//*[@id="accordionDrawing-down"]`)
	openBtn.Click()

	// 検索条件の作業Noを入力する
	workNoFld := page.FindByXPath(`//*[@id="searchContent"]/div[2]/div[1]/input`)
	if err := workNoFld.Fill(workNum); err != nil {
		p.sugar.Fatal("作業Noの入力に失敗しました", err)
	}
	// time.Sleep(time.Millisecond * 500)
	searchBtn := page.FindByXPath(`//*[@id="searchButton"]/div/button`)
	searchBtn.Click()

	// データ準備まで待つ
	i := 0
	for i < 60 {
		// くるくる回るエフェクトのxpath
		selector := page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div[2]`)
		// 処理中の子要素(DIV)が存在する間はクリックしてもエラーにならない
		if err := selector.Click(); err != nil {
			break
		}
		time.Sleep(time.Second)
		i++
	}

	if i >= 60 {
		p.sugar.Fatal("ダウンロードタイムアウト", i)
	}

	// 該当あるか確認
	td := page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/form/table/tbody/tr/td`)
	if _, err := td.Elements(); err == nil {
		// エラーなしは該当なし
		return false, nil
	}
	return true, nil
}

func (p *PropositionTable) updatedDeliveryDate(page *agouti.Page, w string, d time.Time) error {
	curContentsDom, err := page.HTML()
	if err != nil {
		p.sugar.Fatal("DOMの取得に失敗しました", err)
	}

	readerCurContents := strings.NewReader(curContentsDom)
	// htmlをパースする
	contentsDom, _ := goquery.NewDocumentFromReader(readerCurContents)

	// tbodyを取得して td要素数を取得する
	// TODO: 変数名をどうにかして
	hoge := contentsDom.Find(`#app > div > div.contents-wrapper > div.main-wrapper > div > div > div > form > table > tbody`)
	hoge2 := hoge.Children()
	notes := hoge2.Nodes
	p.sugar.Debug(len(notes))

	// TODO: 作業Noの比較して検索失敗に備える
	for i := 1; i <= len(notes); i += 2 {
		wcell := contentsDom.Find(fmt.Sprintf("#app > div > div.contents-wrapper > div.main-wrapper > div > div > div > form > table > tbody > tr:nth-child(%d) > td:nth-child(2)", i)).Text()
		p.sugar.Debug(wcell)

	}
	// TODO: 更新できた作業No できなかった作業Noが分かるように戻り値考える

	tbody := page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/form/table/tbody`)
	// checkbox := tbody.AllByName("orderList")
	// チェックボックスの数で1ページ当たりの行数を取得する
	anchor := tbody.AllByLink("詳細")

	rowCount, err := anchor.Count()
	if err != nil {
		p.sugar.Fatal("チェックボックス(orderList)のカウントに失敗しました", err)
	}
	p.sugar.Debug(rowCount)
	// wsel := tbody.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/form/table/tbody/tr[1]/td[2]`)
	// hoge, _ := wsel.FindForAppium(`//*[@id="app"]/div/div[2]/div[2]/div/div/div/form/table/tbody/tr[1]/td[2]`)
	// p.sugar.Debug(hoge)
	// ek, oh := checkbox.Elements()
	// p.sugar.Debug(oh)
	// for k, v := range ek {
	// 	p.sugar.Debug(k, v)
	// }
	// e, err := tbody.Elements()
	// p.sugar.Debug(e, err)
	// for _, v := range e {
	// 	p.sugar.Debug(v)
	// }
	return nil
}
