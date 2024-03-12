declare global {
  interface Window {
    __LAUNCHPAD_TARGET__?: string
  }
}

export async function request(input: string, init?: RequestInit) {
  const prefix = window.__LAUNCHPAD_TARGET__ ?? ""
  return await fetch(prefix + input, init)
}
