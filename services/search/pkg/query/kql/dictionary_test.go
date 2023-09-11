package kql_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	tAssert "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast/test"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var mustParseTime = func(t *testing.T, ts string) time.Time {
	tp, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		t.Fatalf("time.Parse(...) error = %v", err)
	}

	return tp
}

var mustJoin = func(v []string) string {
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
	tests := []struct {
		name          string
		skip          bool
		givenQuery    string
		expectedAst   *ast.Ast
		expectedError error
	}{
		// SPEC //////////////////////////////////////////////////////////////////////////////
		// https://msopenspecs.azureedge.net/files/MS-KQL/%5bMS-KQL%5d.pdf
		//
		// 2.1.2 AND Operator
		{
			name: `cat AND dog`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Value: "dog"},
				},
			},
		},
		{
			name:          `AND`,
			expectedError: errors.New(""),
		},
		{
			name:          `AND cat AND dog`,
			expectedError: errors.New(""),
		},
		// 3.1.11 Implicit Operator
		{
			name: `cat dog`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "dog"},
				},
			},
		},
		{
			name: `cat AND (dog OR fox)`,
			expectedAst: &ast.Ast{
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
			name: `cat (dog OR fox)`,
			expectedAst: &ast.Ast{
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
		// 3.1.12 Parentheses
		{
			name: `(cat OR dog) AND fox`,
			expectedAst: &ast.Ast{
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
			name: `author:"John Smith" filetype:docx`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Key: "filetype", Value: "docx"},
				},
			},
		},
		{
			name: `author:"John Smith" AND filetype:docx`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Key: "filetype", Value: "docx"},
				},
			},
		},
		{
			name: `author:"John Smith" author:"Jane Smith"`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Key: "author", Value: "Jane Smith"},
				},
			},
		},
		{
			name: `author:"John Smith" OR author:"Jane Smith"`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Key: "author", Value: "Jane Smith"},
				},
			},
		},
		{
			name: `cat filetype:docx`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Key: "filetype", Value: "docx"},
				},
			},
		},
		{
			name: `cat AND filetype:docx`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Key: "filetype", Value: "docx"},
				},
			},
		},
		// 3.3.1.1.1 Implicit AND Operator
		{
			name: `cat +dog`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Value: "dog"},
				},
			},
		},
		{
			name: `cat AND dog`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Value: "dog"},
				},
			},
		},
		{
			name: `cat -dog`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.OperatorNode{Value: kql.BoolNOT},
					&ast.StringNode{Value: "dog"},
				},
			},
		},
		{
			name: `cat AND NOT dog`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.OperatorNode{Value: kql.BoolNOT},
					&ast.StringNode{Value: "dog"},
				},
			},
		},
		{
			name: `cat +dog -fox`,
			expectedAst: &ast.Ast{
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
			name: `cat AND dog AND NOT fox`,
			expectedAst: &ast.Ast{
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
			name: `cat dog +fox`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "dog"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Value: "fox"},
				},
			},
		},
		{
			name: `fox OR (fox AND (cat OR dog))`,
			expectedAst: &ast.Ast{
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
			name: `cat dog -fox`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "cat"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "dog"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.OperatorNode{Value: kql.BoolNOT},
					&ast.StringNode{Value: "fox"},
				},
			},
		},
		{
			name: `(NOT fox) AND (cat OR dog)`,
			expectedAst: &ast.Ast{
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
			name: `cat +dog -fox`,
			expectedAst: &ast.Ast{
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
			name: `(NOT fox) AND (dog OR (dog AND cat))`,
			expectedAst: &ast.Ast{
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
		//////////////////////////////////////////////////////////////////////////////////////
		// everything else
		{
			name:       "FullDictionary",
			givenQuery: mustJoin(FullDictionary),
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{Value: "federated"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "search"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "federat*"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "search"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "search"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "fed*"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Key: "author", Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Key: "filetype", Value: "docx"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Key: "filename", Value: "budget.xlsx"},
					&ast.OperatorNode{Value: kql.BoolAND},
					&ast.StringNode{Value: "author"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "author"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "author"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "author"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "John Smith"},
					&ast.OperatorNode{Value: kql.BoolOR},
					&ast.StringNode{Value: "author"},
					&ast.OperatorNode{Value: kql.BoolOR},
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
							&ast.OperatorNode{Value: kql.BoolOR},
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
			name: "Group",
			givenQuery: mustJoin([]string{
				`(name:"moby di*" OR tag:bestseller) AND tag:book NOT tag:read`,
				`author:("John Smith" Jane)`,
				`author:("John Smith" OR Jane)`,
			}),
			expectedAst: &ast.Ast{
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
							&ast.OperatorNode{Value: kql.BoolOR},
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
			name: `author:("John Smith" Jane) author:"Jack" AND author:"Oggy"`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.GroupNode{
						Key: "author",
						Nodes: []ast.Node{
							&ast.StringNode{Value: "John Smith"},
							&ast.OperatorNode{Value: kql.BoolOR},
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
			name: `author:("John Smith" OR Jane)`,
			expectedAst: &ast.Ast{
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
			name: `NOT "John Smith" NOT Jane`,
			expectedAst: &ast.Ast{
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
			name: `NOT author:"John Smith" NOT author:"Jane Smith" NOT tag:sifi`,
			expectedAst: &ast.Ast{
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
			name: `scope:"<uuid>/new folder/subfolder" file`,
			expectedAst: &ast.Ast{
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
			name: `	😂 "*😀 😁*" name:😂💁👌🎍😍 name:😂💁👌 😍`,
			expectedAst: &ast.Ast{
				Nodes: []ast.Node{
					&ast.StringNode{
						Value: "😂",
					},
					&ast.OperatorNode{Value: kql.BoolOR},
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
			name: "DateTimeRestrictionNode",
			givenQuery: mustJoin([]string{
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
			expectedAst: &ast.Ast{
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
			name: "id",
			givenQuery: mustJoin([]string{
				`id:b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c`,
				`ID:b27d3bf1-b254-459f-92e8-bdba668d6d3f$d0648459-25fb-4ed8-8684-bc62c7dca29c!d0648459-25fb-4ed8-8684-bc62c7dca29c`,
			}),
			expectedAst: &ast.Ast{
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
	}

	assert := tAssert.New(t)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}

			q := tt.name

			if tt.givenQuery != "" {
				q = tt.givenQuery
			}

			parsedAST, err := kql.Builder{}.Build(q)

			if tt.expectedError != nil {
				if tt.expectedError.Error() != "" {
					assert.Equal(err, tt.expectedError)
				} else {
					assert.NotNil(err)
				}

				return
			}

			if diff := test.DiffAst(tt.expectedAst, parsedAST); diff != "" {
				t.Fatalf("AST mismatch \nquery: '%s' \n(-want +got): %s", q, diff)
			}
		})
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
