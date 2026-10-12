# Screenshots

How to update `docs/os-screenshots.webp`, the animated hero image in `README.md`.

## Founders preferences

These are requirements. Follow them without asking.

- Format: animated WebP, 1280x720 (16:9 widescreen). GitHub renders it in the README content column.
- File size: aim for 400-500 KB. A smaller file loses quality. A larger file takes longer to load.
- Duration: 13-15 seconds, looping. Hard cuts between scenes are fine.
- Order: WebTerminal running OpenCode first, then the Marketplace install modal, then the File Manager list at `/app/html`, then the File Manager editor, then the databases scenes in this order: MariaDB list, PostgreSQL tab, Create Database User modal. The terminal opens the animation.
- Skip boring pages. Do not show Cron Jobs, Mappings, Runtime Settings, Accounts, or SSL pages.
- Pacing: cut dead time. The OpenCode segment runs at 2x speed. Never open on a long empty terminal. The first visible action must start within about half a second.
- No cursor overlay in the recording.
- Voice: the animation is a product tour for junior devs and vibe coders. Keep scenes clean and readable.

## 1. Start a disposable instance

Requires Go, mise, and podman. See `docs/DEVELOPMENT.md` for the full setup.

From the repository root, run the dev build detached. It builds the image and rebuilds on every source change.

```
setsid bash dev-build.sh http-unpriv --pid-only >/tmp/opencode/dev-build.log 2>&1 < /dev/null &
```

Wait for port 1618, then log in at `https://127.0.0.1:1618/` with `dev` / `abc123!`.

Stop it later with `kill -TERM "$(cat logs/dev-build.pid)"`. Use SIGTERM, never SIGINT (see `AGENTS.override.md`).

## 2. Install OpenCode in the container

Terminal sessions run as `nobody` with `HOME=/app` (`src/infra/terminalSession/terminalMultiplexerClient.go`). Put the binary in a shared path. Give `nobody` write access to the opencode directories under `/app`.

```
podman exec os bash -lc 'curl -fsSL https://opencode.ai/install | bash'
podman exec os bash -c 'install -m 0755 /root/.opencode/bin/opencode /usr/local/bin/opencode'
podman exec os bash -c 'mkdir -p /app/.local/share /app/.config /app/.cache && chown -R nobody:nogroup /app/.local /app/.config /app/.cache'
```

No credentials are needed. The free Big Pickle model works out of the box. Verify:

```
podman exec -u nobody -e HOME=/app os bash -c 'cd /app && timeout 90 opencode run "Reply with exactly: ready"'
```

## 3. Seed demo data

The Marketplace scenes need WordPress installed, and the Databases scenes need content in both engines.

```
podman exec os bash -c 'os db create -t mariadb -n wordpress_db'
podman exec os bash -c 'os db create-user -t mariadb -n wordpress_db -u wp_user -p "abc123!" -r "select,insert,update,delete"'
podman exec os bash -c 'os services create-installable -n postgresql -v 16 -p 5432/tcp -a false'
podman exec os bash -c 'os db create -t postgresql -n app_db'
podman exec os bash -c 'os db create-user -t postgresql -n app_db -u app_user -p "abc123!" -r "select,insert,update,delete"'
```

Install WordPress from the Marketplace if it is not already installed:

```
podman exec os bash -c 'os mktplace install -s wp -f "adminUsername:admin" -f "adminPassword:abc123!" -f "adminMailAddress:user@example.com"'
```

## 4. Capture the stills

Use agent-browser with a worktree-scoped session and a 1680x944 viewport. Every command needs `--ignore-https-errors` because the dashboard serves a self-signed certificate.

```
agent-browser session id --scope worktree --prefix "ui-verify-screenshots"
agent-browser --session <id> --ignore-https-errors set viewport 1680 944
agent-browser --session <id> --ignore-https-errors open https://127.0.0.1:1618/overview/
mkdir -p /tmp/opencode/anim/raw
```

The scenes, in order. Each screenshot command saves the raw capture under `/tmp/opencode/anim/raw/`:

1. Marketplace install modal. On the Overview page, click `install catalog item`, click WordPress in the carousel, wait one second. The Install Parameters accordion opens by default.
   `agent-browser --session <id> --ignore-https-errors screenshot /tmp/opencode/anim/raw/1-marketplace.png`
2. File Manager list. Open `/file-manager/`, type `/app/html` into the path input, press Enter.
   `agent-browser --session <id> --ignore-https-errors screenshot /tmp/opencode/anim/raw/2-file-list.png`
3. File Manager editor. Click the `index.php` row, click `Edit`, wait for the editor.
   `agent-browser --session <id> --ignore-https-errors screenshot /tmp/opencode/anim/raw/3-file-editor.png`
4. Databases MariaDB tab. Open `/databases/`.
   `agent-browser --session <id> --ignore-https-errors screenshot /tmp/opencode/anim/raw/4-databases-mariadb.png`
5. Databases PostgreSQL tab. Click the `POSTGRESQL` tab.
   `agent-browser --session <id> --ignore-https-errors screenshot /tmp/opencode/anim/raw/5-databases-postgresql.png`
6. Create Database User modal. Click `create user`, select `wordpress_db`, fill a username and password, click the `Privileges` accordion, toggle SELECT, INSERT and DELETE.
   `agent-browser --session <id> --ignore-https-errors screenshot /tmp/opencode/anim/raw/6-databases-create-user.png`

Tab clicks need an eval because the tabs are `li` elements: find the `li` whose text matches and call `click()`.

## 5. Record the terminal segment

Create a session from the Overview page (`new session`). The Terminal Sessions modal opens and attaches. Click `.xterm-screen` to focus, then type with real keystrokes.

Prepare a clean take. Exit any running TUI with `Control+c` twice, then clear:

```
agent-browser --session <id> --ignore-https-errors keyboard type "rm -f hello.php && clear"
agent-browser --session <id> --ignore-https-errors press Enter
```

Record:

```
agent-browser --session <id> --ignore-https-errors record start /tmp/opencode/terminal-take.webm
agent-browser --session <id> --ignore-https-errors keyboard type "opencode"
agent-browser --session <id> --ignore-https-errors press Enter
# wait about 5 seconds for the TUI
agent-browser --session <id> --ignore-https-errors keyboard type "Create a hello.php page that prints Hello from Infinite OS"
agent-browser --session <id> --ignore-https-errors press Enter
# wait for the response to finish, then
agent-browser --session <id> --ignore-https-errors record stop
```

The recording will contain a long dead lead. Agent inference delay puts the `opencode` command around the 6 second mark. Trim it in assembly.

## 6. Assemble

Check a frame sheet to find where the typing starts. Set `-ss` to cut everything before it, then extract the motion frames at 2x speed. The last take used `-ss 5.5`.

```
mkdir -p /tmp/opencode/anim/motion
ffmpeg -hide_banner -loglevel error -ss 5.5 -t 10 -i /tmp/opencode/terminal-take.webm \
  -vf "setpts=PTS/2,fps=10,scale=1280:720:flags=lanczos" /tmp/opencode/anim/motion/c_%04d.png
```

Scale the stills to the same canvas. The names carry the scene, so the assembly order stays easy to check:

```
cd /tmp/opencode/anim
magick raw/1-marketplace.png -resize 1280x720! still1-marketplace.png
magick raw/2-file-list.png -resize 1280x720! still2-file-list.png
magick raw/3-file-editor.png -resize 1280x720! still3-file-editor.png
magick raw/4-databases-mariadb.png -resize 1280x720! still4-databases-mariadb.png
magick raw/5-databases-postgresql.png -resize 1280x720! still5-databases-postgresql.png
magick raw/6-databases-create-user.png -resize 1280x720! still6-databases-create-user.png
```

Join everything with per-frame delays. Stills hold 1.4-1.6 seconds; motion frames run at 10 fps (delay 10). Quality 78 lands in the size target; raise or lower it and re-measure.

```
cd /tmp/opencode/anim
magick -delay 10 motion/c_*.png \
  -delay 160 still1-marketplace.png -delay 160 still2-file-list.png -delay 160 still3-file-editor.png \
  -delay 140 still4-databases-mariadb.png -delay 140 still5-databases-postgresql.png -delay 160 still6-databases-create-user.png \
  -loop 0 -define webp:method=6 -quality 78 /tmp/opencode/os-screenshots.webp
ls -la /tmp/opencode/os-screenshots.webp
```

## 7. Verify and replace

Coalesce the output and inspect a montage. Check the first frames are not empty, the typing is visible, and the scene order matches the preferences.

```
mkdir -p /tmp/opencode/check
magick /tmp/opencode/os-screenshots.webp -coalesce /tmp/opencode/check/f_%03d.png
magick montage /tmp/opencode/check/f_000.png /tmp/opencode/check/f_004.png /tmp/opencode/check/f_012.png \
  -tile 3x1 -geometry 420x236+2+2 -background black /tmp/opencode/check/review.png
```

Open the file in agent-browser to prove it decodes and plays. Then replace the asset:

```
cp /tmp/opencode/os-screenshots.webp docs/os-screenshots.webp
```

Keep the README alt text in sync with the scene order. It currently reads: "Infinite OS demo: running OpenCode in the web terminal, installing WordPress, and managing files and databases".
