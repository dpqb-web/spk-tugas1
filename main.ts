import { getPort } from "jsr:@openjs/port-free";
import { Webview } from "jsr:@webview/webview";

import fetch from "./server.ts";

const PORT_SERVE: number = await getPort({ port: 8000 });
// const APP: Webview = new Webview();

Deno.serve({
  hostname: "127.0.0.1",
  port: PORT_SERVE,
  onListen: ({ port, hostname }) => {
    console.log(`http://${hostname}:${port}`);
    // APP.navigate(`http://${hostname}:${port}`);
    // APP.run();
  },
}, fetch.fetch);
