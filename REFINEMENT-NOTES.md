# Springfield Protocol: Refinement Direction

**Status:** v0.1 complete with 9 characters. Ready for refinement toward lean, capability-based model.

**Date:** 2026-02-16

---

## Intent: From Roles to Capabilities

The current Springfield Protocol uses 9 distinct characters (Wiggum, Troy, Frink, Marge, Lisa, Ralph, Bart, Herb, Lovejoy) to represent specialized roles in a development workflow.

**Key Insight for Refinement:**

The framework should NOT prescribe an org chart or require specialized roles. Instead, it should enable **autonomous teams** where:

- ✅ **Capabilities are distributed** (architecture, DevOps, QA are skills the team possesses, not dedicated roles)
- ✅ **Aligned yet autonomous** (team members exercise their capabilities within shared principles)
- ✅ **Principles drive behavior** (not role titles)
- ✅ **Lean and minimal** (only what's essential to the loops)

---

## Current State (v0.1)

**9 Characters = Enterprise Org Chart:**
- Wiggum (Triage)
- Troy (Discovery)
- Frink (Architecture)
- Marge (Empathy/Gates)
- Lisa (Planning)
- Ralph (TDD Execution)
- Bart (Adversarial Review)
- Herb (Quality)
- Lovejoy (Release Ceremony)

**Problem:** Implies you need to hire or be these specific roles to use the framework.

---

## Target State (v0.2)

**Lean Core: Principles + Loops**

**Two Guiding Characters:**
1. **Ralph Wiggum** - The Loop itself (represents TDD, stateless resampling, persistent iteration)
2. **Troy McClure** - The Discovery Mindset (represents investigation, learning, pragmatic decision-making)

**Essential Capabilities (Not Roles):**
1. **Architecture Validation** - Someone on team applies this capability
2. **Quality Verification** - Someone on team applies this capability (95%+ coverage + adversarial review)
3. **Stakeholder Alignment & Gating** - Someone on team owns this
4. **TDD Execution & Planning** - Team executes this collaboratively

**Three Core Workflows:**
1. **Discovery Flow** - How teams investigate problems (Troy mindset: Five Whys, Gemba walks, unknowns mapping)
2. **Delivery Flow** - How teams execute (Ralph loop: TDD → review → verify → learn)
3. **Learning Loop** - How teams monitor and adjust (during delivery: capture signals, update strategy)

---

## Refinement Roadmap

### Phase 1: Simplify Core
- [ ] Reduce from 9 to 2-3 guiding characters (Ralph, Troy, possibly one bridge)
- [ ] Remove role-based characters (Wiggum, Lisa, Bart, Herb, Lovejoy become capabilities, not characters)
- [ ] Document as **Principles + Capabilities**, not Org Chart

### Phase 2: Restructure Documentation
- [ ] Replace character profiles with capability specifications
- [ ] Create workflow docs focused on loops, not roles
- [ ] Show how **any team structure** can implement these principles

### Phase 3: Add Tooling Notes
- [ ] How to apply framework with different team sizes (2-person startup to large org)
- [ ] Decision trees for capability assignment ("Who does architecture review?")
- [ ] Integration with actual PLAN.md, TODO.md, AGENTS.md files

### Phase 4: Validate Against Philosophy
- [ ] ✅ Plan before you build (discovery loop)
- [ ] ✅ Steer as you go (learning loop)
- [ ] ✅ Accept incompleteness (unknowns mapping, decision gates)
- [ ] ✅ Autonomous teams (capabilities, not roles)

---

## Key Principles to Preserve

**From Current v0.1:**
- ✅ The Ralph Wiggum Loop (stateless resampling, persistent iteration)
- ✅ Discovery mindset (Five Whys, Gemba walks, narrative synthesis)
- ✅ Explicit uncertainty capture (unknowns map, decision gates)
- ✅ Learning during delivery (Troy monitors, Lisa adjusts)
- ✅ Quality gates (architecture, coverage, user fit)
- ✅ "Plan before you build, steer as you go"

**To Remove:**
- ❌ Role-based thinking (Wiggum = triage person, Lisa = planner person)
- ❌ Specialized characters (Lovejoy as a person, not a process)
- ❌ Org chart implications (implies you need 9 different people/skills)

---

## Questions for Refinement

1. **Should Ralph and Troy be the only "characters"?** Or should we have 3-4 guiding principles without characters?

2. **How do we represent capabilities?** Should they be:
   - Specific "hats" a team member wears during different phases?
   - Skill sets documented separately?
   - Decision checkpoints in the workflow?

3. **What about Marge (the bridge/gatekeeper)?** Does this capability need its own character, or is it implicit in "stakeholder alignment"?

4. **For release ceremony:** Should this be a process/checklist, or keep Lovejoy as a lightweight character representing "release hygiene"?

5. **Should we create a "team composition guide"?** Examples like:
   - "2-person team: Both do discovery, both do execution, both verify"
   - "5-person team: One leads discovery, three execute, one verifies"
   - "Enterprise: Discovery team, Delivery team, Quality team—aligned through principles"

---

## Files to Refactor

**Keep:**
- ✅ `README.md` - Update navigation
- ✅ `core-principles.md` - Core to everything
- ✅ `troy-mcclure.md` - Guiding mindset (keep refined)
- ✅ `ralph-wiggum-loop.md` (to be created) - Core engine

**Consolidate/Remove:**
- ❌ `character-map.md` - Replace with "Capabilities & Workflows"
- ❌ `wiggum.md` - Become "Triage Capability"
- ❌ `frink.md` - Become "Architecture Validation Capability"
- ❌ `lisa.md` - Become "Planning & Orchestration Capability"
- ❌ `marge.md` - Become "Stakeholder Alignment Capability"
- ❌ `ralph.md` - Keep as "Ralph Loop" concept, not role
- ❌ `bart.md` - Become "Adversarial Review Practice"
- ❌ `herb.md` - Become "Quality Verification Capability"
- ❌ `lovejoy.md` - Become "Release Ceremony Process"

**New Files to Create:**
- 📄 `ralph-wiggum-loop.md` - Core execution engine explained
- 📄 `discovery-workflow.md` - How to investigate (Troy mindset)
- 📄 `delivery-workflow.md` - How to execute (Ralph loop)
- 📄 `capabilities.md` - Essential team capabilities (not roles)
- 📄 `team-compositions.md` - How different team structures use this

---

## Next Session

When ready to refactor, create a **proper git repository** separate from the Obsidian vault:

```bash
git clone <repo>
cd springfield-protocol
git checkout -b refactor/v0.2-lean-capabilities
# Make changes
git commit -am "refactor: shift from roles to capabilities"
git push origin refactor/v0.2-lean-capabilities
```

This allows:
- ✅ Clean versioning of the framework itself
- ✅ Collaboration (PRs, reviews)
- ✅ History tracking (why decisions were made)
- ✅ Separation from project-specific vault

---

## Summary

**Current:** 9-character org chart (enterprise model)  
**Target:** 2 guiding characters + Capabilities-based workflows (lean, autonomous teams)  
**Benefit:** Teams of any size can use the framework without needing specialized roles

Ready to refactor when you are. 🚀
