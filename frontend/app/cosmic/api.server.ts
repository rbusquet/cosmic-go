import { z } from "zod";

const dateSchema = z.preprocess((arg) => {
  if (typeof arg == "string" || arg instanceof Date) return new Date(arg);
}, z.date());
const castToInt = z.preprocess((arg) => {
  if (typeof arg == "string") return parseInt(arg, 10);
  if (typeof arg == "number") return arg;
}, z.number().int());

let StockSchema = z.object({
  reference: z.string(),
  sku: z.string(),
  purchased_quantity: castToInt,
  eta: dateSchema,
});

export function addBatch(data: unknown) {
  console.log({ data });
  let stock = StockSchema.parse(data);
  let headers = new Headers();
  headers.set("content-type", "application/json");
  return fetch("http://cosmicgo:8080/stock", {
    method: "POST",
    headers,
    body: JSON.stringify(stock),
  });
}
