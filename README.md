# llm-go-friend

A Go linter built for LLM coding agents.

LLMs write Go that compiles and passes tests but is structurally hostile to the next agent (or human) that has to read it. 500-line files, 8-level nesting, functions with 12 parameters. The code works, but it's a maintenance nightmare.

`llm-go-friend` catches these structural problems at the AST level and outputs machine-readable violations in [TOON](https://github.com/toon-format/spec) format, so agents can parse the results and fix their own mess.

## Claude Code Skill

The repo includes a Claude Code skill at `skills/llm-go-friend/SKILL.md` that pairs the tool with manual judgment checks an LLM reviewer should perform (stale comments, excessive indirection, implicit behavior, inconsistent patterns, missing type context). The skill runs `llm-go-friend` for the mechanical checks and guides the reviewer through the things the tool can't detect.

## Install

```
go install github.com/tom-pang/llm-go-friend/cmd/llm-go-friend@latest
```

## Usage

```
llm-go-friend file.go [file.go ...]
```

Exit codes:
- `0` — no violations
- `1` — violations found (printed to stdout as TOON)
- `2` — usage error or parse failure

## Checks

| Check | Threshold | What it catches |
|-------|-----------|-----------------|
| `file_length` | 500 lines | Files too large for an agent to hold in context |
| `func_length` | 50 lines | Functions that should be decomposed |
| `nesting_depth` | 4 levels | Deeply nested control flow |
| `cyclomatic_complexity` | 10 | Functions with too many branching paths |
| `param_count` | 5 params | Function signatures that need a struct |
| `bare_interface` | 0 | Exported functions accepting `interface{}` or `any` |

## Output

Violations are TOON-encoded tables. An agent or script can parse these directly:

```
violations[1]{file,line,check,name,value,threshold}:
  testdata/long.go,1,file_length,"",501,500
```

Each violation includes the file, line number, which check fired, the measured value, and the threshold it exceeded. The `name` field contains the function name for function-scoped checks (empty for file-level checks).

## Tradeoffs

The thresholds are opinionated. They're set for LLM-generated code, which tends to be verbose — a 50-line function limit is generous by human standards but catches the 200-line monsters that agents produce. There's no configuration. If you need different thresholds, fork it.

The bare interface check only fires on exported functions. Unexported functions accepting `any` is a code smell but not the same category of problem — exported APIs with bare interfaces make the whole package harder for agents to use correctly.

No `//nolint` or suppression mechanism. This is intentional. Suppression comments are the first thing an agent adds when it can't figure out how to fix the real problem.
