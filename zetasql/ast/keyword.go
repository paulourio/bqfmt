package ast

import (
	"strings"
)

func IsReservedKeyword(word string) bool {
	return reservedKeywordsMap[strings.ToUpper(word)]
}

var reservedKeywordsList = []string{
	"ALL",
	"AND",
	"ANY",
	"ARRAY",
	"AS",
	"ASC",
	"AT",
	"BETWEEN",
	"BY",
	"CASE",
	"CAST",
	"COLLATE",
	"CREATE",
	"CROSS",
	"CURRENT",
	"DEFAULT",
	"DEFINE",
	"DESC",
	"DISTINCT",
	"ELSE",
	"END",
	"ENUM",
	"EXCEPT",
	"EXISTS",
	"EXTRACT",
	"FALSE",
	"FOLLOWING",
	"FROM",
	"FULL",
	"GROUP",
	"GROUPING",
	"HASH",
	"HAVING",
	"IF",
	"IGNORE",
	"IN",
	"INNER",
	"INTERSECT",
	"INTERVAL",
	"INTO",
	"IS",
	"JOIN",
	"LEFT",
	"LIKE",
	"LIMIT",
	"LOOKUP",
	"MERGE",
	"NATURAL",
	"NEW",
	"NO",
	"NOT",
	"NULL",
	"NULLS",
	"ON",
	"OR",
	"ORDER",
	"OUTER",
	"OVER",
	"PARTITION",
	"PRECEDING",
	"PROTO",
	"RANGE",
	"RECURSIVE",
	"RESPECT",
	"RIGHT",
	"ROLLUP",
	"ROWS",
	"SELECT",
	"SET",
	"STRUCT",
	"TABLESAMPLE",
	"THEN",
	"TO",
	"TRUE",
	"UNBOUNDED",
	"UNION",
	"USING",
	"WHEN",
	"WHERE",
	"WINDOW",
	"WITH",
	"UNNEST",
	"CONTAINS",
	"CUBE",
	"ESCAPE",
	"EXCLUDE",
	"FETCH",
	"FOR",
	"GROUPS",
	"LATERAL",
	"OF",
	"SOME",
	"TREAT",
	"WITHIN",
	"QUALIFY",
}

var reservedKeywordsMap map[string]bool

func initKeywords() {
	reservedKeywordsMap = make(map[string]bool, len(reservedKeywordsList))
	for _, kw := range reservedKeywordsList {
		reservedKeywordsMap[kw] = true
	}
}

func init() {
	initKeywords()
}
