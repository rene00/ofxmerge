package merge

import (
	"bytes"
	"fmt"

	"github.com/aclindsa/ofxgo"
)

type Merger interface {
	Add([]byte) error
	Merge() ([]byte, error)
}

func NewOFXMerger() Merger {
	return &OFXMerger{}
}

type OFXMerger struct {
	ofxFiles ofxFiles
}

func ccAccountsMatch(c1, c2 ofxgo.CCAcct) bool {
	if c1.AcctID != c2.AcctID {
		return false
	}
	if c1.AcctKey != c2.AcctKey {
		return false
	}

	return true

}

func bankAccountsMatch(b1, b2 ofxgo.BankAcct) bool {
	if b1.BankID != b2.BankID {
		return false
	}
	if b1.BranchID != b2.BranchID {
		return false
	}
	if b1.AcctID != b2.AcctID {
		return false
	}
	if b1.AcctType != b2.AcctType {
		return false
	}
	if b1.AcctKey != b2.AcctKey {
		return false
	}
	return true
}

// Add accepts a slice of bytes, parses bytes into a response and adds response
// to ofxFiles.
func (o *OFXMerger) Add(b []byte) error {
	var err error

	resp, err := ofxgo.ParseResponse(bytes.NewReader(b))
	if err != nil {
		return err
	}

	o.ofxFiles.statementType, err = getStatementType(resp)
	if err != nil {
		return err
	}

	ofxFile := ofxFile{resp: resp}

	if err = o.ofxFiles.validate(resp); err != nil {
		return err
	}

	switch o.ofxFiles.statementType {
	case statementTypeBank:
		if stmt, ok := resp.Bank[0].(*ofxgo.StatementResponse); ok {
			o.ofxFiles.currSymbol = stmt.CurDef
			if o.ofxFiles.bankAccount != nil {
				if !bankAccountsMatch(*o.ofxFiles.bankAccount, stmt.BankAcctFrom) {
					return fmt.Errorf("bank accounts do not match")
				}
			} else {
				o.ofxFiles.bankAccount = &stmt.BankAcctFrom
			}
		}
	case statementTypeCreditCard:
		if stmt, ok := resp.CreditCard[0].(*ofxgo.CCStatementResponse); ok {
			o.ofxFiles.currSymbol = stmt.CurDef

			if o.ofxFiles.ccAccount != nil {
				if !ccAccountsMatch(*o.ofxFiles.ccAccount, stmt.CCAcctFrom) {
					return fmt.Errorf("credit card accounts do not match")
				}
			} else {
				o.ofxFiles.ccAccount = &stmt.CCAcctFrom
			}
		}
	}

	o.ofxFiles.files = append(o.ofxFiles.files, ofxFile)
	return nil
}

func newStatementResponse() (ofxgo.StatementResponse, error) {
	var err error
	var resp ofxgo.StatementResponse
	resp.Status = ofxgo.Status{
		Code:     ofxgo.Int(0),
		Severity: ofxgo.String("INFO"),
	}

	trnUID, err := ofxgo.RandomUID()
	if err != nil {
		return resp, err
	}
	resp.TrnUID = *trnUID

	return resp, nil
}

func (o *OFXMerger) Merge() ([]byte, error) {
	var buf *bytes.Buffer

	switch o.ofxFiles.statementType {
	case statementTypeBank:
		stmt, err := newStatementResponse()
		if err != nil {
			return buf.Bytes(), err
		}

		stmt.CurDef = o.ofxFiles.currSymbol
		stmt.DtAsOf = o.ofxFiles.dtAsOf()
		stmt.BankAcctFrom = *o.ofxFiles.bankAccount
		stmt.BankTranList = &ofxgo.TransactionList{
			Transactions: []ofxgo.Transaction{},
			DtEnd:        o.ofxFiles.dtEnd(),
			DtStart:      o.ofxFiles.dtStart(),
		}
		for _, ofxFile := range o.ofxFiles.files {
			ofxFileStmt, ok := ofxFile.resp.Bank[0].(*ofxgo.StatementResponse)
			if !ok {
				return buf.Bytes(), fmt.Errorf("unable to process bank statement")
			}
			for _, i := range ofxFileStmt.BankTranList.Transactions {
				stmt.BankTranList.Transactions = append(stmt.BankTranList.Transactions, i)
			}
		}
		resp := ofxgo.Response{
			Bank: []ofxgo.Message{&stmt},
			Signon: ofxgo.SignonResponse{
				Language: ofxgo.String("ENG"),
				Status: ofxgo.Status{
					Code:     ofxgo.Int(0),
					Severity: ofxgo.String("INFO"),
				},
			},
		}

		buf, err = resp.Marshal()
		if err != nil {
			return buf.Bytes(), err
		}
		return buf.Bytes(), nil
	default:
	}

	return buf.Bytes(), nil
}
