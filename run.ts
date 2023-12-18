#!/usr/bin/env -S deno run --unstable -A

const dylib = Deno.dlopen("./libteleport.so", {
  Run2: { parameters: ["buffer", "buffer"], result: "void" },
  Run: { parameters: ["buffer"], result: "void", nonblocking: true },
});

function s2b(s) {
  const buffer = new TextEncoder().encode(s);
  return buffer;
}

function a2b(a) {
  return new Uint8Array([...a.map((s) => [...s2b(s), 0]).flat()]);
}

export function runGoFunction2() {
  dylib.symbols.Run2(
    s2b("https://ufo.k0s.io"),
    s2b("https://k0s.io"),
  );
}

export async function runGoFunction(n) {
  await dylib.symbols.Run(
    a2b(n),
  );
}

console.log("Deno.args:", Deno.args);

// Example usage
// await runGoFunction2();
// await runGoFunctionN(["https://ufo.k0s.io", "https://k0s.io"]);
// await runGoFunction(["https://ufo.k0s.io/test/test1/test2/test3", "https://k0s.io"]);
// runGoFunction(["https://ufo.k0s.io", "https://k0s.io"]);
await runGoFunction(Deno.args);

console.log("EOF");
