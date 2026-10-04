# Releasing DB Viewer

This guide describes the current [GitHub release workflow](.github/workflows/release.yml).

## What gets published

Pushing a version tag runs **Release Linux AppImage** and publishes these files to the repository's GitHub Releases page:

- `db-viewer-v0.1.0-linux-x86_64.AppImage` (with the actual tag in the filename).
- `SHA256SUMS`, containing the AppImage's SHA-256 checksum.
- Automatically generated release notes.

The workflow currently produces Linux Intel/AMD 64-bit packages only. It does not build ARM64, Windows, or macOS packages. Ordinary branch pushes do not publish releases.

## Prepare a release

1. Choose an unused version, such as `v0.1.0` or `v0.1.1`. Use `v0.2.0-beta.1` for a prerelease.
2. Commit the application changes, workflow, documentation, and dependency lock files you want included. A release builds the tagged commit, so uncommitted changes will not be included.
3. Review the changes and run the local checks:

   ```bash
   go test -tags gtk3 ./internal/...
   npm --prefix frontend ci
   npm --prefix frontend run build
   ```

   The Go checks require the Linux GTK3/WebKit2GTK 4.1 development libraries. The workflow installs these dependencies on its runner.

4. Push your release commit to GitHub using your normal branch or pull request process. Confirm that the intended commit contains `.github/workflows/release.yml`.

## Tag and publish

From the commit you want to release, run:

```bash
git status --short
git log -1 --oneline
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

Replace `v0.1.0` with the new version in both commands. Tags must use `vMAJOR.MINOR.PATCH`, optionally followed by a prerelease suffix. Tags containing a hyphen are published as prereleases.

Open **GitHub → Actions → Release Linux AppImage** to monitor the run. The workflow:

1. Checks out and validates the tag.
2. Sets up Go from `go.mod`, Node.js 22, and Linux build dependencies.
3. Installs the Wails CLI version specified by the project's dependency.
4. Runs backend tests, installs frontend dependencies from `package-lock.json`, and regenerates bindings.
5. Builds the frontend and production Linux executable.
6. Packages the AppImage, checks that it extracts, and generates checksums.
7. Creates a draft release, uploads the files, and publishes the release.

The workflow uses GitHub's automatic token with `contents: write` permission. No personal access token or additional secret is needed. Repository or organization policies must allow the workflow to run and write releases.

## Verify the published release

Open the repository's **Releases** page and confirm that the version, prerelease status, AppImage, and `SHA256SUMS` are present. Download both files into the same directory and run:

```bash
sha256sum --check SHA256SUMS
chmod +x db-viewer-v0.1.0-linux-x86_64.AppImage
./db-viewer-v0.1.0-linux-x86_64.AppImage
```

If FUSE mounting is unavailable, try:

```bash
./db-viewer-v0.1.0-linux-x86_64.AppImage --appimage-extract-and-run
```

Check that the app starts, can connect to a test database, and can display query results. The workflow checks package extraction, but does not launch the graphical application or verify every target distribution.

## Retry or fix a release

- **Temporary build or upload failure:** open the failed Actions run and select **Re-run jobs**.
- **Manual retry:** select **Run workflow** and choose the existing version tag. Selecting a branch skips the release job.
- **Draft left behind:** rerunning the same tag uploads the assets again and completes publication.
- **Release already published:** reruns leave it unchanged. Commit fixes and publish a new version tag.
- **Application or workflow bug:** commit the fix, push it, and use a new version tag. Rerunning an old tag builds its original source and workflow.
- **Permission failure:** check the Actions log and repository or organization policies for restrictions on the automatic token's release write access.

## Build environment and compatibility

Release builds use Ubuntu 22.04, CGO, and the `production,gtk3` build tags. The AppImage generator bundles runtime dependencies for GTK3/WebKit2GTK 4.1. The Ubuntu baseline targets modern glibc-based distributions; older distributions, Alpine/musl, and NixOS need separate compatibility testing.

The AppImage contains the application, not a user's saved connections or databases. Internal app data remains in the user's configuration directory, normally `~/.config/db-viewer/app.db` on Linux, or `$XDG_CONFIG_HOME/db-viewer/app.db` when that environment variable is set.

When changing the Go, Node.js, Wails, GTK, or Ubuntu versions, review the workflow and test a prerelease before publishing a stable release. Full Ubuntu packaging must be verified by a successful GitHub-hosted run.
