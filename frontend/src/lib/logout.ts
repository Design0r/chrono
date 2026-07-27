// Eigenes Modul, damit die API-Schicht nicht auf auth.tsx zeigt (Zirkelimport).
let logoutFn: (() => Promise<void>) | null = null;

export function registerLogout(fn: () => Promise<void>) {
  logoutFn = fn;
}

export async function logoutOutsideReact() {
  if (logoutFn) {
    await logoutFn();
  }
}
