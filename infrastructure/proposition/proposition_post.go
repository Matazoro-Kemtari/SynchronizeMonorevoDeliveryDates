package proposition

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
	"golang.org/x/net/html"
)

func (p *PropositionTable) PostRange(postablePropositions []monorevo.DifferentProposition) ([]monorevo.UpdatedProposition, error) {
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

	var editedPropositions []monorevo.UpdatedProposition
	for _, v := range postablePropositions {
		d := time.Date(
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			0, 0, 0, 0, time.UTC)
		if v.UpdatedDeliveryDate.Before(d) {
			// 現在日より過去日は処理しない ものレボが受け付けない
			editedPropositions = append(
				editedPropositions,
				*monorevo.NewUpdatedProposition(
					v.WorkedNumber,
					v.DET,
					false,
					v.DeliveryDate,
					v.UpdatedDeliveryDate,
				))
			p.sugar.Errorf(
				"現在日(%v)より過去の納期(%v)は受付できない",
				d,
				v.UpdatedDeliveryDate)
			continue
		}

		// 案件検索をする
		if r, err := p.searchPropositionTable(page, v); err != nil {
			editedPropositions = append(
				editedPropositions,
				*monorevo.NewUpdatedProposition(
					v.WorkedNumber,
					v.DET,
					false,
					v.DeliveryDate,
					v.UpdatedDeliveryDate,
				))
			p.sugar.Error(
				"案件検索ができなかった",
				v.WorkedNumber,
				v.DET,
				err)
			continue
		} else if !r {
			editedPropositions = append(
				editedPropositions,
				*monorevo.NewUpdatedProposition(
					v.WorkedNumber,
					v.DET,
					false,
					v.DeliveryDate,
					v.UpdatedDeliveryDate,
				))
			p.sugar.Errorf(
				"作業No(%v),DET番号(%v)の該当がなかった",
				v.WorkedNumber,
				v.DET)
			continue
		}

		// 納期を更新する
		successful, err := p.updatedDeliveryDate(page, v)
		if successful == unspecified && err != nil {
			editedPropositions = append(
				editedPropositions,
				*monorevo.NewUpdatedProposition(
					v.WorkedNumber,
					v.DET,
					false,
					v.DeliveryDate,
					v.UpdatedDeliveryDate,
				))
			p.sugar.Error(
				"納期の更新ができませんでした",
				v.DeliveryDate,
				v.DET,
				err)
			continue
		}
		editedPropositions = append(
			editedPropositions,
			*monorevo.NewUpdatedProposition(
				v.WorkedNumber,
				v.DET,
				(successful == success),
				v.DeliveryDate,
				v.UpdatedDeliveryDate,
			))
	}

	return editedPropositions, nil
}

type hasRecord bool

func (p *PropositionTable) searchPropositionTable(page *agouti.Page, proposition monorevo.DifferentProposition) (hasRecord, error) {
	// 検索条件を開く
	openBtn := page.FindByXPath(`//*[@id="accordionDrawing-down"]`)
	openBtn.Click()

	// **検索条件**
	// 作業Noを入力する
	workNoFld := page.FindByXPath(`//*[@id="searchContent"]/div[2]/div[1]/input`)
	if err := workNoFld.Fill(proposition.WorkedNumber); err != nil {
		p.sugar.Debug("作業Noの入力に失敗しました", err)
		return false, fmt.Errorf("作業Noの入力に失敗しました error: %v", err)
	}
	// DET番号を入力する
	detFld := page.FindByXPath(`//*[@id="searchContent"]/div[2]/div[2]/input`)
	if err := detFld.Fill(proposition.DET); err != nil {
		p.sugar.Debug("DET番号の入力に失敗した", err)
		return false, fmt.Errorf("DET番号の入力に失敗した error: %v", err)
	}
	p.sugar.Infof("案件検索: 作業No(%v) DET番号(%v)", proposition.WorkedNumber, proposition.DET)
	time.Sleep(time.Millisecond * 100)
	searchBtn := page.FindByXPath(`//*[@id="searchButton"]/div/button`)
	searchBtn.Click()

	// データ準備まで待つ
	selector := page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div[2]`)
	for i := 0; ; i++ {
		// くるくる回るエフェクトのxpath
		// 処理中の子要素(DIV)が存在する間はクリックしてもエラーにならない
		if err := selector.Click(); err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)

		if i >= 60 {
			p.sugar.Error("検索タイムアウト", i)
			return false, fmt.Errorf("検索タイムアウト count: %v", i)
		}
	}
	time.Sleep(time.Millisecond * 100)

	// 該当あるか確認
	doc, err := p.getWebDocument(page)
	tr := doc.Find(`#app > div > div.contents-wrapper > div.main-wrapper > div > div > div > form > table > tbody > tr`)
	trs := tr.Nodes
	if err != nil || len(trs) < 2 {
		// 2行(trが2つより少ない)場合は該当なし
		msg := fmt.Sprintf(
			"作業No(%v):DET番号(%v)は該当案件がありません",
			proposition.WorkedNumber,
			proposition.DET,
		)
		p.sugar.Errorf(msg)
		return false, errors.New(msg)
	}
	p.sugar.Infof("案件該当: 作業No(%v) DET番号(%v) nodes: %v", proposition.WorkedNumber, proposition.DET, len(trs))
	return true, nil
}

type successful int

const (
	success successful = iota
	failure
	unspecified
)

func (p *PropositionTable) updatedDeliveryDate(
	page *agouti.Page,
	diff monorevo.DifferentProposition,
) (successful, error) {
	// htmlをパースする
	contentsDom, err := p.getWebDocument(page)
	if err != nil {
		return unspecified, fmt.Errorf("htmlをパースする error: %v", err)
	}

	// tbodySelectionを取得して td要素数を取得する
	// 1Recordにつき2行なので倍になっている
	rows := p.getSearchResults(contentsDom)
	if len(rows) > 2 {
		return unspecified, fmt.Errorf(
			"作業No(%v)とDET(%v)で検索した結果が2レコード以上あるため中止する record: %v",
			diff.WorkedNumber,
			diff.DET,
			(len(rows) / 2))
	}

	// 表をループして納期を更新する
	// 作業NoとDETで検索しているので 原則1レコードだけど
	for i := 1; i <= len(rows); i += 2 {
		// 表中の作業No
		wk := contentsDom.Find(fmt.Sprintf("#app > div > div.contents-wrapper > div.main-wrapper > div > div > div > form > table > tbody > tr:nth-child(%d) > td:nth-child(2)", i)).Text()
		// 表中のDET番号
		dt := contentsDom.Find(fmt.Sprintf("#app > div > div.contents-wrapper > div.main-wrapper > div > div > div > form > table > tbody > tr:nth-child(%d) > td:nth-child(1)", i+1)).Text()
		p.sugar.Infof("処理中の案件: 作業No(%v) DET番号(%v)", wk, dt)

		if diff.WorkedNumber != wk && diff.DET != dt {
			// たまに検索に失敗していることがあったので保険的に比較する
			msg := fmt.Sprintf("検索失敗 期待の作業No(%v) 表示の作業No(%v)が相違している", diff.WorkedNumber, wk)
			p.sugar.Errorf(msg)
			return unspecified, errors.New(msg)
		}

		// 詳細画面を開く
		if err := p.openPropositionDETail(page, i); err != nil {
			return failure,
				fmt.Errorf("案件詳細が開けませんでした error: %v", err)
		}

		// 案件編集ウィンドウを開く
		if err := p.openEditableProposition(page); err != nil {
			return failure,
				fmt.Errorf("案件編集ウィンドウが開けませんでした error: %v", err)
		}

		// 編集する
		updatedDeliveryDateStr := diff.UpdatedDeliveryDate.Format("2006/01/02")
		if err := p.editProposition(page, updatedDeliveryDateStr); err != nil {
			return failure,
				fmt.Errorf(
					"作業No(%v) DET番号(%v)の編集ができませんでした error: %v",
					diff.WorkedNumber,
					diff.DET,
					err,
				)
		}
		p.sugar.Infof(
			"更新: 作業No(%v) DET番号(%v): 納期 %v -> %v (サンドボックスモード: %v)",
			diff.WorkedNumber,
			diff.DET,
			diff.DeliveryDate,
			diff.UpdatedDeliveryDate,
			p.sandboxMode,
		)
		time.Sleep(time.Millisecond * 50)

		// エラー表示を確認
		parent := page.FindByXPath(`/html/body/div[2]`)
		pid, _ := parent.Attribute("id") // idが動的に変わる

		dlg := page.FindByXPath(`/html/body/div[2]/div`)
		for i := 0; ; i++ {
			if v, err := dlg.Visible(); err != nil {
				p.sugar.Info("更新結果ダイアログが閉じたのを確認")
				break
			} else if v {
				// ダイアログ表示された
				doc, err := p.getWebDocument(page)
				if err != nil {
					return failure, fmt.Errorf("ドキュメントの取得に失敗しました error: %v", err)
				}

				sel := doc.Find(fmt.Sprintf("#%v > div", pid))
				msg := sel.Text()
				if msg != "データの登録が完了しました" {
					return failure, fmt.Errorf("更新に失敗しました message: %v", msg)
				}
				break
			}
			p.sugar.Infof("更新結果ダイアログ消失待ち %v * 100ミリ秒", i+1)
			time.Sleep(time.Millisecond * 100)

			if i >= 600 {
				p.sugar.Error("更新結果ダイアログ消失待ち タイムアウト", i)
				return failure, fmt.Errorf("更新結果ダイアログ消失待ちタイムアウト error: %v", i)
			}
		}
	}

	return success, nil
}

func (p *PropositionTable) editProposition(page *agouti.Page, updatedDeliveryDateStr string) error {
	// 納期に入力
	deliveryDateFld := page.FindByXPath(`//*[@id="deliveryDate"]/div[2]/div/input`)
	deliveryDateFld.Fill(updatedDeliveryDateStr)

	if p.sandboxMode {
		// サンドボックスモードのときは バツボタンを押して終わる
		xBtn := page.FindByXPath(`/html/body/div[3]/div[1]/div/div/header/button`) // idが動的に変わる
		xBtn.Click()
		time.Sleep(time.Millisecond * 100)
		return nil
	}

	// 登録して案件一覧に移動ボタンを押す
	entryNextBtn := page.FindByXPath(`//*[@id="smlot-detail"]/div/div/div/form/div[4]/div/button[2]`)
	entryNextBtn.Click()

	time.Sleep(time.Second * 2)
	// くるくる回るエフェクトのxpath
	selector := page.FindByXPath(`//*[@id="app"]/div/div[2]/div[2]/div/div[2]`)
	for i := 0; ; i++ {
		// 処理中の子要素(DIV)が存在する間はクリックしてもエラーにならない
		if err := selector.Click(); err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)

		if i >= 60 {
			p.sugar.Error("検索タイムアウト", i)
			return fmt.Errorf("検索タイムアウト error: %v", i)
		}
	}
	return nil
}

func (p *PropositionTable) openEditableProposition(page *agouti.Page) error {
	// 計画変更ボタンを押す
	updPlanBtn := page.FindByXPath(`//*[@id="smlot-detail"]/div/div/div/div/div[1]/div[1]/button[1]`)
	updPlanBtn.Click()

	entBtn := page.FindByXPath(`//*[@id="smlot-detail"]/div/div/div/form/div[4]/div/button[4]`)
	for i := 0; ; i++ {
		if _, err := entBtn.Enabled(); err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)

		if i >= 60 {
			p.sugar.Error("案件編集を開くタイムアウト", i)
			return fmt.Errorf("案件編集を開くタイムアウト count: %v", i)
		}
	}
	return nil
}

func (p *PropositionTable) openPropositionDETail(page *agouti.Page, row int) error {
	// 詳細ボタンを押す
	xpath := `//*[@id="app"]/div/div[2]/div[2]/div/div/div/form/table/tbody/` +
		fmt.Sprintf("tr[%d]", row) +
		`/td[10]/a`
	detailBtn := page.FindByXPath(xpath)
	detailBtn.Click()

	// 詳細が開くまで待つ
	detailEffect := page.FindByXPath(`//*[@id="smlot-detail"]/div/div/div/div/div[9]`)
	for j := 0; j < 60; j++ {
		// くるくる回るエフェクトのxpath
		err := detailEffect.Click()
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)

		if j >= 60 {
			p.sugar.Error("詳細を開くタイムアウト", j)
			return fmt.Errorf("詳細を開くタイムアウト count: %v", j)
		}
	}
	return nil
}

func (p *PropositionTable) getSearchResults(contentsDom *goquery.Document) []*html.Node {
	tbodySelection := contentsDom.Find(`#app > div > div.contents-wrapper > div.main-wrapper > div > div > div > form > table > tbody`)
	rowSelection := tbodySelection.Children()

	// 1Recordにつき2行なので倍になっている
	rows := rowSelection.Nodes
	p.sugar.Debugf("案件一覧テーブル %vレコード", (len(rows) / 2))
	return rows
}

func (p *PropositionTable) getWebDocument(page *agouti.Page) (*goquery.Document, error) {
	curContentsDom, err := page.HTML()
	if err != nil {
		p.sugar.Error("DOMの取得に失敗しました", err)
		return nil, fmt.Errorf("DOMの取得に失敗しました error: %v", err)
	}

	readerCurContents := strings.NewReader(curContentsDom)

	contentsDom, err := goquery.NewDocumentFromReader(readerCurContents)
	if err != nil {
		p.sugar.Error("パースに失敗しました", err)
		return nil, fmt.Errorf("パースに失敗しました error: %v", err)
	}
	return contentsDom, nil
}
