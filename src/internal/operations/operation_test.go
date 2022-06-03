package operations

import (
	"testing"

	"github.com/alcionai/corso/internal/kopia"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type OperationSuite struct {
	suite.Suite
}

func TestOperationSuite(t *testing.T) {
	suite.Run(t, new(OperationSuite))
}

func (suite *OperationSuite) TestNewOperation() {
	t := suite.T()
	op := newOperation(OperationOpts{}, nil)
	assert.NotNil(t, op.Errors)
	assert.NotNil(t, op.ID)
}

func (suite *OperationSuite) TestOperation_Validate() {
	table := []struct {
		name     string
		kw       *kopia.KopiaWrapper
		errCheck assert.ErrorAssertionFunc
	}{
		{"good", new(kopia.KopiaWrapper), assert.NoError},
		{"missing kopia", nil, assert.Error},
	}
	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			op := newOperation(OperationOpts{}, test.kw)
			test.errCheck(t, op.validate())
		})
	}
}
