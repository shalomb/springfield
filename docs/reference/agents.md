# Agent & Skills Reference

The Springfield Protocol uses a **5-Agent "Single Pizza" Team** model. Each agent represents a focused mindset and exercises specific standardized skills.

---

## The 5-Agent Team

| Agent | Mindset | Focus | Primary Skills | Personas (Personas) |
| :--- | :--- | :--- | :--- | :--- |
| **Product** | Empathy | What & Why | `discovery`, `triage` | Troy, Marge |
| **Planning** | Logic | How & Structure | `planning`, `architecture` | Lisa, Frink |
| **Build** | Optimism | Doing | `implementation`, `testing` | Ralph, Homer |
| **Quality** | Pessimism | Critiquing | `review`, `verification` | Bart, Herb |
| **Release** | Ceremony | Shipping | `release`, `learning` | Lovejoy |

---

## 1. Product Agent ("The What & Why")
- **Role:** Investigates user needs, defines the problem, and enforces the "Definition of Ready."
- **Output:** Feature Briefs (`Feature.md`), Problem Statements.
- **Key Skills:**
    - `discovery-skill`: Interviews, Five Whys, Gemba walks.
    - `triage-skill`: Categorizing and prioritizing requests.

## 2. Planning Agent ("The How & Structure")
- **Role:** Breaks down features into tasks, validates architectural fit, and plans dependencies.
- **Output:** `PLAN.md`, `TODO.md`, `ADRs`.
- **Key Skills:**
    - `planning-skill`: Task breakdown, estimation, sequencing.
    - `architecture-skill`: Design decision records (ADRs), pattern validation.

## 3. Build Agent ("The Doer")
- **Role:** Optimistic implementation of tasks following strict TDD practices.
- **Output:** Tested code, infrastructure config.
- **Key Skills:**
    - `implementation-skill`: Writing Red-Green-Refactor code.
    - `testing-skill`: Writing mocks and unit tests.

## 4. Quality Agent ("The Critic")
- **Role:** Pessimistic adversarial reviewer focused on finding what the Build Agent missed.
- **Output:** `FEEDBACK.md`, Gate sign-offs.
- **Key Skills:**
    - `review-skill`: Security scans, edge case analysis, adversarial questioning.
    - `verification-skill`: Coverage checks (>95%), performance regression checks.

## 5. Release Agent ("The Shipper")
- **Role:** Ceremonial master of publishing and organizational learning.
- **Output:** `CHANGELOG.md`, Release tags, KEDB entries.
- **Key Skills:**
    - `release-skill`: Versioning, changelog automation, deployment ceremony.
    - `learning-skill`: Capturing insights and updating the Known Error Database.

---

## Persona Mapping (The Simpsons Mnemonics)

Characters are useful mental models for how an agent should behave:
- **Ralph Wiggum:** Persistent effort, Red-Green-Refactor.
- **Bart Simpson:** Adversarial thinking, "I'm going to break this."
- **Lisa Simpson:** Strategic, logical, planning everything.
- **Troy McClure:** Investigations, storytelling, synthesis, triage, "Definition of Ready."
- **Herb Powell:** Quality obsession, high standards.
- **Frink:** Scientific method, architectural patterns.
- **Marge Simpson:** User empathy, the "voice of reason."
- **Rev. Lovejoy:** Release ceremony, public announcements.

For detailed persona instructions, see individual profiles in [.github/agents/](../../.github/agents/).
