# Manual Database Reset Instructions

If the automated reset script fails due to connectivity issues, you can manually reset the database through CockroachDB Cloud console:

## Option 1: CockroachDB Cloud Console

1. **Access your CockroachDB Cloud cluster:**
   - Go to https://cockroachlabs.cloud/
   - Navigate to your cluster: `hiring-test-stock-cluster-13493`

2. **Open SQL Console:**
   - Click on "SQL" in the cluster dashboard
   - This opens the web-based SQL console

3. **Run the reset commands one by one:**

```sql
-- Step 1: Drop all tables
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS chat_sessions CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS subscriptions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS ingestion_logs CASCADE;
DROP TABLE IF EXISTS recommendations CASCADE;
DROP TABLE IF EXISTS stocks CASCADE;
DROP TABLE IF EXISTS brokers CASCADE;

-- Step 2: Drop migration tracking
DROP TABLE IF EXISTS schema_migrations CASCADE;

-- Step 3: Drop custom types
DROP TYPE IF EXISTS subscription_plan CASCADE;
DROP TYPE IF EXISTS subscription_status CASCADE;
DROP TYPE IF EXISTS user_tier CASCADE;

-- Step 4: Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Step 5: Verify cleanup
SELECT 'Database reset complete' AS status;
```

4. **Run migrations from your local machine:**
```bash
# This should work now that the database is clean
go run cmd/migrator/main.go -direction=up
```

## Option 2: Local psql (if working)

If your local psql connection is working:

```bash
# Load environment and run reset
export $(grep -v '^#' .env | xargs)
psql "$DATABASE_URL" -f scripts/reset_database.sql

# Then run migrations
go run cmd/migrator/main.go -direction=up
```

## Verification

After either method, verify success:

```bash
# Check migration version (should be 4)
make migrate-status

# List all tables
export $(grep -v '^#' .env | xargs)
psql "$DATABASE_URL" -c "\dt"

# Test migration 004 reset
make migrate-004-reset
```

## Expected Tables After Reset

Your database should have:
- ✅ `brokers` (with initial data)
- ✅ `stocks` (empty, ready for ingestion)
- ✅ `recommendations` (empty)
- ✅ `ingestion_logs` (empty)
- ✅ `users` (empty, ready for authentication)
- ✅ `subscriptions` (empty)
- ✅ `sessions` (empty)
- ✅ `chat_sessions` (empty, ready for chat feature)
- ✅ `chat_messages` (empty, ready for chat feature)
- ✅ `schema_migrations` (tracking table, version 4) 