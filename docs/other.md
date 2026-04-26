plaintext# CANDY LANGUAGE - MISSING FEATURES ANALYSIS
# What Should We Add to Make Candy Complete?

================================================================================
## CRITICAL MISSING FEATURES (MUST HAVE)
================================================================================

### 1. SWITCH/MATCH STATEMENT

**Why it's needed:** Complex conditionals are cleaner with switch

```candy
// Current way (verbose)
if value == 1 {
  doOne()
} else if value == 2 {
  doTwo()
} else if value == 3 {
  doThree()
} else {
  doDefault()
}

// Better way
switch value {
  case 1:
    doOne()
  case 2:
    doTwo()
  case 3:
    doThree()
  default:
    doDefault()
}

// Or even simpler pattern matching
match value {
  1 => doOne()
  2 => doTwo()
  3 => doThree()
  _ => doDefault()
}
```

**Add to checklist:**
☐ Switch/Match statement with fall-through control


### 2. MULTIPLE RETURN VALUES

**Why it's needed:** Common in game dev and C libraries

```candy
// Return multiple values
fun getPosition() {
  return x, y, z
}

// Destructure on assignment
x, y, z = getPosition()

// Or with parentheses (more explicit)
fun getPlayerInfo() {
  return (name, health, score)
}

name, health, score = getPlayerInfo()

// Ignore unwanted values
name, _, score = getPlayerInfo()  // Skip health
```

**Add to checklist:**
☐ Multiple return values
☐ Tuple destructuring
☐ Underscore for ignored values


### 3. DEFAULT PARAMETER VALUES

**Why it's needed:** Makes APIs much friendlier

```candy
// Function with defaults
fun createWindow(width = 800, height = 600, title = "My Game") {
  // code
}

// Call with all defaults
createWindow()

// Call with some overrides
createWindow(1024, 768)
createWindow(title = "Custom Title")
```

**Add to checklist:**
☐ Default parameter values
☐ Named parameters (optional)


### 4. RANGE OPERATOR

**Why it's needed:** Cleaner array slicing and iteration

```candy
// Create range
numbers = 1..10        // [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
numbers = 1..<10       // [1, 2, 3, 4, 5, 6, 7, 8, 9] (exclusive)

// Array slicing
items = [10, 20, 30, 40, 50]
subset = items[1..3]   // [20, 30, 40]

// Iterate over range
for i in 1..10 {
  print(i)
}

// Step range
for i in 0..100 step 10 {  // 0, 10, 20, 30...
  print(i)
}
```

**Add to checklist:**
☐ Range operator (..)
☐ Exclusive range (..<)
☐ Array slicing


### 5. STRING METHODS (NOT JUST FUNCTIONS)

**Why it's needed:** More intuitive object-oriented style

```candy
// Current (function style)
name = upper("hello")
length = len("world")

// Better (method style)
name = "hello".upper()
length = "world".length

// Chaining
result = "  hello  ".trim().upper().replace("H", "J")
// "JELLO"

// Both styles should work
upper("hello")      // Function style
"hello".upper()     // Method style
```

**Add to checklist:**
☐ Method syntax for strings
☐ Method chaining
☐ Additional string methods:
  - .trim() - Remove whitespace
  - .startsWith(prefix)
  - .endsWith(suffix)
  - .indexOf(substring)
  - .substring(start, end)


### 6. ARRAY/LIST METHODS

**Why it's needed:** Modern languages expect these

```candy
numbers = [1, 2, 3, 4, 5]

// Functional operations
doubled = numbers.map(x => x * 2)          // [2, 4, 6, 8, 10]
evens = numbers.filter(x => x % 2 == 0)    // [2, 4]
sum = numbers.reduce((a, b) => a + b, 0)   // 15

// Checking
hasThree = numbers.contains(3)              // true
allPositive = numbers.all(x => x > 0)      // true
anyEven = numbers.any(x => x % 2 == 0)     // true

// Finding
found = numbers.find(x => x > 3)           // 4
index = numbers.indexOf(3)                 // 2

// Sorting
sorted = numbers.sort()                    // Ascending
sorted = numbers.sort((a, b) => b - a)     // Descending

// Transforming
reversed = numbers.reverse()
unique = [1, 2, 2, 3, 3].unique()         // [1, 2, 3]
```

**Add to checklist:**
☐ .map() - Transform each element
☐ .filter() - Keep matching elements
☐ .reduce() - Combine to single value
☐ .contains() - Check if element exists
☐ .find() - Find first matching
☐ .indexOf() - Get index of element
☐ .all() - Check if all match
☐ .any() - Check if any match
☐ .sort() - Sort array
☐ .reverse() - Reverse array
☐ .unique() - Remove duplicates
☐ .join() - Join to string


### 7. DEFER STATEMENT (RESOURCE CLEANUP)

**Why it's needed:** Simpler than 'with', more flexible

```candy
fun processFile() {
  file = openFile("data.txt")
  defer closeFile(file)  // Always runs at function end
  
  texture = loadTexture("img.png")
  defer unloadTexture(texture)
  
  // Even if error occurs, defers still run
  if errorCondition {
    return  // Cleanup happens automatically
  }
  
  // Normal processing
  processData(file)
}
// closeFile and unloadTexture called here automatically
```

**Add to checklist:**
☐ Defer statement
☐ LIFO order execution (last defer runs first)


### 8. ASSERTION/DEBUG TOOLS

**Why it's needed:** Catch bugs early in development

```candy
// Assert conditions
assert(player.health > 0, "Player health must be positive")
assert(score >= 0)

// Debug-only code
debug {
  print("Debug mode: player at {player.x}, {player.y}")
  drawDebugBox(player.x, player.y, 32, 32)
}

// Compile-time checks
static_assert(SCREEN_WIDTH > 0, "Screen width must be positive")
```

**Add to checklist:**
☐ assert() function
☐ debug { } blocks (removed in release builds)
☐ static_assert for compile-time checks


### 9. MODULES/NAMESPACES

**Why it's needed:** Organize larger projects

```candy
// Define module
module Math {
  fun add(a, b) {
    return a + b
  }
  
  fun multiply(a, b) {
    return a * b
  }
  
  const PI = 3.14159
}

// Use module
result = Math.add(5, 3)
circle = 2 * Math.PI * radius

// Import specific items
from Math import add, PI

result = add(5, 3)
```

**Add to checklist:**
☐ Module declaration
☐ Module member access
☐ Selective imports


### 10. ELVIS OPERATOR / NULL COALESCING

**Why it's needed:** Cleaner null handling

```candy
// Current
value = getValue()
if value == null {
  value = defaultValue
}

// Better with 'or' (already have this)
value = getValue() or defaultValue

// But also need null-safe access
name = player?.name or "Unknown"    // If player is null, skip .name
health = player?.stats?.health or 100

// Safe call chain
result = obj?.method1()?.method2()?.value
```

**Add to checklist:**
☐ Safe navigation operator (?.)
☐ Null coalescing is already covered by 'or'


================================================================================
## NICE-TO-HAVE FEATURES (OPTIONAL)
================================================================================

### 11. LIST COMPREHENSIONS

```candy
// Instead of
evens = []
for i in numbers {
  if i % 2 == 0 {
    evens.add(i)
  }
}

// Simpler
evens = [i for i in numbers if i % 2 == 0]

// Nested
pairs = [(x, y) for x in 1..3 for y in 1..3]
```

**Add to checklist (optional):**
☐ List comprehensions


### 12. SPREAD OPERATOR

```candy
// Spread array
arr1 = [1, 2, 3]
arr2 = [4, 5, 6]
combined = [...arr1, ...arr2]  // [1, 2, 3, 4, 5, 6]

// Spread in function calls
numbers = [1, 2, 3]
result = add(...numbers)  // Same as add(1, 2, 3)

// Spread in objects
defaults = {x: 0, y: 0}
player = {...defaults, health: 100}
```

**Add to checklist (optional):**
☐ Spread operator (...)


### 13. TERNARY OPERATOR

```candy
// Instead of
if score > 100 {
  message = "High Score!"
} else {
  message = "Keep trying"
}

// Simpler
message = score > 100 ? "High Score!" : "Keep trying"
```

**Add to checklist (optional):**
☐ Ternary operator (? :)


### 14. IN OPERATOR

```candy
// Check if value in array
if 3 in [1, 2, 3, 4, 5] {
  print("Found!")
}

// Check if key in object
if "name" in player {
  print(player.name)
}

// Not in
if value not in excludeList {
  process(value)
}
```

**Add to checklist (optional):**
☐ 'in' operator
☐ 'not in' operator


### 15. GLOBAL KEYWORD (FOR CLOSURES)

```candy
score = 0

fun incrementScore() {
  global score  // Modify outer variable
  score = score + 1
}

// Or implicit (capture by reference)
fun incrementScore() {
  score = score + 1  // Automatically captures
}
```

**Add to checklist (optional):**
☐ Global keyword for clarity
☐ Automatic closure capture


### 16. STATIC VARIABLES (PERSIST ACROSS CALLS)

```candy
fun getID() {
  static counter = 0  // Initialized once
  counter = counter + 1
  return counter
}

id1 = getID()  // 1
id2 = getID()  // 2
id3 = getID()  // 3
```

**Add to checklist (optional):**
☐ Static local variables


### 17. GOTO (FOR LOW-LEVEL CONTROL)

```candy
// Only for specific use cases (state machines, cleanup)
retry:
  result = attemptConnection()
  if result.error {
    retries = retries + 1
    if retries < MAX_RETRIES {
      goto retry
    }
  }

cleanup:
  closeConnection()
  freeResources()
```

**Add to checklist (optional, controversial):**
☐ Goto statement (label + goto)


================================================================================
## UPDATED COMPLETE CHECKLIST
================================================================================

### CRITICAL ADDITIONS:

☐ **Switch/Match Statement**
  - Pattern matching
  - Fall-through control
  - Default case

☐ **Multiple Return Values**
  - Tuple syntax
  - Destructuring assignment
  - Underscore for ignored values

☐ **Default Parameters**
  - Optional named parameters
  - Parameter defaults

☐ **Range Operator**
  - Inclusive range (1..10)
  - Exclusive range (1..<10)
  - Array slicing [arr[1..3]]
  - Range iteration

☐ **String Methods**
  - .upper(), .lower()
  - .trim(), .split()
  - .startsWith(), .endsWith()
  - .indexOf(), .substring()
  - .replace(), .contains()
  - Method chaining

☐ **Array Methods**
  - .map(), .filter(), .reduce()
  - .contains(), .indexOf()
  - .find(), .all(), .any()
  - .sort(), .reverse()
  - .unique(), .join()

☐ **Defer Statement**
  - Automatic cleanup
  - LIFO execution order

☐ **Assert/Debug Tools**
  - assert() function
  - debug { } blocks
  - Compile-time assertions

☐ **Modules/Namespaces**
  - Module declaration
  - Selective imports
  - Namespace access

☐ **Safe Navigation**
  - ?. operator for null-safe access
  - Null coalescing (already have 'or')


### OPTIONAL ADDITIONS (NICE TO HAVE):

☐ List comprehensions
☐ Spread operator (...)
☐ Ternary operator (? :)
☐ 'in' operator
☐ Global keyword
☐ Static variables
☐ Goto statement


================================================================================
## REVISED FEATURE COUNT
================================================================================

**Original Core Features:** 85
**Critical Additions:** +35
**Optional Additions:** +20

**New Total:** ~140 features for a complete modern language

**Revised Implementation Time:**
- Phase 1 (MVP): 2-4 weeks
- Phase 2 (Enhanced): 2-3 weeks
- Phase 3 (C Interop): 2-3 weeks
- Phase 4 (Critical Additions): 3-4 weeks  ← NEW
- Phase 5 (Standard Library): 1-2 weeks
- Phase 6 (Optimizations): 1-2 weeks
- Phase 7 (Optional Features): 2-3 weeks  ← NEW

**Total: 13-21 weeks for complete implementation**


================================================================================
## PRIORITY RANKING
================================================================================

**HIGHEST PRIORITY (Implement First):**
1. Switch/Match statement - Essential for clean code
2. String methods - Makes strings usable
3. Array methods (.map, .filter, etc.) - Modern expectation
4. Default parameters - API friendliness
5. Defer statement - Resource management

**HIGH PRIORITY:**
6. Multiple return values - Common pattern
7. Range operator - Very useful
8. Assert/Debug tools - Development aid
9. Safe navigation (?.) - Null safety

**MEDIUM PRIORITY:**
10. Modules - Code organization
11. Ternary operator - Convenience
12. 'in' operator - Readability
13. List comprehensions - Elegance

**LOW PRIORITY:**
14. Spread operator - Nice but not critical
15. Global keyword - Can work around
16. Static variables - Niche use case
17. Goto - Controversial, rarely needed


================================================================================
## RECOMMENDATION
================================================================================

**For MVP (Phase 1-3):** Keep original 85 features

**For Production-Ready (Phase 4):** Add these 10 critical features:
1. Switch/Match
2. String methods
3. Array methods (.map, .filter, .reduce, .contains)
4. Default parameters
5. Defer
6. Multiple return values
7. Range operator
8. Assert
9. Safe navigation (?.)
10. Ternary operator

**This gives you a modern, complete language that:**
- ✓ Is still simple enough for beginners
- ✓ Has features programmers expect
- ✓ Can build real applications
- ✓ Interops well with C
- ✓ Feels modern and productive

**Total MVP + Production: ~95 features, 13-14 weeks**

🍬 **Candy would be competitive with modern scripting languages!** 🍬