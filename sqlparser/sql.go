//line ./sqlparser/sql.y:20
package sqlparser

import __yyfmt__ "fmt"

//line ./sqlparser/sql.y:20

import "bytes"

func SetParseTree(yylex interface{}, stmt Statement) {
	yylex.(*Tokenizer).ParseTree = stmt
}

func SetAllowComments(yylex interface{}, allow bool) {
	yylex.(*Tokenizer).AllowComments = allow
}

func ForceEOF(yylex interface{}) {
	yylex.(*Tokenizer).ForceEOF = true
}

var (
	SHARE        = []byte("share")
	MODE         = []byte("mode")
	IF_BYTES     = []byte("if")
	VALUES_BYTES = []byte("values")
)

//line ./sqlparser/sql.y:45
type yySymType struct {
	yys         int
	empty       struct{}
	statement   Statement
	selStmt     SelectStatement
	byt         byte
	bytes       []byte
	bytes2      [][]byte
	str         string
	selectExprs SelectExprs
	selectExpr  SelectExpr
	columns     Columns
	colName     *ColName
	tableExprs  TableExprs
	tableExpr   TableExpr
	smTableExpr SimpleTableExpr
	tableName   *TableName
	indexHints  *IndexHints
	expr        Expr
	boolExpr    BoolExpr
	valExpr     ValExpr
	tuple       Tuple
	valExprs    ValExprs
	values      Values
	subquery    *Subquery
	caseExpr    *CaseExpr
	whens       []*When
	when        *When
	orderBy     OrderBy
	order       *Order
	limit       *Limit
	insRows     InsertRows
	updateExprs UpdateExprs
	updateExpr  *UpdateExpr
}

const LEX_ERROR = 57346
const SELECT = 57347
const INSERT = 57348
const UPDATE = 57349
const DELETE = 57350
const FROM = 57351
const WHERE = 57352
const GROUP = 57353
const HAVING = 57354
const ORDER = 57355
const BY = 57356
const LIMIT = 57357
const FOR = 57358
const ALL = 57359
const DISTINCT = 57360
const AS = 57361
const EXISTS = 57362
const NULL = 57363
const ASC = 57364
const DESC = 57365
const VALUES = 57366
const INTO = 57367
const DUPLICATE = 57368
const KEY = 57369
const DEFAULT = 57370
const SET = 57371
const LOCK = 57372
const ID = 57373
const STRING = 57374
const NUMBER = 57375
const VALUE_ARG = 57376
const COMMENT = 57377
const UNION = 57378
const MINUS = 57379
const EXCEPT = 57380
const INTERSECT = 57381
const JOIN = 57382
const STRAIGHT_JOIN = 57383
const LEFT = 57384
const RIGHT = 57385
const INNER = 57386
const OUTER = 57387
const CROSS = 57388
const NATURAL = 57389
const USE = 57390
const FORCE = 57391
const ON = 57392
const OR = 57393
const AND = 57394
const NOT = 57395
const BETWEEN = 57396
const CASE = 57397
const WHEN = 57398
const THEN = 57399
const ELSE = 57400
const LE = 57401
const GE = 57402
const NE = 57403
const NULL_SAFE_EQUAL = 57404
const IS = 57405
const LIKE = 57406
const IN = 57407
const UNARY = 57408
const END = 57409
const BEGIN = 57410
const START = 57411
const TRANSACTION = 57412
const COMMIT = 57413
const ROLLBACK = 57414
const NAMES = 57415
const REPLACE = 57416
const OFFSET = 57417
const COLLATE = 57418
const CREATE = 57419
const ALTER = 57420
const DROP = 57421
const RENAME = 57422
const TABLE = 57423
const INDEX = 57424
const VIEW = 57425
const TO = 57426
const IGNORE = 57427
const IF = 57428
const UNIQUE = 57429
const USING = 57430
const TRUNCATE = 57431

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LEX_ERROR",
	"SELECT",
	"INSERT",
	"UPDATE",
	"DELETE",
	"FROM",
	"WHERE",
	"GROUP",
	"HAVING",
	"ORDER",
	"BY",
	"LIMIT",
	"FOR",
	"ALL",
	"DISTINCT",
	"AS",
	"EXISTS",
	"NULL",
	"ASC",
	"DESC",
	"VALUES",
	"INTO",
	"DUPLICATE",
	"KEY",
	"DEFAULT",
	"SET",
	"LOCK",
	"ID",
	"STRING",
	"NUMBER",
	"VALUE_ARG",
	"COMMENT",
	"'('",
	"'~'",
	"UNION",
	"MINUS",
	"EXCEPT",
	"INTERSECT",
	"','",
	"JOIN",
	"STRAIGHT_JOIN",
	"LEFT",
	"RIGHT",
	"INNER",
	"OUTER",
	"CROSS",
	"NATURAL",
	"USE",
	"FORCE",
	"ON",
	"OR",
	"AND",
	"NOT",
	"BETWEEN",
	"CASE",
	"WHEN",
	"THEN",
	"ELSE",
	"'='",
	"'<'",
	"'>'",
	"LE",
	"GE",
	"NE",
	"NULL_SAFE_EQUAL",
	"IS",
	"LIKE",
	"IN",
	"'|'",
	"'&'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'^'",
	"'.'",
	"UNARY",
	"END",
	"BEGIN",
	"START",
	"TRANSACTION",
	"COMMIT",
	"ROLLBACK",
	"NAMES",
	"REPLACE",
	"OFFSET",
	"COLLATE",
	"CREATE",
	"ALTER",
	"DROP",
	"RENAME",
	"TABLE",
	"INDEX",
	"VIEW",
	"TO",
	"IGNORE",
	"IF",
	"UNIQUE",
	"USING",
	"TRUNCATE",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 215
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 693

var yyAct = [...]int{

	104, 343, 379, 303, 101, 173, 175, 337, 112, 70,
	256, 176, 3, 133, 215, 266, 251, 211, 102, 90,
	72, 188, 67, 35, 36, 37, 38, 91, 388, 150,
	149, 18, 200, 58, 84, 275, 276, 277, 278, 279,
	95, 280, 281, 388, 388, 74, 61, 111, 79, 73,
	117, 81, 144, 77, 312, 85, 144, 75, 108, 109,
	110, 51, 127, 115, 45, 144, 47, 356, 265, 96,
	48, 50, 244, 51, 53, 54, 55, 126, 136, 132,
	242, 355, 159, 89, 118, 354, 78, 140, 325, 327,
	245, 390, 123, 80, 146, 135, 52, 329, 75, 286,
	113, 114, 56, 334, 260, 177, 389, 387, 148, 178,
	141, 142, 129, 122, 196, 333, 181, 311, 393, 295,
	128, 252, 74, 298, 184, 74, 73, 116, 293, 73,
	195, 186, 252, 192, 193, 243, 131, 326, 83, 172,
	174, 150, 149, 185, 197, 190, 149, 360, 221, 195,
	96, 124, 219, 351, 210, 71, 338, 225, 206, 222,
	230, 231, 338, 234, 235, 236, 237, 238, 239, 240,
	241, 226, 220, 204, 229, 232, 207, 158, 157, 160,
	161, 162, 163, 164, 159, 96, 96, 228, 227, 261,
	223, 224, 150, 149, 86, 361, 139, 259, 336, 247,
	249, 262, 353, 253, 162, 163, 164, 159, 255, 352,
	233, 319, 323, 74, 74, 263, 320, 73, 271, 317,
	322, 321, 269, 268, 318, 189, 219, 124, 189, 143,
	258, 244, 285, 362, 288, 289, 272, 346, 203, 205,
	202, 275, 276, 277, 278, 279, 287, 280, 281, 49,
	292, 359, 373, 218, 96, 74, 254, 273, 217, 73,
	124, 306, 144, 213, 301, 268, 302, 297, 294, 300,
	372, 158, 157, 160, 161, 162, 163, 164, 159, 364,
	365, 18, 219, 219, 371, 212, 310, 127, 315, 316,
	191, 66, 213, 299, 182, 180, 332, 157, 160, 161,
	162, 163, 164, 159, 335, 179, 340, 218, 119, 339,
	341, 344, 217, 74, 59, 284, 147, 347, 75, 330,
	328, 345, 35, 36, 37, 38, 331, 283, 59, 158,
	157, 160, 161, 162, 163, 164, 159, 357, 385, 308,
	307, 209, 358, 158, 157, 160, 161, 162, 163, 164,
	159, 208, 386, 374, 196, 187, 192, 369, 68, 367,
	137, 134, 130, 82, 377, 366, 375, 376, 344, 121,
	120, 378, 380, 380, 380, 381, 382, 87, 18, 19,
	20, 21, 291, 198, 74, 18, 138, 64, 73, 394,
	368, 62, 370, 391, 395, 248, 396, 107, 111, 290,
	304, 117, 22, 350, 267, 305, 257, 349, 94, 108,
	109, 110, 314, 99, 115, 69, 158, 157, 160, 161,
	162, 163, 164, 159, 32, 107, 111, 189, 392, 117,
	383, 18, 40, 98, 17, 118, 94, 108, 109, 110,
	16, 99, 115, 158, 157, 160, 161, 162, 163, 164,
	159, 113, 114, 92, 15, 14, 27, 28, 13, 29,
	30, 98, 31, 118, 12, 23, 24, 26, 25, 18,
	160, 161, 162, 163, 164, 159, 88, 33, 116, 113,
	114, 92, 246, 199, 107, 111, 46, 39, 117, 264,
	201, 76, 270, 384, 363, 75, 108, 109, 110, 342,
	99, 115, 348, 313, 296, 183, 116, 41, 42, 43,
	44, 250, 107, 111, 106, 103, 117, 105, 309, 57,
	98, 60, 118, 75, 108, 109, 110, 100, 99, 115,
	151, 111, 97, 324, 117, 216, 274, 214, 113, 114,
	93, 75, 108, 109, 110, 282, 127, 115, 98, 145,
	118, 63, 34, 65, 11, 10, 9, 8, 7, 6,
	5, 4, 2, 194, 1, 116, 113, 114, 118, 0,
	0, 0, 0, 0, 0, 0, 111, 0, 0, 117,
	0, 0, 0, 125, 113, 114, 75, 108, 109, 110,
	0, 127, 115, 116, 111, 0, 0, 117, 0, 0,
	0, 0, 0, 0, 75, 108, 109, 110, 0, 127,
	115, 116, 0, 118, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 113,
	114, 118, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 113, 114, 0,
	0, 0, 0, 0, 0, 0, 116, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 153,
	155, 0, 0, 0, 116, 165, 166, 167, 168, 169,
	170, 171, 156, 154, 152, 158, 157, 160, 161, 162,
	163, 164, 159,
}
var yyPact = [...]int{

	373, -1000, -1000, 284, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -32, -27, 0, -22, -1000, 17, -1000,
	-1000, -1000, 283, -1000, 426, 374, -1000, -1000, -1000, 369,
	-1000, -39, 327, 406, 67, -48, -11, 283, -1000, -3,
	283, -1000, 332, -67, 283, -67, -1000, 352, -1000, -1000,
	-13, -1000, -1000, 405, -1000, 273, 345, 340, 33, 327,
	185, 555, -1000, 58, -1000, 32, 331, 80, 283, -1000,
	330, -1000, -21, 329, 366, 143, 283, 327, 327, -1000,
	220, -1000, -1000, 297, 28, 87, 613, -1000, 492, 464,
	-1000, -1000, -1000, 573, 269, 259, -1000, 258, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 573, -1000,
	327, 287, 324, 417, 287, -1000, 199, 26, 510, 323,
	-1000, 363, -71, -1000, 145, -1000, 320, -1000, -1000, 310,
	-1000, 256, -1000, 222, 405, -1000, -1000, 283, 83, 492,
	492, 573, 251, 117, 573, 573, 154, 573, 573, 573,
	573, 573, 573, 573, 573, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 613, -25, 30, -15, 613, -1000, 377,
	405, -1000, 426, 73, 371, 227, 218, -1000, 393, 492,
	-1000, 573, 371, 371, -1000, -1000, 24, -1000, -1000, 136,
	283, -1000, -31, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 380, 287, 287, 215, 198, 296, 276, 19, -1000,
	-1000, -1000, -1000, -1000, 91, 371, -1000, 251, 573, 573,
	371, 344, -1000, 361, 396, 224, -1000, 128, 128, 3,
	3, 3, -1000, -1000, 573, -1000, -1000, 23, 405, 14,
	62, -1000, 492, 380, 287, 393, 385, 391, 87, 371,
	283, 309, -1000, -1000, 308, -1000, -1000, 251, 284, 185,
	12, -1000, 401, 222, 222, -1000, -1000, 176, 168, 178,
	177, 169, 37, -1000, 289, -8, 288, -1000, 371, 271,
	573, -1000, 371, -1000, 10, -1000, 21, -1000, 573, 138,
	103, 109, 385, -1000, 573, 573, -1000, -1000, -1000, 195,
	-1000, -1000, 287, 395, 389, 198, 100, -1000, 166, -1000,
	159, -1000, -1000, -1000, -1000, -12, -16, -30, -1000, -1000,
	-1000, 573, 371, -1000, -1000, 371, 573, -1000, 225, -1000,
	-1000, 105, 191, -1000, 257, -1000, 251, -1000, 393, 492,
	573, 492, -1000, -1000, 248, 234, 216, 371, 371, 326,
	573, 573, 573, -1000, -1000, -1000, -1000, 385, 87, 189,
	87, 283, 283, 283, 423, 371, 371, -1000, 322, 2,
	-1000, 1, -14, 287, -1000, 421, 47, -1000, 283, -1000,
	-1000, 185, -1000, 283, -1000, 283, -1000,
}
var yyPgo = [...]int{

	0, 564, 562, 11, 561, 560, 559, 558, 557, 556,
	555, 554, 487, 553, 552, 551, 249, 19, 27, 549,
	545, 540, 537, 14, 536, 535, 22, 533, 2, 21,
	40, 532, 530, 15, 527, 5, 18, 6, 518, 517,
	8, 515, 4, 514, 511, 16, 505, 504, 503, 502,
	10, 499, 1, 494, 3, 493, 17, 492, 7, 9,
	20, 138, 491, 490, 489, 486, 483, 0, 13, 476,
	464, 458, 455, 454, 440, 434, 432,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 3, 3, 4,
	4, 73, 73, 5, 6, 7, 7, 7, 7, 70,
	70, 71, 72, 74, 75, 8, 8, 8, 9, 9,
	9, 10, 11, 11, 11, 76, 12, 13, 13, 14,
	14, 14, 14, 14, 15, 15, 17, 17, 18, 18,
	18, 21, 21, 19, 19, 19, 22, 22, 23, 23,
	23, 23, 20, 20, 20, 24, 24, 24, 24, 24,
	24, 24, 24, 24, 25, 25, 25, 26, 26, 27,
	27, 27, 27, 28, 28, 29, 29, 30, 30, 30,
	30, 30, 31, 31, 31, 31, 31, 31, 31, 31,
	31, 31, 32, 32, 32, 32, 32, 32, 32, 33,
	33, 38, 38, 36, 36, 40, 37, 37, 35, 35,
	35, 35, 35, 35, 35, 35, 35, 35, 35, 35,
	35, 35, 35, 35, 35, 39, 39, 41, 41, 41,
	43, 46, 46, 44, 44, 45, 47, 47, 42, 42,
	42, 34, 34, 34, 34, 48, 48, 49, 49, 50,
	50, 51, 51, 52, 53, 53, 53, 54, 54, 54,
	54, 55, 55, 55, 56, 56, 57, 57, 58, 58,
	59, 59, 60, 60, 61, 61, 62, 62, 16, 16,
	63, 63, 63, 63, 63, 64, 64, 65, 65, 66,
	66, 67, 68, 69, 69,
}
var yyR2 = [...]int{

	0, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 12, 3, 8,
	8, 6, 6, 8, 7, 3, 4, 4, 6, 1,
	2, 1, 1, 2, 4, 5, 8, 4, 6, 7,
	4, 5, 4, 5, 5, 0, 2, 0, 2, 1,
	2, 1, 1, 1, 0, 1, 1, 3, 1, 2,
	3, 1, 1, 0, 1, 2, 1, 3, 3, 3,
	3, 5, 0, 1, 2, 1, 1, 2, 3, 2,
	3, 2, 2, 2, 1, 3, 1, 1, 3, 0,
	5, 5, 5, 1, 3, 0, 2, 1, 3, 3,
	2, 3, 3, 3, 4, 3, 4, 5, 6, 3,
	4, 2, 1, 1, 1, 1, 1, 1, 1, 2,
	1, 1, 3, 3, 1, 3, 1, 3, 1, 1,
	1, 3, 3, 3, 3, 3, 3, 3, 3, 2,
	3, 4, 5, 4, 1, 1, 1, 1, 1, 1,
	5, 0, 1, 1, 2, 4, 0, 2, 1, 3,
	5, 1, 1, 1, 1, 0, 3, 0, 2, 0,
	3, 1, 3, 2, 0, 1, 1, 0, 2, 4,
	4, 0, 2, 4, 0, 3, 1, 3, 0, 5,
	1, 3, 3, 3, 0, 2, 0, 3, 0, 1,
	1, 1, 1, 1, 1, 0, 1, 0, 1, 0,
	2, 1, 0, 0, 1,
}
var yyChk = [...]int{

	-1000, -1, -2, -3, -4, -5, -6, -7, -8, -9,
	-10, -11, -70, -71, -72, -73, -74, -75, 5, 6,
	7, 8, 29, 92, 93, 95, 94, 83, 84, 86,
	87, 89, 51, 104, -14, 38, 39, 40, 41, -12,
	-76, -12, -12, -12, -12, 96, -65, 98, 102, -16,
	98, 100, 96, 96, 97, 98, 85, -12, -67, 31,
	-12, -3, 17, -15, 18, -13, -16, -26, 31, 9,
	-59, 88, -60, -42, -67, 31, -62, 101, 97, -67,
	96, -67, 31, -61, 101, -67, -61, 25, -69, 96,
	-17, -18, 76, -21, 31, -30, -35, -31, 56, 36,
	-34, -42, -36, -41, -67, -39, -43, 20, 32, 33,
	34, 21, -40, 74, 75, 37, 101, 24, 58, 35,
	25, 29, 80, -26, 42, 28, -35, 36, 62, 80,
	31, 56, -67, -68, 31, -68, 99, 31, 20, 53,
	-67, -26, -26, 9, 42, -19, -67, 19, 80, 55,
	54, -32, 71, 56, 70, 57, 69, 73, 72, 79,
	74, 75, 76, 77, 78, 62, 63, 64, 65, 66,
	67, 68, -30, -35, -30, -37, -3, -35, -35, 36,
	36, -40, 36, -46, -35, -26, -59, 31, -29, 10,
	-60, 91, -35, -35, 53, -67, 31, -68, 20, -66,
	103, -63, 95, 93, 28, 94, 13, 31, 31, 31,
	-68, -56, 29, 36, -22, -23, -25, 36, 31, -40,
	-18, -67, 76, -30, -30, -35, -36, 71, 70, 57,
	-35, -35, 21, 56, -35, -35, -35, -35, -35, -35,
	-35, -35, 105, 105, 42, 105, 105, -17, 18, -17,
	-44, -45, 59, -56, 29, -29, -50, 13, -30, -35,
	80, 53, -67, -68, -64, 99, -33, 24, -3, -59,
	-57, -42, -29, 42, -24, 43, 44, 45, 46, 47,
	49, 50, -20, 31, 19, -23, 80, -36, -35, -35,
	55, 21, -35, 105, -17, 105, -47, -45, 61, -30,
	-33, -59, -50, -54, 15, 14, -67, 31, 31, -38,
	-36, 105, 42, -48, 11, -23, -23, 43, 48, 43,
	48, 43, 43, 43, -27, 51, 100, 52, 31, 105,
	31, 55, -35, 105, 82, -35, 60, -58, 53, -58,
	-54, -35, -51, -52, -35, -68, 42, -42, -49, 12,
	14, 53, 43, 43, 97, 97, 97, -35, -35, 26,
	42, 90, 42, -53, 22, 23, -36, -50, -30, -37,
	-30, 36, 36, 36, 27, -35, -35, -52, -54, -28,
	-67, -28, -28, 7, -55, 16, 30, 105, 42, 105,
	105, -59, 7, 71, -67, -67, -67,
}
var yyDef = [...]int{

	0, -2, 1, 2, 3, 4, 5, 6, 7, 8,
	9, 10, 11, 12, 13, 14, 15, 16, 45, 45,
	45, 45, 45, 207, 198, 0, 0, 29, 0, 31,
	32, 45, 0, 45, 0, 49, 51, 52, 53, 54,
	47, 198, 0, 0, 0, 196, 0, 0, 208, 0,
	0, 199, 0, 194, 0, 194, 30, 0, 33, 211,
	213, 18, 50, 0, 55, 46, 0, 0, 87, 0,
	25, 0, 190, 0, 158, 211, 0, 0, 0, 212,
	0, 212, 0, 0, 0, 0, 0, 0, 0, 214,
	0, 56, 58, 63, 211, 61, 62, 97, 0, 0,
	128, 129, 130, 0, 158, 0, 144, 0, 161, 162,
	163, 164, 124, 147, 148, 149, 145, 146, 151, 48,
	0, 0, 0, 95, 0, 26, 27, 0, 0, 0,
	212, 0, 209, 37, 0, 40, 0, 42, 195, 0,
	212, 184, 34, 0, 0, 59, 64, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 112, 113, 114, 115, 116,
	117, 118, 100, 0, 0, 0, 0, 126, 139, 0,
	0, 111, 0, 0, 152, 184, 95, 88, 169, 0,
	191, 0, 126, 192, 193, 159, 211, 35, 197, 0,
	0, 212, 205, 200, 201, 202, 203, 204, 41, 43,
	44, 0, 0, 0, 95, 66, 72, 0, 84, 86,
	57, 65, 60, 98, 99, 102, 103, 0, 0, 0,
	105, 0, 109, 0, 131, 132, 133, 134, 135, 136,
	137, 138, 101, 123, 0, 125, 140, 0, 0, 0,
	156, 153, 0, 0, 0, 169, 177, 0, 96, 28,
	0, 0, 210, 38, 0, 206, 21, 0, 120, 22,
	0, 186, 165, 0, 0, 75, 76, 0, 0, 0,
	0, 0, 89, 73, 0, 0, 0, 104, 106, 0,
	0, 110, 127, 141, 0, 143, 0, 154, 0, 0,
	188, 188, 177, 24, 0, 0, 160, 212, 39, 119,
	121, 185, 0, 167, 0, 67, 70, 77, 0, 79,
	0, 81, 82, 83, 68, 0, 0, 0, 74, 69,
	85, 0, 107, 142, 150, 157, 0, 19, 0, 20,
	23, 178, 170, 171, 174, 36, 0, 187, 169, 0,
	0, 0, 78, 80, 0, 0, 0, 108, 155, 0,
	0, 0, 0, 173, 175, 176, 122, 177, 168, 166,
	71, 0, 0, 0, 0, 179, 180, 172, 181, 0,
	93, 0, 0, 0, 17, 0, 0, 90, 0, 91,
	92, 189, 182, 0, 94, 0, 183,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 78, 73, 3,
	36, 105, 76, 74, 42, 75, 80, 77, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	63, 62, 64, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 79, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 72, 3, 37,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 38, 39, 40, 41, 43, 44,
	45, 46, 47, 48, 49, 50, 51, 52, 53, 54,
	55, 56, 57, 58, 59, 60, 61, 65, 66, 67,
	68, 69, 70, 71, 81, 82, 83, 84, 85, 86,
	87, 88, 89, 90, 91, 92, 93, 94, 95, 96,
	97, 98, 99, 100, 101, 102, 103, 104,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:190
		{
			SetParseTree(yylex, yyDollar[1].statement)
		}
	case 2:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:196
		{
			yyVAL.statement = yyDollar[1].selStmt
		}
	case 17:
		yyDollar = yyS[yypt-12 : yypt+1]
		//line ./sqlparser/sql.y:216
		{
			yyVAL.selStmt = &Select{Comments: Comments(yyDollar[2].bytes2), Distinct: yyDollar[3].str, SelectExprs: yyDollar[4].selectExprs, From: yyDollar[6].tableExprs, Where: NewWhere(AST_WHERE, yyDollar[7].boolExpr), GroupBy: GroupBy(yyDollar[8].valExprs), Having: NewWhere(AST_HAVING, yyDollar[9].boolExpr), OrderBy: yyDollar[10].orderBy, Limit: yyDollar[11].limit, Lock: yyDollar[12].str}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:220
		{
			yyVAL.selStmt = &Union{Type: yyDollar[2].str, Left: yyDollar[1].selStmt, Right: yyDollar[3].selStmt}
		}
	case 19:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line ./sqlparser/sql.y:227
		{
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Ignore: yyDollar[3].str, Table: yyDollar[5].tableName, Columns: yyDollar[6].columns, Rows: yyDollar[7].insRows, OnDup: OnDup(yyDollar[8].updateExprs)}
		}
	case 20:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line ./sqlparser/sql.y:231
		{
			cols := make(Columns, 0, len(yyDollar[7].updateExprs))
			vals := make(ValTuple, 0, len(yyDollar[7].updateExprs))
			for _, col := range yyDollar[7].updateExprs {
				cols = append(cols, &NonStarExpr{Expr: col.Name})
				vals = append(vals, col.Expr)
			}
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Ignore: yyDollar[3].str, Table: yyDollar[5].tableName, Columns: cols, Rows: Values{vals}, OnDup: OnDup(yyDollar[8].updateExprs)}
		}
	case 21:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line ./sqlparser/sql.y:243
		{
			yyVAL.statement = &Replace{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: yyDollar[5].columns, Rows: yyDollar[6].insRows}
		}
	case 22:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line ./sqlparser/sql.y:247
		{
			cols := make(Columns, 0, len(yyDollar[6].updateExprs))
			vals := make(ValTuple, 0, len(yyDollar[6].updateExprs))
			for _, col := range yyDollar[6].updateExprs {
				cols = append(cols, &NonStarExpr{Expr: col.Name})
				vals = append(vals, col.Expr)
			}
			yyVAL.statement = &Replace{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: cols, Rows: Values{vals}}
		}
	case 23:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line ./sqlparser/sql.y:260
		{
			yyVAL.statement = &Update{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[3].tableName, Exprs: yyDollar[5].updateExprs, Where: NewWhere(AST_WHERE, yyDollar[6].boolExpr), OrderBy: yyDollar[7].orderBy, Limit: yyDollar[8].limit}
		}
	case 24:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line ./sqlparser/sql.y:266
		{
			yyVAL.statement = &Delete{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Where: NewWhere(AST_WHERE, yyDollar[5].boolExpr), OrderBy: yyDollar[6].orderBy, Limit: yyDollar[7].limit}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:272
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: yyDollar[3].updateExprs}
		}
	case 26:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:276
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: UpdateExprs{&UpdateExpr{Name: &ColName{Name: []byte("names")}, Expr: StrVal("default")}}}
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:280
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: UpdateExprs{&UpdateExpr{Name: &ColName{Name: []byte("names")}, Expr: yyDollar[4].valExpr}}}
		}
	case 28:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line ./sqlparser/sql.y:284
		{
			yyVAL.statement = &Set{
				Comments: Comments(yyDollar[2].bytes2),
				Exprs: UpdateExprs{
					&UpdateExpr{
						Name: &ColName{Name: []byte("names")}, Expr: yyDollar[4].valExpr,
					},
					&UpdateExpr{
						Name: &ColName{Name: []byte("collate")}, Expr: yyDollar[6].valExpr,
					},
				},
			}
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:300
		{
			yyVAL.statement = &Begin{}
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:304
		{
			yyVAL.statement = &Begin{}
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:311
		{
			yyVAL.statement = &Commit{}
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:317
		{
			yyVAL.statement = &Rollback{}
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:323
		{
			yyVAL.statement = &UseDB{DB: string(yyDollar[2].bytes)}
		}
	case 34:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:329
		{
			yyVAL.statement = &Truncate{Comments: Comments(yyDollar[2].bytes2), TableOpt: yyDollar[3].str, Table: yyDollar[4].tableName}
		}
	case 35:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:335
		{
			yyVAL.statement = &DDL{Action: AST_CREATE, NewName: yyDollar[4].bytes}
		}
	case 36:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line ./sqlparser/sql.y:339
		{
			// Change this to an alter statement
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[7].bytes, NewName: yyDollar[7].bytes}
		}
	case 37:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:344
		{
			yyVAL.statement = &DDL{Action: AST_CREATE, NewName: yyDollar[3].bytes}
		}
	case 38:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line ./sqlparser/sql.y:350
		{
			yyVAL.statement = &DDL{Action: AST_ALTER, Ignore: yyDollar[2].str, Table: yyDollar[4].bytes, NewName: yyDollar[4].bytes}
		}
	case 39:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line ./sqlparser/sql.y:354
		{
			// Change this to a rename statement
			yyVAL.statement = &DDL{Action: AST_RENAME, Ignore: yyDollar[2].str, Table: yyDollar[4].bytes, NewName: yyDollar[7].bytes}
		}
	case 40:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:359
		{
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[3].bytes, NewName: yyDollar[3].bytes}
		}
	case 41:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:365
		{
			yyVAL.statement = &DDL{Action: AST_RENAME, Table: yyDollar[3].bytes, NewName: yyDollar[5].bytes}
		}
	case 42:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:371
		{
			yyVAL.statement = &DDL{Action: AST_DROP, Table: yyDollar[4].bytes}
		}
	case 43:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:375
		{
			// Change this to an alter statement
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[5].bytes, NewName: yyDollar[5].bytes}
		}
	case 44:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:380
		{
			yyVAL.statement = &DDL{Action: AST_DROP, Table: yyDollar[4].bytes}
		}
	case 45:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:385
		{
			SetAllowComments(yylex, true)
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:389
		{
			yyVAL.bytes2 = yyDollar[2].bytes2
			SetAllowComments(yylex, false)
		}
	case 47:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:395
		{
			yyVAL.bytes2 = nil
		}
	case 48:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:399
		{
			yyVAL.bytes2 = append(yyDollar[1].bytes2, yyDollar[2].bytes)
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:405
		{
			yyVAL.str = AST_UNION
		}
	case 50:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:409
		{
			yyVAL.str = AST_UNION_ALL
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:413
		{
			yyVAL.str = AST_SET_MINUS
		}
	case 52:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:417
		{
			yyVAL.str = AST_EXCEPT
		}
	case 53:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:421
		{
			yyVAL.str = AST_INTERSECT
		}
	case 54:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:426
		{
			yyVAL.str = ""
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:430
		{
			yyVAL.str = AST_DISTINCT
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:436
		{
			yyVAL.selectExprs = SelectExprs{yyDollar[1].selectExpr}
		}
	case 57:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:440
		{
			yyVAL.selectExprs = append(yyVAL.selectExprs, yyDollar[3].selectExpr)
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:446
		{
			yyVAL.selectExpr = &StarExpr{}
		}
	case 59:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:450
		{
			yyVAL.selectExpr = &NonStarExpr{Expr: yyDollar[1].expr, As: yyDollar[2].bytes}
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:454
		{
			yyVAL.selectExpr = &StarExpr{TableName: yyDollar[1].bytes}
		}
	case 61:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:460
		{
			yyVAL.expr = yyDollar[1].boolExpr
		}
	case 62:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:464
		{
			yyVAL.expr = yyDollar[1].valExpr
		}
	case 63:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:469
		{
			yyVAL.bytes = nil
		}
	case 64:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:473
		{
			yyVAL.bytes = yyDollar[1].bytes
		}
	case 65:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:477
		{
			yyVAL.bytes = yyDollar[2].bytes
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:483
		{
			yyVAL.tableExprs = TableExprs{yyDollar[1].tableExpr}
		}
	case 67:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:487
		{
			yyVAL.tableExprs = append(yyVAL.tableExprs, yyDollar[3].tableExpr)
		}
	case 68:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:493
		{
			yyVAL.tableExpr = &AliasedTableExpr{Expr: yyDollar[1].smTableExpr, As: yyDollar[2].bytes, Hints: yyDollar[3].indexHints}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:497
		{
			yyVAL.tableExpr = &ParenTableExpr{Expr: yyDollar[2].tableExpr}
		}
	case 70:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:501
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr}
		}
	case 71:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:505
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr, On: yyDollar[5].boolExpr}
		}
	case 72:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:510
		{
			yyVAL.bytes = nil
		}
	case 73:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:514
		{
			yyVAL.bytes = yyDollar[1].bytes
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:518
		{
			yyVAL.bytes = yyDollar[2].bytes
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:524
		{
			yyVAL.str = AST_JOIN
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:528
		{
			yyVAL.str = AST_STRAIGHT_JOIN
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:532
		{
			yyVAL.str = AST_LEFT_JOIN
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:536
		{
			yyVAL.str = AST_LEFT_JOIN
		}
	case 79:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:540
		{
			yyVAL.str = AST_RIGHT_JOIN
		}
	case 80:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:544
		{
			yyVAL.str = AST_RIGHT_JOIN
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:548
		{
			yyVAL.str = AST_JOIN
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:552
		{
			yyVAL.str = AST_CROSS_JOIN
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:556
		{
			yyVAL.str = AST_NATURAL_JOIN
		}
	case 84:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:562
		{
			yyVAL.smTableExpr = &TableName{Name: yyDollar[1].bytes}
		}
	case 85:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:566
		{
			yyVAL.smTableExpr = &TableName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:570
		{
			yyVAL.smTableExpr = yyDollar[1].subquery
		}
	case 87:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:576
		{
			yyVAL.tableName = &TableName{Name: yyDollar[1].bytes}
		}
	case 88:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:580
		{
			yyVAL.tableName = &TableName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 89:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:585
		{
			yyVAL.indexHints = nil
		}
	case 90:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:589
		{
			yyVAL.indexHints = &IndexHints{Type: AST_USE, Indexes: yyDollar[4].bytes2}
		}
	case 91:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:593
		{
			yyVAL.indexHints = &IndexHints{Type: AST_IGNORE, Indexes: yyDollar[4].bytes2}
		}
	case 92:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:597
		{
			yyVAL.indexHints = &IndexHints{Type: AST_FORCE, Indexes: yyDollar[4].bytes2}
		}
	case 93:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:603
		{
			yyVAL.bytes2 = [][]byte{yyDollar[1].bytes}
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:607
		{
			yyVAL.bytes2 = append(yyDollar[1].bytes2, yyDollar[3].bytes)
		}
	case 95:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:612
		{
			yyVAL.boolExpr = nil
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:616
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:623
		{
			yyVAL.boolExpr = &AndExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 99:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:627
		{
			yyVAL.boolExpr = &OrExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 100:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:631
		{
			yyVAL.boolExpr = &NotExpr{Expr: yyDollar[2].boolExpr}
		}
	case 101:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:635
		{
			yyVAL.boolExpr = &ParenBoolExpr{Expr: yyDollar[2].boolExpr}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:641
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: yyDollar[2].str, Right: yyDollar[3].valExpr}
		}
	case 103:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:645
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_IN, Right: yyDollar[3].tuple}
		}
	case 104:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:649
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_NOT_IN, Right: yyDollar[4].tuple}
		}
	case 105:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:653
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_LIKE, Right: yyDollar[3].valExpr}
		}
	case 106:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:657
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_NOT_LIKE, Right: yyDollar[4].valExpr}
		}
	case 107:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:661
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: AST_BETWEEN, From: yyDollar[3].valExpr, To: yyDollar[5].valExpr}
		}
	case 108:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line ./sqlparser/sql.y:665
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: AST_NOT_BETWEEN, From: yyDollar[4].valExpr, To: yyDollar[6].valExpr}
		}
	case 109:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:669
		{
			yyVAL.boolExpr = &NullCheck{Operator: AST_IS_NULL, Expr: yyDollar[1].valExpr}
		}
	case 110:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:673
		{
			yyVAL.boolExpr = &NullCheck{Operator: AST_IS_NOT_NULL, Expr: yyDollar[1].valExpr}
		}
	case 111:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:677
		{
			yyVAL.boolExpr = &ExistsExpr{Subquery: yyDollar[2].subquery}
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:683
		{
			yyVAL.str = AST_EQ
		}
	case 113:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:687
		{
			yyVAL.str = AST_LT
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:691
		{
			yyVAL.str = AST_GT
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:695
		{
			yyVAL.str = AST_LE
		}
	case 116:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:699
		{
			yyVAL.str = AST_GE
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:703
		{
			yyVAL.str = AST_NE
		}
	case 118:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:707
		{
			yyVAL.str = AST_NSE
		}
	case 119:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:713
		{
			yyVAL.insRows = yyDollar[2].values
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:717
		{
			yyVAL.insRows = yyDollar[1].selStmt
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:723
		{
			yyVAL.values = Values{yyDollar[1].tuple}
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:727
		{
			yyVAL.values = append(yyDollar[1].values, yyDollar[3].tuple)
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:733
		{
			yyVAL.tuple = ValTuple(yyDollar[2].valExprs)
		}
	case 124:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:737
		{
			yyVAL.tuple = yyDollar[1].subquery
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:743
		{
			yyVAL.subquery = &Subquery{yyDollar[2].selStmt}
		}
	case 126:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:749
		{
			yyVAL.valExprs = ValExprs{yyDollar[1].valExpr}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:753
		{
			yyVAL.valExprs = append(yyDollar[1].valExprs, yyDollar[3].valExpr)
		}
	case 128:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:759
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 129:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:763
		{
			yyVAL.valExpr = yyDollar[1].colName
		}
	case 130:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:767
		{
			yyVAL.valExpr = yyDollar[1].tuple
		}
	case 131:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:771
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITAND, Right: yyDollar[3].valExpr}
		}
	case 132:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:775
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITOR, Right: yyDollar[3].valExpr}
		}
	case 133:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:779
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITXOR, Right: yyDollar[3].valExpr}
		}
	case 134:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:783
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_PLUS, Right: yyDollar[3].valExpr}
		}
	case 135:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:787
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MINUS, Right: yyDollar[3].valExpr}
		}
	case 136:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:791
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MULT, Right: yyDollar[3].valExpr}
		}
	case 137:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:795
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_DIV, Right: yyDollar[3].valExpr}
		}
	case 138:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:799
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MOD, Right: yyDollar[3].valExpr}
		}
	case 139:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:803
		{
			if num, ok := yyDollar[2].valExpr.(NumVal); ok {
				switch yyDollar[1].byt {
				case '-':
					yyVAL.valExpr = append(NumVal("-"), num...)
				case '+':
					yyVAL.valExpr = num
				default:
					yyVAL.valExpr = &UnaryExpr{Operator: yyDollar[1].byt, Expr: yyDollar[2].valExpr}
				}
			} else {
				yyVAL.valExpr = &UnaryExpr{Operator: yyDollar[1].byt, Expr: yyDollar[2].valExpr}
			}
		}
	case 140:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:818
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes}
		}
	case 141:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:822
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Exprs: yyDollar[3].selectExprs}
		}
	case 142:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:826
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Distinct: true, Exprs: yyDollar[4].selectExprs}
		}
	case 143:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:830
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Exprs: yyDollar[3].selectExprs}
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:834
		{
			yyVAL.valExpr = yyDollar[1].caseExpr
		}
	case 145:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:840
		{
			yyVAL.bytes = IF_BYTES
		}
	case 146:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:844
		{
			yyVAL.bytes = VALUES_BYTES
		}
	case 147:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:850
		{
			yyVAL.byt = AST_UPLUS
		}
	case 148:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:854
		{
			yyVAL.byt = AST_UMINUS
		}
	case 149:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:858
		{
			yyVAL.byt = AST_TILDA
		}
	case 150:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:864
		{
			yyVAL.caseExpr = &CaseExpr{Expr: yyDollar[2].valExpr, Whens: yyDollar[3].whens, Else: yyDollar[4].valExpr}
		}
	case 151:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:869
		{
			yyVAL.valExpr = nil
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:873
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 153:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:879
		{
			yyVAL.whens = []*When{yyDollar[1].when}
		}
	case 154:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:883
		{
			yyVAL.whens = append(yyDollar[1].whens, yyDollar[2].when)
		}
	case 155:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:889
		{
			yyVAL.when = &When{Cond: yyDollar[2].boolExpr, Val: yyDollar[4].valExpr}
		}
	case 156:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:894
		{
			yyVAL.valExpr = nil
		}
	case 157:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:898
		{
			yyVAL.valExpr = yyDollar[2].valExpr
		}
	case 158:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:904
		{
			yyVAL.colName = &ColName{Name: yyDollar[1].bytes}
		}
	case 159:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:908
		{
			yyVAL.colName = &ColName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 160:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:912
		{
			yyVAL.colName = &ColName{Qualifier: yyDollar[3].bytes, Name: yyDollar[5].bytes}
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:918
		{
			yyVAL.valExpr = StrVal(yyDollar[1].bytes)
		}
	case 162:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:922
		{
			yyVAL.valExpr = NumVal(yyDollar[1].bytes)
		}
	case 163:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:926
		{
			yyVAL.valExpr = ValArg(yyDollar[1].bytes)
		}
	case 164:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:930
		{
			yyVAL.valExpr = &NullVal{}
		}
	case 165:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:935
		{
			yyVAL.valExprs = nil
		}
	case 166:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:939
		{
			yyVAL.valExprs = yyDollar[3].valExprs
		}
	case 167:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:944
		{
			yyVAL.boolExpr = nil
		}
	case 168:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:948
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 169:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:953
		{
			yyVAL.orderBy = nil
		}
	case 170:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:957
		{
			yyVAL.orderBy = yyDollar[3].orderBy
		}
	case 171:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:963
		{
			yyVAL.orderBy = OrderBy{yyDollar[1].order}
		}
	case 172:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:967
		{
			yyVAL.orderBy = append(yyDollar[1].orderBy, yyDollar[3].order)
		}
	case 173:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:973
		{
			yyVAL.order = &Order{Expr: yyDollar[1].valExpr, Direction: yyDollar[2].str}
		}
	case 174:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:978
		{
			yyVAL.str = AST_ASC
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:982
		{
			yyVAL.str = AST_ASC
		}
	case 176:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:986
		{
			yyVAL.str = AST_DESC
		}
	case 177:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:991
		{
			yyVAL.limit = nil
		}
	case 178:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:995
		{
			yyVAL.limit = &Limit{Rowcount: yyDollar[2].valExpr}
		}
	case 179:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:999
		{
			yyVAL.limit = &Limit{Offset: yyDollar[2].valExpr, Rowcount: yyDollar[4].valExpr}
		}
	case 180:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:1003
		{
			yyVAL.limit = &Limit{Offset: yyDollar[4].valExpr, Rowcount: yyDollar[2].valExpr}
		}
	case 181:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1008
		{
			yyVAL.str = ""
		}
	case 182:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:1012
		{
			yyVAL.str = AST_FOR_UPDATE
		}
	case 183:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line ./sqlparser/sql.y:1016
		{
			if !bytes.Equal(yyDollar[3].bytes, SHARE) {
				yylex.Error("expecting share")
				return 1
			}
			if !bytes.Equal(yyDollar[4].bytes, MODE) {
				yylex.Error("expecting mode")
				return 1
			}
			yyVAL.str = AST_SHARE_MODE
		}
	case 184:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1029
		{
			yyVAL.columns = nil
		}
	case 185:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:1033
		{
			yyVAL.columns = yyDollar[2].columns
		}
	case 186:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1039
		{
			yyVAL.columns = Columns{&NonStarExpr{Expr: yyDollar[1].colName}}
		}
	case 187:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:1043
		{
			yyVAL.columns = append(yyVAL.columns, &NonStarExpr{Expr: yyDollar[3].colName})
		}
	case 188:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1048
		{
			yyVAL.updateExprs = nil
		}
	case 189:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line ./sqlparser/sql.y:1052
		{
			yyVAL.updateExprs = yyDollar[5].updateExprs
		}
	case 190:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1058
		{
			yyVAL.updateExprs = UpdateExprs{yyDollar[1].updateExpr}
		}
	case 191:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:1062
		{
			yyVAL.updateExprs = append(yyDollar[1].updateExprs, yyDollar[3].updateExpr)
		}
	case 192:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:1068
		{
			yyVAL.updateExpr = &UpdateExpr{Name: yyDollar[1].colName, Expr: yyDollar[3].valExpr}
		}
	case 193:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:1072
		{
			yyVAL.updateExpr = &UpdateExpr{Name: yyDollar[1].colName, Expr: StrVal("ON")}
		}
	case 194:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1077
		{
			yyVAL.empty = struct{}{}
		}
	case 195:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:1079
		{
			yyVAL.empty = struct{}{}
		}
	case 196:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1082
		{
			yyVAL.empty = struct{}{}
		}
	case 197:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line ./sqlparser/sql.y:1084
		{
			yyVAL.empty = struct{}{}
		}
	case 198:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1087
		{
			yyVAL.str = ""
		}
	case 199:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1089
		{
			yyVAL.str = AST_IGNORE
		}
	case 200:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1093
		{
			yyVAL.empty = struct{}{}
		}
	case 201:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1095
		{
			yyVAL.empty = struct{}{}
		}
	case 202:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1097
		{
			yyVAL.empty = struct{}{}
		}
	case 203:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1099
		{
			yyVAL.empty = struct{}{}
		}
	case 204:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1101
		{
			yyVAL.empty = struct{}{}
		}
	case 205:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1104
		{
			yyVAL.empty = struct{}{}
		}
	case 206:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1106
		{
			yyVAL.empty = struct{}{}
		}
	case 207:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1109
		{
			yyVAL.empty = struct{}{}
		}
	case 208:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1111
		{
			yyVAL.empty = struct{}{}
		}
	case 209:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1114
		{
			yyVAL.empty = struct{}{}
		}
	case 210:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line ./sqlparser/sql.y:1116
		{
			yyVAL.empty = struct{}{}
		}
	case 211:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1120
		{
			yyVAL.bytes = bytes.ToLower(yyDollar[1].bytes)
		}
	case 212:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1125
		{
			ForceEOF(yylex)
		}
	case 213:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line ./sqlparser/sql.y:1130
		{
			yyVAL.str = ""
		}
	case 214:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line ./sqlparser/sql.y:1134
		{
			yyVAL.str = AST_TABLE
		}
	}
	goto yystack /* stack new state and value */
}
