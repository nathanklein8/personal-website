import { getDb } from '$lib/server/db/client';
import { test } from '$lib/server/db/schema';
import { json, type RequestEvent } from '@sveltejs/kit';
import { eq, sql } from 'drizzle-orm';

export async function POST({ params }: RequestEvent) {
  const db = getDb();
  const id = Number(params.id);
  if (isNaN(id)) {
    return json({ error: 'Invalid ID' }, { status: 400 });
  }

  let updated;

  await db.transaction(async (tx) => {
    // Try to update the row
    const result = await tx
      .update(test)
      .set({ count: sql`${test.count} + 1` })
      .where(eq(test.id, id))
      .returning();

    if (result.length === 0) {
      // Row doesn't exist — insert it with count = 0
      await tx.insert(test).values({ id, count: 0 });
      // Increment the count after inserting
      const newResult = await tx
        .update(test)
        .set({ count: sql`${test.count} + 1` })
        .where(eq(test.id, id))
        .returning();

      updated = newResult[0];
    } else {
      updated = result[0];
    }
  });

  return json(updated);
}
