plaintext# CANDY LANGUAGE SPECIFICATION - COMPILER IMPLEMENTATION CHECKLIST
# Complete Feature List for Compiler Developers

================================================================================
## CORE LANGUAGE FEATURES
================================================================================

### 1. VARIABLES & TYPES

☐ **Variable Declaration (Type Inference)**
  - Syntax: `x = 10`
  - Infer type from value (int, float, string, bool)
  - No explicit type annotations required in user code
  - Support: int, float, string, bool

☐ **Type Inference Rules**
  - Numbers without decimal → int
  - Numbers with decimal → float
  - Quoted text → string
  - true/false → bool
  - [] → array
  - {} → object/struct

☐ **Constants**
  - Syntax: `const NAME = value`
  - Immutable after declaration
  - Support compile-time constant folding


### 2. OPERATORS

☐ **Arithmetic Operators**
  - Addition: `+`
  - Subtraction: `-`
  - Multiplication: `*`
  - Division: `/`
  - Modulo: `%`
  - Increment: `++` (postfix)
  - Decrement: `--` (postfix)

☐ **Comparison Operators**
  - Equal: `==`
  - Not equal: `!=`
  - Less than: `<`
  - Greater than: `>`
  - Less or equal: `<=`
  - Greater or equal: `>=`

☐ **Logical Operators**
  - AND: `and` or `&&`
  - OR: `or` or `||`
  - NOT: `not` or `!`

☐ **Bitwise Operators** (for C interop)
  - OR: `|`
  - AND: `&`
  - XOR: `^`
  - NOT: `~`
  - Left shift: `<<`
  - Right shift: `>>`

☐ **Assignment Operators**
  - Basic: `=`
  - Compound: `+=`, `-=`, `*=`, `/=`


### 3. CONTROL FLOW

☐ **If Statement**
```candy
  if condition {
    // code
  }
```

☐ **If-Else Statement**
```candy
  if condition {
    // code
  } else {
    // code
  }
```

☐ **If-Else-If Chain**
```candy
  if condition1 {
    // code
  } else if condition2 {
    // code
  } else {
    // code
  }
```

☐ **While Loop**
```candy
  while condition {
    // code
  }
```

☐ **Do-While Loop**
```candy
  do {
    // code
  } while condition
```

☐ **For Loop (Range-based)**
```candy
  for i = 0 to 10 {
    // code
  }
  
  for i = 10 to 0 step -1 {
    // code
  }
```
  - Translate to: `for (int i = start; i <= end; i += step)`
  - Default step = 1

☐ **For-Each Loop**
```candy
  for item in collection {
    // code
  }
```

☐ **Repeat Loop**
```candy
  repeat 5 {
    // code
  }
```
  - Translate to: `for (int _i = 0; _i < count; _i++)`

☐ **Loop Control**
  - `break` - Exit loop
  - `continue` - Skip to next iteration

☐ **Infinite Loop**
```candy
  loop {
    // code
  }
```
  - Translate to: `while (1)`


### 4. FUNCTIONS

☐ **Function Declaration**
```candy
  fun functionName() {
    // code
  }
```

☐ **Function with Parameters**
```candy
  fun functionName(param1, param2, param3) {
    // code
  }
```

☐ **Function with Return Value**
```candy
  fun functionName(x, y) {
    return x + y
  }
```

☐ **Lambda/Arrow Functions** (optional syntax)
```candy
  callback = () => {
    // code
  }
  
  onClick = (event) => {
    print(event)
  }
```

☐ **Function Calls**
```candy
  result = myFunction(arg1, arg2)
```


### 5. OBJECTS/STRUCTS

☐ **Object Declaration**
```candy
  object Player {
    x = 0
    y = 0
    health = 100
    
    fun move(dx, dy) {
      x = x + dx
      y = y + dy
    }
  }
```

☐ **Object Instantiation**
```candy
  player = Player()
```

☐ **Property Access**
```candy
  player.x = 100
  value = player.health
```

☐ **Method Calls**
```candy
  player.move(5, 10)
```

☐ **Inline Object/Struct Creation**
```candy
  point = {x: 10, y: 20}
  person = {name: "Alice", age: 25}
```


### 6. ARRAYS/LISTS

☐ **Array Literal**
```candy
  numbers = [1, 2, 3, 4, 5]
  names = ["Alice", "Bob", "Charlie"]
  mixed = [1, "hello", 3.14, true]
```

☐ **Array Access**
```candy
  value = array[index]
  array[index] = value
```

☐ **Array Properties/Methods**
  - `.count` or `.length` - Get size
  - `.add(item)` - Append item
  - `.remove(index)` - Remove at index
  - `.removeLast()` - Remove last item
  - `.clear()` - Remove all items
  - `.empty()` - Check if empty

☐ **Fixed-Size Arrays** (for C interop)
```candy
  buffer = array(100)    // 100-element array
  data = bytes(256)      // 256-byte buffer
```


### 7. STRINGS

☐ **String Literals**
  - Double quotes: `"hello"`
  - Single quotes: `'hello'` (optional)

☐ **String Concatenation**
```candy
  message = "Hello " + name
```

☐ **String Interpolation**
```candy
  message = "Score: {score}, Name: {name}"
```
  - Translate to: sprintf or TextFormat

☐ **String Functions**
  - `upper(string)` - Convert to uppercase
  - `lower(string)` - Convert to lowercase
  - `length(string)` - Get string length
  - `split(string, delimiter)` - Split into array
  - `join(strings...)` - Join strings
  - `contains(string, substring)` - Check contains
  - `replace(string, old, new)` - Replace text


### 8. NULL/NONE HANDLING

☐ **Null Keyword**
```candy
  value = null
```

☐ **Null Checks**
```candy
  if value == null {
    // handle null
  }
  
  if value {
    // value is not null
  }
```

☐ **Default Value Operator**
```candy
  result = getValue() or defaultValue
```
  - If getValue() returns null, use defaultValue


### 9. ERROR HANDLING

☐ **Try-Catch**
```candy
  try {
    riskyOperation()
  } catch {
    print("Error occurred!")
  }
```

☐ **Try-Catch with Error Variable** (optional)
```candy
  try {
    riskyOperation()
  } catch error {
    print("Error: " + error)
  }
```


### 10. RESOURCE MANAGEMENT

☐ **With Statement** (automatic cleanup)
```candy
  with file = openFile("data.txt") {
    content = read(file)
    // file automatically closed here
  }
```

☐ **Manual Delete** (when needed)
```candy
  resource = createResource()
  // use resource
  delete resource
```


### 11. COMMENTS

☐ **Single-Line Comments**
```candy
  // This is a comment
```

☐ **Multi-Line Comments**
```candy
  /*
    This is a
    multi-line comment
  */
```


### 12. ENUMS (Optional)

☐ **Enum Declaration**
```candy
  enum State {
    IDLE,
    RUNNING,
    JUMPING
  }
```

☐ **Enum with Values**
```candy
  enum Color {
    RED = 1,
    GREEN = 2,
    BLUE = 3
  }
```


================================================================================
## C INTEROP FEATURES
================================================================================

### 13. EXTERNAL FUNCTION DECLARATIONS

☐ **Extern Keyword**
```candy
  extern functionName(param1: type, param2: type): returnType
```

☐ **Type Annotations** (only in extern declarations)
  - `int` - 32-bit integer
  - `float` - 32-bit float
  - `double` - 64-bit float
  - `bool` - Boolean
  - `byte` - 8-bit unsigned
  - `cstring` - C string (char*)
  - `pointer` - void*
  - Custom types from C libraries

☐ **Variadic Functions**
```candy
  extern printf(format: cstring, ...args): int
```


### 14. LIBRARY SYSTEM

☐ **Library Declaration**
```candy
  library "libraryname" {
    // external declarations
    // constants
    // helper functions
  }
```

☐ **Import Statement**
```candy
  import "libraryname"
```

☐ **Type Definitions**
```candy
  type StructName {
    field1: type
    field2: type
  }
```


================================================================================
## BUILT-IN FUNCTIONS (STANDARD LIBRARY)
================================================================================

### 15. I/O FUNCTIONS

☐ **Output**
  - `print(value)` - Print with newline
  - `println(value)` - Print with newline (alias)

☐ **Input**
  - `input(prompt)` - Read line from user
  - `readLine()` - Read line without prompt


### 16. MATH FUNCTIONS

☐ **Basic Math**
  - `abs(x)` - Absolute value
  - `sqrt(x)` - Square root
  - `pow(x, y)` - Power
  - `min(a, b)` - Minimum
  - `max(a, b)` - Maximum

☐ **Rounding**
  - `round(x)` - Round to nearest
  - `floor(x)` - Round down
  - `ceil(x)` - Round up

☐ **Trigonometry**
  - `sin(x)` - Sine
  - `cos(x)` - Cosine
  - `tan(x)` - Tangent

☐ **Utility**
  - `random(min, max)` - Random integer
  - `clamp(value, min, max)` - Clamp value
  - `lerp(a, b, t)` - Linear interpolation


### 17. TYPE CONVERSION

☐ **Conversions**
  - `int(value)` or `toInt(value)` - Convert to int
  - `float(value)` or `toFloat(value)` - Convert to float
  - `string(value)` or `toString(value)` - Convert to string
  - `bool(value)` or `toBool(value)` - Convert to bool


### 18. UTILITY FUNCTIONS

☐ **Program Control**
  - `exit()` - Exit program
  - `wait(seconds)` - Pause execution

☐ **Time Functions**
  - `seconds()` - Get elapsed time
  - `deltaTime()` - Get frame delta time


================================================================================
## FILE I/O (OPTIONAL STANDARD LIBRARY)
================================================================================

### 19. FILE OPERATIONS

☐ **File Reading/Writing**
  - `readFile(filename)` - Read entire file
  - `writeFile(filename, content)` - Write file
  - `appendFile(filename, content)` - Append to file

☐ **File Checks**
  - `fileExists(filename)` - Check if file exists

☐ **Data Persistence**
  - `save(key, value)` - Save key-value data
  - `load(key, default)` - Load key-value data


================================================================================
## TRANSPILER REQUIREMENTS
================================================================================

### 20. C CODE GENERATION

☐ **Header Generation**
  - Generate necessary #includes
  - Generate type definitions
  - Generate function prototypes

☐ **Variable Translation**
  - Infer C types from Candy types
  - Generate proper C variable declarations
  - Handle type conversions

☐ **Control Flow Translation**
  - Translate if/else to C
  - Translate loops to C
  - Handle break/continue

☐ **Function Translation**
  - Generate C function definitions
  - Handle return values
  - Generate lambda/closure code

☐ **Object Translation**
  - Translate objects to C structs
  - Generate constructor functions
  - Generate method functions

☐ **Array Translation**
  - Use dynamic arrays (malloc/realloc)
  - Or use static arrays for fixed-size
  - Generate array helper functions

☐ **String Translation**
  - Use C strings (char*)
  - Handle string concatenation
  - Implement string interpolation
  - Use strdup for string copies

☐ **Memory Management**
  - Generate malloc/free calls
  - Implement reference counting (optional)
  - Or use garbage collection (optional)
  - Handle resource cleanup


### 21. ERROR HANDLING IN GENERATED CODE

☐ **Try-Catch Translation**
  - Use setjmp/longjmp
  - Or generate if-else error checks
  - Generate error handling code

☐ **Null Safety**
  - Generate null checks
  - Handle null pointer dereferences


### 22. OPTIMIZATION

☐ **Constant Folding**
  - Evaluate constant expressions at compile time
  - Example: `x = 5 + 3` → `x = 8`

☐ **Dead Code Elimination**
  - Remove unreachable code

☐ **Inline Small Functions**
  - Inline functions < 50 lines


================================================================================
## COMPILER ARCHITECTURE
================================================================================

### 23. LEXER (TOKENIZER)

☐ **Token Types**
  - Keywords (if, else, for, while, fun, object, etc.)
  - Identifiers (variable names, function names)
  - Literals (numbers, strings, booleans)
  - Operators (+, -, *, /, ==, etc.)
  - Delimiters ({, }, (, ), [, ], ,, :, ;)
  - Comments (skip these)

☐ **Lexer Output**
  - Array of tokens with:
    - Type
    - Value/Lexeme
    - Line number
    - Column number


### 24. PARSER

☐ **AST Node Types**
  - Program
  - Variable Declaration
  - Assignment
  - Binary Operation
  - Unary Operation
  - Function Declaration
  - Function Call
  - Object Declaration
  - If Statement
  - Loop Statement
  - Return Statement
  - Literal (number, string, bool)
  - Identifier

☐ **Parser Output**
  - Abstract Syntax Tree (AST)


### 25. SEMANTIC ANALYZER

☐ **Type Checking**
  - Infer and track types
  - Check type compatibility
  - Generate type errors

☐ **Scope Management**
  - Track variable scopes
  - Check for undefined variables
  - Handle shadowing

☐ **Symbol Table**
  - Track all declarations
  - Function signatures
  - Variable types


### 26. CODE GENERATOR

☐ **Generate C Code**
  - Walk AST
  - Generate corresponding C code
  - Handle type conversions
  - Generate helper functions

☐ **Runtime Library**
  - Link with runtime.c (array functions, string functions, etc.)


================================================================================
## BUILD SYSTEM
================================================================================

### 27. COMPILER INVOCATION

☐ **Command Line**
```bash
  candy build program.candy
  candy run program.candy
  candy compile program.candy -o output
```

☐ **Flags**
  - `-o output` - Specify output file
  - `--debug` - Generate debug symbols
  - `--optimize` - Enable optimizations
  - `--verbose` - Show compilation steps


### 28. LINKING

☐ **Static Linking**
  - Link with C runtime library
  - Link with external C libraries
  - Generate single executable


================================================================================
## TESTING & DEBUGGING
================================================================================

### 29. ERROR MESSAGES

☐ **Helpful Error Messages**
  - Show line and column number
  - Show snippet of problematic code
  - Suggest fixes when possible

☐ **Error Types**
  - Syntax errors
  - Type errors
  - Undefined variable errors
  - Function signature mismatches


### 30. DEBUGGING SUPPORT

☐ **Debug Output**
  - Generate debug symbols
  - Support for gdb/lldb
  - Line number mapping


================================================================================
## SUMMARY: MINIMAL VIABLE COMPILER
================================================================================

**Phase 1 - Core Language (MVP):**
✓ Variables (int, float, string, bool)
✓ Operators (+, -, *, /, ==, <, >, &&, ||)
✓ If/Else
✓ While loops
✓ Functions
✓ Print statement
✓ Basic transpiler (Candy → C)

**Phase 2 - Enhanced Features:**
✓ For loops (range-based)
✓ Arrays
✓ Objects/Structs
✓ String interpolation
✓ Error handling (try/catch)

**Phase 3 - C Interop:**
✓ Extern declarations
✓ Library system
✓ Type annotations
✓ candy-bindgen tool

**Phase 4 - Standard Library:**
✓ Math functions
✓ String functions
✓ File I/O
✓ Utility functions

**Phase 5 - Optimizations:**
✓ Constant folding
✓ Dead code elimination
✓ Inline optimization


================================================================================
## TOTAL FEATURE COUNT
================================================================================

**Core Syntax Features:** ~30
**Built-in Functions:** ~25
**Operators:** ~20
**C Interop Features:** ~10

**Total:** ~85 features to implement