# Springfield Protocol: Quick Start Guide

A concise reference for using the Springfield Protocol framework and its loops and skills.

---

## One-Minute Overview

**The Springfield Protocol** is a character-driven framework for Agile Agentic Development. It combines:

1. **The Ralph Wiggum Loop** - A stateless resampling engine that ensures quality through persistent iteration
2. **Character Skills** - Specialized agents (Lisa, Ralph, Bart, Herb, etc.) that handle specific roles
3. **Agentic Loops** - Feedback patterns (Tree of Thoughts, ReAct, Plan-and-Execute, etc.) for different problems

---

## The Core Engine

### Ralph Wiggum Variant: Stateless Resampling Loop

```
┌─────────────────────────────────────────┐
│  Monitor PLAN.json for failed tasks     │
└──────────────┬──────────────────────────┘
               │
        ┌──────▼──────┐
        │ Spawn agent │ (ephemeral context)
        │ in clean    │
        │ worktree    │
        └──────┬──────┘
               │
        ┌──────▼──────────┐
        │ Execute task    │
        │ (Strict TDD)    │
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │ Verify results  │
        │ (Herb/Bart)     │
        └──────┬──────────┘
               │
        ┌──────▼──────────┐
        │ Update PLAN.json│
        └──────┬──────────┘
               │
        [Loop back]
```

---

## Character Quick Reference

| Character | Use When | Output |
|:----------|:---------|:-------|
| **Lisa** | You need to plan work | TODO.md with tasks |
| **Ralph** | You need to implement | Tested code + git commits |
| **Bart** | You need to break things | Security/quality issues |
| **Herb** | You need to verify quality | Coverage report + sign-off |
| **Marge** | You need user alignment | PR feedback + scope notes |
| **Frink** | You need architecture | ADR + pattern spec |
| **Wiggum** | You have an issue to triage | Triaged issue + DoR checklist |
| **Homer** | You need infrastructure import | Terraform modules + validation |
| **Lovejoy** | You need to release | Version tag + CHANGELOG |

---

## Loop Quick Reference

### For Problem-Solving

**Have a vague problem?** → **Tree of Thoughts**
- Generate multiple solution paths
- Evaluate and prune low-scoring ones
- Explore the most promising paths

**Have a specific error?** → **ReAct**
- Verbalize reasoning at each step
- Take concrete actions
- Observe results before proceeding

**Have a clear spec?** → **Plan-and-Execute**
- Break down into tasks upfront
- Execute sequentially
- Validate each step

### For Quality Improvement

**Need to polish output?** → **GECR (Generate → Evaluate → Critique → Refine)**
- Generate candidates
- Score them
- Critique weaknesses
- Refine iteratively

**Need to learn from testing?** → **TALAR (Test → Analyze → Learn → Adjust → Retest)**
- Run experiments
- Analyze results
- Extract insights
- Adjust and retest

### For Coordination

**Multiple agents working together?** → **Manager-Worker Loop**
- Lisa orchestrates
- Workers (Ralph, Bart, Herb) specialize
- Manager aggregates results

**Two agents iterating?** → **Dialogue Loop**
- Developer proposes
- Reviewer critiques
- Iterate to consensus

---

## Discovery Track vs. Delivery Track

### Discovery (Design Thinking)
```
Product Discovery → Define user need
         ↓
    Ideation & prototyping
         ↓
    Feature Brief (validated)
```

**Characters:** Product Discovery, Marge, Frink
**Loops:** Tree of Thoughts, Dialogue, Observe-Hypothesize-Experiment-Conclude

---

### Delivery (Agile)
```
PLAN.md → TODO.md tasks
         ↓
    Ralph Wiggum Loop (implementation)
         ↓
    Verified, tested, quality code
```

**Characters:** Lisa, Ralph, Bart, Herb
**Loops:** Plan-and-Execute, Ralph Wiggum, Dialogue

---

## Common Workflows

### ✅ Implement a Feature (Happy Path)

```
1. Feature Brief arrives (validated)
   ↓
2. @lisa "Break this into tasks" → TODO.md
   ↓
3. @ralph "Implement task 1" (TDD loop)
   ↓
4. @bart "Review for security" (adversarial)
   ↓
5. @herb "Verify coverage" (95%+ check)
   ↓
6. @marge "Check user alignment" (feedback)
   ↓
7. @lovejoy "Release it" (publish + tag)
```

### 🔍 Debug an Issue

```
1. @wiggum "Triage this issue"
   ↓
2. Search kedb for known solutions
   ↓
3. IF found → document & close
   ↓
4. IF not found → go to feature workflow
```

### 🏗️ Import Infrastructure

```
1. @wiggum "Map this AWS account"
   ↓
2. @homer "Create zero-change Terraform"
   ↓
3. @herb "Verify no changes"
   ↓
4. @tfc-api "Inspect the state"
   ↓
5. @lovejoy "Release the module"
```

### 🤔 Decide Architecture

```
1. @frink "Create an ADR for this decision"
   ↓
2. @bart "Poke holes in this design"
   ↓
3. @frink "Refine based on feedback"
   ↓
4. Document the pattern for reuse
```

---

## When to Use Each Loop

| Loop | Problem Type | Time Constraint | Team Size |
|:-----|:-------------|:----------------|:----------|
| Sense-Plan-Act | Real-time decisions | Immediate | 1 |
| ReAct | Debugging errors | Minutes-hours | 1 |
| Tree of Thoughts | Complex decisions | Hours | 1-3 |
| Plan-and-Execute | Clear requirements | Days-weeks | Any |
| Ralph Wiggum | Quality delivery | Days-weeks | Multi-agent |
| Manager-Worker | Parallel work | Days-weeks | 3+ |
| Dialogue | Collaborative refinement | Hours-days | 2 |

---

## File Organization

```
Project Root/
├── PLAN.md                 # High-level roadmap (Lisa)
├── TODO.md                 # Executable tasks (Lisa)
├── FEEDBACK.md             # Review feedback (Marge)
├── docs/
│   ├── ADR-*.md           # Architecture decisions (Frink)
│   └── RUNBOOK.md         # Operations guide
├── .github/
│   ├── agents/            # Agent definitions
│   └── workflows/         # CI/CD
└── src/                   # Implementation
    └── (tests/ + code)
```

---

## Key Principles

1. **Iteration over Perfection**: Ralph Wiggum Loop ensures quality through persistence, not one-shot perfection.

2. **Memorable Personas**: Character-driven approach makes roles intuitive and easy to remember.

3. **Modular & Lean**: Each skill is self-contained to minimize context rot and fit in limited context windows.

4. **Dual-Track Agility**: Discovery (Design Thinking) and Delivery (Agile) move in parallel.

5. **Feedback Loops**: Every phase includes feedback mechanisms to catch issues early.

---

## Invocation Examples

### In Pi CLI

```bash
@lisa "Break down the user authentication feature into tasks"

@ralph "Implement the login endpoint with TDD"

@bart "Review this code for security vulnerabilities"

@herb "Verify test coverage meets 95%"

@marge "Conduct a user feedback review of this PR"

@frink "Document the decision to use Event Sourcing"
```

### In Other Harnesses

Load the skill's instructions from `~/.pi/agent/skills/{character}/SKILL.md` and adapt to your agent's command syntax.

---

## Troubleshooting

### "Coverage is low"
→ Use **Herb** to identify gaps → **Ralph** to add tests

### "Code has security issues"
→ Use **Bart** to find them → **Ralph** to fix → **Herb** to verify

### "I don't know where to start"
→ Use **Product Discovery** to validate → **Lisa** to plan → **Ralph** to execute

### "This feels like busywork"
→ Use **Tree of Thoughts** to explore alternatives → **Frink** to review options

### "I'm stuck in a loop"
→ Switch loops (e.g., ReAct if Tree of Thoughts isn't working) → Escalate to **Marge** for scope review

---

## Next Steps

1. Read `LOOP_CATALOG.md` for detailed loop specifications
2. Read `CHARACTER_SKILLS.md` for detailed skill descriptions
3. Review the full `Simpsons.md` (in the obsidian vault) for framework principles
4. Install skills in `~/.pi/agent/skills/` (or your agent harness equivalent)
5. Start with **Lisa** to plan your first task
6. Use **Ralph** to implement
7. Use **Herb** to verify quality
8. Iterate!

---

## Resources

- **LOOP_CATALOG.md** - Comprehensive loop reference with diagrams
- **CHARACTER_SKILLS.md** - Detailed skill descriptions and interfaces
- **Simpsons.md** - Framework vision and principles (Obsidian vault)
- **~/.pi/agent/skills/** - Installed skill implementations
