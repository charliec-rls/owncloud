package kql_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/now"
	tAssert "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast/test"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var mustParseTime = func(t *testing.T, ts string) time.Time {
	tp, err := now.Parse(ts)
	if err != nil {
		t.Fatalf("time.Parse(...) error = %v", err)
	}

	return tp
}

var setWorldClock = func(t *testing.T, ts string) func() func() {
	return func() func() {
		kql.PatchTimeNow(func() time.Time {
			return mustParseTime(t, ts)
		})

		return func() {
			kql.PatchTimeNow(time.Now)
		}
	}
}

var patchNow = func(t *testing.T, ts string) func() func() {
	return func() func() {
		kql.PatchTimeNow(func() time.Time {
			return mustParseTime(t, ts)
		})

		return func() {
			kql.PatchTimeNow(time.Now)
		}
	}
}

var join = func(v []string) string {
	return strings.Join(v, " ")
}

var FullDictionary = []string{
	`federated search`,
	`federat* search`,
	`search fed*`,
	`author:"John Smith"`,
	`filetype:docx`,
	`filename:budget.xlsx`,
	`author: "John Smith"`,
	`author :"John Smith"`,
	`author : "John Smith"`,
	`author "John Smith"`,
	`author "John Smith"`,
	`author:Shakespear`,
	`author:Paul`,
	`author:Shakesp*`,
	`title:"Advanced Search"`,
	`title:"Advanced Sear*"`,
	`title:"Advan* Search"`,
	`title:"*anced Search"`,
	`author:"John Smith" OR author:"Jane Smith"`,
	`author:"John Smith" AND filetype:docx`,
	`author:("John Smith" "Jane Smith")`,
	`author:("John Smith" OR "Jane Smith")`,
	`(DepartmentId:* OR RelatedHubSites:*) AND contentclass:sts_site NOT IsHubSite:false`,
	`author:"John Smith" (filetype:docx title:"Advanced Search")`,
}

func TestParse(t *testing.T) {
	type tc struct {
		name  string
		query string
		ast   *ast.Ast
		error error
		skip  bool
		patch func() func()
	}

	tests := []struct {
		name  string
		cases []tc
	}{
		{
			// SPEC //////////////////////////////////////////////////////////////////////////////
			//
			// https://msopenspecs.azureedge.net/files/MS-KQL/%5bMS-KQL%5d.pdf
			// https://learn.microsoft.com/en-us/openspecs/sharepoint_protocols/ms-kql/3bbf06cd-8fc1-4277-bd92-8661ccd3c9b0
			// https://learn.microsoft.com/en-us/sharepoint/dev/general-development/keyword-query-language-kql-syntax-reference
			name: `spec`,
			cases: []tc{
				// 2.1.2 AND Operator
				// 3.1.2 AND Operator
				{
					query: `cat AND dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `AND`,
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolAND},
					},
				},
				{
					query: `AND cat AND dog`,
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolAND},
					},
				},
				// 2.1.6 NOT Operator
				// 3.1.6 NOT Operator
				{
					query: `cat NOT dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `NOT dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				// 2.1.8 OR Operator
				// 3.1.8 OR Operator
				{
					query: `cat OR dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `OR`,
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolOR},
					},
				},
				{
					query: `OR cat AND dog`,
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolOR},
					},
				},
				// 3.1.11 Implicit Operator
				{
					query: `cat dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `cat AND (dog OR fox)`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.StringNode{Value: "dog"},
								&ast.OperatorNode{Value: kql.BoolOR},
								&ast.StringNode{Value: "fox"},
							}},
						},
					},
				},
				{
					query: `cat (dog OR fox)`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.StringNode{Value: "dog"},
								&ast.OperatorNode{Value: kql.BoolOR},
								&ast.StringNode{Value: "fox"},
							}},
						},
					},
				},
				// 2.1.12 Parentheses
				// 3.1.12 Parentheses
				{
					query: `(cat OR dog) AND fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.StringNode{Value: "cat"},
								&ast.OperatorNode{Value: kql.BoolOR},
								&ast.StringNode{Value: "dog"},
							}},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				// 3.2.3 Implicit Operator for Property Restriction
				{
					query: `author:"John Smith" filetype:docx`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "filetype", Value: "docx"},
						},
					},
				},
				{
					query: `author:"John Smith" AND filetype:docx`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "filetype", Value: "docx"},
						},
					},
				},
				{
					query: `author:"John Smith" author:"Jane Smith"`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "author", Value: "Jane Smith"},
						},
					},
				},
				{
					query: `author:"John Smith" OR author:"Jane Smith"`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "author", Value: "Jane Smith"},
						},
					},
				},
				{
					query: `cat filetype:docx`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "filetype", Value: "docx"},
						},
					},
				},
				{
					query: `cat AND filetype:docx`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "filetype", Value: "docx"},
						},
					},
				},
				// 3.3.1.1.1 Implicit AND Operator
				{
					query: `cat +dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `cat AND dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `cat -dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `cat AND NOT dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
				{
					query: `cat +dog -fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `cat AND dog AND NOT fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `cat dog +fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `fox OR (fox AND (cat OR dog))`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "fox"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.StringNode{Value: "fox"},
								&ast.OperatorNode{Value: kql.BoolAND},
								&ast.GroupNode{Nodes: []ast.Node{
									&ast.StringNode{Value: "cat"},
									&ast.OperatorNode{Value: kql.BoolOR},
									&ast.StringNode{Value: "dog"},
								}},
							}},
						},
					},
				},
				{
					query: `cat dog -fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `(NOT fox) AND (cat OR dog)`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.OperatorNode{Value: kql.BoolNOT},
								&ast.StringNode{Value: "fox"},
							}},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.StringNode{Value: "cat"},
								&ast.OperatorNode{Value: kql.BoolOR},
								&ast.StringNode{Value: "dog"},
							}},
						},
					},
				},
				{
					query: `cat +dog -fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `(NOT fox) AND (dog OR (dog AND cat))`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.OperatorNode{Value: kql.BoolNOT},
								&ast.StringNode{Value: "fox"},
							}},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{Nodes: []ast.Node{
								&ast.StringNode{Value: "dog"},
								&ast.OperatorNode{Value: kql.BoolOR},
								&ast.GroupNode{Nodes: []ast.Node{
									&ast.StringNode{Value: "dog"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Value: "cat"},
								}},
							}},
						},
					},
				},
				// 2.3.5 Date Tokens
				// 3.3.5 Date Tokens
				{
					query: `Modified:2023-09-05`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Modified",
								Operator: &ast.OperatorNode{Value: ":"},
								Value:    mustParseTime(t, "2023-09-05"),
							},
						},
					},
				},
				{
					query: `Modified:"2008-01-29"`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Modified",
								Operator: &ast.OperatorNode{Value: ":"},
								Value:    mustParseTime(t, "2008-01-29"),
							},
						},
					},
				},
				{
					query: `Modified:today`,
					patch: patchNow(t, "2023-09-10"),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Modified",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-09-10"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Modified",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-09-10 23:59:59.999999999"),
							},
						},
					},
				},
			},
		},
		{
			name: "DateTimeRestrictionNode",
			cases: []tc{
				{
					query: join([]string{
						`Mtime:"2023-09-05T08:42:11.23554+02:00"`,
						`Mtime:2023-09-05T08:42:11.23554+02:00`,
						`Mtime="2023-09-05T08:42:11.23554+02:00"`,
						`Mtime=2023-09-05T08:42:11.23554+02:00`,
						`Mtime<"2023-09-05T08:42:11.23554+02:00"`,
						`Mtime<2023-09-05T08:42:11.23554+02:00`,
						`Mtime<="2023-09-05T08:42:11.23554+02:00"`,
						`Mtime<=2023-09-05T08:42:11.23554+02:00`,
						`Mtime>"2023-09-05T08:42:11.23554+02:00"`,
						`Mtime>2023-09-05T08:42:11.23554+02:00`,
						`Mtime>="2023-09-05T08:42:11.23554+02:00"`,
						`Mtime>=2023-09-05T08:42:11.23554+02:00`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ":"},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ":"},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "="},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "="},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<"},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<"},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">"},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">"},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-09-05T08:42:11.23554+02:00"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - today",
					patch: setWorldClock(t, "2023-09-10"),
					query: join([]string{
						`Mtime:today`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-09-10"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-09-10 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - yesterday",
					patch: setWorldClock(t, "2023-09-10"),
					query: join([]string{
						`Mtime:yesterday`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-09-09"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-09-09 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - yesterday - the beginning of the month",
					patch: setWorldClock(t, "2023-09-01"),
					query: join([]string{
						`Mtime:yesterday`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-08-31"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-08-31 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - this week",
					patch: setWorldClock(t, "2023-09-06"),
					query: join([]string{
						`Mtime:"this week"`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-09-04"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-09-10 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - this month",
					patch: setWorldClock(t, "2023-09-02"),
					query: join([]string{
						`Mtime:"this month"`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-09-01"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-09-30 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - last month",
					patch: setWorldClock(t, "2023-09-02"),
					query: join([]string{
						`Mtime:"last month"`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-08-01"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-08-31 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - last month - the beginning of the year",
					patch: setWorldClock(t, "2023-01-01"),
					query: join([]string{
						`Mtime:"last month"`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2022-12-01"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2022-12-31 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - this year",
					patch: setWorldClock(t, "2023-06-18"),
					query: join([]string{
						`Mtime:"this year"`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2023-01-01"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2023-12-31 23:59:59.999999999"),
							},
						},
					},
				},
				{
					name:  "NaturalLanguage DateTimeNode - last year",
					patch: setWorldClock(t, "2023-01-01"),
					query: join([]string{
						`Mtime:"last year"`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: ">="},
								Value:    mustParseTime(t, "2022-01-01"),
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.DateTimeNode{
								Key:      "Mtime",
								Operator: &ast.OperatorNode{Value: "<="},
								Value:    mustParseTime(t, "2022-12-31 23:59:59.999999999"),
							},
						},
					},
				},
			},
		},
		{
			name: "errors",
			cases: []tc{
				{
					query: "animal:(mammal:cat mammal:dog reptile:turtle)",
					error: kql.NamedGroupInvalidNodesError{
						Node: &ast.StringNode{Key: "mammal", Value: "cat"},
					},
				},
				{
					query: "animal:(cat mammal:dog turtle)",
					error: kql.NamedGroupInvalidNodesError{
						Node: &ast.StringNode{Key: "mammal", Value: "dog"},
					},
				},
				{
					query: "animal:(AND cat)",
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolAND},
					},
				},
				{
					query: "animal:(OR cat)",
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolOR},
					},
				},
				{
					query: "(AND cat)",
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolAND},
					},
				},
				{
					query: "(OR cat)",
					error: kql.StartsWithBinaryOperatorError{
						Node: &ast.OperatorNode{Value: kql.BoolOR},
					},
				},
				{
					query: `cat dog`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
						},
					},
				},
			},
		},
		{
			name: "Random",
			cases: []tc{
				{
					name:  "FullDictionary",
					query: join(FullDictionary),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "federated"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "search"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "federat*"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "search"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "search"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "fed*"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "filetype", Value: "docx"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "filename", Value: "budget.xlsx"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "author"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "author"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "author"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "author"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "author"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "author", Value: "Shakespear"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "author", Value: "Paul"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "author", Value: "Shakesp*"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "title", Value: "Advanced Search"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "title", Value: "Advanced Sear*"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "title", Value: "Advan* Search"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "title", Value: "*anced Search"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "author", Value: "Jane Smith"},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "filetype", Value: "docx"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{
								Key: "author",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "John Smith"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Value: "Jane Smith"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.GroupNode{
								Key: "author",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "John Smith"},
									&ast.OperatorNode{Value: kql.BoolOR},
									&ast.StringNode{Value: "Jane Smith"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{
								Nodes: []ast.Node{
									&ast.StringNode{Key: "DepartmentId", Value: "*"},
									&ast.OperatorNode{Value: kql.BoolOR},
									&ast.StringNode{Key: "RelatedHubSites", Value: "*"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "contentclass", Value: "sts_site"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.BooleanNode{Key: "IsHubSite", Value: false},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{
								Nodes: []ast.Node{
									&ast.StringNode{Key: "filetype", Value: "docx"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Key: "title", Value: "Advanced Search"},
								},
							},
						},
					},
				},
				{
					query: join([]string{
						`(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`,
						`author:("John Smith" Jane)`,
						`author:("John Smith" OR Jane)`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Nodes: []ast.Node{
									&ast.StringNode{Key: "name", Value: "moby di*"},
									&ast.OperatorNode{Value: kql.BoolOR},
									&ast.StringNode{Key: "tag", Value: "bestseller"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "tag", Value: "book"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Key: "tag", Value: "read"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.GroupNode{
								Key: "author",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "John Smith"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Value: "Jane"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.GroupNode{
								Key: "author",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "John Smith"},
									&ast.OperatorNode{Value: kql.BoolOR},
									&ast.StringNode{Value: "Jane"},
								},
							},
						},
					},
				},
				{
					query: `author:("John Smith" Jane) author:"Jack" AND author:"Oggy"`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Key: "author",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "John Smith"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Value: "Jane"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "author", Value: "Jack"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Key: "author", Value: "Oggy"},
						},
					},
				},
				{
					query: `author:("John Smith" OR Jane)`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Key: "author",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "John Smith"},
									&ast.OperatorNode{Value: kql.BoolOR},
									&ast.StringNode{Value: "Jane"},
								},
							},
						},
					},
				},
				{
					query: `NOT "John Smith" NOT Jane`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Value: "Jane"},
						},
					},
				},
				{
					query: `NOT author:"John Smith" NOT author:"Jane Smith" NOT tag:sifi`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Key: "author", Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Key: "author", Value: "Jane Smith"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.OperatorNode{Value: kql.BoolNOT},
							&ast.StringNode{Key: "tag", Value: "sifi"},
						},
					},
				},
				{
					query: `scope:"<uuid>/new folder/subfolder" file`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{
								Key:   "scope",
								Value: "<uuid>/new folder/subfolder",
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{
								Value: "file",
							},
						},
					},
				},
				{
					query: `	😂 "*😀 😁*" name:😂💁👌🎍😍 name:😂💁👌 😍`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{
								Value: "😂",
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{
								Value: "*😀 😁*",
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{
								Key:   "name",
								Value: "😂💁👌🎍😍",
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{
								Key:   "name",
								Value: "😂💁👌",
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{
								Value: "😍",
							},
						},
					},
				},
				{
					query: "animal:(cat dog turtle)",
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Key: "animal",
								Nodes: []ast.Node{
									&ast.StringNode{
										Value: "cat",
									},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{
										Value: "dog",
									},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{
										Value: "turtle",
									},
								},
							},
						},
					},
				},
				{
					query: "(cat dog turtle)",
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Nodes: []ast.Node{
									&ast.StringNode{
										Value: "cat",
									},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{
										Value: "dog",
									},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{
										Value: "turtle",
									},
								},
							},
						},
					},
				},
				{
					query: `cat dog fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{Value: "cat"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "dog"},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `(cat dog) fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Nodes: []ast.Node{
									&ast.StringNode{Value: "cat"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Value: "dog"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `(mammal:cat mammal:dog) fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Nodes: []ast.Node{
									&ast.StringNode{Key: "mammal", Value: "cat"},
									&ast.OperatorNode{Value: kql.BoolOR},
									&ast.StringNode{Key: "mammal", Value: "dog"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `mammal:(cat dog) fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Key: "mammal",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "cat"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Value: "dog"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolAND},
							&ast.StringNode{Value: "fox"},
						},
					},
				},
				{
					query: `mammal:(cat dog) mammal:fox`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Key: "mammal",
								Nodes: []ast.Node{
									&ast.StringNode{Value: "cat"},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.StringNode{Value: "dog"},
								},
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{Key: "mammal", Value: "fox"},
						},
					},
				},
				{
					query: `title:((Advanced OR Search OR Query) -"Advanced Search Query")`,
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.GroupNode{
								Key: "title",
								Nodes: []ast.Node{
									&ast.GroupNode{
										Nodes: []ast.Node{
											&ast.StringNode{Value: "Advanced"},
											&ast.OperatorNode{Value: kql.BoolOR},
											&ast.StringNode{Value: "Search"},
											&ast.OperatorNode{Value: kql.BoolOR},
											&ast.StringNode{Value: "Query"},
										},
									},
									&ast.OperatorNode{Value: kql.BoolAND},
									&ast.OperatorNode{Value: kql.BoolNOT},
									&ast.StringNode{Value: "Advanced Search Query"},
								},
							},
						},
					},
				},
				{
					query: join([]string{
						`id:b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c`,
						`ID:b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c`,
					}),
					ast: &ast.Ast{
						Nodes: []ast.Node{
							&ast.StringNode{
								Key:   "id",
								Value: "b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c",
							},
							&ast.OperatorNode{Value: kql.BoolOR},
							&ast.StringNode{
								Key:   "ID",
								Value: "b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c",
							},
						},
					},
				},
			},
		},
	}

	for _, currentTest := range tests {
		currentTest := currentTest

		t.Run(
			currentTest.name,
			func(t *testing.T) {
				for i, testCase := range currentTest.cases {
					currentTest := currentTest
					subTestName := strconv.Itoa(i)

					if testCase.query != "" {
						subTestName = subTestName + "  -  " + testCase.query
					}

					if testCase.name != "" {
						subTestName = testCase.name
					}

					t.Run(subTestName, func(t *testing.T) {
						if testCase.skip {
							t.Skip()
						}

						if testCase.patch != nil {
							revert := testCase.patch()
							defer revert()
						}

						query := currentTest.name
						if testCase.query != "" {
							query = testCase.query
						}

						astResult, err := kql.Builder{}.Build(query)
						assert := tAssert.New(t)

						if testCase.error != nil {
							if expectedError := testCase.error.Error(); expectedError != "" {
								assert.Equal(err.Error(), expectedError)
							} else {
								assert.NotNil(err)
							}

							return
						}

						if diff := test.DiffAst(testCase.ast, astResult); diff != "" {
							t.Fatalf("AST mismatch \nquery: '%s' \n(-expected +got): %s", query, diff)
						}
					})
				}
			},
		)
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := kql.Parse("", []byte(strings.Join(FullDictionary, " "))); err != nil {
			b.Fatal(err)
		}
	}
}
