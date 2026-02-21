# PLAN.md - Springfield Product Backlog

**Last Updated:** 2026-02-21 19:50 GMT+1  
**Status:** EPIC-005 Phase 2 (Structured Output) Complete; v0.7.0 Planning Active

---

## 🚀 Current Release: v0.6.0-beta (Release in Progress)

### EPIC-005 Phase 2: Robust Structured Output Parsing ✅ COMPLETE
**Status:** Shipped to `feat/epic-005-structured-output` (89 commits)  
**PR:** `feat/epic-005-structured-output` -> `main`

**What Shipped:**
- ✅ Lexical Sanitizer (MarkdownSanitizer) integrated into Agent core
- ✅ Promise Semantic Contract (<promise>COMPLETE</promise>) enforced
- ✅ Real-time output streaming in pi.go
- ✅ 100% ACP compliance and 92%+ test coverage

**Deferred:**
- ⚠️ Task 3: Native JSON Stream Integration (Deferred to v0.7.0 due to pi CLI external dependency)

---

## 📋 Next Release: v0.7.0 (Planning Cycle)

### 🎯 High Priority Path Options (Select One)

#### **Option A: Cost Controls & Model Optimization** (Recommended) ⭐
- **Effort:** 3 weeks | **Value:** HIGH
- **Scope:** Budget Enforcer (per-session/per-day), Model Selection Logic, Cost calculation fixes.
- **Why:** Essential for scaling autonomous agent usage safely.

#### **Option B: Advanced Output Parsing**
- **Effort:** 2 weeks | **Value:** MEDIUM
- **Scope:** Parse DECISION: directives, automated feedback loop for minor issues.

#### **Option C: Enterprise Governance/Audit**
- **Effort:** 4 weeks | **Value:** MEDIUM
- **Scope:** Audit logging, RBAC, Compliance guardrails.

#### **Option D: Production Operations**
- **Effort:** 2 weeks | **Value:** HIGH
- **Scope:** Docker/K8s deployment patterns, monitoring, observability dashboard.

---

## ✅ Completed History

### EPIC-009: Springfield Binary Orchestrator
- **Status:** ✅ Done (2026-02-21)
- **Outcome:** Type-safe Go CLI, td(1) integration, multi-agent orchestration.

### EPIC-2e90ba: Unified Agent Runner Architecture
- **Status:** ✅ Done (2026-02-21)
- **Outcome:** Consolidated specialized runners into data-driven Agent struct.

---

## 🗂️ Backlog (Lower Priority)

### Nice-To-Have Features

| Task | Reason | Status |
|------|--------|--------|
| Temperature parameter support | pi CLI needs --temperature flag | 🔴 DEPRIORITIZED |
| Environment variable overrides | `SPRINGFIELD_MODEL=...` | ⏳ BACKLOG |
| Dynamic model selection | Select model based on task/budget | ⏳ BACKLOG |
| Multi-provider fallback chains | More than 2 fallbacks | ⏳ BACKLOG |
| Agent resource limits | Memory/CPU constraints | ⏳ BACKLOG |

---

## 📊 Success Metrics (v0.6.0)

| Metric | Target | Current | Status |
|--------|--------|---------|--------|
| **Test Coverage** | 90%+ | 92%+ | ✅ |
| **Agents Coordinating** | Lisa→Ralph→Bart→Lovejoy | All 4 working | ✅ |
| **Sanitization** | Zero false tag extractions | Verified | ✅ |
| **Promise Compliance** | Agents must promise | Enforced | ✅ |

---

## 🚦 Release Gating Criteria

**BLOCKERS (must fix before v0.6.0-beta tag):**
- [ ] Lovejoy merge `feat/epic-005-structured-output` -> `main`
- [ ] CHANGELOG.md updated with Phase 2 notes

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
