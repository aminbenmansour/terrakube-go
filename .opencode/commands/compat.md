---
description: Review a change against the Terrakube compatibility contract
agent: plan
---

Review the current diff for compatibility regressions. $ARGUMENTS

1. Read `docs/compat.md`.
2. Check paths, methods, status codes, content types, JSON:API payloads,
   auth, table/column names and migration safety.
3. Compare against the upstream behavior described in the document.
4. Report findings in severity order with file and line references.

Do not edit files.
