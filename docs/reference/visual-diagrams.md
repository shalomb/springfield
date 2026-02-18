# Springfield Protocol: Visual Reference Guide

Diagrams, flowcharts, and ASCII art for understanding the Springfield Protocol at a glance.

---

## 1. The Ralph Wiggum Loop (Core Engine)

```
                    ┌──────────────────────────────────┐
                    │      PLAN.json (Source of Truth) │
                    │  [{ task: "...", passing: false},│
                    │   { task: "...", passing: true } │
                    │                                  │
                    └──────────────┬───────────────────┘
                                   │
                                   │
                    ┌──────────────▼──────────────┐
                    │     CONTROL LOOP (Scheduler)│
                    │  while (exists failing task)│
                    │    Spawn Agent for task[0]  │
                    │                              │
                    └──────────────┬───────────────┘
                                   │
                   ┌───────────────┴──────────────┐
                   │                              │
        ┌──────────▼────────────┐   ┌────────────▼───────────┐
        │  Ephemeral Context    │   │  Verification Loop     │
        │  (Clean Git Worktree) │   │  (Bart Review)         │
        │                       │   │                        │
        │  ┌─────────────────┐  │   │  ┌──────────────────┐  │
        │  │ Red-Green-      │  │   │  │ Test Coverage    │  │
        │  │ Refactor (TDD)  │  │   │  │ Security Scan    │  │
        │  │                 │  │   │  │ Quality Metrics  │  │
        │  │ Atomic Commits  │  │   │  │ Edge Cases       │  │
        │  └─────────────────┘  │   │  └──────────────────┘  │
        │                       │   │                        │
        └───────────┬───────────┘   └────────────┬───────────┘
                    │                             │
                    └────────────┬────────────────┘
                                 │
                       ┌─────────▼─────────┐
                       │ Update PLAN.json  │
                       │ (passing: true/   │
                       │  false)           │
                       └────────┬──────────┘
                                │
                                │ [Loop back]
                                │
                       ┌────────▼─────────┐
                       │ All tasks        │
                       │ passing?         │
                       └────────┬─────────┘
                                │
                        ┌───────┴────────┐
                        │                │
                      No              Yes
                        │                │
                    [LOOP]         ┌─────▼─────┐
                                   │  SUCCESS  │
                                   │  Verified │
                                   │  Feature  │
                                   └───────────┘
```

**Key Innovation:** No persistent agent context. Each iteration starts fresh, preventing hallucination accumulation.

---

## 2. Discovery vs. Delivery Tracks

```
┌─────────────────────────────────────────────────────────────────────┐
│                         PRODUCT LIFECYCLE                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  DISCOVERY TRACK (Design Thinking)         DELIVERY TRACK (Agile) │
│  ════════════════════════════════════════  ═════════════════════════│
│                                                                     │
│  1. Empathize & Research                   1. Plan → TODO.md       │
│     (Product Agent)                           (Lisa)                │
│        ↓                                       ↓                    │
│  2. Define Problem                         2. Architect & ADR       │
│     (Product Agent)                           (Lisa)                │
│        ↓                                       ↓                    │
│  3. Ideate Solutions                       3. Implement & Test      │
│     (Marge + Lisa)                           (Ralph + TDD)         │
│        ↓                                       ↓                    │
│  4. Prototype & Validate                   4. Review & Verify       │
│     (Marge + Ralph)                          (Bart)                │
│        ↓                                       ↓                    │
│  5. Feature Brief (READY!)                 5. User Alignment Check   │
│     ─────────────────────────────────────────────→                │
│                                                   (Marge)           │
│                                                   ↓                 │
│                                           6. Release & Publish      │
│                                              (Lovejoy)             │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 3. Character Skill Map (Who Does What)

```
                      FEATURE REQUEST
                            │
                            ▼
                    ┌───────────────┐
                    │  MARGE (Prod) │ ← Triage & Investigate
                    └───────┬───────┘
                            │
                ┌───────────┴────────────┐
                │                        │
                ▼                        ▼
        ┌──────────────┐        ┌──────────────┐
        │  MARGE       │        │  LISA        │ ← Design & ADR
        │  (Discovery) │        │  (Planning)  │
        └──────┬───────┘        └──────┬───────┘
               │                       │
               └───────────┬───────────┘
                           │
                           ▼
        ┌──────────────────────────────────┐
        │  RALPH WIGGUM LOOP BEGINS        │
        └─────────┬──────────────────────┬─┘
                  │                      │
         ┌────────▼─────┐      ┌────────▼──────┐
         │ RALPH        │      │ Loop for each │
         │ (Build)      │      │ task in       │
         │ TDD Loop     │      │ TODO.md       │
         └────────┬─────┘      └───────────────┘
                  │
                  ▼
         ┌────────────────────┐
         │ BART (Quality)     │ ← Review & Verify
         │ Adversarial/TDD    │
         └────────┬───────────┘
                  │
     ┌────────────┴────────────┐
     │                         │
   FAIL                      PASS
     │                         │
     └─────────┬──────────────┘
               │ [Loop back]
               │
               ▼
        ┌─────────────────────┐
        │ MARGE (Product)     │ ← User alignment
        │ PR Review & Feedback│
        └──────────┬──────────┘
                   │
                   ▼
        ┌──────────────────────┐
        │ LOVEJOY (Release)    │ ← Release & publish
        └──────────────────────┘
```

---

## 4. Agentic Loop Taxonomy

```
                         AGENTIC LOOPS
                              │
                ┌─────────────┬┴──────────────┐
                │             │               │
           FOUNDATIONAL    REFINEMENT    EXPLORATION
           LOOPS           LOOPS         & DECISION
                │             │               │
    ┌───────────┼──┐  ┌──────┼────┐  ┌──────┼────┐
    │           │  │  │      │    │  │      │    │
 Sense-Plan-  ReAct │  GECR TALAR │ ToT   GSPE VS
   Act              │          OHECI│
                    │
                    └─────────────────────┐
                                          │
                                   ORCHESTRATION
                                   & COLLABORATION
                                          │
                    ┌─────────────┬───────┼─────┐
                    │             │       │     │
              Plan-and-      Ralph Wiggum  Manager-  Dialogue
              Execute        Loop        Worker     Loop
```

---

## 5. Loop Decision Tree (Simplified)

```
                         START
                          │
                          ▼
                    "What's the problem?"
                          │
                ┌─────────┬┼─┬───────┬──────┐
                │         ││ │       │      │
            Vague &   Specific Known Clear  High
            Complex   Error    Solution Spec Quality
             Error    ?        Exists  ?    Needed
              │       │         │      │      │
              ▼       ▼         ▼      ▼      ▼
           Tree of  ReAct    Use KEDB Plan-  GECR
           Thoughts (Debug)  or       and-   (Refine
           (Multi-path SEARCH Execute  &
            pruning)                   Validate)
```

---

## 6. Multi-Agent Orchestration Pattern

```
                      ┌─────────────┐
                      │   MANAGER   │
                      │  (e.g., Lisa│
                      │ or Lovejoy) │
                      └────────┬────┘
                               │
                 ┌─────────────┼─────────────┐
                 │             │             │
            ┌────▼───┐    ┌────▼───┐   ┌───▼────┐
            │ Worker │    │ Worker │   │ Worker │
            │ (Ralph)│    │ (Bart) │   │ (Marge)│
            │        │    │        │   │        │
            │ Build  │    │ Quality│   │ Product│
            └────┬───┘    └────┬───┘   └───┬────┘
                 │             │           │
                 └─────────────┬───────────┘
                               │
                    ┌──────────▼──────────┐
                    │ Manager aggregates  │
                    │ and coordinates     │
                    │ next phase          │
                    └────────────────────┘
```

**Parallelism:** Workers can execute in parallel; Manager coordinates next round.

---

## 7. Feature Implementation Workflow

```
┌──────────────────────────────────────────────────────────────────┐
│                    FEATURE IMPLEMENTATION                        │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Feature Brief (Validated)                                   │
│     ─────────────────────────────────────────────┐              │
│                                                  │              │
│  2. LISA: Plan & Task Breakdown                 │              │
│     Outputs: PLAN.md → TODO.md                  │              │
│                                                  │              │
│  3. RALPH WIGGUM LOOP ──┐                       │              │
│     For each task:      │                       │              │
│       ├─ Red (Test)     │                       │              │
│       ├─ Green (Code)   │                       │              │
│       ├─ Refactor       │                       │              │
│       ├─ Commit         │                       │              │
│       └─ Next task ─────┘                       │              │
│                                                  │              │
│  4. BART: Adversarial Review                    │              │
│     └─ Find security holes, edge cases         │              │
│                                                  │              │
│  5. HERB: Quality Verification                 │              │
│     └─ 95%+ coverage, mock-first               │              │
│                                                  │              │
│  6. MARGE: User Alignment                      │              │
│     └─ Feature meets user needs?               │              │
│                                                  │              │
│  7. LOVEJOY: Release & Publish                 │              │
│     └─ Semantic version, CHANGELOG, registry   │              │
│                                                  │              │
│  ✓ DONE: Production-ready feature              │              │
│          with quality, coverage, docs          │              │
│                                                  │              │
└──────────────────────────────────────────────────────────────────┘
```

---

## 8. Feedback Loops & Backpressure

```
                    Quality Metrics
                         │
            ┌────────────┬┴───────────┐
            │            │            │
        Coverage      Tests      Security
         95%+         Pass?        Check
            │            │            │
            └────────────┬────────────┘
                         │
                         ▼
                  ┌──────────────┐
                  │ Quality >    │
                  │ Threshold?   │
                  └──────┬───────┘
                         │
              ┌──────────┴──────────┐
              │                     │
            YES                    NO
              │                     │
              ▼                     ▼
         Continue       ┌───────────────────┐
         Work           │ BACKPRESSURE      │
                        │ - Pause work      │
                        │ - Investigate     │
                        │ - Fix issues      │
                        │ - Retest          │
                        └────────┬──────────┘
                                 │
                                 └─────[LOOP back]
```

---

## 9. Context Degradation Prevention

```
Agent A starts       Agent A (5 iterations)
      │                     │
      ├─→ Task 1            ├─ Hallucinations ↑
      ├─→ Task 2            ├─ Context drift ↑↑
      ├─→ Task 3            └─ Quality ↓↓
      ├─→ Task 4
      └─→ Task 5
          (Poor Quality)

                    VS

Ralph Wiggum Loop
      │
      ├─→ [CLEAN] Task 1 ✓
      │   [Verify & Update]
      │
      ├─→ [CLEAN] Task 2 ✓
      │   [Verify & Update]
      │
      ├─→ [CLEAN] Task 3 ✓
      │   [Verify & Update]
      │
      ├─→ [CLEAN] Task 4 ✓
      │   [Verify & Update]
      │
      └─→ [CLEAN] Task 5 ✓
          [Verified, High Quality]

Key: Ephemeral context + stateless = No degradation!
```

---

## 10. Knowledge Accumulation

```
            Work Happens
                 │
                 ▼
        ┌────────────────┐
        │ Reflect &      │
        │ Extract        │
        │ Insights       │
        └────────┬───────┘
                 │
      ┌──────────┴───────────┐
      │                      │
      ▼                      ▼
   ADR             KEDB (Known Error DB)
   (Architecture   (Troubleshooting)
    Decisions)
      │                      │
      │  ┌──────────────────┐│
      │  │ Future Work Can  ││
      │  │ Reuse Patterns   ││
      └─→│ & Solutions      ││
         │ Faster           ││
         └──────────────────┘│
                             │
                    Organization
                    Learning
                    Accelerates
```

---

## 11. Character Trait Mapping (5-Agent Team)

```
        Marge (Product)         │        Lisa (Planning)
        "The What & Why"        │        "The How & Structure"

        Ralph (Build)           │        Bart (Quality)
        "The Doing & TDD"       │        "The Review & Verify"

                        Lovejoy (Release)
                        "The Ship & Learn"
```

---

## 12. Time Horizons & Loop Types

```
              IMMEDIATE              SHORT-TERM           LONG-TERM
              (minutes-hours)        (hours-days)         (days-weeks+)
                    │                    │                    │
                    │                    │                    │
              Sense-Plan-Act      Plan-and-Execute    Ralph Wiggum Loop
                    │                    │                    │
              ReAct (Debug)         Tree of Thoughts   Manager-Worker
                    │                    │                    │
              Dialogue                  GECR            Multi-epic
              (immediate)                │            Orchestration
                    │              TALAR (Learning)          │
                    │              Ralph variant             │
                    │            (per-task QA)              │
                    │                    │                    │
                └────────────────────────┴────────────────────┘
                        All can be strung together!
```

---

## 13. The Full Picture (System Architecture)

```
┌────────────────────────────────────────────────────────────────────┐
│                    SPRINGFIELD PROTOCOL                           │
│                    (Agile Agentic Development)                    │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  INPUTS                  PROCESS                   OUTPUTS         │
│  ═════════              ═══════════               ════════════     │
│                                                                    │
│  Feature                 Discovery Track          Feature Brief   │
│  Request                 (Design Thinking)        (Validated)     │
│        │                       │                        │         │
│        ├───────────────────────┤                        ▼         │
│        │                       └──→ PLAN.md            │         │
│        │                                               │         │
│        │                       Delivery Track          │         │
│        │                       (Agile + Ralph          │         │
│        │                        Wiggum Loop)           │         │
│        │                            │                 │         │
│        └────────────────────────────┤                 │         │
│                                      │                 │         │
│                              Ralph Wiggum Loop         │         │
│                              [Stateless Resampling]   │         │
│                                      │                 │         │
│                              Verified, Tested Code    │         │
│                                      │                 │         │
│                                      ├────────────────→│         │
│                                      │                 │         │
│                              Released Version         │         │
│                              in Registry              │         │
│                                                        │         │
│  Knowledge Base Updates                                ▼         │
│  (ADRs, KEDB, Runbooks)                          Production!     │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```

---

## 14. Quality Metrics & Progression

```
Progress through the Ralph Wiggum Loop:

Task State                 Quality Metrics              Status
═══════════════════════════════════════════════════════════════════

Task Assigned              ❌ Untested
                          ❌ 0% coverage
                          ❌ No verification
                                              → [In Progress]

Red Phase                  ❌ Tests failing
                          ❌ 0% coverage
                          ⚠️  Test written
                                              → [In Progress]

Green Phase                ✓ Tests passing
                          ⚠️  Low coverage
                          ✓ Code written
                                              → [In Progress]

Refactor Phase             ✓ Tests passing
                          ⚠️  Coverage improving
                          ✓ Code optimized
                                              → [In Progress]

Bart Review                ✓ Tests passing
                          ⚠️  No security issues
                          ⚠️  Edge cases found
                                              → [Needs Fix]
                                              → [In Progress]

Herb Verification          ✓ Tests passing
                          ✓ 95%+ coverage
                          ✓ No mock issues
                                              → [Complete]

✓ Task Complete & Verified
```

---

**Ready to use these? Start with QUICK_START.md and refer back to this guide as needed!**
