-- Supabase Schema Migration Script
-- Auto-generated to match CercaChain FlutterFlow schema

CREATE TABLE IF NOT EXISTS public."analytics" (
  "id" text PRIMARY KEY,
  "user_ref" text,
  "user_objects" text,
  "user_orders" text,
  "object_views_accum" text,
  "user_object_count" text,
  "object_reach" text,
  "user_impressions" text,
  "object_voterate" text,
  "object_pin" text,
  "object_share" text,
  "top_performer" text,
  "order_avg_ref" text,
  "order_avg_time" text,
  "order_avg_review" text,
  "bag_order_ratio" text,
  "order_accum_ref" text,
  "order_accum" text,
  "user_pin" text,
  "market_index" text,
  "performance" text,
  "user_hash" text,
  "objectAvgReference" text,
  "orderAvgReference" text,
  "order_pubrev_list" text,
  "order_user_rev_list" text,
  "pubusers_ref" text,
  "orders_as_maker" text,
  "date_created" text,
  "ratings" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."bag" (
  "id" text PRIMARY KEY,
  "user_ref" text,
  "items" text,
  "in_bag_items" text,
  "public_user_ref" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."catalogue" (
  "id" text PRIMARY KEY,
  "user_ref" text,
  "catalalogue_name" text,
  "catalogue_choicechips" text,
  "catalogue_items" text,
  "total_items" text,
  "catalogue_buzzwords" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."chats" (
  "id" text PRIMARY KEY,
  "users" text,
  "user_a" text,
  "user_b" text,
  "last_message" text,
  "last_message_time" text,
  "last_message_sent_by" text,
  "last_message_seen_by" text,
  "group_chat_id" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."chat_messages" (
  "id" text PRIMARY KEY,
  "user" text,
  "chat" text,
  "text" text,
  "timestamp" text,
  "image" text,
  "video" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."credits" (
  "id" text PRIMARY KEY,
  "userRef" text,
  "f_x" text,
  "generation_i" text,
  "participation_i" text,
  "transition_i" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."interactions" (
  "id" text PRIMARY KEY,
  "user_ref" text,
  "type" text,
  "value" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."order_methods" (
  "id" text PRIMARY KEY,
  "method_poster" text,
  "method_type" text,
  "method_tag" text,
  "method_name" text,
  "method_thread" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."order" (
  "id" text PRIMARY KEY,
  "date" text,
  "ref_value" text,
  "total_ref_value" text,
  "order_completed" text,
  "payment_method" text,
  "order_comlpletion_date" text,
  "user_ref" text,
  "items" text,
  "root_item" text,
  "ordersCompleted" text,
  "items_inorder" text,
  "publicuser_ref" text,
  "order_stats" text,
  "order_method" text,
  "wallet_method" text,
  "username" text,
  "publicusername" text,
  "order_image" text,
  "is_orderimage_uploaded" text,
  "method_order" text,
  "method_wallet" text,
  "is_order_accepted" text,
  "is_orderimage_accepted" text,
  "key_1" text,
  "key_2" text,
  "rev_user" text,
  "rev_pubuser" text,
  "ref_list_orders" text,
  "revby_maker" text,
  "revby_taker" text,
  "order_finished" text,
  "order_closed" text,
  "rating" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."ratings" (
  "id" text PRIMARY KEY,
  "user_ref" text,
  "value" text,
  "comment" text,
  "date" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."submission" (
  "id" text PRIMARY KEY,
  "image" text,
  "date" text,
  "header" text,
  "poster" text,
  "imagesextra" text,
  "type1choice" text,
  "type2choice" text,
  "video" text,
  "audio" text,
  "body" text,
  "upvote" text,
  "downvote" text,
  "refvalue" text,
  "type0choice" text,
  "submitted_date" text,
  "type_order" text,
  "type_activity" text,
  "type_ticket" text,
  "submitter" text,
  "profilePictureObjectPost" text,
  "ups" text,
  "pins" text,
  "object_is_edited" text,
  "shares" text,
  "type0object" text,
  "type_object" text,
  "thread" text,
  "object_thread" text,
  "object_in_bag" text,
  "object_refvalue" text,
  "analytics_ref" text,
  "object_views" text,
  "object_viewers" text,
  "objectViewsS" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."thread" (
  "id" text PRIMARY KEY,
  "user_ref" text,
  "public_ref" text,
  "submission_ref" text,
  "object_thread" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."transactions" (
  "id" text PRIMARY KEY,
  "user_madeby" text,
  "user_madeto" text,
  "items" text,
  "id_transaction" text,
  "date" text,
  "total_ref_value_transaction" text,
  "order_ref" text,
  "w_method" text,
  "o_method" text,
  "UserIn" text,
  "userOut" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."usernames" (
  "id" text PRIMARY KEY,
  "usernames_in_use" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."users" (
  "id" text PRIMARY KEY,
  "bio" text,
  "email" text,
  "location" text,
  "name" text,
  "surname" text,
  "username" text,
  "display_name" text,
  "photo_url" text,
  "uid" text,
  "created_time" text,
  "phone_number" text,
  "last_active_time" text,
  "role" text,
  "title" text,
  "banner" text,
  "profile_picture" text,
  "usernames" text,
  "pinned_users" text,
  "user_pins" text,
  "bag_ref" text,
  "user_occupations" text,
  "user_interests" text,
  "pinned_objects" text,
  "user_objects" text,
  "wallet_methods_user" text,
  "wallet_address" text,
  "user_bag_objects" text,
  "order_ref" text,
  "order_methods_user" text,
  "user_places" text,
  "user_type" text,
  "user_items" text,
  "user_posts" text,
  "user_transactions" text,
  "shortDescription" text,
  "analytics_ref" text,
  "user_verified" text,
  "user_verified_pending" text,
  "user_is_admin" text,
  "u_ratings" text,
  "avg_u_rating" text,
  "ratings" text,
  "user_ratings" text,
  "userCredits" text,
  "following_users" text,
  "users_follwoing_me" text,
  "verification_id_url" text,
  "verification_selfie_url" text,
  "verification_location_url" text,
  "verification_municipality" text,
  "verification_submitted_at" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."user_ratings" (
  "id" text PRIMARY KEY,
  "rated_user" text,
  "value" text,
  "comment" text,
  "date" text,
  "created_at" timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."wallet_methods" (
  "id" text PRIMARY KEY,
  "method_poster" text,
  "method_name" text,
  "method_type" text,
  "method_id" text,
  "method_account" text,
  "method_thread" text,
  "created_at" timestamp with time zone DEFAULT now()
);
