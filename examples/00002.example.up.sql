CREATE TYPE "public"."tunnel_enc" AS ENUM('aes_gcm', 'chacha20', 'rsa_4096');

CREATE TABLE "auth_gateways" (
	"gid" text PRIMARY KEY NOT NULL,
	"endpoint_v4" text NOT NULL,
	"encryption" "tunnel_enc" DEFAULT 'aes_gcm' NOT NULL,
	"max_sessions" integer DEFAULT 1000 NOT NULL,
	"is_active" boolean DEFAULT true NOT NULL,
	"meta_config" jsonb DEFAULT '{"retries": 3}'::jsonb NOT NULL,
	"last_heartbeat" timestamp,
	CONSTRAINT "gateway_ipv4_unique" UNIQUE("endpoint_v4")
);

CREATE TABLE "session_keys" (
    "kid" uuid PRIMARY KEY DEFAULT gen_random_uuid () NOT NULL,
    "gateway_id" text NOT NULL,
    "key_blob" text NOT NULL,
    "expires_at" timestamp NOT NULL,
    "created_at" timestamp DEFAULT now () NOT NULL
);

ALTER TABLE "session_keys"
ADD CONSTRAINT "session_keys_gateway_id_fkey" FOREIGN KEY ("gateway_id") REFERENCES "public"."auth_gateways" ("gid") ON DELETE cascade ON UPDATE no action;

CREATE INDEX "idx_keys_expiry" ON "session_keys" USING btree ("expires_at");

CREATE INDEX "idx_gateways_status" ON "auth_gateways" USING btree ("is_active");