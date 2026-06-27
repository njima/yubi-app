import { readdir, readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";

const currentDir = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(currentDir, "..");
const srcRoot = path.join(root, "src");
const sourceExtensions = new Set([".ts", ".tsx"]);

const rules = [
  {
    name: "app imports are route-local only",
    appliesTo: (file) => !relative(file).startsWith("app/"),
    forbidden: (specifier) =>
      specifier === "@/app" || specifier.startsWith("@/app/"),
  },
  {
    name: "shared must not depend on app, features, or components",
    appliesTo: (file) => relative(file).startsWith("shared/"),
    forbidden: (specifier) =>
      specifier === "@/app" ||
      specifier.startsWith("@/app/") ||
      specifier === "@/features" ||
      specifier.startsWith("@/features/") ||
      specifier === "@/components" ||
      specifier.startsWith("@/components/"),
  },
  {
    name: "ui primitives must not depend on app, features, or lib/api",
    appliesTo: (file) => relative(file).startsWith("components/ui/"),
    forbidden: (specifier) =>
      specifier === "@/app" ||
      specifier.startsWith("@/app/") ||
      specifier === "@/features" ||
      specifier.startsWith("@/features/") ||
      specifier === "@/lib/api" ||
      specifier.startsWith("@/lib/api/"),
  },
  {
    name: "lib/api must not depend on app, features, components, or shared",
    appliesTo: (file) => relative(file).startsWith("lib/api/"),
    forbidden: (specifier) =>
      specifier === "@/app" ||
      specifier.startsWith("@/app/") ||
      specifier === "@/features" ||
      specifier.startsWith("@/features/") ||
      specifier === "@/components" ||
      specifier.startsWith("@/components/") ||
      specifier === "@/shared" ||
      specifier.startsWith("@/shared/"),
  },
];

function relative(file) {
  return path.relative(srcRoot, file).replaceAll(path.sep, "/");
}

async function collectFiles(dir) {
  const entries = await readdir(dir, { withFileTypes: true });
  const files = await Promise.all(
    entries.map(async (entry) => {
      const fullPath = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        if (entry.name === "generated") return [];
        return collectFiles(fullPath);
      }
      if (!sourceExtensions.has(path.extname(entry.name))) return [];
      return [fullPath];
    })
  );
  return files.flat();
}

function extractImportSpecifiers(source) {
  const specifiers = [];
  const patterns = [
    /\bimport\s+(?:type\s+)?(?:[^'"]*?\s+from\s+)?["']([^"']+)["']/g,
    /\bexport\s+(?:type\s+)?(?:[^'"]*?\s+from\s+)["']([^"']+)["']/g,
    /\bimport\s*\(\s*["']([^"']+)["']\s*\)/g,
  ];

  for (const pattern of patterns) {
    for (const match of source.matchAll(pattern)) {
      specifiers.push(match[1]);
    }
  }
  return specifiers;
}

const files = await collectFiles(srcRoot);
const violations = [];

for (const file of files) {
  const source = await readFile(file, "utf8");
  const specifiers = extractImportSpecifiers(source);
  for (const specifier of specifiers) {
    for (const rule of rules) {
      if (rule.appliesTo(file) && rule.forbidden(specifier)) {
        violations.push(
          `${relative(file)} imports ${specifier} (${rule.name})`
        );
      }
    }
  }
}

if (violations.length > 0) {
  console.error("Import boundary violations:");
  for (const violation of violations) {
    console.error(`- ${violation}`);
  }
  process.exit(1);
}

console.log(`Import boundary check passed (${files.length} files scanned).`);
