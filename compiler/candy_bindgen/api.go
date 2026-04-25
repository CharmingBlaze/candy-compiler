package candy_bindgen

type API struct {
	Functions []Function `json:"functions,omitempty"`
	Types     []TypeDef  `json:"types,omitempty"`
	Structs   []Struct   `json:"structs,omitempty"`
	Enums     []Enum     `json:"enums,omitempty"`
	Constants []Constant `json:"constants,omitempty"`
}

type Function struct {
	Name       string      `json:"name"`
	Symbol     string      `json:"symbol,omitempty"`
	ReturnType string      `json:"return_type"`
	Params     []Parameter `json:"params,omitempty"`
	Variadic   bool        `json:"variadic,omitempty"`
	Comment    string      `json:"comment,omitempty"`
}

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type TypeDef struct {
	Name      string `json:"name"`
	CType     string `json:"c_type"`
	CandyType string `json:"candy_type,omitempty"`
}

type Struct struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields,omitempty"`
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Enum struct {
	Name   string      `json:"name"`
	Values []EnumValue `json:"values,omitempty"`
}

type EnumValue struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type Constant struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

type ParseOptions struct {
	IncludeDirs []string
	Defines     []string
	SimpleOnly  bool
	Language    string
	UnsafeABI   bool
}
