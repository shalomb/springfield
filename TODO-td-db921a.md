# TODO-td-db921a.md - Autonomous Control Plane

## Intent
Transform Springfield into a daemonized Control Plane that manages agent lifecycles, environments, and state transitions. Replace loose string signaling with a strict, sentinel-protected CLI callback protocol.

## Approach
1.  **Signaling:** Implement `springfield signal` as the exclusive exit mechanism. Validate Sentinel tokens to prevent unauthorized state changes.
2.  **Sentinel:** Generate cryptographic tokens per session. Inject them into prompts.
3.  **Daemon:** Create a persistent `orchestrate --daemon` mode that polls `td` and manages worker processes.
4.  **Templating:** Convert static prompts to Go templates to support dynamic injection.

## Constraints
- **Atomic Commits:** Follow ACP. One task = one commit.
- **Backward Compatibility:** Legacy agents must still work until migration is complete.
- **Security:** Sentinels must be unpredictable.

## Tasks (in `td`)
1.  `td-b6dae6`: Implement `springfield signal` command.
2.  `td-d59930`: Implement Sentinel Token logic.
3.  `td-e53407`: Migrate prompts to templates.
4.  `td-278c10`: Implement Daemon loop.

## References
- `docs/features/EPIC-010-control-plane.md`
- `docs/standards/signaling-protocol.md`
- `tests/integration/features/control_plane.feature`
