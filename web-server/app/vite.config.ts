import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import {existsSync, mkdirSync} from'fs';

export default ({ mode }) => {
  process.env = {...process.env, ...loadEnv(mode, process.cwd())};
  const defaultOutDir = "dist"
  const defaultBase = "/"
  const prefix = process.env.VITE_URI_PREFIX
  const base = !prefix ? defaultBase : prefix
  const outDir = !prefix? defaultOutDir : `${defaultOutDir}${prefix}`

  return defineConfig({
      plugins: [react()],
      base,
      build: {
        outDir,
      }
  });
}