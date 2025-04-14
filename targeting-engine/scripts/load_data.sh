#!/bin/bash
# Initialize database with sample data
psql $DATABASE_URL -f migrations/001_init_tables.up.sql
go run cmd/seed/main.go