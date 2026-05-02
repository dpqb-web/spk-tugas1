import { resolve } from "jsr:@std/path";
import { Webview } from "jsr:@webview/webview";

import { renderSPA } from "./template.ts";
import * as db from "./db.json.ts";

export const PATH_CACHE: string = resolve("cache");

// NOTE TO SELF: never use Fedora again. can't even do good things. fuck it.
// uncomment below if you didn't use it

// export const app: Webview = new Webview();

// app.navigate(`data:text/html,${
//   encodeURIComponent(renderSPA({
//     name: "Sistem Pendukung Keputusan - Pemilihan Penyedia Layanan Cloud",
//     description:
//       "Aplikasi Sistem Pendukung Keputusan bertema Pemilihan Penyedia Layanan Cloud",
//     author: "Muhammad Rizki Fauzan",
//   }))
// }`);
// app.bind("select", db.select);
// app.bind("insert", db.insert);
// app.bind("update", db.update);
// app.bind("drop", db.drop);
// app.run();
