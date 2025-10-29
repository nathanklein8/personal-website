import { getDb } from '$lib/server/db/client';
import { test } from '$lib/server/db/schema';
import { json, type RequestEvent } from '@sveltejs/kit';
import { eq } from 'drizzle-orm';

const db = getDb();

export async function GET({ params }: RequestEvent) {
  const id = Number(params.id);
  if (isNaN(id)) {
    return json({ error: 'Invalid ID' }, { status: 400 });
  }

  const [counter] = await db.select().from(test).where(eq(test.id, id));

  if (!counter) {
    return json({ error: 'ID not found' }, { status: 404 });
  }

  return json(counter);
}
