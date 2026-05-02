export type typeCriteria = "benefit" | "cost";
interface table {
  name: string;
}

export type nameTable = "alt" | "criteria";
export interface dataCriteria extends table {
  nameL10N: string;
  type: typeCriteria;
  weight: number;
}
interface valueCriteria extends table {
  value: number;
}
export interface dataAlt extends table {
  values: valueCriteria[];
}
