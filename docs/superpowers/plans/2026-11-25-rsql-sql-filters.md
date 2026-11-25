# RSQL filters in SQL — implementation plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Parse Elide/RSQL filter expressions and apply them, plus multi-key sorts, in SQL for the organization, project, tag and workspace lists.

**Architecture:** A dependency-free parser in `internal/rsql` turns the filter string into an AST. A hand-written builder in `internal/store` renders that AST into a parameterized `WHERE` against a whitelisted attribute map per resource, adds `ORDER BY`/`LIMIT`/`OFFSET`, and scans `SELECT *` by position. Handlers parse once and pass the AST down; the twelve paged/count sqlc queries are deleted.

**Tech Stack:** Go 1.27, pgx v5, sqlc (generated code untouched by hand), Postgres 16.

**Spec:** `docs/superpowers/specs/2026-11-25-rsql-sql-filters-design.md`

## Global Constraints

- Upstream Elide 5 semantics: `==`, `!=`, `=in=`, `=out=` case-sensitive; `=ini=`, `=outi=` case-insensitive; `*` is the wildcard.
- `%`, `_` and `\` in a value stay literal (`LIKE … ESCAPE '\'`).
- Malformed syntax or an unparseable value → `400` JSON:API error. An attribute we do not model → the whole filter is ignored (widen, never hide).
- Never edit `internal/store/*.sql.go`; change `internal/store/queries/*.sql` and run `make generate`.
- Migrations are append-only; no migration is needed for this change.
- Every commit builds and passes alone; generated code never shares a commit with hand-written logic.
- Merge with `gh pr merge --merge` only.

## Review Focus

- A filter containing an unmodelled attribute (`workspaceTag.tagId`, `job.status`) must not narrow the result; the whole expression is dropped.
- Quoted values containing `;`, `,`, `(` or `)` must survive the lexer.
- `=in=()` with no values and a dangling `;` must be a `400`, not a panic or a silent full list.
- A `lastJobDate` value that is not a timestamp must be a `400`, not a `500` from Postgres.
- Position scanning must fail loudly if model field order ever diverges from table column order.

---

### Task 1: RSQL parser

**Files:**
- Create: `internal/rsql/rsql.go`
- Test: `internal/rsql/rsql_test.go`

**Interfaces:**
- Produces:
  ```go
  type Op string // "==", "!=", "=in=", "=out=", "=isnull=", "=ini=", "=outi=", "=lt=", "=le=", "=gt=", "=ge="
  type Node interface{ node() }
  type And struct{ Left, Right Node }
  type Or struct{ Left, Right Node }
  type Compare struct{ Attr string; Op Op; Values []string }
  func Parse(raw string) (Node, error)
  var ErrSyntax = errors.New("invalid filter expression")
  ```

- [ ] **Step 1: Write the failing test** — table of the verbatim UI expressions from `docs/upstream/ui-operations.md`, each with the s-expression the test helper prints (`(and (== name prod) (=in= status failed cancelled))`). Include: `name=ini='*prod*'`, `(name=ini='*prod*',description=ini='*prod*')`, `project.id=isnull=true`, `id=in=('<a>','<b>')`, `(lastJobStatus=isnull=true,lastJobStatus=='NeverExecuted')`, `name=="<org>"`, precedence (`a==1;b==2,c==3` → `(or (and (== a 1) (== b 2)) (== c 3))`), a quoted value containing `;`,` and a wildcard, and the error cases ``, `name`, `name==`, `name=in=()`, `(name==x`, `name==x;`.
- [ ] **Step 2: Run it** — `go test ./internal/rsql` fails with "undefined: Parse".
- [ ] **Step 3: Implement the lexer and parser** — scan values respecting quotes and backslash escapes, match operators longest-first, precedence-climbing with `;` binding tighter than `,`, parentheses, and `=in=`/`=out=` value lists.
- [ ] **Step 4: Run the tests** — `go test ./internal/rsql` passes.
- [ ] **Step 5: Commit** — `feat(rsql): parse the upstream filter grammar`.

### Task 2: Store list builder and resource specs

**Files:**
- Create: `internal/store/list.go`, `internal/store/listspecs.go`
- Test: `internal/store/list_test.go` (uses `dbtest.Pool`)

**Interfaces:**
- Consumes: `rsql.Node`, `rsql.Compare`, `rsql.Op`, `rsql.ErrSyntax`.
- Produces:
  ```go
  type Kind int
  const (KindText Kind = iota; KindTime; KindBool)
  type Attr struct{ Expr string; Kind Kind }
  type ListSpec struct {
      Table    string
      Base     string   // e.g. "organization_id = $1 AND deleted = false"
      BaseArgs []any
      Attrs    map[string]Attr
      Sorts    map[string]string
  }
  type SortKey struct{ Column string; Desc bool }
  func (q *Queries) ListPage[T any](ctx context.Context, spec ListSpec, filter rsql.Node, sorts []SortKey, limit *int64, offset int64) ([]T, error)
  func (q *Queries) CountList(ctx context.Context, spec ListSpec, filter rsql.Node) (int64, error)
  var OrganizationList, ProjectList, TagList, WorkspaceList, ProjectWorkspaceList ListSpec
  ```
- Attribute maps: organization `name`, `description`, `id`; project `name`, `description`, `id`; tag `name`, `id`; workspace `name`, `description`, `id`, `lastJobStatus`, `policyComplianceStatus`, `project.id` (→ `project_id`), `lastJobDate`. Sorts: `name` → `lower(name)`, `id` → `id`, `lastJobDate` → `last_job_date`.
- `OrganizationList.BaseArgs` is set by the caller (the membership restriction); specs are copied, maps are read-only.

- [ ] **Step 1: Write the failing tests** — one database test that seeds mixed-case rows and asserts `ListPage` results for: `name==` (case-sensitive), `name=ini='*PROD*'`, `=in=`, `=isnull=`, `!=`, `=gt=` on `lastJobDate`, an unsupported attribute returning the unfiltered list, descending and multi-key sorts, and `CountList` matching the filtered length.
- [ ] **Step 2: Run them** — `go test ./internal/store -run TestList` fails.
- [ ] **Step 3: Implement the renderer** — placeholders continue from `len(BaseArgs)`; `*` → `%` with `%`/`_`/`\` escaped; `=ini=` wraps in `lower()`; `=isnull=` emits `IS NULL`/`IS NOT NULL`; `KindTime` parses RFC3339 and returns `fmt.Errorf("%w: …", rsql.ErrSyntax)`; unsupported attribute drops the whole filter; unknown sort keys are ignored and `id` is appended as the tiebreaker; `ListPage` scans with `pgx.RowToStructByPos[T]`.
- [ ] **Step 4: Add the scan-order test** — reflect over `Organization`, `Project`, `Tag`, `Workspace` field names, snake_case them, and compare with `information_schema.columns` ordered by `ordinal_position`.
- [ ] **Step 5: Run the tests** — `go test ./internal/store` passes.
- [ ] **Step 6: Commit** — `feat(store): filter and sort list queries in SQL`.

### Task 3: Handler cutover

**Files:**
- Modify: `internal/api/api.go` (`listQuery`, `parseList`; delete `nameFilter`, `prefixPattern`)
- Modify: `internal/api/organizations.go`, `projects.go`, `tags.go`, `workspaces.go`
- Delete: `internal/api/filter_sort_page_test.go`
- Modify: `internal/api/organizations_sql_test.go`, `projects_tags_sql_test.go`, `workspaces_sql_test.go`
- Modify: `internal/store/queries/organizations.sql`, `projects.sql`, `tags.sql`, `workspaces.sql` (remove the twelve `*Page*`/`Count*ByName` queries)
- Regenerate: `make generate`

**Interfaces:**
- Consumes: `rsql.Parse`, `store.ListPage`, `store.CountList`, the five specs, `store.SortKey`.
- Produces: `func parseList(r *http.Request, filterType string) (listQuery, error)`; `listQuery{Filter rsql.Node, Sorts []store.SortKey, …}` keeps `page()`/`pageWindow()`.

- [ ] **Step 1: Rewrite the endpoint tests** as expected-result cases over seeded rows: the documented workspace expressions, `sort=name,id`/`-name,-id`, `sort=lastJobDate`, paging, `page[totals]`, wildcard escaping, case sensitivity, and `filter[workspace]=name==` returning `400`.
- [ ] **Step 2: Run them** — they fail against the current handlers.
- [ ] **Step 3: Cut the handlers over** — parse once; `errors.Is(err, rsql.ErrSyntax)` → `jsonapi.WriteError(w, http.StatusBadRequest, "Invalid Filter", err.Error())`; organization copies `store.OrganizationList` and sets `BaseArgs`; the other four use their spec directly; `CountList` only when `q.Totals`.
- [ ] **Step 4: Delete the twelve queries and regenerate** — `make generate`; commit generated files separately.
- [ ] **Step 5: Run the tests** — `go test ./internal/api` against Postgres passes.
- [ ] **Step 6: Commit** — `feat(api): apply RSQL filters and sorts in SQL`, then `chore(store): regenerate sqlc`.

### Task 4: Docs and graph

**Files:**
- Modify: `docs/compat.md` (query conventions; M1 limitations bullet)
- Modify: `AGENTS.md` (list bullet)
- Regenerate: `graphify-out/`

- [ ] **Step 1: Document the operator set, case sensitivity, the ignored-attribute rule, the relationship filters still unsupported, and the `action` exception.**
- [ ] **Step 2: Run `graphify update .` and commit `graphify-out/` alone.**
- [ ] **Step 3: Verify** — `make check`, `make test-race`, `make smoke`.
- [ ] **Step 4: Commit** — `docs: record the RSQL filter support` and `chore(graph): update knowledge graph`.

## Self-review

- Spec coverage: parser (T1), SQL application (T2), handler wiring and 400s (T3), docs (T4). Sorts and totals are in T2/T3.
- The spec's non-goals (relationship filters, collection operators, global `filter=`) are enforced by the ignored-attribute rule and documented in T4.
- Types line up: `rsql.Node` flows from `Parse` into `ListPage`/`CountList`; `store.SortKey` is produced by `parseList`.
