import { execFileSync } from "node:child_process";
import { copyFileSync, mkdirSync, mkdtempSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

const runtimeRoot = process.env.GOPURS_ROOT ?? fileURLToPath(new URL("../../gopurs/", import.meta.url));
const workspace = mkdtempSync(join(tmpdir(), "gopurs-listener-added-"));
try {
  mkdirSync(join(workspace, "gopurs_runtime"));
  copyFileSync(join(runtimeRoot, "runtime/runtime.go"), join(workspace, "gopurs_runtime/runtime.go"));
  copyFileSync(fileURLToPath(new URL("../src/Node/EventEmitter.go", import.meta.url)), join(workspace, "emitter.go"));
  copyFileSync(fileURLToPath(new URL("./listener-added-native_test.go", import.meta.url)), join(workspace, "emitter_test.go"));
  writeFileSync(join(workspace, "go.mod"), "module gopurs/output\n\ngo 1.22\n");
  const race = process.env.GOPURS_NATIVE_TEST_RACE === "1" ? ["-race"] : [];
  execFileSync("go", ["test", ...race, "-count=1", "-timeout=10s", "-v", "."], {
    cwd: workspace,
    stdio: "inherit",
    timeout: 120_000,
    env: {
      ...process.env,
      GOCACHE: process.env.GOCACHE ?? join(workspace, "go-cache"),
      GOWORK: "off",
      GOTOOLCHAIN: "local",
      GOPROXY: "off",
      GOSUMDB: "off",
    },
  });
} finally {
  rmSync(workspace, { recursive: true, force: true });
}
