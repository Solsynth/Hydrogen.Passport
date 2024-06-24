export async function request(input: string, init?: RequestInit) {
  return await fetch(input, init)
}
