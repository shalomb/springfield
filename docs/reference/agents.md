# Agent & Skills Reference

The Springfield Protocol uses a **5-Agent "Single Pizza" Team** model. Each agent represents a focused mindset and exercises specific standardized skills.

---

## The 5-Agent Team

| Agent | Mindset | Focus | Primary Skills | Primary Persona |
| :--- | :--- | :--- | :--- | :--- |
| **Product** | Empathy | What & Why | `discovery`, `triage` | Marge |
| **Planning** | Logic | How & Structure | `planning`, `architecture` | Lisa |
| **Build** | Optimism | Doing | `implementation`, `testing` | Ralph |
| **Quality** | Pessimism | Critiquing | `review`, `verification` | Bart |
| **Release** | Ceremony | Shipping | `release`, `learning` | Lovejoy |

---

## 1. Product Agent (Marge)
- **Mindset:** Empathy & User Alignment.
- **Role:** Investigates user needs, defines the problem, and enforces the "Definition of Ready."
- **Output:** Feature Briefs (`Feature.md`), Problem Statements.
- **Skills:** `discovery`, `triage`.

## 2. Planning Agent (Lisa)
- **Mindset:** Logical Strategy & Structure.
- **Role:** Breaks down features into tasks, validates architectural fit (ADRs), and plans dependencies.
- **Output:** `PLAN.md`, `TODO.md`, `ADRs`.
- **Skills:** `planning`, `architecture`.

## 3. Build Agent (Ralph)
- **Mindset:** Persistent Optimism & TDD.
- **Role:** Implementation of tasks following strict TDD practices.
- **Output:** Tested code, infrastructure config.
- **Skills:** `implementation`, `testing`.

## 4. Quality Agent (Bart)
- **Mindset:** Adversarial Pessimism ("I'm going to break this").
- **Role:** Finding what the Build Agent missed; enforcing standards.
- **Output:** `FEEDBACK.md`, Gate sign-offs.
- **Skills:** `review`, `verification`.

## 5. Release Agent (Lovejoy)
- **Mindset:** Ceremony & Learning.
- **Role:** Master of publishing and organizational learning.
- **Output:** `CHANGELOG.md`, Release tags.
- **Skills:** `release`, `learning`.

---

## Retired Personas
The following personas have been consolidated into the core 5-agent team to simplify coordination:
- **Troy McClure & Wiggum:** Merged into **Product (Marge)**.
- **Frink:** Merged into **Planning (Lisa)**.
- **Homer:** Merged into **Build (Ralph)**.
- **Herb Powell:** Merged into **Quality (Bart)**.

For detailed persona instructions, see individual profiles in [.github/agents/](../../.github/agents/).
