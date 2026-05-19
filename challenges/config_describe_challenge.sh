#!/usr/bin/env bash
# config_describe_challenge.sh
#
# Round-253 paired-mutation deep-doc challenge for digital.vasic.config.
#
# Validates that:
#   1. The deep-doc ledger (docs/test-coverage.md) lists every exported
#      symbol from pkg/{config,env,validator,storageconfig}.
#   2. The bilingual fixture (tests/fixtures/i18n/payloads.json) parses
#      and contains at least 3 locales.
#   3. The bilingual runner (challenges/runner/main.go) builds and runs,
#      byte-preserving non-ASCII input through SaveFile/LoadFile,
#      env.Load, storageconfig.MarshalConfig + every protocol Unmarshal,
#      and validator.Validate.
#   4. The README enumerates the four packages and the round-253
#      anti-bluff guarantees section.
#
# Paired-mutation invariant (CONST-035 + CONST-050(B)):
#   With --anti-bluff-mutate the script plants a deliberate symbol-rename
#   mutation in the ledger (in a tmp copy), reruns validation, and asserts
#   the gate FAILS with exit 99. This proves the gate actually catches
#   ledger-vs-source drift instead of rubber-stamping it.
#
# Exit codes:
#   0  — gate PASS on clean tree
#   1  — gate FAIL on clean tree (real failure to fix)
#   99 — paired-mutation correctly detected (good — proves anti-bluff)
#   2  — usage / environment error

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULE_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

MUTATE=0
for arg in "$@"; do
    case "$arg" in
        --anti-bluff-mutate) MUTATE=1 ;;
        --help|-h)
            sed -n '1,30p' "$0"
            exit 0
            ;;
        *)
            echo "unknown argument: $arg" >&2
            exit 2
            ;;
    esac
done

PASS=0
FAIL=0
TOTAL=0

pass() { PASS=$((PASS+1)); TOTAL=$((TOTAL+1)); echo "  PASS: $1"; }
fail() { FAIL=$((FAIL+1)); TOTAL=$((TOTAL+1)); echo "  FAIL: $1"; }

LEDGER="${MODULE_DIR}/docs/test-coverage.md"
FIXTURE="${MODULE_DIR}/tests/fixtures/i18n/payloads.json"
RUNNER="${MODULE_DIR}/challenges/runner/main.go"
README="${MODULE_DIR}/README.md"

LEDGER_WORK="${LEDGER}"
TMP_LEDGER=""
if [ "${MUTATE}" -eq 1 ]; then
    TMP_LEDGER="$(mktemp)"
    cp "${LEDGER}" "${TMP_LEDGER}"
    # Plant a rename: LoadFile -> LoadBogus_MUTATED
    sed -i 's/LoadFile/LoadBogus_MUTATED/g' "${TMP_LEDGER}"
    LEDGER_WORK="${TMP_LEDGER}"
    echo "=== Config Describe Challenge (anti-bluff-mutate mode) ==="
else
    echo "=== Config Describe Challenge (clean mode) ==="
fi
echo ""

# Section 1: ledger presence and freshness
echo "Section 1: docs/test-coverage.md ledger"
if [ ! -f "${LEDGER_WORK}" ]; then
    fail "ledger missing at ${LEDGER_WORK}"
else
    pass "ledger present"
    if grep -q "round-253" "${LEDGER_WORK}"; then
        pass "ledger marked round-253"
    else
        fail "ledger missing round-253 marker"
    fi
    if grep -q "execution of tests and Challenges MUST guarantee" "${LEDGER_WORK}"; then
        pass "ledger carries Article XI §11.9 mandate"
    else
        fail "ledger missing Article XI §11.9 mandate"
    fi
fi

# Section 2: every exported pkg symbol appears in ledger
echo ""
echo "Section 2: exported symbols cross-reference"

extract_symbols() {
    local pkg_dir="$1"
    local files
    files=$(find "${pkg_dir}" -maxdepth 1 -type f -name '*.go' \
        ! -name '*_test.go')
    [ -z "${files}" ] && return 0
    # shellcheck disable=SC2086
    grep -hE '^(func ([A-Z][A-Za-z0-9_]*\()|func \([^)]+\) ([A-Z][A-Za-z0-9_]*\()|type [A-Z][A-Za-z0-9_]* )' \
        ${files} 2>/dev/null \
        | sed -E 's/^func \([^)]+\) ([A-Z][A-Za-z0-9_]*)\(.*$/\1/; s/^func ([A-Z][A-Za-z0-9_]*)\(.*$/\1/; s/^type ([A-Z][A-Za-z0-9_]*).*$/\1/' \
        | sort -u
}

CHECKED=0
MISSING=0
for pkg in config env validator storageconfig; do
    PKG_DIR="${MODULE_DIR}/pkg/${pkg}"
    if [ ! -d "${PKG_DIR}" ]; then
        fail "pkg/${pkg} missing — cannot cross-reference"
        continue
    fi
    while IFS= read -r sym; do
        [ -z "${sym}" ] && continue
        CHECKED=$((CHECKED + 1))
        if grep -qE "\\b${sym}\\b" "${LEDGER_WORK}"; then
            : # symbol cross-referenced
        else
            fail "ledger missing symbol ${pkg}.${sym}"
            MISSING=$((MISSING + 1))
        fi
    done < <(extract_symbols "${PKG_DIR}")
done
if [ "${CHECKED}" -gt 0 ] && [ "${MISSING}" -eq 0 ]; then
    pass "all ${CHECKED} exported symbols cross-referenced in ledger"
fi

# Section 3: bilingual fixture sanity
echo ""
echo "Section 3: bilingual fixture"
if [ ! -f "${FIXTURE}" ]; then
    fail "fixture missing at ${FIXTURE}"
else
    pass "fixture present"
    LOCALE_COUNT=$(grep -oE '"locale":\s*"[^"]+"' "${FIXTURE}" | sort -u | wc -l)
    if [ "${LOCALE_COUNT}" -ge 3 ]; then
        pass "fixture covers ${LOCALE_COUNT} locales (>=3)"
    else
        fail "fixture covers only ${LOCALE_COUNT} locales (<3)"
    fi
fi

# Section 4: runner builds + runs against every package
echo ""
echo "Section 4: bilingual runner build + run (real os + env transport)"
if [ ! -f "${RUNNER}" ]; then
    fail "runner missing at ${RUNNER}"
else
    pass "runner source present"
    cd "${MODULE_DIR}"
    if go build -o /tmp/config_round253_runner ./challenges/runner/ 2>/tmp/cfg_build.log; then
        pass "runner builds"
        if /tmp/config_round253_runner -fixtures "${FIXTURE}" > /tmp/cfg_run.log 2>&1; then
            pass "runner exit 0 across every package + locale"
            if grep -q "PASS: \[config\]\[sr\]" /tmp/cfg_run.log; then
                pass "config pkg Cyrillic (sr) round-trip"
            else
                fail "config pkg Cyrillic (sr) missing from runner output"
            fi
            if grep -q "PASS: \[env\]\[ja\]" /tmp/cfg_run.log; then
                pass "env pkg Japanese (ja) os.Setenv round-trip"
            else
                fail "env pkg Japanese (ja) missing from runner output"
            fi
            if grep -q "PASS: \[storage\]\[WebDAV\]" /tmp/cfg_run.log; then
                pass "storageconfig WebDAV non-ASCII Marshal/Unmarshal"
            else
                fail "storageconfig WebDAV section missing"
            fi
            if grep -q "PASS: \[storage\]\[Git\]" /tmp/cfg_run.log; then
                pass "storageconfig Git non-ASCII Marshal/Unmarshal"
            else
                fail "storageconfig Git section missing"
            fi
            if grep -q "PASS: \[validator\]\[ar\]" /tmp/cfg_run.log; then
                pass "validator Arabic (ar) accept"
            else
                fail "validator Arabic (ar) missing from runner output"
            fi
            if grep -q "PASS: \[validator-failmode\]" /tmp/cfg_run.log; then
                pass "validator failure-mode (3 violations surfaced)"
            else
                fail "validator failure-mode missing from runner output"
            fi
        else
            fail "runner exit non-zero — see /tmp/cfg_run.log"
            sed -n '1,40p' /tmp/cfg_run.log
        fi
    else
        fail "runner build failed — see /tmp/cfg_build.log"
        sed -n '1,40p' /tmp/cfg_build.log
    fi
    rm -f /tmp/config_round253_runner
fi

# Section 5: README round-253 anti-bluff section
echo ""
echo "Section 5: README round-253 anti-bluff section"
if grep -q "Anti-bluff guarantees" "${README}"; then
    pass "README declares Anti-bluff guarantees"
else
    fail "README missing Anti-bluff guarantees section"
fi
if grep -q "round-253" "${README}"; then
    pass "README marked round-253"
else
    fail "README missing round-253 marker"
fi

# Cleanup mutated ledger if any
if [ -n "${TMP_LEDGER}" ]; then
    rm -f "${TMP_LEDGER}"
fi

echo ""
echo "=== Summary: ${PASS}/${TOTAL} PASS, ${FAIL} FAIL ==="

if [ "${MUTATE}" -eq 1 ]; then
    if [ "${FAIL}" -gt 0 ]; then
        echo "anti-bluff-mutate: gate correctly detected planted mutation (exit 99)"
        exit 99
    else
        echo "anti-bluff-mutate: gate FAILED to detect planted mutation — bluff!"
        exit 1
    fi
fi

if [ "${FAIL}" -gt 0 ]; then
    exit 1
fi
exit 0
