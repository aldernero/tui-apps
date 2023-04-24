# Rule 30

This terminal app implements a [Rule 30](https://en.wikipedia.org/wiki/Rule_30) cellular automaton. The seed is the 
current Unix time in nanoseconds when the app is started.

## Installation

Download one of the binaries from the release page, or:
```bash
go install github.com/aldernero/tui-apps/rule30@latest
```

## Controls

* `Esc`, `Enter`, `Ctrl+c`, `q` - quit
* `Space` - Toggle between wrapped bounds or random bounds (default is random).
* `p` - cycle through color palettes
* `Up/Down` - increase/decrease speed, there are 20 speeds, starts in the middle
* `Left/Right` - increase/decrease the seed, restarts the simulation