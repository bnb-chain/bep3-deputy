package util

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalcActualOutAmount(t *testing.T) {
	tests := []struct {
		amount   *big.Int
		ratio    *big.Float
		fixedFee *big.Int
		res      *big.Int
	}{
		{
			amount:   big.NewInt(1000),
			ratio:    big.NewFloat(1.0),
			fixedFee: big.NewInt(100),
			res:      big.NewInt(900),
		},
	}

	for _, test := range tests {
		res := CalcActualOutAmount(test.amount, test.ratio, test.fixedFee)
		require.Equal(t, res.Cmp(test.res), 0)
	}
}
