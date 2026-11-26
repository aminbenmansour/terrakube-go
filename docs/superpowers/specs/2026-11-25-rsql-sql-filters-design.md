# RSQL filters in SQL — design

Issue #55, milestone M6. Follows the three slices that moved list paging into
SQL (#179–#184).

## Intent

The management list endpoints accept Elide/RSQL filters but only understand
`name==value` and `name==value*`. Every other term is silently dropped, so the
UI's workspace search (`name=ini='*prod*'`) and status filters return the whole
list. #55's remaining acceptance criterion is "RSQL operators match upstream";
this change parses RSQL and applies it in the query.

## Goals

- Parse the expressions the UI sends, quoted verbatim in
  `docs/upstream/ui-operations.md`: `==`, `!=`, `=in=`, `=out=`, `=isnull=`,
  `=ini=`, `=outi=`, `=lt=`, `=le=`, `=gt=`, `=ge=`, `*` wildcards, single- or
  double-quoted and bare values, `;` (AND), `,` (OR) and parentheses.
- Apply filters on scalar attributes of the organization, project, tag and
  workspace lists in SQL; `page[totals]` counts the filtered set.
- Sort by a whitelisted set of attributes (`name`, `id`, `lastJobDate`),
  multi-key, with `-` for descending and `id` appended as the tiebreaker.
- Malformed filter → `400` with a JSON:API error, as Elide does. An unknown
  attribute is ignored, as today: we model fewer attributes than upstream, and
  failing would break the UI on attributes we have not built.

## Non-goals

- Relationship filters beyond the workspace's `workspaceTag.tagId`,
  `workspaceTag.value` and `job.status` (those three are a follow-up slice and
  use `EXISTS`; `project.id` is a plain column).
- Collection operators (`=hasmember=`, `=isempty=`, `=between=`, `=subsetof=`).
- The untyped global `filter=` parameter.
- The `action` list keeps its dedicated two-term parser: the catalog is static
  and small, so SQL buys nothing there.

## Approaches considered

1. **Dynamic SQL from a whitelisted column map (chosen).** One parser, one
   renderer, four small resource specs. Replaces twelve near-identical
   paged/count queries.
2. **Nullable parameters on the existing sqlc queries.** No parser, no dynamic
   SQL, but OR groups (`(lastJobStatus=isnull=true,lastJobStatus=='NeverExecuted')`)
   and operator combinations do not fit; each new expression becomes a query
   and parameter change. Rejected as a dead end.
3. **Parse in Go, filter in memory.** Rejected: it reintroduces the scaling
   problem #55 exists to fix.

## Design

- `internal/rsql`: lexer and recursive-descent parser producing an AST
  (`And`, `Or`, `Compare{Attr, Op, Values}`). AND binds tighter than OR;
  parentheses override. Operators are matched longest-first (`=isnull=`,
  `=ini=`, `=out=` before `=`). Values strip `'…'`/`"…"` and accept backslash
  escapes.
- `internal/store/list.go`: `ListSpec{Table, BaseWhere, BaseArgs, Attrs, Sorts}`
  with `Attr{Expr, Kind}`. `ListPage[T]` renders a parameterized `WHERE` from
  the AST, appends `ORDER BY`, `LIMIT`/`OFFSET`, runs `SELECT *` and scans with
  `pgx.RowToStructByPos[T]`. `Count` uses the same `WHERE`.
- Position scanning is safe because sqlc emits model fields in table-column
  order, which `ALTER TABLE` preserves by appending. A test asserts the order
  for all four models.
- Semantics follow Elide 5: `==`, `!=`, `=in=`, `=out=` are case-sensitive;
  `=ini=`, `=outi=` are case-insensitive; `*` maps to `%` with `%`, `_` and
  `\` escaped; `=isnull=true|false` is `IS NULL` / `IS NOT NULL`; comparisons
  cast by kind. A timestamp that does not parse is a `400`.
- Sorts: `lower(name)` for `name` (keeps the existing order and index), `id`,
  `last_job_date`. Unknown sort keys are ignored.
- Resource specs live in `internal/store/listspecs.go` and use JSON:API
  attribute names. The organization spec keeps its `allowed_ids` restriction in
  the base `WHERE`.
- Handlers parse once, pass the AST and sort keys down, and return `400` on a
  parse error. The twelve paged/count sqlc queries and their parameter structs
  are deleted.

## Behavior changes

- `name==` becomes case-sensitive, matching upstream; the UI uses `=ini=` when
  it wants a case-insensitive match.
- A malformed filter is a `400` instead of being ignored.
- `sort=id`, `sort=lastJobDate` and multi-key sorts start working.

## Layering

1. `internal/rsql` parser and its table tests — no behavior change.
2. `internal/store` builder, resource specs, scan-order test; remove the sqlc
   list queries.
3. Handler cutover, endpoint tests rewritten as expected-result cases over the
   documented expressions, docs and graph.

## Test plan

- Parser: the verbatim UI expressions, precedence, quoting, wildcards,
  malformed input.
- Endpoint: each documented workspace filter against seeded rows (search,
  status set/null, id set, project id, compliance status), sorts in both
  directions, paging, `page[totals]`, wildcard escaping, case sensitivity, and
  a `400` for bad syntax.
- Scan order: compare `information_schema.columns` order with the model field
  order for organization, project, tag and workspace.
