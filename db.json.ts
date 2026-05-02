import { join, resolve } from "jsr:@std/path";
import { ensureDir } from "jsr:@std/fs";

const PATH_CACHE: string = resolve("cache");
await ensureDir(PATH_CACHE);

const write = async (table: string, data: object[]): Promise<void> =>
  await Deno.writeTextFile(
    join(PATH_CACHE, `${table}.json`),
    JSON.stringify(data),
  );

export const select = async (table: string): Promise<object[]> =>
  JSON.parse(await Deno.readTextFile(join(PATH_CACHE, `${table}.json`)));

export const insert = async (
  table: string,
  ...params: object[]
): Promise<object> => {
  const tmp: object[] = await select(table);
  tmp.push(params);

  await write(table, tmp);
  return tmp;
};

export const update = async (
  table: string,
  filter: (element: object, index: number, array: object[]) => boolean,
  params: object,
): Promise<object | null> => {
  const tmpAll: object[] = await select(table);
  const tmpSel: object[] = tmpAll.filter(filter);
  const tmpUnsel: object[] = tmpAll.filter((val) => !tmpSel.includes(val));

  if (!tmpSel) return null;

  for (const item of tmpSel) {
    Object.assign(item, params);
    Object.assign(tmpUnsel, item);
  }

  await write(table, tmpUnsel);
  return tmpSel;
};

export const drop = async (
  table: string,
  filter: (element: object, index: number, array: object[]) => boolean,
): Promise<object | null> => {
  const tmpAll: object[] = await select(table);
  const tmp2del: object[] = tmpAll.filter(filter);

  if (tmp2del.length == tmpAll.length) return null;

  const tmpDelD: object[] = tmpAll.filter((val) => !tmp2del.includes(val));

  await write(table, tmpDelD);
  return tmp2del;
};
