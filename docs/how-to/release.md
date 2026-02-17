# How to Release a Version

Step-by-step guide for publishing a new version with ceremony.

## The Workflow

```
Feature/fixes completed
    ↓
Gather all changes
    ↓
@lovejoy: Plan the release
    ↓
Update version & CHANGELOG
    ↓
Tag & commit
    ↓
Publish to registry
    ↓
Announce
```

## Step 1: Gather Changes

**Who:** Anyone

**What to do:**
- List all features, fixes, and changes since last release
- Reference commit hashes or PR numbers
- Categorize:
  - Breaking Changes
  - Features (new capabilities)
  - Fixes (bug fixes)
  - Improvements (refactors, perf, etc.)

**Output:**
- Changelog draft (can be informal)

**When you're done:** You have a list of what's going into this release.

---

## Step 2: Determine Version

**Who:** Lovejoy (Release Master)

**What to do:**
- Follow Semantic Versioning: `MAJOR.MINOR.PATCH`
  - `MAJOR` - Breaking changes
  - `MINOR` - New features (backward compatible)
  - `PATCH` - Bug fixes only
- Examples:
  - `1.0.0` → `2.0.0` (breaking)
  - `1.0.0` → `1.1.0` (new feature)
  - `1.0.0` → `1.0.1` (bug fix)

**Output:**
- New version number (e.g., `v1.2.3`)

**When you're done:** Version is decided.

---

## Step 3: Write CHANGELOG

**Who:** Lovejoy (Release Master)

**What to do:**
- Update `CHANGELOG.md` with:
  ```
  ## [1.2.3] - 2026-02-16
  
  ### Added
  - New feature X (PR #123)
  - New feature Y
  
  ### Fixed
  - Bug fix A (Issue #456)
  - Bug fix B
  
  ### Changed
  - Refactored module X
  - Improved performance of Y
  
  ### Breaking
  - Removed deprecated API Z
  ```
- Follow [Keep a Changelog](https://keepachangelog.com/) format
- Be clear and user-focused
- Include PR/issue references

**Output:**
- Updated CHANGELOG.md

**When you're done:** Change history is documented.

---

## Step 4: Update Version Files

**Who:** Lovejoy (Release Master)

**What to do:**
- Update version in all version files:
  - `package.json` (if applicable)
  - `VERSION` file (if you have one)
  - Code version constants
  - Documentation version
- Ensure consistency across all files

**Output:**
- All version files updated to new version

**When you're done:** Version is consistent everywhere.

---

## Step 5: Tag & Commit

**Who:** Lovejoy (Release Master)

**What to do:**
```bash
# Create annotated tag
git tag -a v1.2.3 -m "Release v1.2.3: [brief description]"

# Push tag
git push origin v1.2.3

# Alternative: commit version bump first
git add package.json CHANGELOG.md
git commit -m "chore: bump version to 1.2.3"
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin main
git push origin v1.2.3
```

**Output:**
- git tag created
- Commit recorded

**When you're done:** Version is tagged in git.

---

## Step 6: Publish to Registry

**Who:** Lovejoy (Release Master)

**What to do:**
- For npm packages:
  ```bash
  npm publish
  ```
- For other registries (PyPI, Maven, etc.):
  - Follow registry-specific instructions
  - Verify package published correctly
  - Check that version appears in registry

**Output:**
- Package published
- Version available to users

**When you're done:** Users can download the new version.

---

## Step 7: Announce

**Who:** Lovejoy (Release Master)

**What to do:**
- Announce in appropriate channels:
  - GitHub Releases (copy CHANGELOG)
  - Slack/Discord
  - Email list
  - Twitter/social media (if applicable)
- Highlight:
  - What's new (features)
  - What's fixed (bugs)
  - What's changed (breaking changes)
  - Upgrade instructions (if breaking)

**Output:**
- Release announced
- Users notified

**When you're done:** Release is live and known!

---

## Example: v1.2.3 Release

### Step 1-2: Gather & Version
```
Changes since v1.2.2:
- PR #123: Add user authentication (feature)
- PR #124: Fix password reset bug (fix)
- PR #125: Improve query performance (improvement)
- PR #126: Refactor internal API (change)

Breaking changes? No
New features? Yes (1)
Bug fixes? Yes (1)
Version: 1.2.2 → 1.2.3 (patch, wait—new features mean MINOR)
Actual version: 1.2.2 → 1.3.0
```

### Step 3: CHANGELOG
```
## [1.3.0] - 2026-02-16

### Added
- User authentication system (PR #123)

### Fixed
- Password reset email not sending (PR #124)

### Changed
- Query performance improved by 40% (PR #125)
- Internal API refactored for clarity (PR #126)
```

### Step 4-5: Version & Tag
```bash
# Update package.json, VERSION, etc.
git add package.json CHANGELOG.md
git commit -m "chore: release v1.3.0"
git tag -a v1.3.0 -m "Release v1.3.0: Add authentication, improve perf"
git push origin main v1.3.0
```

### Step 6: Publish
```bash
npm publish
# ✓ Published to npm registry
```

### Step 7: Announce
```
🎉 Springfield Protocol v1.3.0 released!

New features:
✨ User authentication system

Improvements:
⚡ 40% query performance boost
🔧 Cleaner internal APIs

Bug fixes:
🐛 Password reset emails now sending

Upgrade: npm install springfield@1.3.0

Full changelog: https://github.com/.../releases/v1.3.0
```

---

## Troubleshooting

**Version number unclear**
- Refer to Semantic Versioning rules
- When in doubt, ask: "Does this break existing users?"
- Breaking → MAJOR, New feature → MINOR, Fix → PATCH

**Forgot to update CHANGELOG**
- Create a follow-up commit
- Update CHANGELOG with previous release info
- Tag next release with both

**Published wrong version**
- Republish with correct version (npm allows this)
- Mark incorrect version as deprecated on registry
- Announce the correction

---

## Files Involved

- `package.json` (or equivalent version file)
- `CHANGELOG.md`
- git tags
- Registry (npm, PyPI, etc.)

## See Also

- **lovejoy.md** (`.github/agents/`) - Release Master profile
- **Semantic Versioning** - https://semver.org/
- **Keep a Changelog** - https://keepachangelog.com/
