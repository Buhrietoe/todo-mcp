# TODO

## Plan to align Go MCP server with Rust example and go-sdk

1. **Review Rust implementation**
   - Use agent tool to summarize functionality (already done).
2. **Define tool metadata**
   3. **Expose tools via MCP SDK**
   - Ensure tools are registered with `mcp.NewServer` (already done in `main.go`).
4. **Implement read/write semantics**
   - `todo_read` returns current content.
   - `todo_write` overwrites content and returns update message.
5. **Add optional fallback storage**
   - Use in‑memory fallback (already present).
6. **Write tests**
   - Test `handleRead` and `handleWrite` directly.
   - Test end‑to‑end client‑server interaction.
7. **Update documentation**
   - Add description in `README.md`.
   - Document tool usage and limits.
8. **Run lint / vet**
   - Ensure no unused imports.
9. **Commit changes**
   - Follow conventional commit message.

After completing these steps, the Go server will mirror the Rust example's behavior and be ready for production use.