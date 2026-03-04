# llm-go-friend

Static analysis tool for Go that checks code for LLM-friendliness — patterns that make code harder for AI coding agents to understand, modify, and reason about.

## Why

LLMs process code as linear token sequences with finite context windows. Research shows specific, measurable code characteristics that degrade LLM performance on code tasks. Most existing linters check for human readability. This tool checks for machine readability.

## Research basis

The checks are backed by published research:

**File length.** LLM attention degrades linearly across long files. The "Lost in the Middle" paper (Liu et al., Stanford/UC Berkeley) showed 30%+ accuracy drops when relevant information sits in the middle of long contexts. A separate study found 50% of faults correctly localized by LLMs are in the first 25% of a file's lines, while only 18% are detected in the last 25%.

**Function length.** LLM refactoring studies show 50% LOC reduction correlates with better LLM handling. Addy Osmani's LLM coding workflow research emphasizes LLMs work best targeting "one function, fix one bug, add one feature at a time" — functions should be single-purpose and independently comprehensible.

**Nesting depth.** The LM-CC paper (2026) found that compositional nesting depth and branching factor correlate at -0.93 to -0.97 with LLM task performance across repair, translation, and reasoning tasks. Reducing nesting through semantics-preserving rewrites yielded 10-21% relative improvement. Traditional cyclomatic complexity is a weaker signal — it's the depth that kills LLMs, not the number of paths.

**Parameter count.** Large parameter lists are ambiguous at call sites. Dagster's "Dignified Python" rules (applicable across languages) flag functions with more than 5 parameters as a readability hazard for both humans and LLMs.

**Type coverage.** Type-constrained code generation research shows enforcing type constraints during LLM code generation reduces compilation errors by more than half and increases functional correctness by 3.5-5.5%. In Go this means exported functions should have clear, named types — not bare `interface{}` or deeply nested anonymous structs.

**Naming.** "How Does Naming Affect LLMs on Code Analysis Tasks?" found that anonymizing all names dropped code search accuracy from 70.36% to 17.42%. Misleading names performed worse than random gibberish. Single-letter variables outside loop iterators actively hurt LLM comprehension.

**Comments.** Research on comments and LLM comprehension found correct comments improve understanding from 84% to 96%. But when comments are inaccurate, comprehension drops from 73% to 61%. LLMs attempt to reconcile contradictions by treating all comments as simultaneously correct. Stale comments are actively harmful — worse than no comments at all.

## Checks

| Check | Threshold | Detectable from AST |
|-------|-----------|-------------------|
| File length | >500 lines | Yes (`go/token`) |
| Function length | >50 lines | Yes (`go/ast`) |
| Nesting depth | >4 levels | Yes (`go/ast`) |
| Cyclomatic complexity | >10 per function | Yes (`go/ast`) |
| Parameter count | >5 per function | Yes (`go/ast`) |
| Bare `interface{}` / `any` params | Flag on exported functions | Yes (`go/ast`) |
| Single-letter variables | Outside `i,j,k` in loops | Yes (`go/ast`) |

### Not checked (need LLM or human judgment)

- Stale/incorrect comments
- Naming quality beyond length
- Pattern consistency across files
- Excessive indirection depth
- Implicit behavior (init functions doing too much)

These remain in the `llm-friendly-check` skill as a manual review checklist.

## Goals

1. Fast — runs on large codebases in seconds, no LLM calls
2. Useful as a pre-commit hook or CI check
3. Reports violations with file, line, and current value vs threshold
4. Non-zero exit on violations
5. Configurable thresholds via flags or config file
6. Go only — no ambition to be multi-language

## Sources

- Liu et al. "Lost in the Middle: How Language Models Use Long Contexts" (2023)
- "Rethinking Code Complexity Through the Lens of LLMs" — LM-CC paper (2026)
- "How Does Naming Affect LLMs on Code Analysis Tasks?" (2023)
- "Impact of Comments on LLM Comprehension of Legacy Code" (2025)
- "Type-Constrained Code Generation with Language Models" (2025)
- "Code Refactoring with LLM: Comprehensive Evaluation" (2025)
- Osmani, "My LLM Coding Workflow Going Into 2026"
- Crawshaw, "How I Program with Agents" / "Eight More Months of Agents"
- Dagster, "Dignified Python: 10 Rules to Improve Your LLM Agents"
