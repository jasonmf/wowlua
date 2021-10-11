# wowlua

[![Go Reference](https://pkg.go.dev/badge/github.com/jasonmf/wowlua.svg)](https://pkg.go.dev/github.com/jasonmf/wowlua)

Some time around 2015 I wrote tooling to synchronize our WoW guild calendar
to a Google Calendar. This had two parts: one was a WoW UI addon that would retrieve
the guild calendar data and ensure it was written out to disk as Lua data.
The other would read that data and synchronize it to Google Calendar. This
package was written to support the latter. I don't have any of the othe code.

I'm not sure that this package entirely works and I haven't really touched it
since 2015.

It's simple enough to use:

```
table, err := wowlua.ParseLua(string(luaTableDataByteSlice))
```

The top-level data structure must be a table. The package only parses into
types defined by this package, not arbitrary Go types like `encoding/json`. To
use the table you must either node how data is stored with in it or be willing
to inspect the keys and check node types.