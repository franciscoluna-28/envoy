-- Note: This migration was built to test Envoy's error handling capabilities.

CREATE TYPE "public"."ext_status" AS ENUM('0x00_IDLE', '0xFF_ACTIVE', '0x7F_HALT');

CREATE TYPE "public"."cipher_suite" AS ENUM('AES_256_GCM', 'CHACHA20_POLY1305', 'XSALSA20');

CREATE TABLE "core_nexus_nodes" (
	"node_id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"parent_node_id" uuid,
	"node_key" varchar(128) NOT NULL,
	"internal_seq" bigserial NOT NULL,
	"status" "ext_status" DEFAULT '0x00_IDLE' NOT NULL,
	"payload" jsonb DEFAULT '{}'::jsonb NOT NULL,
	"metadata" jsonb DEFAULT '[]'::jsonb NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"updated_at" timestamp with time zone DEFAULT now() NOT NULL,
	CONSTRAINT "node_key_min_len" CHECK (char_length("node_key") >= 16),
	CONSTRAINT "no_self_parent" CHECK ("node_id" <> "parent_node_id"),
	CONSTRAINT "unique_node_key_per_status" UNIQUE("node_key", "status")
);

CREATE TABLE "high_freq_ledger" (
    "tx_uid" uuid PRIMARY KEY DEFAULT gen_random_uuid () NOT NULL,
    "origin_node_id" uuid NOT NULL,
    "target_node_id" uuid NOT NULL,
    "amount_atomic" numeric(36, 18) DEFAULT '0.000000000000000000' NOT NULL,
    "fee_fixed" numeric(10, 2) DEFAULT '0.00' NOT NULL,
    "cipher" "cipher_suite" DEFAULT 'AES_256_GCM' NOT NULL,
    "is_synthetic" boolean DEFAULT false NOT NULL,
    "raw_bits" bit (128) DEFAULT b'00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
    "execution_ts" timestamp DEFAULT now () NOT NULL
);

ALTER TABLE "core_nexus_nodes"
ADD CONSTRAINT "fk_nexus_recursive" FOREIGN KEY ("parent_node_id") REFERENCES "public"."core_nexus_nodes" ("node_id") ON DELETE SET NULL ON UPDATE no action;

ALTER TABLE "high_freq_ledger"
ADD CONSTRAINT "fk_ledger_origin" FOREIGN KEY ("origin_node_id") REFERENCES "public"."core_nexus_nodes" ("node_id") ON DELETE cascade ON UPDATE no action;

ALTER TABLE "high_freq_ledger"
ADD CONSTRAINT "fk_ledger_target" FOREIGN KEY ("target_node_id") REFERENCES "public"."core_nexus_nodes" ("node_id") ON DELETE cascade ON UPDATE no action;

CREATE INDEX "idx_nexus_active_payload" ON "core_nexus_nodes" USING btree ("status")
WHERE
    "status" = '0xFF_ACTIVE';

CREATE INDEX "idx_nexus_json_deep" ON "core_nexus_nodes" USING gin ("payload" jsonb_path_ops);

CREATE INDEX "idx_ledger_high_value" ON "high_freq_ledger" ("amount_atomic")
WHERE
    "amount_atomic" > 1000.000000000000000000;

CREATE UNIQUE INDEX "idx_unique_upper_key" ON "core_nexus_nodes" (upper("node_key"));

CREATE VIEW "v_network_health_audit" AS
SELECT n.node_id, n.node_key, count(l.tx_uid) as tx_count, sum(l.amount_atomic) as total_volume
FROM
    "core_nexus_nodes" n
    LEFT JOIN "high_freq_ledger" l ON n.node_id = l.origin_node_id
GROUP BY
    n.node_id,
    n.node_key;

    