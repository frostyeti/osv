---
type: "docs"
tags: ["docs", "enhancement"]
---

# Proposed New Features

## Implementation Plan

After reviewing the current `osv` (Operating System Vault CLI) codebase, here is a proposal for new features to enhance its functionality and user experience.

1. **Copy to Clipboard Support**
   - **Feature**: Add a `--clip` or `-c` flag to the `osv get` command.
   - **Benefit**: Instead of printing sensitive secrets to stdout where they might be saved in shell history or visible on screen, this feature copies the secret directly to the system clipboard.
   - **Implementation**: Integrate a cross-platform clipboard library (e.g., `github.com/atotto/clipboard`) and invoke it when the flag is specified.

2. **Rename Secret Command**
   - **Feature**: Add an `osv rename <old-key> <new-key>` command.
   - **Benefit**: Simplifies the process of renaming a secret. Currently, users have to manually get the old secret, set the new one, and delete the old one.
   - **Implementation**: Implement a wrapper in `cmd/rename.go` that safely retrieves the old secret, writes it to the new key, and deletes the old key only upon successful creation of the new one.

3. **Interactive TUI Mode**
   - **Feature**: Introduce an `osv ui` command to launch an interactive mode.
   - **Benefit**: Provides a visual way to browse, filter, copy, and manage secrets directly from the terminal.
   - **Implementation**: Use a rich terminal UI library like `github.com/charmbracelet/bubbletea` and `bubbles` to render lists, input fields, and confirmation dialogs.

4. **Bulk Import/Export**
   - **Feature**: Provide `osv export` and `osv import` commands.
   - **Benefit**: Enables easy migration of secrets between machines or environments, as well as backup capabilities.
   - **Implementation**: Export secrets to an encrypted archive format (using `age` or AES-GCM encryption with a user-provided passphrase) and parse them safely during import.
   
5. **Secret Expiration/Tags (Metadata)**
   - **Feature**: Allow attaching arbitrary tags or an expiration date to a secret during `osv set`.
   - **Benefit**: Helps users organize large keyrings and receive warnings for stale secrets (e.g., API keys that should be rotated).
   - **Implementation**: Since standard OS keyrings generally only store strings, we can implement a lightweight local index (e.g., `~/.osv/metadata.json`) mapped to the secret keys to store extra attributes.
