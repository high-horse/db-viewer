# Welcome to Your New Wails3 Project!

Congratulations on generating your Wails3 application! This README will guide you through the next steps to get your project up and running.

## Getting Started

1. Navigate to your project directory in the terminal.

2. To run your application in development mode, use the following command:

   ```
   wails3 dev
   ```

   This will start your application and enable hot-reloading for both frontend and backend changes.

3. To build your application for production, use:

   ```
   wails3 build
   ```

   This will create a production-ready executable in the `build` directory.

## Exploring Wails3 Features

Now that you have your project set up, it's time to explore the features that Wails3 offers:

1. **Check out the examples**: The best way to learn is by example. Visit the `examples` directory in the `v3/examples` directory to see various sample applications.

2. **Run an example**: To run any of the examples, navigate to the example's directory and use:

   ```
   go run .
   ```

   Note: Some examples may be under development during the alpha phase.

3. **Explore the documentation**: Visit the [Wails3 documentation](https://v3.wails.io/) for in-depth guides and API references.

4. **Join the community**: Have questions or want to share your progress? Join the [Wails Discord](https://discord.gg/JDdSxwjhGf) or visit the [Wails discussions on GitHub](https://github.com/wailsapp/wails/discussions).

## Project Structure

Take a moment to familiarize yourself with your project structure:

- `frontend/`: Contains your frontend code (HTML, CSS, JavaScript/TypeScript)
- `main.go`: The entry point of your Go backend
- `app.go`: Define your application structure and methods here
- `wails.json`: Configuration file for your Wails project

## Next Steps

1. Modify the frontend in the `frontend/` directory to create your desired UI.
2. Add backend functionality in `main.go`.
3. Use `wails3 dev` to see your changes in real-time.
4. When ready, build your application with `wails3 build`.

Happy coding with Wails3! If you encounter any issues or have questions, don't hesitate to consult the documentation or reach out to the Wails community.

## Remote databases over SSH

Select PostgreSQL, MySQL, or MongoDB and enable **Connect through SSH tunnel**. Enter the SSH server host, port (usually 22), username, and either a password or a private key. Keys may be pasted as PEM/OpenSSH text or supplied as a local file path, including `~/.ssh/id_ed25519`. Encrypted keys support a passphrase. Use **Edit SSH Configuration** to change these settings.

The database host and port are resolved from the SSH server. For a database running on that server, use `127.0.0.1` and its database port (5432 for PostgreSQL, 3306 for MySQL, or 27017 for MongoDB). The SSH server must allow TCP forwarding.

Server keys are checked against `~/.ssh/known_hosts` using Go's [SSH library](https://pkg.go.dev/golang.org/x/crypto/ssh). Before connecting for the first time, run `ssh -p 22 user@server` and verify its fingerprint before accepting it. Unknown or changed keys are rejected.

**Test Connection** checks SSH authentication and database access. Saved connections retain their SSH settings. Disable **Save connection** for a temporary session. Tunnels close on disconnect, failed database connection, or application shutdown. SQLite remains a local file connection.

## MongoDB

Select **MongoDB**, enter a host and port (default `27017`), and choose a database. Username and password are optional for servers without authentication. Host connections authenticate against `admin`. For other authentication sources, TLS, replica sets, or Atlas, enter a `mongodb://` or `mongodb+srv://` URI in the host field, for example `mongodb://user:password@localhost:27017/?authSource=mydb`. URI credentials and options are used directly; port and separate credential fields are disabled. The Database field selects the database used by commands.

MongoDB also supports the SSH settings above. Use a plain database hostname and port with SSH. Tunnels connect directly to that MongoDB member so replica-set discovery cannot bypass the tunnel; connect to a primary for writes. Saved connections keep the URI or SSH configuration.

Double-click a collection to browse its documents. The query console accepts one MongoDB database command as JSON or Extended JSON (rather than JavaScript shell expressions), following the driver's [database command API](https://www.mongodb.com/docs/drivers/go/v1.x/usage-examples/command/):

```json
{"find":"users","filter":{"active":true},"sort":{"_id":1}}
```

```json
{"aggregate":"users","pipeline":[{"$match":{"active":true}},{"$group":{"_id":"$role","count":{"$sum":1}}}]}
```

```json
{"find":"users","filter":{"_id":{"$oid":"507f1f77bcf86cd799439011"}}}
```

Supported commands include `find`, `aggregate`, `count`, `distinct`, `listCollections`, `listIndexes`, `collStats`, `dbStats`, `ping`, `hello`, `insert`, `update`, `delete`, `findAndModify`, `create`, `createIndexes`, `drop`, and `dropIndexes`. Readonly connections reject writes, including aggregation pipelines using `$out` or `$merge`. Change streams and tailable cursors are not supported.

Documents are displayed as canonical Extended JSON so BSON types and large integers keep their precision. Double-click a document cell to inspect its formatted contents. Cursor results use sequential pages with up to 500 documents per page (100 by default), retain at most 16 open streams, and expire after 10 minutes. Sort using the MongoDB command's `sort` field or an aggregation `$sort` stage. Closing results, disconnecting, and application shutdown release MongoDB cursors.

## Download and run on Linux

Open this repository's **Releases** page and download the Linux x86_64 AppImage and `SHA256SUMS` from the same release. The AppImage is for Intel/AMD 64-bit machines; ARM64 packages are not currently produced.

For example, for version `v0.1.0`, put both files in the same directory and run:

```bash
sha256sum --check SHA256SUMS
chmod +x db-viewer-v0.1.0-linux-x86_64.AppImage
./db-viewer-v0.1.0-linux-x86_64.AppImage
```

If your system cannot mount AppImages through FUSE, run it with extraction instead:

```bash
./db-viewer-v0.1.0-linux-x86_64.AppImage --appimage-extract-and-run
```

The release build uses Ubuntu 22.04 and GTK3/WebKit2GTK 4.1, with runtime libraries bundled by Wails' AppImage generator. It targets modern glibc-based Linux distributions; compatibility with older distributions, Alpine/musl, and NixOS is not guaranteed. App settings remain in the user's configuration directory (`~/.config/db-viewer/app.db` on Linux by default).

## Publish a release

See [release-document.md](release-document.md) for the full release procedure, verification steps, and troubleshooting.

The [release workflow](.github/workflows/release.yml) runs when you push a version tag. Commit and push the workflow and application changes first, then tag the commit you want to release:

```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

The workflow tests the backend, builds the frontend and Linux executable, packages and verifies the AppImage, and publishes it with `SHA256SUMS` and generated release notes. It uses GitHub's automatic token; no personal access token is required. Tags such as `v0.1.0-beta.1` produce prereleases. Ordinary branch pushes do not publish a release.

To retry a failed run, use **Actions → Release Linux AppImage → Re-run jobs**, or manually run the workflow with an existing version tag selected. An existing draft is completed on retry; a published release is left unchanged. Use a new version tag for subsequent releases. The first GitHub-hosted run is needed to verify the full Ubuntu packaging environment.
