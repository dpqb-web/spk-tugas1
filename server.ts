import { join, resolve } from "jsr:@std/path";
import { Database } from "jsr:@db/sqlite";

import { renderSPA } from "./template.ts";

const PATH_CACHE: string = resolve("cache");

const api = (data: FormData): object => {
  // const db: Database = new Database(join(PATH_CACHE, "db.sqlite"));
  console.log(data);
  // db.close();
  return { message: "hello world" };
};

export default {
  async fetch(req: Request): Promise<Response> {
    const url: URL = new URL(req.url);

    if (url.pathname == "/api") {
      const data: FormData = await req.formData();
      if (!data.has("_method")) {
        data.append("_method", req.method);
      }

      return new Response(JSON.stringify(api(data)), {
        headers: { "content-type": "application/json" },
      });
    }

    return new Response(
      renderSPA({
        name: "text",
        description: "qwerty",
        author: "Muhammad Rizki Fauzan",
      }),
      {
        headers: { "content-type": "text/html" },
      },
    );
  },
} satisfies Deno.ServeDefaultExport;
