package merge

import (
	"fmt"
	"time"

	"github.com/aclindsa/ofxgo"
)

type ofxFiles struct {
	statementType statementType
	currSymbol    ofxgo.CurrSymbol
	bankAccount   *ofxgo.BankAcct
	ccAccount     *ofxgo.CCAcct
	files         []ofxFile
}

func (o ofxFiles) validate(resp *ofxgo.Response) error {
	if len(o.files) == 0 {
		return nil
	}

	respType, err := getStatementType(resp)
	if err != nil {
		return err
	}

	if respType != o.statementType {
		return fmt.Errorf("statement is different type to previous statement (%v, %v)", respType, o.statementType)
	}

	return nil
}

func (o ofxFiles) dtAsOf() ofxgo.Date {
	d := ofxgo.Date{}
	for _, ofxFile := range o.files {
		switch o.statementType {
		case statementTypeBank:
			if len(ofxFile.resp.Bank) == 0 {
				continue
			}
			stmt, ok := ofxFile.resp.Bank[0].(*ofxgo.StatementResponse)
			if !ok {
				continue
			}
			if d.Time.IsZero() {
				d = stmt.DtAsOf
				continue
			}
			if stmt.DtAsOf.Time.Before(d.Time) {
				d = stmt.DtAsOf
				continue
			}
		case statementTypeCreditCard:
			if len(ofxFile.resp.CreditCard) == 0 {
				continue
			}
			stmt, ok := ofxFile.resp.CreditCard[0].(*ofxgo.StatementResponse)
			if !ok {
				continue
			}
			if d.Time.IsZero() {
				d = stmt.DtAsOf
				continue
			}
			if stmt.DtAsOf.Time.Before(d.Time) {
				d = stmt.DtAsOf
				continue
			}
		}
	}

	if d.Time.IsZero() {
		now := time.Now()
		return *ofxgo.NewDate(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second(),
			now.Nanosecond(),
			now.Location(),
		)
	}
	return d
}

func (o ofxFiles) dtStart() ofxgo.Date {
	d := time.Time{}
	for _, ofxFile := range o.files {
		switch o.statementType {
		case statementTypeBank:
			if len(ofxFile.resp.Bank) == 0 {
				continue
			}
			stmt, ok := ofxFile.resp.Bank[0].(*ofxgo.StatementResponse)
			if !ok {
				continue
			}
			if stmt.BankTranList == nil {
				continue
			}
			if d.IsZero() {
				d = stmt.BankTranList.DtStart.Time
				continue
			}
			if stmt.BankTranList.DtStart.Time.After(d) {
				d = stmt.BankTranList.DtStart.Time
				continue
			}
		case statementTypeCreditCard:
			if len(ofxFile.resp.CreditCard) == 0 {
				continue
			}
			stmt, ok := ofxFile.resp.CreditCard[0].(*ofxgo.CCStatementResponse)
			if !ok {
				continue
			}
			if stmt.BankTranList == nil {
				continue
			}
			if d.IsZero() {
				d = stmt.BankTranList.DtStart.Time
				continue
			}
			if stmt.BankTranList.DtStart.Time.After(d) {
				d = stmt.BankTranList.DtStart.Time
				continue
			}
		}
	}
	if d.IsZero() {
		d = time.Now()
	}
	return *ofxgo.NewDate(
		d.Year(),
		d.Month(),
		d.Day(),
		d.Hour(),
		d.Minute(),
		d.Second(),
		d.Nanosecond(),
		d.Location(),
	)
}

func (o ofxFiles) dtEnd() ofxgo.Date {
	d := time.Time{}
	for _, ofxFile := range o.files {
		switch o.statementType {
		case statementTypeBank:
			if len(ofxFile.resp.Bank) == 0 {
				continue
			}
			stmt, ok := ofxFile.resp.Bank[0].(*ofxgo.StatementResponse)
			if !ok {
				continue
			}
			if stmt.BankTranList == nil {
				continue
			}
			if d.IsZero() {
				d = stmt.BankTranList.DtEnd.Time
				continue
			}
			if stmt.BankTranList.DtEnd.Time.Before(d) {
				d = stmt.BankTranList.DtEnd.Time
				continue
			}
		case statementTypeCreditCard:
			if len(ofxFile.resp.Bank) == 0 {
				continue
			}
			stmt, ok := ofxFile.resp.Bank[0].(*ofxgo.StatementResponse)
			if !ok {
				continue
			}
			if stmt.BankTranList == nil {
				continue
			}
			if d.IsZero() {
				d = stmt.BankTranList.DtEnd.Time
				continue
			}
			if stmt.BankTranList.DtEnd.Time.Before(d) {
				d = stmt.BankTranList.DtEnd.Time
				continue
			}
		}
	}
	if d.IsZero() {
		d = time.Now()
	}
	return *ofxgo.NewDate(
		d.Year(),
		d.Month(),
		d.Day(),
		d.Hour(),
		d.Minute(),
		d.Second(),
		d.Nanosecond(),
		d.Location(),
	)
}

type ofxFile struct {
	resp *ofxgo.Response
}
