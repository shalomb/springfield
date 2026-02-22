# PLAN.md - Springfield Product Backlog

**Last Updated:** 2026-02-22 11:15 GMT+1  
**Status:** EPIC-005 Phase 3 (Cost Controls & Model Optimization) Active; Phase 2 Complete

---

## 🚀 Current Release: v0.7.0 (Development Cycle)

### EPIC-005 Phase 3: Agent Cost Controls & Model Optimization ⭐ ACTIVE
**td:** td-e1fd16
**Status:** In Progress
**Priority:** P1

**Objective:** Implement budget enforcement and optimize model selection to ensure sustainable autonomous agent operations.

**Scope:**
- [ ] **td-7d836a**: Per-Session Budget Enforcement (Hard limits per session)
- [ ] **td-752039**: Per-Day Budget Tracking (Aggregation and daily caps)
- [ ] **td-9177a6**: Fix cost calculation precision and aggregation
- [ ] **td-aa88b1**: Native JSON Stream Integration (Prerequisite for reliable enforcement)
- [ ] **td-19b078**: Model Selection per Agent (Tuning roles for cost/performance)

**Acceptance Criteria:**
- [ ] Every LLM call is logged with accurate token count and cost.
- [ ] Session is terminated gracefully when budget is exceeded.
- [ ] Users can define budget overrides in `.springfield.yaml`.

---

## ✅ Completed History

### EPIC-005 Phase 2: Robust Structured Output Parsing
**td:** td-7b0cb8
- **Status:** ✅ Done (2026-02-21)
- **Outcome:** Lexical Sanitizer (MarkdownSanitizer) integrated. Promise Semantic Contract enforced. 92%+ test coverage.

### EPIC-009: Springfield Binary Orchestrator & td(1) Integration
**td:** td-3cc3c3
- **Status:** ✅ Done (2026-02-21)
- **Outcome:** Type-safe Go CLI, td(1) integration, multi-agent orchestration.

### EPIC: Unified Agent Runner Architecture
**td:** td-2e90ba
- **Status:** ✅ Done (2026-02-21)
- **Outcome:** Consolidated specialized runners into data-driven Agent struct.

---

## 🗂️ Backlog (Lower Priority)

### Phase 4: Production Operations (v0.8.0 Choice Candidate) 🔍 DISCOVERY
- **Effort:** 3 weeks | **Value:** HIGH
- **Scope:** Docker-based sandboxing (EPIC-004), observability dashboard, persistent storage for agent state.
- **Why:** Moving from local execution to a robust, scalable environment.

### Phase 5: Advanced Intelligence & Governance
- **Option B: Advanced Output Parsing**: Parse DECISION: directives, automated feedback loop.
- **Option C: Enterprise Governance/Audit**: Audit logging, RBAC, Compliance guardrails.
- **Option E: Multi-Agent Swarms**: Dynamic agent creation and coordination.

---

## 📊 Success Metrics (v0.7.0)

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| **Cost Accuracy** | 100% vs Provider API | 85% (estimate) | ⚠️ |
| **Budget Enforcement** | 0 violations | Untested | ⏳ |
| **Test Coverage** | 90%+ | 92%+ | ✅ |

---

## 🚦 Release Gating Criteria (v0.7.0)

**BLOCKERS:**
- [ ] Budget enforcement logic must be unit tested.
- [ ] JSON stream parsing must handle partial/malformed chunks.

---

## 📝 Retrospective: EPIC-005 Phase 2

**Completed Work:**
- ✅ Integrated MarkdownSanitizer prevents code block tag "leaks"
- ✅ Promise contract ensures deterministic loop termination
- ✅ Rebase madness resolved - system stabilized

**Learning:**
- Rebase operations on complex branches with manual resolutions are high-risk; use atomic commits to minimize merge surface.
- Agent prompts need explicit mission-critical instructions preserved across refactors.

---

*Maintained by Lisa Simpson (Planning Agent) with input from the team.*
