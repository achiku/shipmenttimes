package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	dojo = "--- 同上 ---"
)

var (
	summaryHeader = []string{
		"注文番号",
		"注文時間",
		"郵便番号",
		"都道府県",
		"住所1",
		"住所2",
		"氏名",
		"商品",
		"総商品数",
	}
	clickpostHeader = []string{
		"お届け先郵便番号",
		"お届け先氏名",
		"お届け先敬称",
		"お届け先住所1行目",
		"お届け先住所2行目",
		"お届け先住所3行目",
		"お届け先住所4行目",
		"内容品",
	}
)

// BaseOrderRaw thebase.in stupid csv format
type BaseOrderRaw struct {
	OrderID           string // 1. 注文ID
	OrderTime         string // 2. 注文日時
	FirstName         string // 3. 氏(配送先)
	LastName          string // 4. 名(配送先)
	DeliveryCompanyID string // 5. 配送業者ID
	OrderNumber       string // 6. 伝票番号
	MessageTemplateID string // 7. メッセージテンプレートID
	PostalCode        string // 8. 郵便番号(配送先)
	Prefecture        string // 9. 都道府県(配送先)
	Address1          string // 10. 住所(配送先)
	Address2          string // 11. 住所2(配送先)
	PhoneNumber       string // 12. 電話番号(配送先)
	Note              string // 13. 備考
	ItemCode          string // 14. 商品コード
	ItemName          string // 15. 商品名
	ItemType          string // 16. バリエーション
	ItemQuantity      string // 17. 数量
}

// BaseOrder thebase.in order
type BaseOrder struct {
	OrderID       string
	OrderTime     string
	OrderNumber   string
	FirstName     string
	LastName      string
	PostalCode    string
	Prefecture    string
	Address1      string
	Address2      string
	TotalQuantity int
	OrderItems    []BaseItem
}

// ClickpostFormat clickpost format
// 1. お届け先郵便番号
// 2. お届け先氏名
// 3. お届け先敬称
// 4. お届け先住所1行目
// 5. お届け先住所2行目
// 6. お届け先住所3行目
// 7. お届け先住所4行目
// 8. 内容品
func (od *BaseOrder) ClickpostFormat() []string {
	return []string{
		od.PostalCode,
		fmt.Sprintf("%s%s", od.LastName, od.FirstName),
		"様",
		od.Prefecture,
		od.Address1,
		od.Address2,
		"",
		"CD",
	}
}

// SummaryFormat summary format
func (od BaseOrder) SummaryFormat() []string {
	var items string
	for _, itm := range od.OrderItems {
		items = fmt.Sprintf("%s %s", items, itm.String())
	}
	return []string{
		od.OrderID,
		od.OrderTime,
		od.PostalCode,
		od.Prefecture,
		od.Address1,
		od.Address2,
		fmt.Sprintf("%s%s", od.LastName, od.FirstName),
		items,
		fmt.Sprintf("%d", od.TotalQuantity),
	}
}

// WriteSummaryFormat  write summary
func WriteSummaryFormat(w io.Writer, odrs []BaseOrder, osType string) error {
	var writer *csv.Writer
	if osType == "win" {
		writer = csv.NewWriter(
			transform.NewWriter(w, japanese.ShiftJIS.NewEncoder()))
		writer.UseCRLF = true
	} else {
		writer = csv.NewWriter(w)
	}
	if err := writer.Write(summaryHeader); err != nil {
		return err
	}
	for _, odr := range odrs {
		if err := writer.Write(odr.SummaryFormat()); err != nil {
			return errors.Wrapf(err, "%s", odr.SummaryFormat())
		}
	}
	writer.Flush()
	return nil
}

// WriteClickpostFormat  write summary
func WriteClickpostFormat(w io.Writer, odrs []BaseOrder, osType string) error {
	var writer *csv.Writer
	if osType == "win" {
		writer = csv.NewWriter(
			transform.NewWriter(w, japanese.ShiftJIS.NewEncoder()))
		writer.UseCRLF = true
	} else {
		writer = csv.NewWriter(w)
	}
	if err := writer.Write(clickpostHeader); err != nil {
		return err
	}
	for _, odr := range odrs {
		if err := writer.Write(odr.ClickpostFormat()); err != nil {
			return errors.Wrapf(err, "%s", odr.ClickpostFormat())
		}
	}
	writer.Flush()
	return nil
}

// BaseItem base item
type BaseItem struct {
	ItemCode     string
	ItemName     string
	ItemType     string
	ItemQuantity int
}

func (bi *BaseItem) String() string {
	return fmt.Sprintf("%s %d個/", bi.ItemName, bi.ItemQuantity)
}

// TransformBaseOrder transform base order
func TransformBaseOrder(ods []BaseOrderRaw) ([]BaseOrder, error) {
	var bos []BaseOrder
	for i, od := range ods {
		// 要素がヘッダーの場合スキップ
		if i == 0 {
			continue
		}
		q, err := strconv.Atoi(od.ItemQuantity)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse %s", od.ItemQuantity)
		}
		// 同上の場合商品情報を前の行とマージする
		if od.OrderID == dojo {
			item := BaseItem{
				ItemCode:     od.ItemCode,
				ItemName:     od.ItemName,
				ItemType:     od.ItemType,
				ItemQuantity: q,
			}
			b := bos[len(bos)-1]
			b.TotalQuantity = b.TotalQuantity + q
			b.OrderItems = append(b.OrderItems, item)
			bos[len(bos)-1] = b
		} else {
			itms := []BaseItem{
				BaseItem{
					ItemCode:     od.ItemCode,
					ItemName:     od.ItemName,
					ItemType:     od.ItemType,
					ItemQuantity: q,
				},
			}
			order := BaseOrder{
				OrderID:       od.OrderID,
				OrderTime:     od.OrderTime,
				OrderNumber:   od.OrderNumber,
				FirstName:     od.FirstName,
				LastName:      od.LastName,
				PostalCode:    od.PostalCode,
				Prefecture:    od.Prefecture,
				Address1:      od.Address1,
				Address2:      od.Address2,
				TotalQuantity: q,
				OrderItems:    itms,
			}
			bos = append(bos, order)
		}
	}
	return bos, nil
}

// ParseBaseCSV parse base csv
func ParseBaseCSV(in io.Reader) ([]BaseOrderRaw, error) {
	reader := csv.NewReader(
		transform.NewReader(in, japanese.ShiftJIS.NewDecoder()))
	var ls []BaseOrderRaw
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// 空行を削除
		if rec[1] == "" {
			continue
		}
		l := BaseOrderRaw{
			OrderID:           rec[0],
			OrderTime:         rec[1],
			FirstName:         rec[2],
			LastName:          rec[3],
			DeliveryCompanyID: rec[4],
			OrderNumber:       rec[5],
			MessageTemplateID: rec[6],
			PostalCode:        rec[7],
			Prefecture:        rec[8],
			Address1:          rec[9],
			Address2:          rec[10],
			PhoneNumber:       rec[11],
			Note:              rec[12],
			ItemCode:          rec[13],
			ItemName:          rec[14],
			ItemType:          rec[15],
			ItemQuantity:      rec[16],
		}
		ls = append(ls, l)
	}
	return ls, nil
}

// QuantityFilter quantity filter
func QuantityFilter(ods []BaseOrder, sep int) ([]BaseOrder, []BaseOrder) {
	var a, b []BaseOrder
	for _, od := range ods {
		if od.TotalQuantity <= sep {
			a = append(a, od)
		} else {
			b = append(b, od)
		}
	}
	return a, b
}
