//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Test = newTestTable("public", "test", "")

type testTable struct {
	postgres.Table

	// Columns
	ID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TestTable struct {
	testTable

	EXCLUDED testTable
}

// AS creates new TestTable with assigned alias
func (a TestTable) AS(alias string) *TestTable {
	return newTestTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TestTable with assigned schema name
func (a TestTable) FromSchema(schemaName string) *TestTable {
	return newTestTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TestTable with assigned table prefix
func (a TestTable) WithPrefix(prefix string) *TestTable {
	return newTestTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TestTable with assigned table suffix
func (a TestTable) WithSuffix(suffix string) *TestTable {
	return newTestTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTestTable(schemaName, tableName, alias string) *TestTable {
	return &TestTable{
		testTable: newTestTableImpl(schemaName, tableName, alias),
		EXCLUDED:  newTestTableImpl("", "excluded", ""),
	}
}

func newTestTableImpl(schemaName, tableName, alias string) testTable {
	var (
		IDColumn       = postgres.IntegerColumn("id")
		allColumns     = postgres.ColumnList{IDColumn}
		mutableColumns = postgres.ColumnList{}
	)

	return testTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID: IDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
