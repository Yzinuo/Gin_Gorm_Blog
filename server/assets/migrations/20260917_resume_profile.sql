-- Apply once before deploying the résumé editor if DbAutoMigrate is false.
-- With DbAutoMigrate enabled, GORM creates the equivalent table at startup.
CREATE TABLE IF NOT EXISTS resume_profile (
    id INTEGER NOT NULL PRIMARY KEY,
    document LONGTEXT NOT NULL
);
