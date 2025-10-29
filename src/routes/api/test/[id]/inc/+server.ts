import { getDb } from '$lib/server/db/client';
import { test } from '$lib/server/db/schema';
import { json, type RequestEvent } from '@sveltejs/kit';
import { eq, sql } from 'drizzle-orm';

const db = getDb();

export async function POST({ params }: RequestEvent) {
  const id = Number(params.id);
  if (isNaN(id)) {
    return json({ error: 'Invalid ID' }, { status: 400 });
  }

  // increment the count atomically
  const [updated] = await db
    .update(test)
    .set({ count: sql`${test.count} + 1 ` })
    .where(eq(test.id, id))
    .returning();

  if (!updated) {
    return json({ error: 'ID not found' }, { status: 404 });
  }

  return json(updated);
}
