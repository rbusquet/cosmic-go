import { ActionFunction, Form, redirect, useActionData, useTransition } from "remix";
import { addBatch } from "~/cosmic/api.server";

export let action: ActionFunction = async ({ request }) => {
  let data = await request.formData();
  return await addBatch(Object.fromEntries(data.entries()));
};

export default function Index() {
  const transition = useTransition();
  const data = useActionData();
  console.log({data})
  return (
    <div style={{ fontFamily: "system-ui, sans-serif", lineHeight: "1.4" }}>
      <Form method="post">
        <fieldset disabled={transition.state === "submitting"}>
          <input type="text" name="reference" />
          <br />
          <input type="text" name="sku" />
          <br />
          <input type="number" name="purchased_quantity" />
          <br />
          <input type="date" name="eta" />
          <br />
          <button type="submit">Add stock</button>
        </fieldset>
      </Form>
    </div>
  );
}
