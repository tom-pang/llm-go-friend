---
name: llm-go-friend
description: Use when reviewing code for LLM-friendliness, during code review, or when asked to check if code is easy for AI agents to understand and modify. Use when files feel too large, functions too complex, or codebase hard for agents to navigate.
---

# LLM-Friendly Code Check

## Overview

Review code for patterns that make it harder for LLMs to understand, modify, and reason about. Based on published research showing specific, measurable characteristics that degrade LLM performance on code tasks.

Core principle: LLMs process code as linear token sequences with finite context. Anything that increases the distance between related concepts, hides information behind indirection, or creates gaps between source and runtime behavior makes LLMs worse at their job.

## When to Use

- Code review (part of review checklist)
- Evaluating a codebase for AI-assisted development readiness
- After refactoring, to verify improvements
- When agents keep making mistakes in specific files

## Step 1: Run llm-go-friend (Go code)

For Go code, run the `llm-go-friend` tool FIRST. It checks the mechanical/measurable thresholds automatically:

```bash
llm-go-friend ./...
```

This covers: file length, function length, nesting depth, cyclomatic complexity, parameter count, and bare `interface{}`/`any` on exported functions.

**DO NOT manually check these things yourself.** The tool handles them. Do not eyeball file lengths, count function lines, or estimate nesting depth. If `llm-go-friend` passes, those checks pass. If it reports violations, include its output in your review. Your job is the judgment calls below — not re-doing what the tool already does.

## Step 2: Judgment Checks (reviewer does these)

These require reading and understanding the code. The tool cannot detect them.

| # | Check | What to look for |
|---|-------|-----------------|
| 1 | **Stale/wrong comments** | Comments that contradict the code. LLMs trust all comments equally — incorrect comments are **actively harmful**, dropping comprehension 12+ points. Worse than no comments. Zero tolerance. |
| 2 | **Excessive indirection** | More than 2 levels of call chain to trace. Single-implementation interfaces that add indirection without value. Deep call chains require holding multiple file contexts simultaneously. |
| 3 | **Implicit behavior** | `init()` functions doing non-trivial work, global mutable state modified at runtime, heavy use of `reflect`. Creates gap between source code and runtime behavior. LLMs are stateless — they can only reason about what the source says. |
| 4 | **Inconsistent patterns** | Same concern handled different ways in different places. Error handling via error returns in one place, panics in another. When patterns are inconsistent, LLMs pick randomly between styles. |
| 6 | **Missing type context** | Exported functions returning bare `interface{}` or accepting untyped params where a named type would clarify intent. Types compensate for missing context — they constrain the LLM's search space. |

## How to Report

1. Include `llm-go-friend` output verbatim for any violations.
2. For judgment checks, report as: file:line, which check, what's wrong, suggested fix.

## What Good Looks Like

- `llm-go-friend` passes clean
- No stale comments
- Flat call chains — can understand a function without reading 3 other files
- Explicit behavior — no surprises from `init()` or `reflect`
- One pattern per concern across the codebase
- Named types on public APIs

## What NOT to Check

- Formatting/style (that's what `gofmt` is for)
- Test coverage percentage (separate concern)
- Performance characteristics (orthogonal)
- Single-letter variable names (idiomatic Go — `r` for reader, `w` for writer, `ctx` for context, `err` for error — this is fine)
- Anything `llm-go-friend` already checks (do not duplicate its work)

## Quick Reference: The Three Worst Things for LLMs

1. **Stale comments** — LLMs will confidently follow wrong instructions in comments
2. **1000+ line files** — attention literally cannot cover the whole file
3. **Deep nesting + branching** — exponential blowup in reasoning paths
