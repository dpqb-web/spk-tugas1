import { resolve } from "jsr:@std/path";

import {
  compileFile as compileFilePug,
  compileTemplate as compileTemplatePug,
} from "npm:pug";
import {
  CompileResult as CompileResultSass,
  compileString as compileStringSass,
} from "npm:sass";
import { compile as compileCoffee } from "npm:coffeescript";
import { minify_sync, MinifyOutput } from "npm:terser";

export const renderSPA = (locals?: object): string => {
  const tmp: compileTemplatePug = compileFilePug(resolve("index.pug"), {
    filters: {
      sass: (txt: string): string => {
        const css: CompileResultSass = compileStringSass(txt, {
          style: "compressed",
          syntax: "indented",
        });
        return "<style>" + css.css + "</style>";
      },
      coffee: (txt: string): string => {
        const js: MinifyOutput = minify_sync(compileCoffee(txt), {
          toplevel: true,
          nameCache: {},
        });
        return "<script>" + js.code + "</script>";
      },
    },
  });
  return tmp(locals);
};
