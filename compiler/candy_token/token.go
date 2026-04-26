package candy_token

import "strings"

type TokenType string

// Token is a lexeme; Line/Col are 1-based after lexing.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
	Offset  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"
	FLOAT = "FLOAT"
	STR   = "STR" // string literal

	ASSIGN   = "="
	PLUS_ASSIGN = "+="
	MINUS_ASSIGN = "-="
	STAR_ASSIGN = "*="
	SLASH_ASSIGN = "/="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	PERCENT  = "%"

	INC = "++"
	DEC = "--"

	LT = "<"
	GT = ">"
	LTE = "<="
	GTE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	QUESTION  = "?"
	OR_ASSIGN = "||="
	NULL_COALESCE_ASSIGN = "??="

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACK = "["
	RBRACK = "]"
	DOT    = "."

	FUNCTION  = "FUNCTION"
	VAL       = "VAL"
	VAR       = "VAR"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	IMPORT    = "IMPORT"
	FROM      = "FROM"
	STRUCT    = "STRUCT"
	NULL      = "NULL"
	NEW       = "NEW"
	WHEN      = "WHEN"
	MAP       = "MAP"
	MATCH     = "MATCH"
	IS        = "IS"
	CLASS     = "CLASS"
	OBJECT    = "OBJECT"
	INTERFACE = "INTERFACE"
	TRAIT     = "TRAIT"
	EXTENDS   = "EXTENDS"
	SEALED    = "SEALED"
	SUSPEND   = "SUSPEND"
	EXTERN    = "EXTERN"
	PACKAGE   = "PACKAGE"
	DIM       = "DIM"
	THEN      = "THEN"
	END       = "END"
	SUB       = "SUB"
	AS        = "AS"
	FOR       = "FOR"
	TO        = "TO"
	STEP      = "STEP"
	NEXT      = "NEXT"
	WHILE     = "WHILE"
	DO        = "DO"
	SWITCH    = "SWITCH"
	DEFAULT   = "DEFAULT"
	WEND      = "WEND"
	SELECT    = "SELECT"
	CASE      = "CASE"
	ELSEIF    = "ELSEIF"
	GOTO      = "GOTO"
	MOD       = "MOD"
	VOID      = "VOID"
	CHAR      = "CHAR"
	OVERRIDE  = "OVERRIDE"
	DEFER     = "DEFER"
	ARROW     = "=>"
	CONST     = "CONST"
	PRINT     = "PRINT"
	IN        = "IN"
	REF       = "REF"
	SHARED    = "SHARED"

	MODULE   = "MODULE"
	EXPORT   = "EXPORT"
	ENUM     = "ENUM"
	TRY      = "TRY"
	CATCH    = "CATCH"
	FINALLY  = "FINALLY"
	ASYNC    = "ASYNC"
	AWAIT    = "AWAIT"
	RUN      = "RUN"
	TYPEOF   = "TYPEOF"
	MAYBE    = "MAYBE"
	OPERATOR = "OPERATOR"
	PRIVATE  = "PRIVATE"

	AMPERSAND     = "&"
	BIT_OR        = "|"
	BIT_XOR       = "^"
	BIT_NOT       = "~"
	SHL           = "<<"
	SHR           = ">>"
	SAFE_DOT      = "?."
	NULL_COALESCE = "??"
	USING         = "USING"
	WITH          = "WITH"
	RANGE         = ".."
	RANGE_EXCL    = "..<"
	EACH          = "EACH"
	AND           = "AND"
	OR            = "OR"
	NOT           = "NOT"

	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	FOREACH  = "FOREACH"
	DELETE   = "DELETE"

	AND_AND = "&&"
	OR_OR   = "||"

	REPEAT = "REPEAT"
	LOOP   = "LOOP"
	SUPER  = "SUPER"
	PIPELINE = "|>"
)

var keywords = map[string]TokenType{
	"fun":       FUNCTION,
	"func":      FUNCTION,
	"function":  FUNCTION,
	"val":       VAL,
	"var":       VAR,
	"true":      TRUE,
	"false":     FALSE,
	"if":        IF,
	"else":      ELSE,
	"return":    RETURN,
	"import":    IMPORT,
	"from":      FROM,
	"struct":    STRUCT,
	"new":       NEW,
	"null":      NULL,
	"when":      WHEN,
	"map":       MAP,
	"match":     MATCH,
	"is":        IS,
	"class":     CLASS,
	"object":    OBJECT,
	"interface": INTERFACE,
	"trait":     TRAIT,
	"extends":   EXTENDS,
	"sealed":    SEALED,
	"suspend":   SUSPEND,
	"extern":    EXTERN,
	"package":   PACKAGE,
	"dim":       DIM,
	"then":      THEN,
	"end":       END,
	"sub":       SUB,
	"as":        AS,
	"for":       FOR,
	"to":        TO,
	"step":      STEP,
	"next":      NEXT,
	"while":     WHILE,
	"do":        DO,
	"switch":    SWITCH,
	"default":   DEFAULT,
	"wend":      WEND,
	"select":    SELECT,
	"case":      CASE,
	"elseif":    ELSEIF,
	"goto":      GOTO,
	"mod":       MOD,
	"void":      VOID,
	"char":      CHAR,
	"override":  OVERRIDE,
	"defer":     DEFER,
	"const":     CONST,
	"print":     PRINT,
	"in":        IN,
	"ref":       REF,
	"shared":    SHARED,
	"module":    MODULE,
	"export":    EXPORT,
	"enum":      ENUM,
	"try":       TRY,
	"catch":     CATCH,
	"finally":   FINALLY,
	"async":     ASYNC,
	"await":     AWAIT,
	"run":       RUN,
	"typeof":    TYPEOF,
	"maybe":     MAYBE,
	"operator":  OPERATOR,
	"private":   PRIVATE,
	"using":     USING,
	"with":      WITH,
	"each":      EACH,
	"and":       AND,
	"or":        OR,
	"not":       NOT,
	"break":     BREAK,
	"continue":  CONTINUE,
	"foreach":   FOREACH,
	"delete":    DELETE,
	"repeat":    REPEAT,
	"loop":      LOOP,
	"super":     SUPER,
}

func LookupIdent(ident string) TokenType {
	// Keywords are ASCII case-insensitive (lexer also normalizes, but keep this safe for all callers).
	key := strings.ToLower(ident)
	if tok, ok := keywords[key]; ok {
		return tok
	}
	return IDENT
}
