BEGIN;

DROP TABLE IF EXISTS "families";

DROP TYPE IF EXISTS user_role CASCADE;

DROP TABLE IF EXISTS "users";

DROP TYPE IF EXISTS list_type CASCADE;

DROP TYPE IF EXISTS list_visibility CASCADE;

DROP TABLE IF EXISTS "lists";

DROP TYPE IF EXISTS item_lists_status CASCADE;

DROP TABLE IF EXISTS "item_lists";

DROP TYPE IF EXISTS wishlist_status CASCADE;

DROP TABLE IF EXISTS "wishlists";

DROP TABLE IF EXISTS "diary_items";

DROP TABLE IF EXISTS "locations";

DROP TABLE IF EXISTS "notifications";

DROP TABLE IF EXISTS "chats";

DROP TABLE IF EXISTS "messages";

DROP TABLE IF EXISTS "chat_participants";

DROP TABLE IF EXISTS "user_rewards";

DROP TABLE IF EXISTS "rewards";

DROP TABLE IF EXISTS "reward_redemptions";

COMMIT;