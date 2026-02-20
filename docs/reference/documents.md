# The 7 Core Documents

These files are the only thing that matters. If an agent learns something and doesn't write it down here, it didn't happen.

---

## 1. PLAN.md (The Roadmap)
**Owner:** @Lisa
**Lifecycle:** Permanent

This is the big picture. It lists the Epics (large chunks of value) we plan to deliver.
*   **Contains:** Feature goals, status, and high-level requirements.
*   **Don't touch if:** You are Ralph. Only Lisa and Marge should be messing with the roadmap.

## 2. TODO.md (The Task List)
**Owner:** @Ralph
**Lifecycle:** Ephemeral (Lives for a sprint, dies at release)

This is the immediate work. It breaks down one Epic into tiny, bite-sized tasks.
*   **Contains:** Checkboxes. Lots of them.
*   **Rule:** If a task takes more than an hour, break it down.
*   **Death:** When the Epic is done, `TODO.md` is deleted.

## 3. Feature.md (The Brief)
**Owner:** @Marge
**Lifecycle:** Permanent (per feature)

The source of truth for "Why" we are doing this.
*   **Contains:** The Problem, The User, The Requirements, and The Unknowns.
*   **Rule:** If you don't understand the user, read this.

## 4. FEEDBACK.md (The Roast)
**Owner:** @Bart
**Lifecycle:** Ephemeral (Lives during review, dies at success)

This is where Bart yells at Ralph.
*   **Contains:** Bugs, security holes, code smell, and general complaints.
*   **Rule:** You cannot ship if this file exists and has content.

## 5. CHANGELOG.md (The History)
**Owner:** @Lovejoy
**Lifecycle:** Permanent

The official record of what we shipped and when.
*   **Contains:** Version numbers, dates, and a summary of value delivered.
*   **Rule:** Follow [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## 6. ADRs (The Decisions)
**Owner:** @Lisa
**Lifecycle:** Permanent

Architectural Decision Records. The "Why" behind the "How."
*   **Location:** `docs/adr/XXXX-title.md`
*   **Contains:** Context, Decision, Consequences.
*   **Rule:** If you make a major technical choice (like "Use Postgres"), write an ADR.

## 7. BDD Specs (The Truth)
**Owner:** @Bart / @Marge
**Lifecycle:** Permanent

Executable requirements written in Gherkin.
*   **Location:** `features/*.feature`
*   **Contains:** `Given`, `When`, `Then` scenarios.
*   **Rule:** These are the definition of "Done." If these pass, the feature works.
