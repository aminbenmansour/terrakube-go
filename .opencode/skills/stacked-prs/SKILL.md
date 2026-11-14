---
name: stacked PRs
description: Prepare focused commits and stacked pull requests without rewriting shared history
---

# Stacked PRs

## Rules

- Never commit directly to `main`.
- One logical change per commit; use conventional commit subjects.
- Each PR builds on the previous branch (`feat/01-...` → `feat/02-...`).
- Do not force-push or rewrite shared history without approval.
- Keep generated code, docs and behavior changes separable.

## Workflow

1. Verify the current branch with `make check`.
2. Split the diff into commits by concern.
3. Create a branch from the previous stack tip.
4. Push only the feature branch (`git push -u origin <branch>`).
5. Open a PR with base = previous branch, or `main` for the first PR.
6. Include Summary, Test plan and compatibility notes.
7. After the base merges, retarget the next PR to `main`.

## PR template

```markdown
## Summary
...

## Test plan
- [ ] `make check`
- [ ] `scripts/smoke.sh` (API changes)

## Compatibility
...
```
