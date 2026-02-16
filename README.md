# Springfield Protocol

A character-driven agentic development framework for Lean, iterative product development. Inspired by *The Simpsons*, this protocol combines memorable personas with proven Agile and design thinking practices to create a system that favors learning, feedback, and rapid iteration.

---

## Quick Start

**New to Springfield Protocol?**

1. Start here: [[core-principles.md]] - Vision, philosophy, dual-track framework
2. Meet the team: [[character-map.md]] - All 9 characters and their roles
3. Deep dive: [[troy-mcclure.md]] - Understand the discovery process

**Already know the basics?**

- Want to understand discovery? → [[discovery-track.md]]
- Want to understand delivery? → [[delivery-track.md]]
- Need templates? → [[discovery-brief-template.md]] or [[feature-brief-template.md]]

---

## The Framework at a Glance

### Philosophy
**"Plan before you build, steer as you go."**

Gather available data to reduce uncertainty before committing to build. But stay open to discovery during execution, learn from implementation, and adjust strategy based on evidence.

### Two Tracks

```
DISCOVERY TRACK (Design Thinking)              DELIVERY TRACK (Agile)
─────────────────────────────────              ─────────────────────────────────
Issue/Request                                  Feature Brief
      ↓                                              ↓
Wiggum: Triage & Definition of Ready           Lisa: Planning (PLAN.md)
      ↓                                              ↓
Troy: Gather Intelligence                      Ralph: TDD Execution
      ↓                                              ↓
Frink: Validate Architecture                   Bart: Adversarial Review
      ↓                                              ↓
Marge: Gate & Stakeholder Alignment            Herb: Quality Validation
      ↓                                              ↓
Feature Brief                                  Marge: Merge Gate
                                                    ↓
                                               Lovejoy: Release Ceremony
```

### Core Characters

| Character | Role | Track | Gate? |
|-----------|------|-------|-------|
| **Troy** | Chief Discovery Officer | Discovery | No |
| **Frink** | Architect & Patterns | Discovery | Yes |
| **Marge** | Empathy & Guardrails | Discovery + Merge | Yes |
| **Wiggum** | Triage Officer | Entry | Yes |
| **Lisa** | Strategic Planner | Delivery | No |
| **Ralph** | TDD Executor | Delivery | No |
| **Bart** | Adversarial Reviewer | Delivery | No |
| **Herb** | Quality Engineer | Delivery | Yes |
| **Lovejoy** | Ceremony Master | Release | Soft |

---

## Core Concepts

### The Ralph Wiggum Loop
A stateless resampling loop that iterates toward quality through persistent feedback, not one-shot perfection.

- **Stateless:** Each iteration starts fresh (no context rot)
- **Resampling:** Failures are just prompts for the next iteration
- **Persistent:** Keep iterating until the plan is satisfied

### Plan Before You Build
Conduct discovery first:
1. **Understand the problem** - Use Five Whys, interviews, Gemba walks
2. **Validate assumptions** - Identify unknowns and risks explicitly
3. **Make informed decisions** - Proceed with available data, document uncertainty

### Steer As You Go
Stay engaged during delivery:
1. **Monitor assumptions** - Are they holding up in practice?
2. **Capture learning** - What's different from what we predicted?
3. **Adjust strategy** - If assumptions break, update the plan

---

## Files in This Directory

```
Springfield Protocol/
├── README.md                          ← You are here
├── core-principles.md                 # Vision, philosophy, dual-track framework
├── ralph-wiggum-loop.md               # Core execution engine (coming soon)
├── character-map.md                   # All characters, roles, interactions
├── troy-mcclure.md                    # Chief Discovery Officer (detailed)
├── frink.md                           # Architect & Patterns (stub)
├── lisa.md                            # Strategic Planner (stub)
├── marge.md                           # Empathy & Guardrails (stub)
├── ralph.md                           # TDD Executor (stub)
├── bart.md                            # Adversarial Reviewer (stub)
├── herb.md                            # Quality Engineer (stub)
├── lovejoy.md                         # Ceremony Master (stub)
├── wiggum.md                          # Triage Officer (stub)
├── discovery-track.md                 # Discovery workflow (coming soon)
├── delivery-track.md                  # Delivery workflow (coming soon)
├── discovery-brief-template.md        # Troy's deliverable (coming soon)
└── feature-brief-template.md          # Feature specification (coming soon)
```

**Status:** Core characters are defined. Workflow documents and templates coming soon.

---

## How to Use This Framework

### I Have a Feature Request

1. **Wiggum** triages it and enforces Definition of Ready
2. **Troy** conducts discovery (interviews, Gemba walk, narrative synthesis)
3. **Frink** validates it fits your architecture
4. **Marge** confirms user fit and roadmap alignment
5. **Lisa** plans the delivery
6. **Ralph** executes with TDD
7. **Bart** reviews adversarially
8. **Herb** validates quality
9. **Marge** gates the merge
10. **Lovejoy** releases with ceremony

### I Found a Bug

1. **Wiggum** triages (is this critical or routine?)
2. **Lisa** plans (if straightforward) OR **Troy** investigates (if complex root cause)
3. **Ralph** fixes with TDD
4. **Bart** and **Herb** validate
5. **Lovejoy** releases

### I Want Better Architecture

1. **Frink** leads architecture discussion
2. Propose ADR (Architecture Decision Record)
3. **Marge** gates stakeholder alignment
4. Implement with **Ralph**
5. Validate with **Bart** and **Herb**

---

## Key Principles

### 1. Problems Before Solutions
Discovery focuses on "why" before "what" or "how". Troy uses Five Whys to uncover root causes, not solution-biased requests.

### 2. Explicit Uncertainty
All unknowns are documented. We proceed with available data but acknowledge what we're betting on that might be wrong.

### 3. Learning Loops During Build
Troy stays engaged during delivery. When assumptions break, we update the plan and communicate changes.

### 4. Memorable Personas
Every character has personality, traits, flaws, and decision authority. This makes workflows intuitive and memorable.

### 5. Gatekeeping at Critical Points
Not everyone blocks progress. Specific characters (Wiggum, Frink, Marge, Herb) are gates. Others are advisors/executors.

### 6. Test-Driven Everything
Ralph's TDD discipline, Bart's adversarial review, Herb's coverage enforcement. Quality is built in, not added later.

### 7. Iteration > Perfection
The Ralph Wiggum Loop embraces failure as feedback. Better to iterate quickly than plan perfectly.

---

## Integration with Existing Systems

The Springfield Protocol is designed to integrate with:

- **GitHub Issues** - Issue intake and triage
- **PLAN.md** - Strategic planning and epic tracking
- **TODO.md** - Task lists for Ralph
- **ADRs** - Architecture Decision Records (Frink's domain)
- **Feature Briefs** - Problem specification (Troy's output)
- **Code Reviews** - Bart's adversarial review, Herb's coverage validation
- **git/GitHub** - Version control and release ceremony (Lovejoy)

---

## Next Steps

### To Get Started
1. Read [[core-principles.md]] to understand the philosophy
2. Review [[character-map.md]] to see how characters interact
3. Deep dive into [[troy-mcclure.md]] to understand discovery

### To Implement
1. Define your Definition of Ready (Wiggum's gate)
2. Set up PLAN.md and TODO.md formats (Lisa's domain)
3. Establish TDD practices and coverage requirements (Ralph, Herb)
4. Create issue templates and triage process (Wiggum)
5. Train team on character roles and interactions

### To Expand
1. Complete the stub files (frink.md, lisa.md, etc.)
2. Create workflow documents (discovery-track.md, delivery-track.md)
3. Build templates (discovery-brief, feature-brief)
4. Document ADR decision tree (Frink)
5. Create integration guides for your tools/systems

---

## Questions?

This is a living framework. If something is unclear or needs refinement, that's expected. The Springfield Protocol is designed to evolve as you use it.

**Key contacts by topic:**
- **Discovery:** Troy (see [[troy-mcclure.md]])
- **Architecture:** Frink (see [[frink.md]])
- **Planning:** Lisa (see [[lisa.md]])
- **Quality:** Herb (see [[herb.md]])
- **Release:** Lovejoy (see [[lovejoy.md]])
- **Triage:** Wiggum (see [[wiggum.md]])

---

## Version

**Springfield Protocol v0.1** - Core characters defined, workflows and templates in progress.

Last updated: 2026-02-16

*You might remember me from such frameworks as Agile Development, Design Thinking, and Quality Engineering!*
