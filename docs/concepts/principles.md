# Core Principles

The foundational ideas behind the Springfield Protocol.

## Five Core Principles

### 1. Plan Before You Build
**Gather available data to reduce uncertainty before committing to build.**

- Conduct discovery (interviews, Gemba walks, root-cause analysis)
- Document unknowns explicitly
- Make informed decisions with available information
- Don't try to anticipate unknowns—capture them as you discover them

**In Practice:**
- Troy investigates and synthesizes information
- Frink validates architectural fit
- Marge confirms user alignment
- Then you commit to building

### 2. Steer As You Go
**Stay connected to execution and adjust strategy based on evidence.**

- During delivery, monitor whether assumptions hold
- Capture learning (what's harder than expected? easier?)
- If unknowns become known problems, update the plan
- Treat implementation as a form of discovery

**In Practice:**
- Ralph executes with TDD
- Troy monitors for assumption breakage
- Lisa updates PLAN.md with learnings
- Team adjusts strategy mid-sprint if needed

### 3. Iteration Over Perfection
**Embrace the Ralph Wiggum Loop: fail fast, iterate persistently, ship quality.**

- Don't aim for one-shot perfection
- Each iteration starts clean (ephemeral context)
- Failures are just signals for the next iteration
- Persistent iteration beats brilliant first attempt

**In Practice:**
- Ralph implements, Herb verifies
- If coverage is low, loop back (Ralph adds tests)
- If Bart finds security issues, loop back
- Keep going until PLAN.md is satisfied

### 4. Explicit Uncertainty
**Document what you don't know. Don't pretend to have certainty you don't have.**

- Create "Unknowns Map" during discovery
- Mark decision gates (points where unknowns must be resolved)
- Say "we don't know this, but here's how we'll find out"
- Accept incompleteness; iterate toward clarity

**In Practice:**
- Troy documents unknowns in Feature Brief
- Team agrees on decision gates
- During delivery, unknowns become research tasks
- Update understanding as you learn

### 5. Memorable Personas Drive Behavior
**Agents represent focused modes of thinking, not just roles.**

- **Build Agent** (Optimistic) vs **Quality Agent** (Pessimistic)
- **Product Agent** (Why) vs **Planning Agent** (How)
- Character names (Ralph, Bart) are useful mnemonics for specific skills

### 6. Avoid the Distracted Agent
**Keep agent context focused and expedient.**

- **Don't overload context:** An agent trying to be Architect, Coder, and QA simultaneously performs poorly.
- **Split by purpose:** Use specialized agents (Product, Planning, Build, Quality, Release) to keep prompts lean.
- **Bootstrap context:** Provide sufficient pointers (PLAN.md, Feature.md) without dumping the entire repo.
- **Reference:** [Distracted Agent Anti-Pattern](https://lexler.github.io/augmented-coding-patterns/anti-patterns/distracted-agent/)

---

## Derived Principles

### Problem First, Solution Second
Discovery focuses on **why** before **what** or **how**.
- Troy uses Five Whys to uncover root causes
- Avoid solution-biased thinking
- Understand the problem thoroughly before proposing solutions

### Lean Architecture
No architecture is "perfect." Make decisions that:
- Solve the current problem
- Allow future change
- Are understandable to the team
- Have documented trade-offs

### Quality Built In, Not Added Later
- Ralph uses TDD (tests written first)
- Bart reviews for edge cases and security
- Herb enforces coverage (95%+)
- Quality isn't a phase—it's how you work

### Gatekeeping at Critical Points
Not everyone blocks progress. Specific characters gate decisions:
- **Wiggum** gates issue triage (is this ready to work on?)
- **Frink** gates architecture (does this fit our patterns?)
- **Marge** gates user alignment (does this meet user needs?)
- **Herb** gates quality (is this tested and covered?)
- Others are advisors/executors, not blockers

### Learning Loops During Build
Discovery doesn't end when delivery starts.
- Troy stays engaged during delivery
- Team captures signals: "This is harder than expected"
- If assumptions break, you update the plan
- Learning flows back to refine strategy

---

## What We Don't Do

### We Don't Require a Dedicated Org Chart
The framework describes capabilities, not necessarily separate people. A 2-person team can:
- Both do discovery (wearing Frink/Troy hats)
- Both implement (wearing Ralph hat)
- Both verify (wearing Herb hat)
- Marge's role (alignment) is everyone's responsibility

### We Don't Plan Perfectly Upfront
We gather available data, then iterate. Not waterfall.
- PLAN.md is a living document
- TODO.md changes as learning happens
- Unknowns get resolved during delivery
- Pivot when evidence demands it

### We Don't Ship Perfect Code
We ship tested, verified, quality code that meets the spec.
- Good enough beats perfect
- Technical debt is OK if documented and managed
- Refactor when it matters, not always
- Pragmatism over dogmatism

---

## Principles in Action

### Feature Request Arrives
1. Wiggum triages (is it clear? ready?)
2. Troy investigates (why do users need this?)
3. Frink validates (does it fit our architecture?)
4. Marge confirms (does it align with user needs and roadmap?)
5. **Then** you build

### Bug Is Found
1. Wiggum triages (critical? what type?)
2. Search KEDB (is this a known issue?)
3. If not, Bart or Ralph investigates
4. Fix implemented with test
5. Herb verifies (is it covered? are there regressions?)
6. Lovejoy releases
7. Document in KEDB for next time

### Architecture Decision Needed
1. Frink proposes (researching alternatives)
2. Bart challenges (poke holes, ask tough questions)
3. Frink refines (address concerns)
4. Create ADR (document for posterity)
5. Implement with guardrails (tests that prevent violations)

---

## See Also

- **philosophy.md** - Deeper dive into "Plan before you build, steer as you go"
- **discovery-vs-delivery.md** - The two-track framework explained
- **ralph-wiggum-loop.md** - The core execution engine
- **core-principles.md** (root) - Original detailed principles doc
