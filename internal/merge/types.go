package merge

import (
	"fmt"

	"github.com/aclindsa/ofxgo"
)

const (
	statementTypeUnsupported statementType = iota
	statementTypeBank
	statementTypeCreditCard
)

type statementType int

func (s statementType) String() string {
	switch s {
	case statementTypeBank:
		return "bank"
	case statementTypeCreditCard:
		return "creditcard"
	default:
		return "unsupported"
	}
}

func getStatementType(resp *ofxgo.Response) (statementType, error) {
	if resp == nil {
		return statementTypeUnsupported, fmt.Errorf("nil response")
	}

	if len(resp.Bank) > 0 {
		if _, ok := resp.Bank[0].(*ofxgo.StatementResponse); ok {
			return statementTypeBank, nil
		}
	}

	if len(resp.CreditCard) > 0 {
		if _, ok := resp.CreditCard[0].(*ofxgo.CCStatementResponse); ok {
			return statementTypeCreditCard, nil
		}
	}

	return statementTypeUnsupported, fmt.Errorf("unsupported statement response")
}
