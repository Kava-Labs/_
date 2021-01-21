package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Deposit defines an amount of coins deposited into a hard module account
type Deposit struct {
	Depositor sdk.AccAddress        `json:"depositor" yaml:"depositor"`
	Amount    sdk.Coins             `json:"amount" yaml:"amount"`
	Index     SupplyInterestFactors `json:"index" yaml:"index"`
}

// NewDeposit returns a new deposit
func NewDeposit(depositor sdk.AccAddress, amount sdk.Coins, indexes SupplyInterestFactors) Deposit {
	return Deposit{
		Depositor: depositor,
		Amount:    amount,
		Index:     indexes,
	}
}

// Validate deposit validation
func (d Deposit) Validate() error {
	if d.Depositor.Empty() {
		return fmt.Errorf("Depositor cannot be empty")
	}
	if !d.Amount.IsValid() {
		return fmt.Errorf("Invalid deposit coins: %s", d.Amount)
	}

	if err := d.Index.Validate(); err != nil {
		return err
	}

	return nil
}

func (d Deposit) String() string {
	return fmt.Sprintf(`Deposit:
	%s
	%s
	%s
	`, d.Depositor, d.Amount, d.Index)
}

// Deposits is a slice of Deposit
type Deposits []Deposit

// Validate validates Deposits
func (ds Deposits) Validate() error {
	depositDupMap := make(map[string]Deposit)
	for _, d := range ds {
		if err := d.Validate(); err != nil {
			return err
		}
		dup, ok := depositDupMap[d.Depositor.String()]
		if ok {
			return fmt.Errorf("duplicate depositor: %s\n%s", d, dup)
		}
		depositDupMap[d.Depositor.String()] = d
	}
	return nil
}

// SupplyInterestFactor defines an individual borrow interest factor
type SupplyInterestFactor struct {
	Denom string  `json:"denom" yaml:"denom"`
	Value sdk.Dec `json:"value" yaml:"value"`
}

// NewSupplyInterestFactor returns a new SupplyInterestFactor instance
func NewSupplyInterestFactor(denom string, value sdk.Dec) SupplyInterestFactor {
	return SupplyInterestFactor{
		Denom: denom,
		Value: value,
	}
}

// Validate validates SupplyInterestFactor values
func (sif SupplyInterestFactor) Validate() error {
	if strings.TrimSpace(sif.Denom) == "" {
		return fmt.Errorf("supply interest factor denom cannot be empty")
	}
	if sif.Value.IsNegative() {
		return fmt.Errorf("supply interest factor value cannot be negative: %s", sif)

	}
	return nil
}

func (sif SupplyInterestFactor) String() string {
	return fmt.Sprintf(`[%s,%s]
	`, sif.Denom, sif.Value)
}

// SupplyInterestFactors is a slice of SupplyInterestFactor, because Amino won't marshal maps
type SupplyInterestFactors []SupplyInterestFactor

// Validate validates SupplyInterestFactors
func (sifs SupplyInterestFactors) Validate() error {
	for _, sif := range sifs {
		if err := sif.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (sifs SupplyInterestFactors) String() string {
	out := ""
	for _, sif := range sifs {
		out += sif.String()
	}
	return out
}
