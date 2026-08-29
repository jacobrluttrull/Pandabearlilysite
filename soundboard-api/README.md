# soundboard-api

Go API + SQLite behind the site's soundboard, plus the `cli` tool used to manage clips.

Run everything from **inside this directory** — the default paths (`data/soundboard.db`,
`clips`) are relative, so running from the repo root puts the database somewhere
unexpected.

```powershell
cd "E:\Coding Projects\pandabearlily\soundboard-api"
```

## Everyday tasks

**Add one new clip** — the common case. Length is measured from the file and the label is
derived from the filename, so this is all you need:

```powershell
go run .\cmd\cli upload "C:\path\to\new-clip.mp3"
```

Override either default when you want to:

```powershell
go run .\cmd\cli upload "C:\path\to\clip.mp3" -name "the good one" -date-made 2026-08-27
```

**Add a batch of clips** — drop them in a folder, then:

```powershell
go run .\cmd\cli names -dir "..\Pandalily Soundbites"   # seed labels for the new files
notepad names.json                                       # fix any that read badly
go run .\cmd\cli import -dir "..\Pandalily Soundbites"   # copy + insert
```

`import` is safe to re-run: clips already stored are skipped, so an interrupted run just
needs running again.

**Fix a label:**

```powershell
go run .\cmd\cli rename "clip.mp3" "a better name"      # one clip
# ...or edit names.json for several, then:
go run .\cmd\cli apply-names -dry-run                   # preview
go run .\cmd\cli apply-names                            # apply
```

**Clear duplicates** — catches byte-identical copies (`clip (1).mp3`) that imported as
separate clips. Play counts from the copies are merged into the survivor:

```powershell
go run .\cmd\cli dedupe -dry-run
go run .\cmd\cli dedupe
```

Delete the duplicate *source* files from your clips folder afterwards, or the next
import adds them back.

**Before launching** — audit for drift, then zero out any test plays:

```powershell
go run .\cmd\cli check
go run .\cmd\cli reset-plays -all
```

## All commands

| Command | What it does |
| --- | --- |
| `upload <file>` | Add one clip. Measures length, derives the label. |
| `import -dir <folder>` | Bulk-add a folder of `.mp3`s. Idempotent. |
| `names -dir <folder>` | Seed `names.json` with derived labels. Never overwrites existing entries. |
| `apply-names` | Push edited `names.json` into the database. |
| `rename <file> "<name>"` | Relabel one clip; also records it in `names.json`. |
| `set-date <file> <YYYY-MM-DD>` | Record when a clip was made. `-clear` removes it. |
| `dedupe` | Find byte-identical clips, keep one, merge its copies' play counts. |
| `remove <file>...` | Delete clips — row, audio file, and `names.json` entry. |
| `reset-plays -all` | Zero play tallies. Or pass filenames to reset just those. |
| `check` | Audit database vs. audio files vs. `names.json` for drift. |
| `list` | Show every stored clip with length, plays, and date. |

Most commands take `-dry-run`. All take `-h`.

## How clips are stored

- **Database**: the `soundbites` table. Locally that lives in `data/soundboard.db`,
  which is gitignored. In production it is a Turso database instead — see
  [Connecting to Turso](#connecting-to-turso). Same table either way; only the DSN
  changes.
- **Audio**: `clips/`, copied in from wherever you staged it. **Committed** — it ships
  inside the Docker image, and it is the only backup of the audio that exists.
- **Labels**: `names.json`, a `filename -> display name` map. **Committed** — it is
  hand-written content, and losing it means retyping every label.

Filenames are the key across all three. A clip whose filename does not derive into a
readable label (`asseatenbythesebitches.mp3`) gets its real name from `names.json`;
anything not listed there falls back to the derived name.

Your staged source files are never modified or deleted by any command.

## Schema changes

Schema lives in goose migrations under `internal/db/migrations/`; queries are compiled by
sqlc from `internal/db/queries/` into `internal/db/gen/`. After editing a query or
migration:

```powershell
sqlc generate
```

Migrations run automatically whenever the API or CLI opens the database — there is no
separate migrate step.

Clip data is **not** migration material. Adding, renaming, or removing clips is data, and
goes through the commands above.

## Running the API

```powershell
go run .\cmd\api        # listens on :8080
```

| Variable | Default | Purpose |
| --- | --- | --- |
| `SOUNDBOARD_DB_PATH` | `data/soundboard.db` | Local SQLite file. Used only when `TURSO_DATABASE_URL` is empty. |
| `TURSO_DATABASE_URL` | *(empty)* | Remote libSQL database, e.g. `libsql://soundboard-you.turso.io`. Set it and `SOUNDBOARD_DB_PATH` is ignored. |
| `TURSO_AUTH_TOKEN` | *(empty)* | Credential for `TURSO_DATABASE_URL`. Never logged. |
| `SOUNDBOARD_AUDIO_DIR` | `clips` | Where clip audio lives |
| `SOUNDBOARD_ADDR` | `:8080` | Listen address. Ignored when `PORT` is set. |
| `PORT` | *(empty)* | Takes precedence over `SOUNDBOARD_ADDR`. Railway injects this and health-checks the port it assigned, so the service must honour it. |
| `SOUNDBOARD_ALLOWED_ORIGIN` | *(empty)* | `Access-Control-Allow-Origin`. Empty sends no CORS headers, which is correct while the API and the site share one origin. |
| `SOUNDBOARD_STATIC_DIR` | *(empty)* | Built frontend to serve alongside the API. Empty serves the API alone. |
| `SOUNDBOARD_AUTH_PASSWORD` | *(empty)* | Password for the whole site. **Empty disables the prompt entirely** — fine locally, never in a deployment. |
| `SOUNDBOARD_AUTH_USER` | `panda` | Username shown alongside the password prompt. |

The two `TURSO_*` names are Turso's own spelling, so the values copy straight out of its
dashboard or CLI without renaming. Both are unset by default, which is what keeps a fresh
checkout zero-setup: clone, `go run .\cmd\api`, get a file on disk.

| Endpoint | Purpose |
| --- | --- |
| `GET /soundbites` | List every clip as JSON |
| `GET /soundbites/{id}/audio` | Stream a clip for playback (supports range requests) |
| `GET /soundbites/{id}/download` | Same bytes, sent as a file to save |
| `POST /soundbites/{id}/play` | Record one play, returns the new total |

`/download` exists as its own route because the HTML `download` attribute is ignored
cross-origin, and the site is served from a different origin than this API. It sends
`Content-Disposition: attachment` with a hyphenated form of the clip's display name, so a
visitor saves `ass-eaten-by-these-bitches.mp3` rather than the storage filename. The audio
itself is never modified — both routes stream the same bytes off disk.

The frontend needs no environment variables. It calls `/api` on its own origin — served
by this binary in production, proxied to `localhost:8080` by Vite in development (see
`web/vite.config.ts`). There is no API URL to configure in either case.

## Keeping it private

This is an unofficial fan site for someone else's character, published without her
endorsement, so it is deliberately not a public website. Three things enforce that, and
all three matter:

**A password over everything.** Setting `SOUNDBOARD_AUTH_PASSWORD` puts HTTP Basic Auth in
front of every route — pages, audio, and the JSON API alike. A soundbite list is not public
just because it is JSON.

```powershell
$env:SOUNDBOARD_AUTH_PASSWORD = "something-long"
go run .\cmd\api
```

`/api/health` is the one exemption. Railway and the Dockerfile `HEALTHCHECK` probe it
without credentials, and requiring a password there would mark every deploy unhealthy and
roll it back. It reveals nothing but that the process is up.

The API states which mode it is in at boot, next to the database line:

```
auth: password required (user "panda")
auth: DISABLED — every visitor gets in. Set SOUNDBOARD_AUTH_PASSWORD before deploying.
```

**No search indexing.** `web/static/robots.txt` disallows everything, and the server sends
`X-Robots-Tag: noindex, nofollow, noarchive` on every response — robots.txt is only a
request a crawler may ignore, and it is not part of the response at all.

**A disclaimer on the page.** The footer states the site is unofficial and unaffiliated, so
anyone who does reach it is not misled about whose site it is.

## Connecting to Turso

Production runs against [Turso](https://turso.tech/) (managed libSQL) rather than a file
on a disk that has to survive every deploy. It is still SQLite dialect, so the goose
migrations and the sqlc-generated queries are unchanged — what changes is the driver, the
DSN, and the connection pool.

**The DSN decides the driver.** There is no mode switch to set: `db.Open` is handed one
string, and its shape picks the path.

| DSN starts with | Driver | Means |
| --- | --- | --- |
| `libsql://`, `wss://`, `ws://`, `https://`, `http://` | `libsql` (libsql-client-go) | A libSQL server over the network |
| anything else | `sqlite` (`modernc.org/sqlite`) | A local file path, with the WAL and `busy_timeout` pragmas appended |

Both drivers are pure Go, so `CGO_ENABLED=0` still builds a static binary and Windows
still needs no C toolchain.

**Point it at Turso** by setting the two variables for the session:

```powershell
$env:TURSO_DATABASE_URL = "libsql://soundboard-you.turso.io"
$env:TURSO_AUTH_TOKEN   = "<token from the Turso dashboard>"
go run .\cmd\api
```

The CLI reads the same two variables, so `go run .\cmd\cli list` in that same shell talks
to the remote database rather than your local file. Unset them to go back:

```powershell
Remove-Item Env:TURSO_DATABASE_URL, Env:TURSO_AUTH_TOKEN
```

Setting them in one PowerShell window affects only that window. To make it stick across
new shells, use `[Environment]::SetEnvironmentVariable("TURSO_DATABASE_URL", "libsql://...", "User")`
— but prefer the per-session form, so you cannot forget you left the CLI pointed at
production.

**Testing the remote path — `turso dev` is not available on Windows.** `turso dev` runs a
local libSQL server speaking the real protocol, which would exercise the remote driver
without credentials. It ships only in the Turso CLI, and that CLI has no native Windows
build — its install docs require WSL. Installing WSL just for this is not worth it.

Use a throwaway Turso database instead. The free tier allows plenty, and nothing in one
needs preserving — `seed.Clips` rebuilds every row on boot:

```powershell
$env:TURSO_DATABASE_URL = "libsql://scratch-you.turso.io"
$env:TURSO_AUTH_TOKEN   = "<token for the scratch database>"
go run .\cmd\api
```

`http://` and `ws://` remain accepted remote schemes, so `turso dev` still works if you
already have WSL — that case, never production.

**Which database am I on?** The API prints it as its first line at boot:

```
database: local file data/soundboard.db
database: remote libSQL at libsql://soundboard-you.turso.io
```

The fallback to a local file is silent by design — that is right in development and
dangerous in production, where play counts would accumulate into a container filesystem
that gets thrown away. The log line makes the choice visible immediately rather than at
the first lost tally. The auth token is never part of it.

**Connection pooling differs by target.** Local keeps a single connection, because SQLite
allows one writer and a bigger pool just buys `SQLITE_BUSY` errors. Remote uses a pool of
10 with a 5-minute idle timeout: Turso serialises writes on its own side, so the limit
there is network latency rather than lock contention, and idle sockets to a remote host go
stale silently.

**No data migration is needed** to move to Turso. `seed.Clips` rebuilds every row from
`clips/` plus the committed `names.json` on first boot, so an empty Turso database fills
itself. The only thing in the local file that is not reproducible is play counts, and those
are test data — clear them with `reset-plays -all` before launch either way.

### Status

The Turso path is **implemented and verified against a live Turso database** (2026-08-29).
Confirmed working end to end: both goose migrations applying on boot to an empty database,
`seed.Clips` creating all 56 rows from `clips/` plus `names.json`, `GET /api/soundbites`
returning them, a `POST .../play` write, and `cli list` / `cli import -dry-run` / `cli check`.

**Transactions were the open question and they work.** `internal/db/migrate.go` runs goose
on every boot and goose wraps each migration in a transaction; `cmd/cli/import.go` opens its
own `BeginTx` around a batch import. An interactive transaction over libSQL's HTTP protocol
is state held on the server across requests rather than a lock on a file handle, so this was
not safe to assume — both paths were exercised against a real server before this note was
written.

Still unexercised: a Railway deploy. The container has never been built, because Docker is
not installed on the development machine.
