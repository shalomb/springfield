# The Master Model

This is the brain of the operation.

The Springfield Protocol isn't just a random collection of scripts. It's a way of working designed to keep AI agents from hallucinating, getting distracted, or breaking things.

---

## 1. Core Principles

We built this on nine simple rules.

### 1. See the whole board
Don't just write code. Understand where it fits. We plan before we build so we don't build the wrong thing.

### 2. No hand-offs
We don't pass batons. We leave notes. Agents read the "Blackboard" (our core docs) to know what to do next. This lets everyone work at their own speed.

### 3. Fast feedback
We find bugs immediately, not next week. The **Ralph Wiggum Loop** ensures that for every builder, there's a critic watching closely.

### 4. Treat everything as an experiment
We don't "implement requirements." We test hypotheses. Every feature is a bet. If it works, we keep it. If not, we learn and try again.

### 5. Perfection through iteration
We don't expect to get it right the first time. We expect to get it right eventually. We start fresh often to avoid getting stuck in a bad path.

### 6. Admit what you don't know
If you're unsure, say so. We document "Unknowns" in `Feature.md` instead of guessing. A documented risk is better than a confident hallucination.

### 7. Stay focused
Agents get distracted easily. We split roles (e.g., Build vs. Quality) so no single agent has to hold the entire world in its head at once.

### 8. Use different brains
We use different personas to solve problems. Marge thinks about users. Lisa thinks about systems. By using different "biases," we get better solutions than one generic AI could produce.

### 9. If it's not written down, it doesn't exist
Markdown files are our brain. If an agent learns something but doesn't write it to a file, that knowledge dies when the agent shuts down.

---

## 2. The Team (The 5 Agents)

We split the work into five distinct roles.

| Who | Role | Mindset | What they care about |
| :--- | :--- | :--- | :--- |
| **Product** | @Marge | Empathetic | The User. The "Why." |
| **Planning** | @Lisa | Logical | The System. The "How." |
| **Build** | @Ralph | Optimistic | The Code. Getting it done. |
| **Quality** | @Bart | Pessimistic | The Bugs. Breaking things. |
| **Release** | @Lovejoy | Ceremonial | The Ship. Telling the world. |

---

## 3. The Workflow (Think, then Build)

We work in two phases. We call them "Diamonds."

### Phase 1: Discovery (The "What")
*   **Marge** looks at the problem and asks "Why?"
*   **Lisa** looks at the options and picks the best one.
*   **Result:** A **Feature Brief** that says exactly what we're going to build.

### Phase 2: Delivery (The "How")
*   **Ralph** writes the code.
*   **Bart** tries to break it.
*   **Result:** Working, tested software.

---

## 4. The Engine (Ralph Wiggum Loop)

How do we actually execute this?

1.  **Check the Plan:** We look at `PLAN.md` to see what's next.
2.  **Start Fresh:** We spawn an agent in a clean environment. No baggage.
3.  **Do the Work:** The agent writes code or runs tests.
4.  **Write it Down:** The agent updates the docs (`TODO.md`, `FEEDBACK.md`).
5.  **Repeat:** If it failed, we try again.

---

## 5. The Brain (Shared State)

These 7 files are the only things that matter. They keep us aligned.

1.  **PLAN.md:** The Roadmap.
2.  **TODO.md:** The Task List.
3.  **Feature.md:** The Spec.
4.  **ADRs:** The Decisions.
5.  **BDD Specs:** The Tests.
6.  **FEEDBACK.md:** The Review.
7.  **CHANGELOG.md:** The History.

---

## 6. The Big Picture

```
Problem → [Think] → Feature Brief → [Build] → Release
         (Marge/Lisa)              (Ralph/Bart)   (Lovejoy)
```

Want to see how it all connects? Check out the [Visual Diagrams](../reference/visual-diagrams.md).
