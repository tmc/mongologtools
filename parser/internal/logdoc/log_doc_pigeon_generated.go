package logdoc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ifs(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

type docElem struct {
	Field string
	Value interface{}
}

var g = &grammar{
	rules: []*rule{
		{
			name: "LogDoc",
			pos:  position{line: 18, col: 1, offset: 194},
			expr: &actionExpr{
				pos: position{line: 18, col: 11, offset: 204},
				run: (*parser).callonLogDoc1,
				expr: &seqExpr{
					pos: position{line: 18, col: 11, offset: 204},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 18, col: 11, offset: 204},
							label: "doc",
							expr: &ruleRefExpr{
								pos:  position{line: 18, col: 15, offset: 208},
								name: "Doc",
							},
						},
						&notExpr{
							pos: position{line: 18, col: 19, offset: 212},
							expr: &anyMatcher{
								line: 18, col: 20, offset: 213,
							},
						},
					},
				},
			},
		},
		{
			name: "Doc",
			pos:  position{line: 22, col: 1, offset: 262},
			expr: &actionExpr{
				pos: position{line: 22, col: 8, offset: 269},
				run: (*parser).callonDoc1,
				expr: &seqExpr{
					pos: position{line: 22, col: 8, offset: 269},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 22, col: 8, offset: 269},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 22, col: 12, offset: 273},
							label: "fields",
							expr: &zeroOrOneExpr{
								pos: position{line: 22, col: 19, offset: 280},
								expr: &seqExpr{
									pos: position{line: 22, col: 20, offset: 281},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 22, col: 20, offset: 281},
											name: "DocElem",
										},
										&zeroOrMoreExpr{
											pos: position{line: 22, col: 28, offset: 289},
											expr: &seqExpr{
												pos: position{line: 22, col: 29, offset: 290},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 22, col: 29, offset: 290},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 22, col: 33, offset: 294},
														name: "DocElem",
													},
												},
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 22, col: 46, offset: 307},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "DocElem",
			pos:  position{line: 38, col: 1, offset: 629},
			expr: &actionExpr{
				pos: position{line: 38, col: 12, offset: 640},
				run: (*parser).callonDocElem1,
				expr: &seqExpr{
					pos: position{line: 38, col: 12, offset: 640},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 38, col: 12, offset: 640},
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 12, offset: 640},
								name: "S",
							},
						},
						&labeledExpr{
							pos:   position{line: 38, col: 15, offset: 643},
							label: "field",
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 21, offset: 649},
								name: "Field",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 38, col: 27, offset: 655},
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 27, offset: 655},
								name: "S",
							},
						},
						&litMatcher{
							pos:        position{line: 38, col: 30, offset: 658},
							val:        ":",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 38, col: 34, offset: 662},
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 34, offset: 662},
								name: "S",
							},
						},
						&labeledExpr{
							pos:   position{line: 38, col: 37, offset: 665},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 43, offset: 671},
								name: "Value",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 38, col: 49, offset: 677},
							expr: &ruleRefExpr{
								pos:  position{line: 38, col: 49, offset: 677},
								name: "S",
							},
						},
					},
				},
			},
		},
		{
			name: "List",
			pos:  position{line: 42, col: 1, offset: 743},
			expr: &actionExpr{
				pos: position{line: 42, col: 9, offset: 751},
				run: (*parser).callonList1,
				expr: &seqExpr{
					pos: position{line: 42, col: 9, offset: 751},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 42, col: 9, offset: 751},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 42, col: 13, offset: 755},
							label: "vals",
							expr: &zeroOrOneExpr{
								pos: position{line: 42, col: 18, offset: 760},
								expr: &seqExpr{
									pos: position{line: 42, col: 19, offset: 761},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 42, col: 19, offset: 761},
											name: "ListElem",
										},
										&zeroOrMoreExpr{
											pos: position{line: 42, col: 28, offset: 770},
											expr: &seqExpr{
												pos: position{line: 42, col: 29, offset: 771},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 42, col: 29, offset: 771},
														val:        ",",
														ignoreCase: false,
													},
													&ruleRefExpr{
														pos:  position{line: 42, col: 33, offset: 775},
														name: "ListElem",
													},
												},
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 42, col: 47, offset: 789},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ListElem",
			pos:  position{line: 56, col: 1, offset: 1047},
			expr: &actionExpr{
				pos: position{line: 56, col: 13, offset: 1059},
				run: (*parser).callonListElem1,
				expr: &seqExpr{
					pos: position{line: 56, col: 13, offset: 1059},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 56, col: 13, offset: 1059},
							expr: &ruleRefExpr{
								pos:  position{line: 56, col: 13, offset: 1059},
								name: "S",
							},
						},
						&labeledExpr{
							pos:   position{line: 56, col: 16, offset: 1062},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 56, col: 20, offset: 1066},
								name: "Value",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 56, col: 26, offset: 1072},
							expr: &ruleRefExpr{
								pos:  position{line: 56, col: 26, offset: 1072},
								name: "S",
							},
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 60, col: 1, offset: 1097},
			expr: &actionExpr{
				pos: position{line: 60, col: 10, offset: 1106},
				run: (*parser).callonField1,
				expr: &labeledExpr{
					pos:   position{line: 60, col: 10, offset: 1106},
					label: "fieldName",
					expr: &oneOrMoreExpr{
						pos: position{line: 60, col: 20, offset: 1116},
						expr: &ruleRefExpr{
							pos:  position{line: 60, col: 20, offset: 1116},
							name: "fieldChar",
						},
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 64, col: 1, offset: 1160},
			expr: &choiceExpr{
				pos: position{line: 64, col: 11, offset: 1170},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 64, col: 11, offset: 1170},
						name: "Doc",
					},
					&ruleRefExpr{
						pos:  position{line: 65, col: 11, offset: 1184},
						name: "List",
					},
					&ruleRefExpr{
						pos:  position{line: 66, col: 11, offset: 1199},
						name: "Numeric",
					},
					&ruleRefExpr{
						pos:  position{line: 67, col: 11, offset: 1217},
						name: "Boolean",
					},
					&ruleRefExpr{
						pos:  position{line: 68, col: 11, offset: 1235},
						name: "String",
					},
					&ruleRefExpr{
						pos:  position{line: 69, col: 11, offset: 1252},
						name: "Null",
					},
					&ruleRefExpr{
						pos:  position{line: 70, col: 11, offset: 1267},
						name: "ObjectID",
					},
					&ruleRefExpr{
						pos:  position{line: 71, col: 11, offset: 1286},
						name: "Date",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 11, offset: 1301},
						name: "BinData",
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 11, offset: 1319},
						name: "TimestampVal",
					},
					&ruleRefExpr{
						pos:  position{line: 74, col: 11, offset: 1342},
						name: "Regex",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 11, offset: 1358},
						name: "NumberLong",
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 11, offset: 1379},
						name: "Undefined",
					},
					&ruleRefExpr{
						pos:  position{line: 77, col: 11, offset: 1399},
						name: "MinKey",
					},
					&ruleRefExpr{
						pos:  position{line: 78, col: 11, offset: 1416},
						name: "MaxKey",
					},
				},
			},
		},
		{
			name: "Numeric",
			pos:  position{line: 81, col: 1, offset: 1434},
			expr: &actionExpr{
				pos: position{line: 81, col: 12, offset: 1445},
				run: (*parser).callonNumeric1,
				expr: &seqExpr{
					pos: position{line: 81, col: 12, offset: 1445},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 81, col: 12, offset: 1445},
							expr: &litMatcher{
								pos:        position{line: 81, col: 12, offset: 1445},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 81, col: 17, offset: 1450},
							expr: &charClassMatcher{
								pos:        position{line: 81, col: 17, offset: 1450},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 81, col: 24, offset: 1457},
							expr: &litMatcher{
								pos:        position{line: 81, col: 24, offset: 1457},
								val:        ".",
								ignoreCase: false,
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 81, col: 29, offset: 1462},
							expr: &charClassMatcher{
								pos:        position{line: 81, col: 29, offset: 1462},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "Boolean",
			pos:  position{line: 84, col: 1, offset: 1522},
			expr: &choiceExpr{
				pos: position{line: 84, col: 12, offset: 1533},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 84, col: 12, offset: 1533},
						run: (*parser).callonBoolean2,
						expr: &litMatcher{
							pos:        position{line: 84, col: 12, offset: 1533},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 84, col: 42, offset: 1563},
						run: (*parser).callonBoolean4,
						expr: &litMatcher{
							pos:        position{line: 84, col: 42, offset: 1563},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 85, col: 1, offset: 1593},
			expr: &actionExpr{
				pos: position{line: 85, col: 12, offset: 1604},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 85, col: 12, offset: 1604},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 85, col: 12, offset: 1604},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 85, col: 16, offset: 1608},
							label: "str",
							expr: &ruleRefExpr{
								pos:  position{line: 85, col: 20, offset: 1612},
								name: "stringChars",
							},
						},
						&litMatcher{
							pos:        position{line: 85, col: 32, offset: 1624},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "stringChars",
			pos:  position{line: 86, col: 1, offset: 1657},
			expr: &actionExpr{
				pos: position{line: 86, col: 16, offset: 1672},
				run: (*parser).callonstringChars1,
				expr: &oneOrMoreExpr{
					pos: position{line: 86, col: 16, offset: 1672},
					expr: &ruleRefExpr{
						pos:  position{line: 86, col: 16, offset: 1672},
						name: "stringChar",
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 87, col: 1, offset: 1715},
			expr: &actionExpr{
				pos: position{line: 87, col: 9, offset: 1723},
				run: (*parser).callonNull1,
				expr: &litMatcher{
					pos:        position{line: 87, col: 9, offset: 1723},
					val:        "null",
					ignoreCase: false,
				},
			},
		},
		{
			name: "Date",
			pos:  position{line: 88, col: 1, offset: 1750},
			expr: &actionExpr{
				pos: position{line: 88, col: 9, offset: 1758},
				run: (*parser).callonDate1,
				expr: &seqExpr{
					pos: position{line: 88, col: 9, offset: 1758},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 88, col: 9, offset: 1758},
							expr: &litMatcher{
								pos:        position{line: 88, col: 9, offset: 1758},
								val:        "new ",
								ignoreCase: false,
							},
						},
						&litMatcher{
							pos:        position{line: 88, col: 17, offset: 1766},
							val:        "Date(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 88, col: 25, offset: 1774},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 88, col: 27, offset: 1776},
								name: "number",
							},
						},
						&litMatcher{
							pos:        position{line: 88, col: 34, offset: 1783},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "number",
			pos:  position{line: 91, col: 1, offset: 1833},
			expr: &actionExpr{
				pos: position{line: 91, col: 11, offset: 1843},
				run: (*parser).callonnumber1,
				expr: &seqExpr{
					pos: position{line: 91, col: 11, offset: 1843},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 91, col: 11, offset: 1843},
							expr: &litMatcher{
								pos:        position{line: 91, col: 11, offset: 1843},
								val:        "-",
								ignoreCase: false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 91, col: 16, offset: 1848},
							expr: &charClassMatcher{
								pos:        position{line: 91, col: 16, offset: 1848},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "ObjectID",
			pos:  position{line: 92, col: 1, offset: 1886},
			expr: &actionExpr{
				pos: position{line: 92, col: 13, offset: 1898},
				run: (*parser).callonObjectID1,
				expr: &seqExpr{
					pos: position{line: 92, col: 13, offset: 1898},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 92, col: 13, offset: 1898},
							val:        "ObjectId(",
							ignoreCase: false,
						},
						&charClassMatcher{
							pos:        position{line: 92, col: 25, offset: 1910},
							val:        "['\"]",
							chars:      []rune{'\'', '"'},
							ignoreCase: false,
							inverted:   false,
						},
						&labeledExpr{
							pos:   position{line: 92, col: 30, offset: 1915},
							label: "hex",
							expr: &ruleRefExpr{
								pos:  position{line: 92, col: 34, offset: 1919},
								name: "hexChars",
							},
						},
						&charClassMatcher{
							pos:        position{line: 92, col: 43, offset: 1928},
							val:        "['\"]",
							chars:      []rune{'\'', '"'},
							ignoreCase: false,
							inverted:   false,
						},
						&litMatcher{
							pos:        position{line: 92, col: 48, offset: 1933},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "hexChars",
			pos:  position{line: 95, col: 1, offset: 1989},
			expr: &actionExpr{
				pos: position{line: 95, col: 13, offset: 2001},
				run: (*parser).callonhexChars1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 95, col: 13, offset: 2001},
					expr: &ruleRefExpr{
						pos:  position{line: 95, col: 13, offset: 2001},
						name: "hexChar",
					},
				},
			},
		},
		{
			name: "BinData",
			pos:  position{line: 98, col: 1, offset: 2042},
			expr: &actionExpr{
				pos: position{line: 98, col: 12, offset: 2053},
				run: (*parser).callonBinData1,
				expr: &seqExpr{
					pos: position{line: 98, col: 12, offset: 2053},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 98, col: 12, offset: 2053},
							val:        "BinData(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 98, col: 23, offset: 2064},
							label: "bd",
							expr: &ruleRefExpr{
								pos:  position{line: 98, col: 26, offset: 2067},
								name: "binData",
							},
						},
						&litMatcher{
							pos:        position{line: 98, col: 34, offset: 2075},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "binData",
			pos:  position{line: 101, col: 1, offset: 2137},
			expr: &actionExpr{
				pos: position{line: 101, col: 12, offset: 2148},
				run: (*parser).callonbinData1,
				expr: &oneOrMoreExpr{
					pos: position{line: 101, col: 12, offset: 2148},
					expr: &charClassMatcher{
						pos:        position{line: 101, col: 12, offset: 2148},
						val:        "[^)]",
						chars:      []rune{')'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "Regex",
			pos:  position{line: 104, col: 1, offset: 2178},
			expr: &actionExpr{
				pos: position{line: 104, col: 10, offset: 2187},
				run: (*parser).callonRegex1,
				expr: &seqExpr{
					pos: position{line: 104, col: 10, offset: 2187},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 104, col: 10, offset: 2187},
							val:        "/",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 104, col: 14, offset: 2191},
							label: "rebody",
							expr: &ruleRefExpr{
								pos:  position{line: 104, col: 21, offset: 2198},
								name: "regexBody",
							},
						},
					},
				},
			},
		},
		{
			name: "TimestampVal",
			pos:  position{line: 107, col: 1, offset: 2260},
			expr: &actionExpr{
				pos: position{line: 107, col: 18, offset: 2277},
				run: (*parser).callonTimestampVal1,
				expr: &labeledExpr{
					pos:   position{line: 107, col: 18, offset: 2277},
					label: "ts",
					expr: &choiceExpr{
						pos: position{line: 107, col: 22, offset: 2281},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 107, col: 22, offset: 2281},
								name: "timestampParen",
							},
							&ruleRefExpr{
								pos:  position{line: 107, col: 39, offset: 2298},
								name: "timestampPipe",
							},
						},
					},
				},
			},
		},
		{
			name: "timestampParen",
			pos:  position{line: 110, col: 1, offset: 2333},
			expr: &actionExpr{
				pos: position{line: 110, col: 19, offset: 2351},
				run: (*parser).callontimestampParen1,
				expr: &seqExpr{
					pos: position{line: 110, col: 19, offset: 2351},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 110, col: 19, offset: 2351},
							val:        "Timestamp(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 110, col: 32, offset: 2364},
							label: "ts",
							expr: &ruleRefExpr{
								pos:  position{line: 110, col: 35, offset: 2367},
								name: "charsTillParen",
							},
						},
						&litMatcher{
							pos:        position{line: 110, col: 50, offset: 2382},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "timestampPipe",
			pos:  position{line: 113, col: 1, offset: 2438},
			expr: &actionExpr{
				pos: position{line: 113, col: 18, offset: 2455},
				run: (*parser).callontimestampPipe1,
				expr: &seqExpr{
					pos: position{line: 113, col: 18, offset: 2455},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 113, col: 18, offset: 2455},
							val:        "Timestamp ",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 113, col: 31, offset: 2468},
							label: "ts",
							expr: &ruleRefExpr{
								pos:  position{line: 113, col: 34, offset: 2471},
								name: "timestampPipeChars",
							},
						},
					},
				},
			},
		},
		{
			name: "timestampPipeChars",
			pos:  position{line: 116, col: 1, offset: 2542},
			expr: &actionExpr{
				pos: position{line: 116, col: 23, offset: 2564},
				run: (*parser).callontimestampPipeChars1,
				expr: &oneOrMoreExpr{
					pos: position{line: 116, col: 23, offset: 2564},
					expr: &choiceExpr{
						pos: position{line: 116, col: 24, offset: 2565},
						alternatives: []interface{}{
							&charClassMatcher{
								pos:        position{line: 116, col: 24, offset: 2565},
								val:        "[0-9]",
								ranges:     []rune{'0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
							&litMatcher{
								pos:        position{line: 116, col: 32, offset: 2573},
								val:        "|",
								ignoreCase: false,
							},
						},
					},
				},
			},
		},
		{
			name: "NumberLong",
			pos:  position{line: 119, col: 1, offset: 2611},
			expr: &actionExpr{
				pos: position{line: 119, col: 15, offset: 2625},
				run: (*parser).callonNumberLong1,
				expr: &seqExpr{
					pos: position{line: 119, col: 15, offset: 2625},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 119, col: 15, offset: 2625},
							val:        "NumberLong(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 119, col: 29, offset: 2639},
							label: "n",
							expr: &ruleRefExpr{
								pos:  position{line: 119, col: 31, offset: 2641},
								name: "charsTillParen",
							},
						},
						&litMatcher{
							pos:        position{line: 119, col: 46, offset: 2656},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "MinKey",
			pos:  position{line: 123, col: 1, offset: 2713},
			expr: &actionExpr{
				pos: position{line: 123, col: 11, offset: 2723},
				run: (*parser).callonMinKey1,
				expr: &litMatcher{
					pos:        position{line: 123, col: 11, offset: 2723},
					val:        "MinKey",
					ignoreCase: false,
				},
			},
		},
		{
			name: "MaxKey",
			pos:  position{line: 124, col: 1, offset: 2787},
			expr: &actionExpr{
				pos: position{line: 124, col: 11, offset: 2797},
				run: (*parser).callonMaxKey1,
				expr: &litMatcher{
					pos:        position{line: 124, col: 11, offset: 2797},
					val:        "MaxKey",
					ignoreCase: false,
				},
			},
		},
		{
			name: "Undefined",
			pos:  position{line: 125, col: 1, offset: 2878},
			expr: &actionExpr{
				pos: position{line: 125, col: 14, offset: 2891},
				run: (*parser).callonUndefined1,
				expr: &litMatcher{
					pos:        position{line: 125, col: 14, offset: 2891},
					val:        "undefined",
					ignoreCase: false,
				},
			},
		},
		{
			name: "hexChar",
			pos:  position{line: 127, col: 1, offset: 2956},
			expr: &choiceExpr{
				pos: position{line: 127, col: 12, offset: 2967},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 127, col: 12, offset: 2967},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 127, col: 20, offset: 2975},
						val:        "[a-fA-F]",
						ranges:     []rune{'a', 'f', 'A', 'F'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "regexChar",
			pos:  position{line: 128, col: 1, offset: 2984},
			expr: &charClassMatcher{
				pos:        position{line: 128, col: 14, offset: 2997},
				val:        "[^/]",
				chars:      []rune{'/'},
				ignoreCase: false,
				inverted:   true,
			},
		},
		{
			name: "regexBody",
			pos:  position{line: 129, col: 1, offset: 3002},
			expr: &actionExpr{
				pos: position{line: 129, col: 14, offset: 3015},
				run: (*parser).callonregexBody1,
				expr: &seqExpr{
					pos: position{line: 129, col: 14, offset: 3015},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 129, col: 14, offset: 3015},
							expr: &ruleRefExpr{
								pos:  position{line: 129, col: 14, offset: 3015},
								name: "regexChar",
							},
						},
						&litMatcher{
							pos:        position{line: 129, col: 25, offset: 3026},
							val:        "/",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 129, col: 29, offset: 3030},
							expr: &charClassMatcher{
								pos:        position{line: 129, col: 29, offset: 3030},
								val:        "[gims]",
								chars:      []rune{'g', 'i', 'm', 's'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "charsTillParen",
			pos:  position{line: 133, col: 1, offset: 3071},
			expr: &actionExpr{
				pos: position{line: 133, col: 19, offset: 3089},
				run: (*parser).calloncharsTillParen1,
				expr: &oneOrMoreExpr{
					pos: position{line: 133, col: 19, offset: 3089},
					expr: &charClassMatcher{
						pos:        position{line: 133, col: 19, offset: 3089},
						val:        "[^)]",
						chars:      []rune{')'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "stringChar",
			pos:  position{line: 136, col: 1, offset: 3127},
			expr: &choiceExpr{
				pos: position{line: 136, col: 15, offset: 3141},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 136, col: 15, offset: 3141},
						val:        "[^\"\\\\]",
						chars:      []rune{'"', '\\'},
						ignoreCase: false,
						inverted:   true,
					},
					&seqExpr{
						pos: position{line: 136, col: 24, offset: 3150},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 136, col: 24, offset: 3150},
								val:        "\\",
								ignoreCase: false,
							},
							&charClassMatcher{
								pos:        position{line: 136, col: 29, offset: 3155},
								val:        "[\"\\\\]",
								chars:      []rune{'"', '\\'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "fieldChar",
			pos:  position{line: 137, col: 1, offset: 3161},
			expr: &choiceExpr{
				pos: position{line: 137, col: 14, offset: 3174},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 137, col: 14, offset: 3174},
						val:        "[a-zA-Z]",
						ranges:     []rune{'a', 'z', 'A', 'Z'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 137, col: 25, offset: 3185},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
					&charClassMatcher{
						pos:        position{line: 137, col: 33, offset: 3193},
						val:        "[_$.*]",
						chars:      []rune{'_', '$', '.', '*'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "S",
			pos:  position{line: 139, col: 1, offset: 3201},
			expr: &litMatcher{
				pos:        position{line: 139, col: 6, offset: 3206},
				val:        " ",
				ignoreCase: false,
			},
		},
	},
}

func (c *current) onLogDoc1(doc interface{}) (interface{}, error) {
	return doc.(map[string]interface{}), nil
}

func (p *parser) callonLogDoc1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLogDoc1(stack["doc"])
}

func (c *current) onDoc1(fields interface{}) (interface{}, error) {
	result := map[string]interface{}{}
	fieldsSl := ifs(fields)
	if len(fieldsSl) == 0 {
		return result, nil
	}
	de := fieldsSl[0].(docElem)
	result[de.Field] = de.Value
	restSl := ifs(fieldsSl[1])
	for _, field := range restSl {
		de := ifs(field)[1].(docElem)
		result[de.Field] = de.Value
	}
	return result, nil
}

func (p *parser) callonDoc1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoc1(stack["fields"])
}

func (c *current) onDocElem1(field, value interface{}) (interface{}, error) {
	return docElem{Field: field.(string), Value: value}, nil
}

func (p *parser) callonDocElem1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDocElem1(stack["field"], stack["value"])
}

func (c *current) onList1(vals interface{}) (interface{}, error) {
	result := []interface{}{}
	valsSl := ifs(vals)
	if len(valsSl) == 0 {
		return result, nil
	}
	result = append(result, valsSl[0])
	restSl := valsSl[1]
	for _, val := range ifs(restSl) {
		result = append(result, ifs(val)[1])
	}
	return result, nil
}

func (p *parser) callonList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onList1(stack["vals"])
}

func (c *current) onListElem1(val interface{}) (interface{}, error) {
	return val, nil
}

func (p *parser) callonListElem1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListElem1(stack["val"])
}

func (c *current) onField1(fieldName interface{}) (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["fieldName"])
}

func (c *current) onNumeric1() (interface{}, error) {
	return new(LogDoc).Numeric(string(c.text)), nil
}

func (p *parser) callonNumeric1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumeric1()
}

func (c *current) onBoolean2() (interface{}, error) {
	return true, nil
}

func (p *parser) callonBoolean2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolean2()
}

func (c *current) onBoolean4() (interface{}, error) {
	return false, nil
}

func (p *parser) callonBoolean4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolean4()
}

func (c *current) onString1(str interface{}) (interface{}, error) {
	return str.(string), nil
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1(stack["str"])
}

func (c *current) onstringChars1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonstringChars1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onstringChars1()
}

func (c *current) onNull1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNull1()
}

func (c *current) onDate1(n interface{}) (interface{}, error) {
	return new(LogDoc).Date(n.(string)), nil
}

func (p *parser) callonDate1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDate1(stack["n"])
}

func (c *current) onnumber1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonnumber1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onnumber1()
}

func (c *current) onObjectID1(hex interface{}) (interface{}, error) {
	return new(LogDoc).ObjectId(hex.(string)), nil
}

func (p *parser) callonObjectID1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onObjectID1(stack["hex"])
}

func (c *current) onhexChars1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonhexChars1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onhexChars1()
}

func (c *current) onBinData1(bd interface{}) (interface{}, error) {
	return new(LogDoc).Bindata(string(bd.([]byte))), nil
}

func (p *parser) callonBinData1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinData1(stack["bd"])
}

func (c *current) onbinData1() (interface{}, error) {
	return c.text, nil
}

func (p *parser) callonbinData1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onbinData1()
}

func (c *current) onRegex1(rebody interface{}) (interface{}, error) {
	return new(LogDoc).Regex(rebody.(string)), nil
}

func (p *parser) callonRegex1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRegex1(stack["rebody"])
}

func (c *current) onTimestampVal1(ts interface{}) (interface{}, error) {
	return ts, nil
}

func (p *parser) callonTimestampVal1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimestampVal1(stack["ts"])
}

func (c *current) ontimestampParen1(ts interface{}) (interface{}, error) {
	return new(LogDoc).Timestamp(ts.(string)), nil
}

func (p *parser) callontimestampParen1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ontimestampParen1(stack["ts"])
}

func (c *current) ontimestampPipe1(ts interface{}) (interface{}, error) {
	return new(LogDoc).Timestamp(ts.(string)), nil
}

func (p *parser) callontimestampPipe1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ontimestampPipe1(stack["ts"])
}

func (c *current) ontimestampPipeChars1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callontimestampPipeChars1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.ontimestampPipeChars1()
}

func (c *current) onNumberLong1(n interface{}) (interface{}, error) {
	return new(LogDoc).Numberlong(n.(string)), nil
}

func (p *parser) callonNumberLong1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumberLong1(stack["n"])
}

func (c *current) onMinKey1() (interface{}, error) {
	return new(LogDoc).Minkey(), nil
}

func (p *parser) callonMinKey1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMinKey1()
}

func (c *current) onMaxKey1() (interface{}, error) {
	return new(LogDoc).Maxkey(), nil
}

func (p *parser) callonMaxKey1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMaxKey1()
}

func (c *current) onUndefined1() (interface{}, error) {
	return new(LogDoc).Undefined(), nil
}

func (p *parser) callonUndefined1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUndefined1()
}

func (c *current) onregexBody1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonregexBody1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onregexBody1()
}

func (c *current) oncharsTillParen1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) calloncharsTillParen1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.oncharsTillParen1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
