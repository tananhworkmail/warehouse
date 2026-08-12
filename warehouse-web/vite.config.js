import { fileURLToPath, URL } from "node:url";
import fs from "node:fs";
import path from "node:path";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";

const certDir = fileURLToPath(new URL("./certs", import.meta.url));
const httpsKeyPath = path.join(certDir, "dev-server-key.pem");
const httpsCertPath = path.join(certDir, "dev-server-cert.pem");
const httpsPfxPath = path.join(certDir, "dev-server.pfx");
const hasHttpsPemCert = fs.existsSync(httpsKeyPath) && fs.existsSync(httpsCertPath);
const hasHttpsPfxCert = fs.existsSync(httpsPfxPath);
const httpsOptions = hasHttpsPfxCert
  ? {
      pfx: fs.readFileSync(httpsPfxPath),
      passphrase: "warehouse-dev"
    }
  : hasHttpsPemCert
    ? {
        key: fs.readFileSync(httpsKeyPath),
        cert: fs.readFileSync(httpsCertPath)
      }
    : undefined;

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()]
    }),
    Components({
      resolvers: [ElementPlusResolver()]
    })
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url))
    }
  },
  server: {
    host: "0.0.0.0", 
    port: 8084, 
    strictPort: true, 
    open: false,
    https: httpsOptions,
    proxy: {
      "/api/v1": {
        target: "http://127.0.0.1:3084",
        changeOrigin: true
      }
    }
  }
});
