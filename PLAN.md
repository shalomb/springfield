I have translated the high-level intent from `PLAN.md` into actionable state in `td` and prepared the environment for Ralph.

### Summary of Changes:
1.  **Plan Update**: Updated `PLAN.md` to include `td` IDs for historical tracking of **EPIC-005 Phase 2** (`td-7b0cb8`) and the active **Phase 3** (`td-e1fd16`).
2.  **State Management**:
    *   Closed **Phase 2** epic (`td-7b0cb8`) as it is shipped and complete.
    *   Activated **Phase 3** epic (`td-e1fd16`) and moved it to `in_progress`.
    *   Expanded **Phase 3** scope by moving the deferred **JSON Stream Integration** task (`td-aa88b1`) into its backlog and creating a new task for **Cost Calculation Fixes** (`td-9177a6`).
3.  **Handoff Preparation**:
    *   Updated `worktrees/epic-td-e1fd16/TODO-td-e1fd16.md` with the refined approach and full task list.
    *   Deposited and committed the updated handoff in the active worktree.

### Next Steps:
Ralph should focus on the **Native JSON Stream Integration** (`td-aa88b1`) as it is a prerequisite for reliable budget enforcement.

<action>
td log td-e1fd16 lisa_planned --decision
</action>
