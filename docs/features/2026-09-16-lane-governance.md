# Feature Plan: Lane Governance & Zero-Trust Memory

## 1. Goals
- Implement a Context Firewall at Layer 1 (Substrate) to prevent "Context Bleeding" between different workspaces (lanes).
- Enable seamless "Model Hot-swapping" (e.g., Claude to Gemini) by ensuring context resides in the OS ledger, not the client UI.
- Build the foundation for the "Abstraction Bridge" (Layer 3), allowing the Lexis Engine to anonymize and promote design patterns from private lanes to a shared `commons` lane.

## 2. Acceptance Criteria (ACs)
- **AC1:** The system must recognize `lane_id` defined in `.cog/config.yaml`.
- **AC2:** The Context Engine (`internal/engine/context.go`) must filter out any `CogDoc` that does not match the current `lane_id` or `commons` before passing documents to the Mamba SSM.
- **AC3:** A global `.cog_global/` directory must be supported and mounted for the `commons` lane.
- **AC4:** CLI command `cos lane switch <lane_id>` must update the active context without losing historical turns from the ledger.
- **AC5:** A manual approval CLI/MCP tool must be provided for the Lexis Engine to quarantine and approve abstracted patterns before promoting them to `commons`.

## 3. Data Model Changes
- Add `Lane` (string) field to `CogDoc` struct in `pkg/cogfield`.
- Add `Lane` (string) field to `EventEnvelope` in `pkg/cogblock`.

## 4. API & Protocol Surface
- Add `cog_list_lanes`, `cog_switch_lane` MCP tools.
- Modify `cog_assemble_context` MCP tool to accept/respect the active lane.

## 5. Implementation Steps
1. **Schema Update:** Update `pkg/cogfield` and `pkg/cogblock` to include the `Lane` attribute.
2. **Directory Mounting:** Update daemon boot logic to mount both local `.cog/` and global `.cog_global/`.
3. **Context Engine Middleware:** Inject a pre-scoring filter in `internal/engine/context.go` to enforce Zero-Trust Memory.
4. **CLI & MCP:** Expose lane management via CLI commands and MCP tools.
5. **Lexis Engine Quarantine (Future Phase):** Build the quarantine table for pattern abstraction.
