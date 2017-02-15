package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func convDayOfWeek(str string) string {
	r := strings.NewReplacer("曜日", "", "－", "")
	return r.Replace(str)
}

func convTwiceAWeek(str string) string {
	r := strings.NewReplacer("曜日", "", "・", " ", "－", "")
	return r.Replace(str)
}

func convTwiceAMonth(str string) string {
	r := regexp.MustCompile(`第(\d)・(\d)\s([月火水木金土日])`)
	return r.ReplaceAllString(convDayOfWeek(str), "$3$1 $3$2")

}

func mkAreadays() error {
	url := "https://www.city.edogawa.tokyo.jp/gomi_recycle/yobihyo.html"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
    return err
	}
	defer resp.Body.Close()

	decodedReader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	doc, _ := goquery.NewDocumentFromReader(decodedReader)

	fmt.Println("地区,センター,燃やすごみ,燃やさないごみ,資源")
	doc.Find("table.table01").Each(func(_ int, table *goquery.Selection) {
		area1 := ""
		table.Find("tr").Each(func(_ int, tr *goquery.Selection) {
			data := tr.Find("td")
			if data.Length() == 0 {
				return
			}
			base := 0
			if data.Length() == 6 {
				area1 = data.Eq(0).Text()
				base = 1
			}
			area2 := data.Eq(base).Text()
			recyclable := convDayOfWeek(strings.Replace(data.Eq(base+1).Text(), "資源", "", 1))
			burnable := convTwiceAWeek(strings.Replace(data.Eq(base+2).Text(), "燃やすごみ", "", 1))
			unburnable := convTwiceAMonth(strings.Replace(data.Eq(base+3).Text(), "燃やさないごみ", "", 1))
			center := strings.Replace(data.Eq(base+4).Text(), "管轄", "", 1)

			if burnable == "" && unburnable == "" && recyclable == "" {
				return
			}
			fmt.Printf("%s%s,%s,%s,%s,%s\n", area1, area2, center, burnable, unburnable, recyclable)
		})
	})
  return nil
}
