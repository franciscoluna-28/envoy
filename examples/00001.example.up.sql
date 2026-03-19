CREATE TYPE "public"."vec_protocol" AS ENUM('v_alpha', 'v_beta', 'v_gamma');

CREATE TYPE "public"."node_health" AS ENUM('stable', 'degraded', 'critical');

CREATE TABLE "vector_shards" (
	"sid" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"p_id" uuid,
	"v_hash" text NOT NULL,
	"protocol" "vec_protocol" DEFAULT 'v_alpha' NOT NULL,
	"is_replicated" boolean DEFAULT false NOT NULL,
	"raw_payload" jsonb DEFAULT '{}'::jsonb NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "v_hash_unique" UNIQUE("v_hash")
);

CREATE TABLE "node_metrics" (
    "mid" text PRIMARY KEY NOT NULL,
    "shard_id" uuid NOT NULL,
    "health_status" "node_health" DEFAULT 'stable' NOT NULL,
    "latency_ms" integer DEFAULT 0 NOT NULL,
    "throughput_kb" numeric(12, 4) DEFAULT 0.0000 NOT NULL,
    "captured_at" timestamp
    with
        time zone DEFAULT now () NOT NULL
);

ALTER TABLE "node_metrics"
ADD CONSTRAINT "node_metrics_shard_id_fkey" FOREIGN KEY ("shard_id") REFERENCES "public"."vector_shards" ("sid") ON DELETE cascade ON UPDATE no action;

CREATE INDEX "idx_shards_protocol" ON "vector_shards" USING btree ("protocol");

CREATE INDEX "idx_metrics_health" ON "node_metrics" USING btree ("health_status")
WHERE
    "health_status" = 'critical';