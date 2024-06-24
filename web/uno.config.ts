import { defineConfig, presetAttributify, presetTypography, presetUno } from "unocss"

export default defineConfig({
  presets: [presetAttributify(), presetTypography(), presetUno({ preflight: false })],
})
