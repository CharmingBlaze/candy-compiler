# Standard Library 📚

Candy comes with a built-in library of functions to handle common tasks.

## Core Builtins

These functions are available globally without any imports.

-   `print(args...)`: Prints values joined by spaces, followed by a newline.
-   `println(args...)`: Same as `print`.
-   `printf(format, args...)`: Alias to `println` (useful for C-style habit).
-   `len(value)`: Returns the length of a list or string.
-   `type(value)`: Returns the type name of a value as a string.
-   `sleep(ms)`: Pauses execution for the specified milliseconds.
-   `readLine()`: Reads one line from standard input (no arguments).
-   `parseInt(s)`: Parses a decimal integer from a string (whitespace trimmed).
-   `toUpper(s)` / `toLower(s)`: Same as `string.upper` / `string.lower`.
-   `rand(min, max)`: Inclusive random integer (alias of the `random(min, max)` builtin).

## Math Module (`math`)

Use these for mathematical operations.

-   `math.sqrt(x)`: Square root.
-   `math.pow(base, exp)`: Power.
-   `math.abs(x)`: Absolute value.
-   `math.floor(x)`, `math.ceil(x)`, `math.round(x)`.
-   `math.sin(x)`, `math.cos(x)`, `math.tan(x)`.
-   `math.pi`: The value of π.

## Random Module (`random`)

-   `random.int(min, max)`: Random integer between min and max (inclusive).
-   `random.float(min, max)`: Random float in `[min, max)`.
-   `random.pick(list)`: Returns a random item from a list.
-   `random.shuffle(list)`: Shuffles the list in-place.
-   `random.seed(n)`: Sets the random seed.

## Time Module (`time`)

-   `time.now()`: Returns the current time in milliseconds.
-   `time.sleep(ms)`: Pauses execution.
-   `time.sleep_sec(s)`: Pauses execution for whole/float seconds.

## JSON Module (`json`)

-   `json.parse(string)`: Converts a JSON string to a Candy list or map.
-   `json.stringify(value)`: Converts a Candy value to a JSON string.
-   `json.load(file)`: Loads and parses JSON from a file.
-   `json.save(file, value)`: Saves a value to a file as JSON.

## File Module (`file`)

-   `file.read(path)`: Reads the entire file as a string.
-   `file.write(path, content)`: Writes a string to a file.
-   `file.exists(path)`: Returns true if the file exists.
-   `file.read_lines(path)`: Returns a list of lines from the file.
-   `file.delete(path)`: Deletes a file.
-   `file.list(dir)`: Lists directory entries.

## String Module (`string`)

-   `string.trim(s)`: Removes leading/trailing whitespace.
-   `string.split(s, delim)`: Splits a string into a list.
-   `string.join(list, delim)`: Joins list items into a string.
-   `string.replace(s, old, new)`: Replaces all occurrences.
-   `string.lower(s)`, `string.upper(s)`.
-   `string.starts_with(s, prefix)`, `string.ends_with(s, suffix)`, `string.contains(s, sub)`.

## OS Module (`os`)

-   `os.cwd()`: Returns current working directory.
-   `os.chdir(path)`: Changes process working directory.
-   `os.env(name)`: Returns env variable value or `null`.
-   `os.run(cmd)`: Runs a shell command and returns combined output text.
-   `os.mkdir(path)`: Creates a directory (and parents if needed).
-   `os.rmdir(path)`: Removes a directory recursively.

## Path Module (`path`)

-   `path.join(a, b, ...)`: Joins path parts with platform separators.
-   `path.basename(p)`: Last element of path.
-   `path.dirname(p)`: Parent directory.
-   `path.ext(p)`: File extension (example: `.txt`).
-   `path.normalize(p)`: Cleans `.`/`..` and duplicated separators.

## Collections Module (`collections`)

`collections` provides lightweight constructors using Candy core values.

-   `collections.set([items])`: Returns a map-backed set-like value (`item -> true`).
-   `collections.queue([items])`: Returns list-backed queue storage.
-   `collections.stack([items])`: Returns list-backed stack storage.
-   `collections.deque([items])`: Returns list-backed deque storage.
-   `collections.priority_queue([items])`: Returns a sorted list-backed priority queue seed.

## Color Module (`color`)

-   `color.rgb(r, g, b)`: Creates `{r,g,b,a}` with `a=255`.
-   `color.rgba(r, g, b, a)`: Creates `{r,g,b,a}`.
-   `color.hex("#RRGGBB")` / `color.hex("#RRGGBBAA")`: Parses hex color.
-   `color.lerp(a, b, t)`: Interpolates between two color maps (`t` in `[0,1]`).

## ENET Module (`enet`)

Networking extension with normalized events and peer/packet handles.

-   Lifecycle: `enet.init()`, `enet.deinit()`, `enet.version()`, `enet.backend()`.
-   Address: `enet.address(host, port)` -> `{host,port}` map.
-   Host: `enet.host_create(addressOrNull, peers, channels, inBandwidth, outBandwidth)`, `enet.host_destroy(hostId)`, `enet.host_service(hostId, timeoutMs)`, `enet.host_flush(hostId)`, `enet.host_bandwidth_limit(...)`, `enet.host_channel_limit(...)`, `enet.host_compress_range_coder(hostId)`.
-   Peer: `enet.host_connect(hostId, address, channels, data)`, `enet.peer_disconnect(peerId, data)`, `enet.peer_disconnect_now(peerId, data)`, `enet.peer_ping(peerId)`, `enet.peer_timeout(peerId, limit, min, max)`, `enet.peer_reset(peerId)`.
-   Packet: `enet.packet_create(dataString, flags)`, `enet.packet_destroy(packetId)`, `enet.peer_send(peerId, channel, packetId)`.
-   Event constants: `enet.EVENT_NONE`, `enet.EVENT_CONNECT`, `enet.EVENT_DISCONNECT`, `enet.EVENT_RECEIVE`.
-   Packet flag constants: `enet.PACKET_RELIABLE`, `enet.PACKET_UNSEQUENCED`, `enet.PACKET_NO_ALLOCATE`, `enet.PACKET_UNRELIABLE`.

`enet.host_service(...)` returns:
- `{"type","hostId","peerId","channel","data","packet","address"}`
- `packet` is `{"id","flags","data"}` when `type == "receive"`.
