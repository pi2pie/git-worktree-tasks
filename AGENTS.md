# AGENTS.md

This document provides essential context for any agent working in this repository.

---

## Development Environment

- Language: Go (see `go.mod` for version and module path)
- Entry point: `main.go`
- Local tools: `cli/`, examples in `examples/`
- Workspace root: this repo directory (no external deps needed unless noted)

---

## Skills

Use the following skill for this repository. Skills are defined under `.agents/skills/*`. How they’re discovered/loaded depends on the agent tool.

- `golang-expert` — apply for all Go code changes, refactors, reviews, testing, or Go best‑practice questions.
  - Skill file: `.agents/skills/golang-expert/SKILL.md`
  - Load only the relevant reference files from the skill when needed.

---

## Documentation Conventions

All meaningful agent work SHOULD be documented.

Optional metadata:
- If you update an existing plan, research doc, or job record, you MAY add a `modified-date: YYYY-MM-DD` field to the front-matter.
- Keep the original `date` value unchanged; `modified-date` is for the latest update.

### Timezone

- Unless a document or request specifies otherwise, record dates in UTC.

### Plan Documents

Location:
```text
docs/plans/plan-YYYY-MM-DD-<short-title>.md
```

Notes:
- Do not create or edit `docs/plan.md`.
- Use the date for when the plan is created and a short, kebab-case title.

Front-matter format:
```yaml
---
title: "<plan title>"
date: YYYY-MM-DD
status: draft | active | completed
agent: <agent name>
---
```

---

### Research Documents

Use research docs for exploratory work that is not yet ready for a plan but may inform one.

Location:
```text
docs/research-YYYY-MM-DD-<short-title>.md
```

Notes:
- Use the date the research starts and a short, kebab-case title.
- Keep scope focused on a single topic or question.
- If research becomes actionable, create a plan doc and link to it.

Front-matter format:
```yaml
---
title: "<research title>"
date: YYYY-MM-DD
status: draft | in-progress | completed
agent: <agent name>
---
```

Suggested sections:
- Goal
- Key Findings
- Implications or Recommendations
- Open Questions (optional)
- References (use footnote-style links)

Traceability:
- Research docs should include a short "Related Plans" section when applicable, with links to plan docs.
- Plan docs should include a short "Related Research" section when applicable, with links to research docs.
- Use those exact section titles for consistency.
- Omit the section if there are no relevant links.

---

### Job Records

For concrete tasks or implementations, create a job record.

Location:
```text
docs/plans/jobs/YYYY-MM-DD-<short-title>.md
```

Front-matter format:
```yaml
---
title: "<job title>"
date: YYYY-MM-DD
status: draft | in-progress | completed | blocked
agent: <agent name>
---
```

---

### Status Meanings

- `draft` — idea or exploration, not executed
- `active` — current plan being worked on
- `in-progress` — task implementation ongoing
- `completed` — work finished
- `blocked` — waiting on decision or dependency

---

## Writing Guidelines

- Prefer clarity over verbosity
- Record *what changed* and *why*
- Avoid repeating information already in other documents
- Assume future agents will read this without prior context

---

## Philosophy

> This file exists to reduce guesswork for the next agent.
