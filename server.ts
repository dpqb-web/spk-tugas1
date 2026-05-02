import { renderSPA } from "./template.ts";

export default {
  fetch(_req: Request): Response {
    return new Response(
      renderSPA({
        name: "Sistem Pendukung Keputusan - Pemilihan Penyedia Layanan Cloud",
        description:
          "Aplikasi Sistem Pendukung Keputusan bertema Pemilihan Penyedia Layanan Cloud",
        author: "Muhammad Rizki Fauzan",
      }),
      {
        headers: { "content-type": "text/html" },
      },
    );
  },
} satisfies Deno.ServeDefaultExport;
