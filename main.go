package main

import (
	"flag"
	"fmt"
	"ofxmerge/internal/merge"
	"os"

	"github.com/aclindsa/ofxgo"
	"github.com/aclindsa/xml"
)

type OfxBank struct {
	XMLName      xml.Name               `xml:"OFX"`
	Status       ofxgo.Status           `xml:"STATUS"`
	CurDef       ofxgo.CurrSymbol       `xml:"STMTRS>CURDEF"`
	BankTranList *ofxgo.TransactionList `xml:"STMTRS>BANKTRANLIST,omitempty"`
}

func NewOfxBank() *OfxBank {
	transactionList := ofxgo.TransactionList{
		Transactions: []ofxgo.Transaction{},
	}
	return &OfxBank{
		BankTranList: &transactionList,
	}
}

func main() {
	cmd := flag.NewFlagSet("ofxmerge", flag.ExitOnError)
	cmd.Parse(os.Args[1:])

	o := merge.NewOFXMerger()

	for _, filename := range cmd.Args() {
		b, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
			os.Exit(1)
		}

		if err := o.Add(b); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", filename, err)
			os.Exit(1)
		}
	}

	b, err := o.Merge()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", b)
}
