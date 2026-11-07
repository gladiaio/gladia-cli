#!/bin/sh
# Regression tests for install.sh: serves fake releases from a local HTTP
# server and checks that the installer verifies checksums, rejects unsafe
# archives and tags, and never runs a 'gladia' found earlier on PATH.
# Requires: sh, curl, tar, python3.
set -eu

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
INSTALLER="${ROOT}/install.sh"

work="$(mktemp -d)"
server_pid=""
cleanup() {
  if [ -n "$server_pid" ]; then kill "$server_pid" 2>/dev/null || true; fi
  rm -rf "$work"
}
trap cleanup EXIT INT HUP TERM

case "$(uname -s)" in
  Darwin) os=darwin ;;
  Linux) os=linux ;;
  *) echo "skip: unsupported OS"; exit 0 ;;
esac
case "$(uname -m)" in
  x86_64 | amd64) arch=amd64 ;;
  aarch64 | arm64) arch=arm64 ;;
  armv7l | armv6l) arch=armv7 ;;
  i386 | i686) arch=386 ;;
  *) echo "skip: unsupported arch"; exit 0 ;;
esac

sha256() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{ print $1 }'
  else
    shasum -a 256 "$1" | awk '{ print $1 }'
  fi
}

srv="${work}/srv"
mkdir -p "$srv"
port_file="${work}/port"
python3 -I -c '
import functools, http.server, socketserver, sys
class Quiet(http.server.SimpleHTTPRequestHandler):
    def log_message(self, *args):
        pass
handler = functools.partial(Quiet, directory=sys.argv[1])
with socketserver.TCPServer(("127.0.0.1", 0), handler) as s:
    open(sys.argv[2], "w").write(str(s.server_address[1]))
    s.serve_forever()
' "$srv" "$port_file" &
server_pid=$!
i=0
while [ ! -s "$port_file" ]; do
  i=$((i + 1))
  if [ "$i" -gt 100 ]; then echo "server did not start" >&2; exit 1; fi
  sleep 0.1
done
base="http://127.0.0.1:$(cat "$port_file")"

# publish TAG [MODE]: write latest.json and a release for TAG.
# MODE: ok | badsum | nosum | traversal | symlink
publish() {
  tag="$1"
  mode="${2:-ok}"
  version="${tag#v}"
  archive="gladia_${version}_${os}_${arch}.tar.gz"
  rel="${srv}/${tag}"
  stage="${work}/stage"
  rm -rf "$rel" "$stage"
  mkdir -p "$rel" "$stage"
  printf '{"tag_name": "%s"}\n' "$tag" > "${srv}/latest.json"

  printf '#!/bin/sh\necho "gladia %s"\n' "$version" > "${stage}/gladia"
  chmod 755 "${stage}/gladia"

  case "$mode" in
    traversal | symlink)
      python3 -I - "$stage/gladia" "${rel}/${archive}" "$mode" <<'PY'
import io, sys, tarfile
src, dst, mode = sys.argv[1:4]
with tarfile.open(dst, "w:gz") as t:
    t.add(src, arcname="gladia")
    if mode == "traversal":
        data = b"pwned\n"
        ti = tarfile.TarInfo("../../pwned")
        ti.size = len(data)
        t.addfile(ti, io.BytesIO(data))
    else:
        ti = tarfile.TarInfo("evil-link")
        ti.type = tarfile.SYMTYPE
        ti.linkname = "/etc/passwd"
        t.addfile(ti)
PY
      ;;
    *)
      tar -czf "${rel}/${archive}" -C "$stage" gladia
      ;;
  esac

  sum="$(sha256 "${rel}/${archive}")"
  case "$mode" in
    nosum) printf '%s  %s\n' "$sum" "some-other-file.tar.gz" > "${rel}/checksums.txt" ;;
    *) printf '%s  %s\n' "$sum" "$archive" > "${rel}/checksums.txt" ;;
  esac
  if [ "$mode" = badsum ]; then
    printf 'tampered' >> "${rel}/${archive}"
  fi
}

failures=0
pass() { echo "ok   - $1"; }
fail() { echo "FAIL - $1"; failures=$((failures + 1)); }

# run_installer: runs install.sh into a fresh dir, output in $out, status in $rc.
run_installer() {
  dest="${work}/dest"
  rm -rf "$dest"
  mkdir -p "$dest"
  set +e
  out="$(cd "$work" && env \
    GLADIA_INSTALL_API_URL="${base}/latest.json" \
    GLADIA_INSTALL_DOWNLOAD_BASE="${base}/" \
    GLADIA_INSTALL_DIR="$dest" \
    GLADIA_NO_COMPLETION_PROMPT=1 \
    PATH="${PATH_PREFIX:-}${PATH}" \
    sh "$INSTALLER" 2>&1 </dev/null)"
  rc=$?
  set -e
}

expect_failure() {
  name="$1"
  pattern="$2"
  run_installer
  if [ "$rc" -ne 0 ] && printf '%s' "$out" | grep -q "$pattern" && [ ! -e "${dest}/gladia" ]; then
    pass "$name"
  else
    fail "$name (rc=$rc)"
    printf '%s\n' "$out" | sed 's/^/    /'
  fi
}

# 1. Happy path.
publish v9.9.9 ok
run_installer
if [ "$rc" -eq 0 ] && [ -x "${dest}/gladia" ] && printf '%s' "$out" | grep -q "Checksum verified"; then
  pass "installs a release with a valid checksum"
else
  fail "installs a release with a valid checksum (rc=$rc)"
  printf '%s\n' "$out" | sed 's/^/    /'
fi

# 2. Tampered archive.
publish v9.9.9 badsum
expect_failure "rejects an archive whose checksum does not match" "checksum mismatch"

# 3. Archive missing from checksums.txt.
publish v9.9.9 nosum
expect_failure "rejects an archive missing from checksums.txt" "no sha256 entry"

# 4. Path traversal entry.
publish v9.9.9 traversal
expect_failure "rejects an archive with a '..' entry" "unsafe entries"
if [ -e "${work}/pwned" ] || [ -e "${srv}/pwned" ]; then
  fail "traversal entry was written outside the temp dir"
fi

# 5. Symlink entry.
publish v9.9.9 symlink
expect_failure "rejects an archive with a symlink entry" "unsafe entries"

# 6. Malformed tags are rejected before use.
marker="${work}/tag-marker"
printf '{"tag_name": "v1.0;touch %s"}\n' "$marker" > "${srv}/latest.json"
expect_failure "rejects a tag with shell metacharacters" "unexpected tag format"
if [ -e "$marker" ]; then fail "tag content was executed"; fi
printf '{"tag_name": ".."}\n' > "${srv}/latest.json"
expect_failure "rejects a '..' tag" "unexpected tag format"

# 7. A 'gladia' earlier on PATH is not executed after install.
publish v9.9.9 ok
shim="${work}/shim"
mkdir -p "$shim"
printf '#!/bin/sh\necho SHIM-RAN\n' > "${shim}/gladia"
chmod 755 "${shim}/gladia"
PATH_PREFIX="${shim}:"
run_installer
PATH_PREFIX=""
if [ "$rc" -eq 0 ] && ! printf '%s' "$out" | grep -q SHIM-RAN && printf '%s' "$out" | grep -q "gladia 9.9.9"; then
  pass "uses the installed binary, not one found on PATH"
else
  fail "uses the installed binary, not one found on PATH (rc=$rc)"
  printf '%s\n' "$out" | sed 's/^/    /'
fi

if [ "$failures" -ne 0 ]; then
  echo "${failures} installer test(s) failed"
  exit 1
fi
echo "all installer tests passed"
