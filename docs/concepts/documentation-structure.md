# Documentation Structure

Springfield Protocol documentation is organized for discoverability and clarity.

---

## 📍 Start Here

**→ [INDEX.md](INDEX.md)** - Complete navigation guide and overview

---

## 📚 Root Level (6 Files)

Quick reference files for immediate lookup:

| File | Purpose | Time |
|------|---------|------|
| **INDEX.md** | Navigation & documentation index | 5 min |
| **QUICK_START.md** | Workflows, examples, troubleshooting | 20 min |
| **LOOP_CATALOG.md** | All 16+ agentic loops with specifications | 30 min |
| **CHARACTER_SKILLS.md** | All 9 character skill descriptions | 20 min |
| **VISUAL_REFERENCE.md** | 14 ASCII diagrams explaining the system | 15 min |
| **README.md** | High-level framework overview | 10 min |

---

## 🎭 Agent Profiles (`.github/agents/`)

Individual character profiles:

- **lisa.md** - Strategic Planner
- **ralph.md** - TDD Executor
- **bart.md** - Adversarial Reviewer
- **herb.md** - Quality Engineer
- **marge.md** - Empathy & Guardrails
- **frink.md** - Architect
- **wiggum.md** - Triage Officer
- **lovejoy.md** - Release Master
- **troy-mcclure.md** - Chief Discovery Officer

(Linked from [docs/reference/agents.md](docs/reference/agents.md))

---

## 📖 Detailed Documentation (`docs/`)

Organized by Diataxis principles: **Tutorials, How-To Guides, Reference, Explanation**.

### `docs/how-to/` — Goal-Oriented Guides

End-to-end workflows for solving specific problems:

- **implement-feature.md** - Feature from planning to release
- **debug-issue.md** - Investigate and fix problems
- **design-architecture.md** - Make and document architectural decisions
- **release.md** - Publish and announce versions

### `docs/reference/` — Look-Up Information

Quick reference materials:

- **loops.md** - Loop selection guide and quick reference
- **agents.md** - Agent profiles index (links to `.github/agents/`)
- **glossary.md** - Terminology and definitions

### `docs/concepts/` — Understanding & Philosophy

Deep conceptual dives:

- **principles.md** - Core principles explained (5 main + derived)
- **ralph-wiggum-loop.md** - Stateless resampling execution engine explained

---

## 🔍 Navigation Patterns

### "I want to solve a problem"
→ [docs/how-to/](docs/how-to/) - Pick your use case

### "I want to understand the framework"
→ [INDEX.md](INDEX.md) → [QUICK_START.md](QUICK_START.md) → [docs/concepts/](docs/concepts/)

### "I want to look something up"
→ [docs/reference/](docs/reference/) → Pick your type (loops, agents, glossary)

### "I want character details"
→ [docs/reference/agents.md](docs/reference/agents.md) → Pick a character

---

## 🗂️ File Organization Philosophy

**Root Level:** Essential references you return to frequently (6 files)
- Quick to scan
- High signal-to-noise ratio
- Links to deeper documentation

**`.github/agents/`:** Agent definitions (9 files)
- Separate concern (character profiles)
- Linked from documentation
- Can be referenced by agent harnesses

**`docs/`:** Detailed guides (10 files across 3 directories)
- Organized by use case (how-to)
- Organized by lookup type (reference)
- Organized by topic (concepts)
- Diataxis-aligned for clarity

---

## 📊 At a Glance

| Category | Files | Purpose |
|----------|-------|---------|
| Root Quick Ref | 6 | Essential lookups |
| Agent Profiles | 9 | Character definitions |
| How-To Guides | 4 | Problem-solving |
| Reference | 3 | Look-up information |
| Concepts | 2 | Understanding |
| **Total** | **24** | Complete framework |

---

## 🚀 Typical User Journey

```
1. User arrives → reads INDEX.md (overview)
   ↓
2. User reads QUICK_START.md (20 min, 80% of knowledge)
   ↓
3. User needs to solve a problem → goes to docs/how-to/
   ↓
4. User wants to understand something → goes to docs/concepts/
   ↓
5. User wants to look something up → goes to docs/reference/
   ↓
6. User returns to root files for quick reference
```

---

## 🔗 Key Files

| If You Want To... | Go Here |
|------------------|---------|
| Get started | [INDEX.md](INDEX.md) |
| Understand basics | [QUICK_START.md](QUICK_START.md) |
| Learn loops | [LOOP_CATALOG.md](LOOP_CATALOG.md) |
| Meet characters | [CHARACTER_SKILLS.md](CHARACTER_SKILLS.md) |
| See diagrams | [VISUAL_REFERENCE.md](VISUAL_REFERENCE.md) |
| Solve a problem | [docs/how-to/](docs/how-to/) |
| Look something up | [docs/reference/](docs/reference/) |
| Understand concepts | [docs/concepts/](docs/concepts/) |
| Check agent details | [docs/reference/agents.md](docs/reference/agents.md) |

---

## 🎯 Design Decisions

1. **Flat over nested** - Easier navigation, fewer clicks
2. **Linked not copied** - Agent profiles live in `.github/agents/`, referenced from `docs/`
3. **Diataxis alignment** - How-to (goal), Reference (look-up), Concepts (understanding)
4. **Root = quick ref** - 6 essential files you return to frequently
5. **`docs/` = detailed** - Full guides for deeper learning
6. **`.github/` = metadata** - Profiles and planning docs

---

## 📝 Moving Cleaned Up Files

These files were moved to maintain a clean structure:

- `REFINEMENT-NOTES.md` → `.github/REFINEMENT_ROADMAP.md` (evolution planning)
- `START_HERE.md` → Merged into `INDEX.md` (consolidated navigation)
- `core-principles.md` → `docs/concepts/principles.md` (conceptual docs)
- `character-map.md` → `docs/reference/` (reference material)

---

## ✅ Structure Checklist

- ✓ Root level: 6 essential quick reference files
- ✓ Agents: 9 profiles in `.github/agents/`
- ✓ How-to: 4 goal-oriented guides
- ✓ Reference: 3 look-up resources
- ✓ Concepts: 2 philosophical deep-dives
- ✓ Navigation: INDEX.md as single entry point
- ✓ Diataxis: Aligned to how-to / reference / explanation
- ✓ Clean structure: Flat, navigable, discoverable

---

**Ready to explore? Start with [INDEX.md](INDEX.md)** 🚀
