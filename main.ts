import { Webview } from "jsr:@webview/webview";

import { SPA } from "./server.ts";
import * as db from "./db.json.ts";

// NOTE TO SELF: never use Fedora again. can't even do good things. fuck it.

export const app: Webview = new Webview();

app.navigate(`data:text/html,${encodeURIComponent(SPA)}`);
app.bind("select", db.select);
app.bind("insert", db.insert);
app.bind("update", db.update);
app.bind("drop", db.drop);
app.run();
