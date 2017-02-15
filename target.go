package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func getSpan(s *goquery.Selection, t string) int {
	span := 1
	str, exist := s.Attr(t + "span")
	if exist {
		span, _ = strconv.Atoi(str)
	}
	return span
}

func getText(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
}

func mkTarget() error {
	url := "https://www.city.edogawa.tokyo.jp/gomi_recycle/hinmoku/"
	files := []string{"a.html", "ka.html", "sa.html", "ta.html",
		"na.html", "ha.html", "ma.html", "ya.html", "ra.html"}

	fmt.Println("type,name,notice,furigana")

	for _, fileName := range files {
		resp, err := http.Get(url + fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return err
		}
		defer resp.Body.Close()

		decodedReader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
		doc, _ := goquery.NewDocumentFromReader(decodedReader)

		doc.Find("table.table01").Each(func(_ int, table *goquery.Selection) {
			furigana := ""
			name := ""
			notice := ""
			nameRowspan := 0
			noticeRowspan := 0
			hasName2 := false
			table.Find("tr").Each(func(r int, tr *goquery.Selection) {
				c := 0
				name2 := ""
				kind := ""
				nameColspan := 1
				if r == 0 { // 品目が２カラムかチェック
					tr.Find("th[colspan='2']").Each(func(_ int, s *goquery.Selection) {
						hasName2 = true
					})
					return // ヘッダ行を無視
				}

				// ふりがなを取得
				tr.Find("th > strong").Each(func(_ int, s *goquery.Selection) {
					furigana = getText(s)
				})
				// thでなくtdが使われている場合(ta.html)のため
				tr.Find("td > strong").Each(func(_ int, s *goquery.Selection) {
					furigana = getText(s)
					c++
				})

				data := tr.Find("td")

				// 名前の取得
				if nameRowspan == 0 {
					name = getText(data.Eq(c))
					nameRowspan = getSpan(data.Eq(c), "row")
					nameColspan = getSpan(data.Eq(c), "col")
					c++
				}
				nameRowspan--

				// 名前の詳細取得
				if hasName2 && nameColspan == 1 {
					name2 = "（" + getText(data.Eq(c)) + "）"
					c++
				}

				// 種別の取得
				kind = getText(data.Eq(c))
				c++
				// 注意書きの取得 (乾電池のところにrowspan=4あり)
				if noticeRowspan == 0 {
					notice = getText(data.Eq(c))
					noticeRowspan = getSpan(data.Eq(c), "row")
				}
				noticeRowspan--
				if kind == "燃やすごみ" || kind == "燃やさないごみ" || kind == "資源" {
          // FIXME: データがきれいに並んでないところがある、個別に修正
          if (name2 == "（ペットボトル）") {
            notice += "（キャップ、外装ラベル→容リプラ回収）"
          }
          if (name == "コーヒー用ミルクのパック") {
            notice += "（ふた等裏側が銀色のものは燃やすごみ）"
          }
					fmt.Printf("%s,%s%s,%s,%s\n", kind, name, name2, notice, furigana)
				}
			})
		})
	}
	return nil
}
