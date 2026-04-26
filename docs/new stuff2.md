plaintext# CANDY LANGUAGE - ULTIMATE SIMPLIFICATION & POWER ADDITIONS
# Making Programming Feel Like Magic

================================================================================
## PART 1: RADICAL SIMPLIFICATION IDEAS
================================================================================

### 1. IMPLICIT RETURNS (Last Expression)

**Problem:** Always typing `return`
```candy
// Current
fun add(a, b) {
  return a + b
}

// Simpler - last expression is auto-returned
fun add(a, b) {
  a + b
}

// Still allow explicit return for early exit
fun findItem(array, value) {
  for item in array {
    if item == value {
      return item  // Early exit
    }
  }
  null  // Implicit return
}
```

**Add to checklist:**
☐ Implicit return from last expression
☐ Explicit `return` for early exit


### 2. PIPELINE OPERATOR (Chain Operations)

**Problem:** Nested function calls are hard to read
```candy
// Current (hard to read)
result = round(sqrt(abs(add(5, multiply(3, 2)))))

// Better (left to right)
result = 5
  |> add(multiply(3, 2))
  |> abs()
  |> sqrt()
  |> round()

// Real example
enemies
  |> filter(e => e.alive)
  |> map(e => e.position)
  |> sort((a, b) => a.x - b.x)
  |> first()
```

**Add to checklist:**
☐ Pipeline operator `|>`
☐ Passes result of left as first argument to right


### 3. OPTIONAL CHAINING (Deep Safe Access)

**Problem:** Multiple null checks are verbose
```candy
// Current
if player != null {
  if player.inventory != null {
    if player.inventory.weapon != null {
      damage = player.inventory.weapon.damage
    }
  }
}

// Better
damage = player?.inventory?.weapon?.damage or 0
```

**Already mentioned, but emphasize:**
☐ Safe navigation `?.`
☐ Works with arrays: `arr?.[index]`
☐ Works with calls: `obj?.method?.()`


### 4. PATTERN MATCHING (More Powerful Switch)

**Problem:** Switch is limited to simple values
```candy
// Current switch
switch playerState {
  case "idle":
    // ...
  case "running":
    // ...
}

// Better - pattern matching
match player {
  {health: 0} => gameOver()
  {health: h} if h < 20 => showWarning()
  {position: {y}} if y < -100 => respawn()
  {velocity: v} if v.length() > 10 => showSpeedLines()
  _ => // default
}

// Match on types
match value {
  int => print("Number: {value}")
  string => print("Text: {value}")
  [first, ...rest] => print("Array starting with {first}")
  {x, y} => print("Point at {x}, {y}")
  _ => print("Something else")
}
```

**Add to checklist:**
☐ Pattern matching with `match`
☐ Destructuring in patterns
☐ Guard clauses with `if`
☐ Type matching
☐ Rest patterns `...rest`


### 5. AUTO-VIVIFICATION (Create on Access)

**Problem:** Initializing nested structures is tedious
```candy
// Current
if !enemies[level] {
  enemies[level] = []
}
enemies[level].add(newEnemy)

// Better - auto-create on access
enemies[level] ||= []
enemies[level].add(newEnemy)

// Or even simpler - auto-vivify arrays
enemies[level].add(newEnemy)  // Creates array if needed
```

**Add to checklist:**
☐ `||=` operator (assign if null)
☐ Auto-vivification for arrays/objects


### 6. TEMPLATE LITERALS (Multi-line Strings)

**Problem:** Multi-line strings are awkward
```candy
// Current
message = "Line 1\n" +
          "Line 2\n" +
          "Line 3"

// Better with template literals
message = ```
  Line 1
  Line 2
  Line 3
```

// With interpolation
instructions = ```
  Player Health: {player.health}
  Position: ({player.x}, {player.y})
  Status: {player.status}

**Add to checklist:**
☐ Triple-backtick ``` for multi-line strings
☐ Auto-trim leading/trailing whitespace
☐ Support interpolation


### 7. SHORTHAND PROPERTY SYNTAX

**Problem:** Repeating variable names
```candy
// Current
player = {
  x: x,
  y: y,
  health: health,
  name: name
}

// Better - shorthand
player = {x, y, health, name}

// Mixed shorthand + explicit
player = {
  x, y,  // Shorthand
  health: maxHealth,  // Explicit
  name
}
```

**Add to checklist:**
☐ Shorthand object properties
☐ Mixed shorthand + explicit


### 8. RANGE-BASED OPERATIONS (More Powerful)

**Problem:** Limited range features
```candy
// Array slicing (already mentioned)
subset = arr[1..5]
last3 = arr[-3..]  // Last 3 elements
allButFirst = arr[1..]  // Skip first
allButLast = arr[..-1]  // Skip last

// String slicing
substr = text[0..5]  // First 6 characters
lastChar = text[-1]  // Last character

// In-place modification
arr[2..4] = [10, 20, 30]  // Replace range
```

**Add to checklist:**
☐ Negative indices (from end)
☐ Open-ended ranges (1.., ..-1)
☐ String slicing
☐ Range assignment


### 9. INFINITY LITERAL

**Problem:** No way to represent infinity
```candy
// Add infinity
maxValue = infinity
minValue = -infinity

// Use cases
clampedValue = clamp(value, -infinity, maxValue)

// Check
if distance == infinity {
  print("Unreachable")
}
```

**Add to checklist:**
☐ `infinity` keyword
☐ `-infinity` for negative infinity
☐ `isInfinite()` helper


### 10. LABEL SYNTAX (Named Blocks)

**Problem:** Can't break from nested loops
```candy
// Current - need flags
found = false
for row in grid {
  for col in row {
    if col == target {
      found = true
      break
    }
  }
  if found { break }
}

// Better - labeled break
outer: for row in grid {
  for col in row {
    if col == target {
      break outer
    }
  }
}

// Works with while too
searching: while true {
  for item in items {
    if done {
      break searching
    }
  }
}
```

**Add to checklist:**
☐ Label syntax `label: loop`
☐ `break label` to break from labeled loop
☐ `continue label` to continue labeled loop


================================================================================
## PART 2: CONVENIENCE ADDITIONS
================================================================================

### 11. INLINE CONDITIONALS (Guard Clauses)

```candy
// Early return with guard
fun process(value) {
  guard value != null else return
  guard value > 0 else return
  
  // Process value
}

// Or simpler
fun process(value) {
  return unless value != null
  return unless value > 0
  
  // Process value
}
```

**Add to checklist:**
☐ `guard condition else action`
☐ `return unless condition`
☐ `continue if condition`
☐ `break if condition`


### 12. WITH EXPRESSION (Scope Builder)

```candy
// Create and initialize in one go
player = Player() with {
  x = 100
  y = 200
  health = 100
  name = "Hero"
}

// Equivalent to
player = Player()
player.x = 100
player.y = 200
player.health = 100
player.name = "Hero"
```

**Add to checklist:**
☐ `with` expression for initialization


### 13. CASE EXPRESSIONS (Not Just Statements)

```candy
// Use match as expression
color = match health {
  h if h > 75 => "green"
  h if h > 25 => "yellow"
  _ => "red"
}

// One-liner
message = match state {
  "win" => "Victory!"
  "lose" => "Defeat!"
  _ => "Playing..."
}
```

**Add to checklist:**
☐ Match/switch as expressions (return values)


### 14. LAZY EVALUATION

```candy
// Don't evaluate unless needed
lazy expensiveResult = calculateComplexThing()

// Only computed when accessed
if needResult {
  print(expensiveResult)  // Calculated here
}

// Useful for default parameters
fun loadTexture(path, fallback = lazy defaultTexture()) {
  // fallback only created if needed
}
```

**Add to checklist:**
☐ `lazy` keyword for deferred evaluation


### 15. TUPLE TYPES (Lightweight Data)

```candy
// Tuples - immutable, ordered collections
position = (100, 200, 0)
color = ("red", 255, 0, 0)

// Access by index
x = position[0]
y = position[1]

// Destructure
(x, y, z) = position
(name, r, g, b) = color

// Return multiple values (already mentioned, but use tuples)
fun getMinMax(array) {
  (min(array), max(array))
}

min, max = getMinMax(numbers)
```

**Add to checklist:**
☐ Tuple syntax with parentheses
☐ Immutable by default
☐ Index access
☐ Destructuring


### 16. SET TYPE (Unique Collections)

```candy
// Sets - unique values only
visited = set()
visited.add(position)
visited.add(position)  // No duplicates
print(visited.size)  // 1

// Set literals
primes = {2, 3, 5, 7, 11}

// Set operations
a = {1, 2, 3}
b = {2, 3, 4}
union = a | b        // {1, 2, 3, 4}
intersection = a & b // {2, 3}
difference = a - b   // {1}

// Contains check (fast)
if 3 in primes {
  print("Is prime")
}
```

**Add to checklist:**
☐ Set type with unique values
☐ Set literals `{1, 2, 3}`
☐ Set operations (|, &, -)
☐ Fast `in` operator


### 17. MAP/DICTIONARY METHODS

```candy
// Map with any key type
scores = map()
scores[player1] = 100
scores[player2] = 95

// Or map literal
config = {
  "width": 800,
  "height": 600,
  "title": "My Game"
}

// Methods
config.keys()     // ["width", "height", "title"]
config.values()   // [800, 600, "My Game"]
config.entries()  // [[key, value], ...]
config.has("width")  // true

// Iterate
for key, value in config {
  print("{key}: {value}")
}
```

**Add to checklist:**
☐ Map type (already have objects, but explicit map)
☐ .keys(), .values(), .entries()
☐ .has(key) method
☐ Foreach with key, value


### 18. ASYNC/AWAIT (Optional, Advanced)

```candy
// For loading assets, network calls
async fun loadAssets() {
  texture = await loadTexture("player.png")
  sound = await loadSound("jump.wav")
  music = await loadMusic("theme.mp3")
  return {texture, sound, music}
}

// Use
assets = await loadAssets()
print("Assets loaded!")

// Or with error handling
try {
  assets = await loadAssets()
} catch error {
  print("Failed to load: {error}")
}
```

**Add to checklist (OPTIONAL):**
☐ `async` keyword for async functions
☐ `await` for waiting on async operations
☐ Promise-like behavior


### 19. FIRST-CLASS FUNCTIONS (More Features)

```candy
// Partial application
fun add(a, b) { a + b }
add5 = add(5, _)  // Partial - leave second arg
result = add5(10)  // 15

// Function composition
double = x => x * 2
increment = x => x + 1
doubleAndIncrement = double >> increment
result = doubleAndIncrement(5)  // 11

// Or reverse composition
incrementAndDouble = double << increment
result = incrementAndDouble(5)  // 12
```

**Add to checklist:**
☐ Partial application with `_`
☐ Function composition `>>` and `<<`


### 20. CASTING SHORTCUTS

```candy
// Current
x = toInt("123")
y = toFloat("3.14")
s = toString(42)

// Better - type as function
x = int("123")
y = float("3.14")
s = string(42)

// Also works for checking
if value is int {
  print("Is a number")
}

// Type assertions
num = value as int  // Error if not int
num = value as? int  // null if not int
```

**Add to checklist:**
☐ Type as function: `int()`, `float()`, etc.
☐ Type checking: `is` operator
☐ Type casting: `as` operator
☐ Safe casting: `as?` operator


================================================================================
## PART 3: QUALITY OF LIFE FEATURES
================================================================================

### 21. DO EXPRESSION

```candy
// Execute block and return last value
result = do {
  x = calculate()
  y = transform(x)
  finalize(y)  // This value returned
}

// Useful for complex initialization
config = do {
  base = loadDefaults()
  base.width = 1920
  base.height = 1080
  base  // Return modified
}
```

**Add to checklist:**
☐ `do { }` expression returns last value


### 22. SINGLE-EXPRESSION FUNCTIONS

```candy
// Super short syntax for simple functions
add = (a, b) => a + b
double = x => x * 2
isPositive = x => x > 0

// No braces needed for single expression
numbers.map(x => x * 2)
numbers.filter(x => x > 0)
```

**Add to checklist:**
☐ Single-expression lambdas without braces


### 23. STRING METHODS (Complete Set)

```candy
// Add missing string methods
text = "  Hello World  "

text.trim()           // "Hello World"
text.trimLeft()       // "Hello World  "
text.trimRight()      // "  Hello World"
text.padLeft(20)      // "       Hello World  "
text.padRight(20)     // "  Hello World       "
text.repeat(3)        // "  Hello World    Hello World    Hello World  "
text.reverse()        // "  dlroW olleH  "
text.chars()          // ['  ', 'H', 'e', 'l', ...]

// Character access
text[0]               // " "
text[-1]              // " "
```

**Add to checklist:**
☐ .trim(), .trimLeft(), .trimRight()
☐ .padLeft(), .padRight()
☐ .repeat()
☐ .reverse()
☐ .chars() - return array of characters
☐ Character indexing


### 24. ARRAY METHODS (Complete Set)

```candy
// Add missing array methods
numbers = [1, 2, 3, 4, 5]

numbers.first()           // 1
numbers.last()            // 5
numbers.take(3)           // [1, 2, 3]
numbers.drop(2)           // [3, 4, 5]
numbers.takeWhile(x => x < 4)  // [1, 2, 3]
numbers.dropWhile(x => x < 3)  // [3, 4, 5]
numbers.partition(x => x % 2 == 0)  // [[2, 4], [1, 3, 5]]
numbers.chunk(2)          // [[1, 2], [3, 4], [5]]
numbers.flatten()         // For nested arrays
numbers.flatMap(x => [x, x * 2])  // [1, 2, 2, 4, 3, 6, ...]
numbers.zip([10, 20, 30]) // [[1, 10], [2, 20], [3, 30]]
```

**Add to checklist:**
☐ .first(), .last()
☐ .take(n), .drop(n)
☐ .takeWhile(fn), .dropWhile(fn)
☐ .partition(fn)
☐ .chunk(n)
☐ .flatten()
☐ .flatMap(fn)
☐ .zip(otherArray)


### 25. RANGE METHODS

```candy
// Ranges as first-class values
r = 1..10
r.contains(5)      // true
r.step(2)          // 1, 3, 5, 7, 9
r.reverse()        // 10, 9, 8, ..., 1
r.toArray()        // [1, 2, 3, ..., 10]

// Generate sequences
evens = (0..).step(2).take(10)  // [0, 2, 4, ..., 18]
```

**Add to checklist:**
☐ Range as value
☐ .contains(), .step(), .reverse(), .toArray()
☐ Infinite ranges with `..`


### 26. CONSOLE/DEBUG HELPERS

```candy
// Better debugging
debug(player)  // Pretty-print with type info
dump(enemies)  // Deep inspection with nested data

// Conditional debug
assert(health > 0, "Health must be positive")
debug(position) if DEBUG_MODE

// Trace execution
trace("Entered combat") {
  // Code here logged with timing
  fight()
}
```

**Add to checklist:**
☐ debug() - pretty print
☐ dump() - deep inspection
☐ assert() - runtime assertions
☐ trace() - execution timing


### 27. RANDOM HELPERS

```candy
// Current
dice = random(1, 6)

// More helpers
randomFloat(0.0, 1.0)  // Float between 0 and 1
randomBool()            // true or false
randomChoice([1, 2, 3]) // Pick random from array
shuffle(array)          // Randomize order

// Seeded random (reproducible)
rng = Random(seed: 12345)
val1 = rng.next(1, 6)
val2 = rng.nextFloat()
```

**Add to checklist:**
☐ randomFloat(), randomBool()
☐ randomChoice(array)
☐ shuffle(array)
☐ Seeded Random class


### 28. COLOR TYPE

```candy
// Built-in color type
red = Color(255, 0, 0)
blue = Color.fromHex("#0000FF")
green = Color.fromHSL(120, 1, 0.5)

// Operations
darker = red.darken(0.2)
lighter = red.lighten(0.3)
mixed = red.mix(blue, 0.5)  // Purple

// Properties
r, g, b, a = red.rgba()
h, s, l = red.hsl()
hex = red.toHex()  // "#FF0000"
```

**Add to checklist:**
☐ Color type
☐ .fromHex(), .fromHSL()
☐ .darken(), .lighten(), .mix()
☐ .rgba(), .hsl(), .toHex()


### 29. TIME/TIMER HELPERS

```candy
// Built-in timers
timer = Timer(5.0)  // 5 second timer

loop {
  timer.update(dt)
  
  if timer.finished {
    spawnEnemy()
    timer.reset()
  }
  
  print("Time left: {timer.remaining}")
}

// Stopwatch
watch = Stopwatch()
watch.start()
// ... do work ...
watch.stop()
print("Took: {watch.elapsed} seconds")
```

**Add to checklist:**
☐ Timer class
☐ Stopwatch class
☐ .update(), .finished, .remaining, .reset()
☐ .start(), .stop(), .elapsed


### 30. MATH EXTENSIONS

```candy
// Add useful math
lerp(0, 100, 0.5)      // 50
inverseLerp(0, 100, 50)  // 0.5
remap(5, 0, 10, 0, 100)  // 50 (map from one range to another)
smoothstep(0, 1, 0.5)    // Smooth interpolation
ping pong(value, length) // Bounce value between 0 and length

// Angle helpers
degrees(radians)
radians(degrees)
normalizeAngle(angle)  // Wrap to 0-360 or -180 to 180

// Vector shortcuts
vec2.fromAngle(radians)
vec2.fromAngle(radians, length)
```

**Add to checklist:**
☐ inverseLerp(), remap()
☐ smoothstep(), pingpong()
☐ degrees(), radians()
☐ normalizeAngle()
☐ vec2.fromAngle()


================================================================================
## SUMMARY OF ALL NEW ADDITIONS
================================================================================

**SIMPLIFICATION FEATURES:**
1. ✓ Implicit returns
2. ✓ Pipeline operator |>
3. ✓ Optional chaining ?.
4. ✓ Pattern matching
5. ✓ Auto-vivification
6. ✓ Template literals ```
7. ✓ Shorthand properties
8. ✓ Enhanced ranges
9. ✓ Infinity literal
10. ✓ Labeled loops

**CONVENIENCE FEATURES:**
11. ✓ Guard clauses
12. ✓ With expression
13. ✓ Match as expression
14. ✓ Lazy evaluation
15. ✓ Tuples
16. ✓ Sets
17. ✓ Map methods
18. ○ Async/await (optional)
19. ✓ Partial application
20. ✓ Type casting shortcuts

**QUALITY OF LIFE:**
21. ✓ Do expressions
22. ✓ Single-expression lambdas
23. ✓ Complete string methods
24. ✓ Complete array methods
25. ✓ Range methods
26. ✓ Debug helpers
27. ✓ Random helpers
28. ✓ Color type
29. ✓ Timer/Stopwatch
30. ✓ Math extensions

**NEW TOTAL FEATURES: ~205**

**But Candy is SIMPLER because:**
- Less boilerplate (implicit returns, shorthand)
- More intuitive (pipeline, pattern matching)
- Better defaults (auto-vivification, safe navigation)
- Complete standard library (no missing gaps)

🍬 **Candy: Maximum power, minimum complexity!** 🍬