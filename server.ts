import { fetchBase64, renderSPA } from "./template.ts";

const fonts: string[] = [
  await fetchBase64("http://rsms.me/inter/font-files/InterVariable.woff2"),
  await fetchBase64(
    "http://rsms.me/inter/font-files/InterVariable-Italic.woff2",
  ),
];

export default {
  fetch(_req: Request): Response {
    return new Response(
      renderSPA({
        names: [
          "Sistem Pendukung Keputusan",
          "Pemilihan Penyedia Layanan Cloud",
        ],
        author: "Muhammad Rizki Fauzan",
        fonts,
      }),
      {
        headers: { "content-type": "text/html" },
      },
    );
  },
} satisfies Deno.ServeDefaultExport;
