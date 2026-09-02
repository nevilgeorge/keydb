# KeyDB

A small append-only key/value store in Go, with an interactive REPL.

Records are appended to a log file as `index<TAB>key<TAB>value` lines. An in-memory
map tracks the byte offset of each key's most recent record, so a read is a single
`Seek` plus one line read.

## Requirements

- Go 1.25.1 or newer

## Running

The database file lives at `db/local.db`, so make sure the `db/` directory exists:

```sh
mkdir -p db
go run .
```

Or build a binary:

```sh
go build -o keydb .
./keydb
```

## Commands

| Command            | Description                              |
| ------------------ | ---------------------------------------- |
| `put <key>:<value>` | Append a record and index its offset     |
| `get <key>`         | Print the value for a key                |
| `exit` / `quit`     | Close the file handles and exit          |

Example session:

```
Hi from KeyDB.
> put name:nevil
Record stored: name: nevil
> get name
nevil
> exit
```

## Layout

- `main.go` — REPL loop and command parsing
- `record_store.go` — `RecordStore` interface plus the file-backed reader/writer

## Limitations

- The key/offset index is only built as you write, so keys from earlier runs are
  not readable until the log is replayed on startup.
- Keys and values cannot contain `:` (the command delimiter) or tabs and newlines
  (the record delimiters).
- No deletes and no log compaction — every `put` grows the file.
