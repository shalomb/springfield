# Springfield Protocol: Documentation Index

Welcome to the Springfield Protocol v0.2 framework for Agile Agentic Development.

---

## ⚡ 60-Second Overview

Springfield Protocol gives you:

1. **5 Specialized Agents** - Product, Planning, Build, Quality, Release (Single Pizza Team)
2. **9 Standardized Skills** - Reusable capabilities (discovery, architecture, planning, implementation, testing, review, verification, release, learning)
3. **2 Flow Patterns** - Discovery Diamond (Design Thinking) + Delivery Diamond (Agile)
4. **1 Execution Engine** - Ralph Wiggum Loop (stateless resampling for quality)
5. **7 Core Documents** - Shared state: PLAN.md, TODO.md, Feature.md, ADRs, BDD specs, FEEDBACK.md, CHANGELOG.md
6. **Just CLI** - Simple commands: impersonate, utilize, flow, gate, loop

**Why it works:** Keeps agent contexts focused ("Avoid the Distracted Agent"). Scales from 1 person (wearing 5 hats) to large teams.

---

## 🚀 Getting Started

### Option 1: Fast Track (15 minutes)
1. Read **QUICK_START.md** (this file, 80% of what you need)
2. Browse **docs/reference/visual-diagrams.md** (14 ASCII diagrams)
3. Done!

### Option 2: Structured Learning (45 minutes)
1. Read **QUICK_START.md** (workflows)
2. Skim **docs/reference/loop-catalog.md** (sections 1-3)
3. Review **docs/reference/character-skills.md** (overview)

### Option 3: Complete Deep Dive (2-3 hours)
1. **docs/how-to/getting-started.md** - Implementation guide
2. **docs/concepts/model.md** - Complete v0.2 model
3. **docs/concepts/architecture.md** - System architecture
4. **docs/concepts/essential-documents.md** - The 7 documents
5. **docs/reference/loop-catalog.md** - All loops explained
6. **docs/reference/character-skills.md** - All skills detailed
7. **docs/concepts/principles.md** - Core philosophy

---

## 📚 Root-Level Files (Only 3)

| File | Purpose |
|------|---------|
| **INDEX.md** | You are here - navigation hub |
| **README.md** | Framework high-level overview |
| **QUICK_START.md** | Quick reference for common workflows |

---

## 📖 Documentation Structure (Diataxis-Aligned)

All detailed documentation lives in `docs/` organized by purpose:

### `docs/how-to/` — Goal-Oriented Guides
- **getting-started.md** ⭐ - Quick setup + implementation guide
- **implement-feature.md** - Feature workflow (discovery → delivery → release)
- **debug-issue.md** - Debugging workflow
- **design-architecture.md** - Architecture workflow
- **release.md** - Release workflow

### `docs/reference/` — Look-Up Information
- **loop-catalog.md** - Complete agentic loop reference (was LOOP_CATALOG.md)
- **character-skills.md** - All 9 skills detailed (was CHARACTER_SKILLS.md)
- **visual-diagrams.md** - 14 ASCII diagrams (was VISUAL_REFERENCE.md)
- **documents.md** - The 7 core documents (extended reference)
- **loops.md** - Loop selection guide (quick reference)
- **agents.md** - Agent profiles index
- **quick-reference.md** - Quick tips and commands
- **glossary.md** - Terminology definitions

### `docs/concepts/` — Understanding & Philosophy
- **model.md** - Complete v0.2 model (5 agents, 9 skills, 2 diamonds, 1 loop, 7 docs, 1 CLI)
- **architecture.md** - System architecture + data flows + team scaling
- **essential-documents.md** - The 7 core documents explained
- **documentation-structure.md** - This structure explained (was STRUCTURE.md)
- **principles.md** - 5 core principles + derived principles
- **ralph-wiggum-loop.md** - Execution engine deep dive

---

## 🎯 Navigation by Intent

### "I'm brand new, where do I start?"
→ Read **QUICK_START.md** (this file)
→ Then **docs/how-to/getting-started.md**

### "I want to understand the complete model"
→ Read **docs/concepts/model.md** (start here)
→ Then **docs/concepts/architecture.md** (see how it connects)

### "I want to implement this in my project"
→ Read **docs/how-to/getting-started.md** (step-by-step)
→ Follow the setup instructions

### "I need to solve a specific problem"
→ Pick from **docs/how-to/** (implement feature, debug issue, design architecture, release)

### "I want a quick reference"
→ Read **docs/reference/quick-reference.md** or **docs/reference/loops.md**

### "I want to understand a specific concept"
→ Pick from **docs/concepts/** (model, architecture, principles, documents)

### "I want detailed reference information"
→ Pick from **docs/reference/** (loop catalog, character skills, diagrams, glossary)

---

## ⚡ Quick Reference

### The 5 Agents
- **Product Agent** - Discovery, Triage, Problem Definition
- **Planning Agent** - Structure, Architecture, Task Breakdown
- **Build Agent** - Implementation, TDD, Infrastructure
- **Quality Agent** - Adversarial Review, Verification Gates
- **Release Agent** - Ceremony, Changelogs, Learning

### The 9 Skills
- discovery-skill, architecture-skill, planning-skill
- implementation-skill, testing-skill, review-skill
- verification-skill, release-skill, learning-skill

### The 2 Diamonds
- **Discovery Diamond** (Design Thinking) - Problem → Investigate → Validate → Feature Brief
- **Delivery Diamond** (Agile) - Spec → Plan & Build → Verify → Release

### The 7 Documents
- **PLAN.md** (epic roadmap), **TODO.md** (sprint tasks)
- **Feature.md** (requirements), **ADRs** (decisions)
- **scenarios.feature** (acceptance criteria), **FEEDBACK.md** (review results)
- **CHANGELOG.md** (release history)

### The Just CLI
```bash
just impersonate {agent}    # Load agent definition
just utilize {skill}        # Exercise a skill
just flow {phase}           # Navigate discovery/delivery/devops
just gate {checkpoint}      # Check quality gate
just loop                   # Run Ralph Wiggum Loop
just skills|agents|phases   # List available items
```

---

## 🔑 Key Concepts

### Avoid the Distracted Agent
Agents have focused contexts. A Build Agent doesn't worry about release ceremony. This prevents hallucination and context overload.

### Skills are Standardized & Composable
9 reusable capabilities that follow agentskills.io spec. Same skills work across discovery, delivery, and DevOps.

### Documents ARE Shared State
Markdown files (PLAN.md, Feature.md, ADRs, etc.) are the single source of truth. No database, no message queue.

### Ralph Wiggum Loop is the Engine
Stateless resampling: each iteration starts fresh. Prevents context rot, enables learning through iteration.

### Same Framework for Everything
Discovery + Delivery + DevOps all use same 2 diamonds, same 9 skills, same 7 documents.

---

## 📍 File Map (Where Everything Is)

```
.
├── INDEX.md                      (you are here)
├── README.md                     (framework overview)
├── QUICK_START.md                (quick reference)
│
├── docs/
│   ├── how-to/                   (goal-oriented guides)
│   │   ├─ getting-started.md     ⭐ START HERE
│   │   ├─ implement-feature.md
│   │   ├─ debug-issue.md
│   │   ├─ design-architecture.md
│   │   └─ release.md
│   │
│   ├── reference/                (look-up information)
│   │   ├─ loop-catalog.md        (complete reference)
│   │   ├─ character-skills.md    (all skills)
│   │   ├─ visual-diagrams.md     (14 diagrams)
│   │   ├─ quick-reference.md     (quick tips)
│   │   ├─ documents.md
│   │   ├─ loops.md
│   │   ├─ agents.md
│   │   └─ glossary.md
│   │
│   └── concepts/                 (understanding & philosophy)
│       ├─ model.md               (complete v0.2 model)
│       ├─ architecture.md        (system architecture)
│       ├─ essential-documents.md (the 7 docs)
│       ├─ documentation-structure.md
│       ├─ principles.md
│       └─ ralph-wiggum-loop.md
│
├── .github/
│   ├── agents/                   (agent definitions)
│   ├── skills/                   (skill definitions - to create)
│   ├── CHANGES_SUMMARY.md        (project history)
│   └── REFINEMENT_ROADMAP.md     (framework evolution)
│
├── PLAN.md, TODO.md, CHANGELOG.md, FEEDBACK.md
└── justfile (to be created)
```

---

## ❓ Common Questions

**Q: Why only 3 files at root?**
A: Clean, discoverable. Everything detailed lives in `docs/` organized by Diataxis (how-to, reference, concepts).

**Q: What's the difference between 5 agents and 9 skills?**
A: Agents are the "Who" (context/focus). Skills are the "What" (capability). An agent exercises 1-2 skills at a time.

**Q: Do I need 5 people?**
A: No. One person can wear all 5 hats. The separation is for mental context switching and agent context efficiency.

**Q: What makes the Ralph Wiggum Loop special?**
A: Stateless resampling. Each iteration starts fresh (clean context). Prevents hallucination/error accumulation. Enables learning through persistent iteration.

**Q: Can I use just part of this framework?**
A: Absolutely. Use just the 7 documents. Use just the Ralph Wiggum Loop. Use just the skills. Mix and match.

**Q: How does this compare to v0.1?**
A: v0.2 simplifies from 9 specialized characters to 5 focused agents. Scales better. Prevents "distracted agent" anti-pattern.

---

## 🚀 Next Steps

1. **Read:** [QUICK_START.md](QUICK_START.md) (20 min)
2. **Understand:** [docs/concepts/model.md](docs/concepts/model.md) (30 min)
3. **Implement:** [docs/how-to/getting-started.md](docs/how-to/getting-started.md) (follow steps)
4. **Reference:** Use [docs/reference/](docs/reference/) as needed

---

## 📚 Learning Paths by Role

### Software Developer
1. QUICK_START.md
2. docs/how-to/implement-feature.md
3. docs/reference/character-skills.md
4. docs/reference/loop-catalog.md

### Product Manager
1. QUICK_START.md
2. docs/concepts/model.md
3. docs/how-to/getting-started.md
4. docs/concepts/essential-documents.md

### Architect
1. docs/concepts/model.md
2. docs/concepts/architecture.md
3. docs/how-to/design-architecture.md
4. docs/reference/loop-catalog.md

### Team Lead
1. INDEX.md (this page)
2. docs/concepts/model.md
3. docs/concepts/architecture.md (team scaling section)
4. docs/how-to/getting-started.md

---

**Ready? Start with [QUICK_START.md](QUICK_START.md) or [docs/how-to/getting-started.md](docs/how-to/getting-started.md)** 🚀
