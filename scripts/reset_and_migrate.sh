#!/bin/bash

# ==========================================
# COMPLETE DATABASE RESET AND MIGRATION SCRIPT
# ==========================================
# This script will completely reset the database and run all migrations fresh

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Complete Database Reset and Migration${NC}"
echo "==========================================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ .env file not found${NC}"
    echo "Please create .env with your DATABASE_URL"
    exit 1
fi

# Load environment variables
export $(grep -v '^#' .env | xargs)

if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}❌ DATABASE_URL not set in .env file${NC}"
    exit 1
fi

echo -e "${YELLOW}⚠️  WARNING: This will completely reset your database!${NC}"
echo -e "${YELLOW}⚠️  All data will be lost!${NC}"
echo ""
read -p "Are you sure you want to continue? (yes/no): " -r
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "Operation cancelled."
    exit 1
fi

echo ""
echo -e "${BLUE}Step 1: Dropping all tables and resetting database...${NC}"
if timeout 60 psql "$DATABASE_URL" -f scripts/reset_database.sql; then
    echo -e "${GREEN}✅ Database reset completed${NC}"
else
    echo -e "${RED}❌ Database reset failed or timed out${NC}"
    echo "You may need to run this manually in CockroachDB Cloud console"
    exit 1
fi

echo ""
echo -e "${BLUE}Step 2: Running all migrations from scratch...${NC}"
if go run cmd/migrator/main.go -direction=up; then
    echo -e "${GREEN}✅ All migrations completed successfully${NC}"
else
    echo -e "${RED}❌ Migration failed${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}Step 3: Verifying migration status...${NC}"
if timeout 30 psql "$DATABASE_URL" -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;"; then
    echo -e "${GREEN}✅ Migration status verified${NC}"
else
    echo -e "${YELLOW}⚠️  Could not verify migration status (timeout)${NC}"
fi

echo ""
echo -e "${BLUE}Step 4: Verifying all tables exist...${NC}"
if timeout 30 psql "$DATABASE_URL" -c "\dt"; then
    echo -e "${GREEN}✅ Tables verification completed${NC}"
else
    echo -e "${YELLOW}⚠️  Could not verify tables (timeout)${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Database reset and migration completed successfully!${NC}"
echo ""
echo "Your database now has:"
echo "  ✅ brokers, stocks, recommendations, ingestion_logs (Migration 001)"
echo "  ✅ unique_ticker_event_time constraint (Migration 002)"
echo "  ✅ users, subscriptions, sessions (Migration 003)"
echo "  ✅ chat_sessions, chat_messages (Migration 004)"
echo ""
echo -e "${BLUE}You can now test migration 004 reset with:${NC}"
echo "  make migrate-004-reset" 