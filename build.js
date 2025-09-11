import { build } from "esbuild";
import {
  rmSync,
  existsSync,
  cpSync,
  lstatSync,
  readlinkSync,
  readdirSync,
} from "fs";
import { resolve, dirname } from "path";

const builder = async () => {
  let distPath = "./dist";

  // If dist is a symlink, get the target directory
  if (existsSync(distPath) && lstatSync(distPath).isSymbolicLink()) {
    const linkTarget = readlinkSync(distPath);
    // Resolve relative symlinks
    distPath = resolve(dirname(distPath), linkTarget);
  }

  // Clean the target directory contents, not the symlink itself
  if (existsSync(distPath)) {
    // Remove contents of directory, not the directory itself
    const files = readdirSync(distPath);
    files.forEach((file) => {
      rmSync(resolve(distPath, file), { recursive: true, force: true });
    });
  }

  // Copy public files to the target directory
  if (existsSync("./public")) {
    cpSync("./public", distPath, { recursive: true });
  }

  await build({
    entryPoints: ["./src/main.js"],
    bundle: true,
    minify: false,
    sourcemap: false,
    target: ["chrome58", "firefox57", "safari11", "edge16"],
    outdir: "./dist",
    define: {
      "process.env.NODE_ENV": JSON.stringify("development"),
      __DEV__: "true",
      __PROD__: "false",
    },
  });
};
builder();
