export interface ClaimType {
  icon: string
  name: string
  description: string
}

export const claims: { [id: string]: ClaimType } = {
  openid: {
    icon: "mdi-identifier",
    name: "Open Identity",
    description: "Allow them to read your personal information.",
  },
}
