import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";
import { quasar, transformAssetUrls } from "@quasar/vite-plugin";
import tailwindcss from "@tailwindcss/vite";
import { fileURLToPath, URL } from "node:url";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  // plugins: [vue(), wails("./bindings")],
  plugins: [
    vue({
      template: {
        transformAssetUrls,
      },
    }),

    // quasar({
    //   sassVariables: "src/quasar-variables.sass",
    // }),
    quasar({
      sassVariables: fileURLToPath(
        new URL("./src/quasar-variables.sass", import.meta.url)
      ),
    }),

    tailwindcss(),

    wails("./bindings"),
  ],
});
