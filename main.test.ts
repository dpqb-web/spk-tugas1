import { join } from "jsr:@std/path";
import { ensureDir, exists } from "jsr:@std/fs";

import { PATH_CACHE } from "./main.ts";
import { dataAlt, dataCriteria, nameTable } from "./db.json.d.ts";

await ensureDir(PATH_CACHE);

const tables: nameTable[] = ["alt", "criteria"];

for (const table of tables) {
  const filepath: string = join(PATH_CACHE, `${table}.json`);
  if (!await exists(filepath)) {
    if (table == "criteria") {
      const data: dataCriteria[] = [
        {
          name: "price",
          nameL10N: "Harga",
          type: "cost",
          weight: 0.2,
        },
        {
          name: "memory",
          nameL10N: "Memori",
          type: "benefit",
          weight: 0.2,
        },
        {
          name: "storage",
          nameL10N: "Penyimpanan",
          type: "benefit",
          weight: 0.2,
        },
        {
          name: "uptime",
          nameL10N: "Ketersediaan",
          type: "benefit",
          weight: 0.2,
        },
        {
          name: "security",
          nameL10N: "Keamanan",
          type: "benefit",
          weight: 0.2,
        },
      ];
      await Deno.writeTextFile(filepath, JSON.stringify(data));
    }
    if (table == "alt") {
      const data: dataAlt[] = [
        {
          name: "Amazon Web Services",
          values: [
            {
              name: "price",
              value: 10,
            },
            {
              name: "memory",
              value: 4,
            },
            {
              name: "storage",
              value: 80,
            },
            {
              name: "uptime",
              value: 100,
            },
            {
              name: "security",
              value: 3,
            },
          ],
        },
        {
          name: "Google Cloud Platform",
          values: [
            {
              name: "price",
              value: 15,
            },
            {
              name: "memory",
              value: 8,
            },
            {
              name: "storage",
              value: 160,
            },
            {
              name: "uptime",
              value: 90,
            },
            {
              name: "security",
              value: 2,
            },
          ],
        },
        {
          name: "Microsoft Azure",
          values: [
            {
              name: "price",
              value: 20,
            },
            {
              name: "memory",
              value: 16,
            },
            {
              name: "storage",
              value: 200,
            },
            {
              name: "uptime",
              value: 80,
            },
            {
              name: "security",
              value: 3,
            },
          ],
        },
      ];
      await Deno.writeTextFile(filepath, JSON.stringify(data));
    }
  }
}
