package main

import (
	"fmt"
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TestCSVSuite struct{}

var _ = Suite(&TestCSVSuite{})

// Test calling createRowToJSON with a values' length higher than values' length
func (s *TestCSVSuite) TestCreateRowToJSON_ColumnLengthMustBeAtLeastEqualToValuesLength(c *C) {

	columns := []string{"c1"}
	values := []string{"v1", "v2"}

	_, err := createRowToJSON(columns, values)
	c.Assert(err, NotNil)
}

// Test calling createRowToJSON with two empty (nil value) slices
func (s *TestCSVSuite) TestCreateRowToJSON_EmptyColumnsAndValues(c *C) {

	var columns, values []string

	ret, err := createRowToJSON(columns, values)

	c.Assert(err, IsNil)
	c.Assert(fmt.Sprintf("%s", ret), Equals, "{}")
}

// Test calling createRowToJSON with a valid slice of columns but an empty slice of values
func (s *TestCSVSuite) TestCreateRowToJSON_ValidButEmptyValuesAndValidColumns(c *C) {

	columns := []string{"c1"}
	var values []string //make([]string, 1)

	ret, err := createRowToJSON(columns, values)

	c.Assert(err, IsNil)
	c.Assert(fmt.Sprintf("%s", ret), Equals, "{}")
}

// Test calling createRowToJSON with a valid slice of columns and a valid slice of values
func (s *TestCSVSuite) TestCreateRowToJSON_ValidValuesAndColumns(c *C) {

	columns := []string{"c1", "c2"}
	values := []string{"v1", "v2"}

	ret, err := createRowToJSON(columns, values)

	c.Assert(err, IsNil)
	c.Assert(fmt.Sprintf("%s", ret), Equals, "{\"c1\":\"v1\",\"c2\":\"v2\"}")
}
