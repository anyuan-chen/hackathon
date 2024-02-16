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

var Tokens = newTokensTable("public", "tokens", "")

type tokensTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnInteger
	BearerToken postgres.ColumnString
	ExpiryTime  postgres.ColumnInteger
	UserID      postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type TokensTable struct {
	tokensTable

	EXCLUDED tokensTable
}

// AS creates new TokensTable with assigned alias
func (a TokensTable) AS(alias string) *TokensTable {
	return newTokensTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TokensTable with assigned schema name
func (a TokensTable) FromSchema(schemaName string) *TokensTable {
	return newTokensTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TokensTable with assigned table prefix
func (a TokensTable) WithPrefix(prefix string) *TokensTable {
	return newTokensTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TokensTable with assigned table suffix
func (a TokensTable) WithSuffix(suffix string) *TokensTable {
	return newTokensTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTokensTable(schemaName, tableName, alias string) *TokensTable {
	return &TokensTable{
		tokensTable: newTokensTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newTokensTableImpl("", "excluded", ""),
	}
}

func newTokensTableImpl(schemaName, tableName, alias string) tokensTable {
	var (
		IDColumn          = postgres.IntegerColumn("id")
		BearerTokenColumn = postgres.StringColumn("bearer_token")
		ExpiryTimeColumn  = postgres.IntegerColumn("expiry_time")
		UserIDColumn      = postgres.IntegerColumn("user_id")
		allColumns        = postgres.ColumnList{IDColumn, BearerTokenColumn, ExpiryTimeColumn, UserIDColumn}
		mutableColumns    = postgres.ColumnList{BearerTokenColumn, ExpiryTimeColumn, UserIDColumn}
	)

	return tokensTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		BearerToken: BearerTokenColumn,
		ExpiryTime:  ExpiryTimeColumn,
		UserID:      UserIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
