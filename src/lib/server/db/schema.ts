import { pgTable, serial, integer } from 'drizzle-orm/pg-core';

export const test = pgTable('test', {
	id: serial('id').primaryKey(),
	count: integer('count').default(0).notNull(),
});
