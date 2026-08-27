# soundboard-api

Go API + SQLite behind the site's soundboard, plus the `cli` tool used to manage clips.

Run everything from **inside this directory** — the default paths (`data/soundboard.db`,
`data/audio`) are relative, so running from the repo root puts the database somewhere
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

- **Database**: `data/soundbites` table, `data/soundboard.db`. Gitignored.
- **Audio**: `data/audio/`, copied in from wherever you staged it. Gitignored.
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
| `SOUNDBOARD_DB_PATH` | `data/soundboard.db` | SQLite file |
| `SOUNDBOARD_AUDIO_DIR` | `data/audio` | Where clip audio lives |
| `SOUNDBOARD_ADDR` | `:8080` | Listen address |
| `SOUNDBOARD_ALLOWED_ORIGIN` | `*` | `Access-Control-Allow-Origin` |

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

The frontend reads `PUBLIC_SOUNDBOARD_API_URL` from `web/.env`.
