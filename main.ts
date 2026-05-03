import { Webview } from "jsr:@webview/webview";

import { fetchBase64, renderSPA } from "./template.ts";
import * as db from "./db.json.ts";

const fonts: string[] = [
  await fetchBase64("http://rsms.me/inter/font-files/InterVariable.woff2"),
  await fetchBase64(
    "http://rsms.me/inter/font-files/InterVariable-Italic.woff2",
  ),
];

// NOTE TO SELF: never use Fedora again. can't even do good things. fuck it.

export const app: Webview = new Webview();

app.navigate(`data:text/html,${
  encodeURIComponent(renderSPA({
    names: [
      "Sistem Pendukung Keputusan",
      "Pemilihan Penyedia Layanan Cloud",
    ],
    author: "Muhammad Rizki Fauzan",
    fonts,
  }))
}`);
app.bind("select", db.select);
app.bind("insert", db.insert);
app.bind("update", db.update);
app.bind("drop", db.drop);
app.run();
