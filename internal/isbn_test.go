package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type tc struct {
	isbn string
	err  error
}

var cases = map[string]tc{
	"isbn-10":            {`2-266-11156-6`, nil},
	"isbn-10-x":          {`0-712-67263-X`, nil},
	"isbn-13":            {`978-5-907488-10-6`, nil},
	"isbn-13-2":          {`978-2-743601-72-0`, nil},
	"isbn-empty":         {``, errEmptyString},
	"isbn-len":           {`978-5-907488-10-`, errUnexpectedFormat},
	"isbn-10-format-err": {`2-fff-11155-6`, errUnexpectedFormat},
	"isbn-13-format-err": {`978-5-907488-ff-6`, errUnexpectedFormat},
	"isbn-10-sum-err":    {`2-266-11155-6`, errInvalidNumber},
	"isbn-13-sum-err":    {`978-5-907488-11-6`, errInvalidNumber},
}

func TestCheckISBN(t *testing.T) {
	t.Parallel()

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := CheckISBN(c.isbn)
			assert.ErrorIs(t, got, c.err)
		})
	}
}
