# something cheatsheet

Below is a handy PostgreSQL cheat sheet geared toward common tasks you might perform when interacting with your database and tables. This assumes you’re working primarily through `psql` in a terminal environment and using a database user that has appropriate permissions.

---

## Connecting to PostgreSQL

**Connect to a Database:**

```bash
psql -U <username> -d <database>
```

*e.g.* `psql -U chilli -d ppbadb`

**Switch Databases inside psql (if allowed):**

```sql
\c <database_name>
```

**Quit psql:**

```sql
\q
```

---

### Viewing Database Objects & Info

**List Databases:**

```sql
\l
```

**List Tables:**

```sql
\dt
```

**Describe a Table’s Structure:**

```sql
\d tablename
```

**Show Current User:**

```sql
SELECT current_user;
```

**Show Current Database:**

```sql
SELECT current_database();
```

---

### Creating & Managing Tables

**Create a Table:**

```sql
CREATE TABLE tablename (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    proof BYTEA,
    helper_data BYTEA
);
```

**Drop a Table:**

```sql
DROP TABLE IF EXISTS tablename;
```

**Add a Column:**

```sql
ALTER TABLE tablename
ADD COLUMN newcolumn TEXT;
```

**Remove a Column:**

```sql
ALTER TABLE tablename
DROP COLUMN column_name;
```

**Change a Column’s Data Type:**

```sql
ALTER TABLE tablename
ALTER COLUMN column_name TYPE new_type
USING column_name::new_type;
```

**Rename a Column:**

```sql
ALTER TABLE tablename
RENAME COLUMN old_name TO new_name;
```

**Rename a Table:**

```sql
ALTER TABLE old_tablename
RENAME TO new_tablename;
```

---

### Inserting, Updating, & Deleting Data

**Insert Data:**

```sql
INSERT INTO tablename (username, proof, helper_data)
VALUES ('user1', '\xDEADBEEF', '\xCAFEBABE');
```

*(For text: `VALUES ('user1', 'some text', 'more text');`)*

**Update Data:**

```sql
UPDATE tablename
SET proof = '\x00FF11', helper_data = NULL
WHERE username = 'user1';
```

**Delete Data:**

```sql
DELETE FROM tablename
WHERE username = 'user1';
```

---

### Querying Data

**Select All Rows:**

```sql
SELECT * FROM tablename;
```

**Select Specific Columns:**

```sql
SELECT username, proof FROM tablename;
```

**Filter Rows (WHERE clause):**

```sql
SELECT * FROM tablename
WHERE username = 'user1';
```

**Order Results:**

```sql
SELECT * FROM tablename
ORDER BY username ASC;
```

**Limit Results:**

```sql
SELECT * FROM tablename
LIMIT 5;
```

**Count Rows:**

```sql
SELECT COUNT(*) FROM tablename;
```

---

### Transactions

**Start a Transaction:**

```sql
BEGIN;
```

**Commit a Transaction:**

```sql
COMMIT;
```

**Rollback a Transaction:**

```sql
ROLLBACK;
```

### (Useful for making multiple changes atomically.)*

---

### Indexes & Performance

**Create an Index (e.g., on username):**

```sql
CREATE INDEX idx_tablename_username ON tablename (username);
```

**Drop an Index:**

```sql
DROP INDEX idx_tablename_username;
```

---

### User & Permission Management (If Needed)

**Create a New User (Role):**

```sql
CREATE ROLE newuser WITH LOGIN PASSWORD 'newpassword';
```

**Grant Privileges:**

```sql
GRANT SELECT, INSERT, UPDATE, DELETE ON tablename TO newuser;
```

---

### Utility Commands & Help

**Show Help for `psql` Commands:**

```sql
\?
```

**Show Help for SQL Commands:**

```sql
\h <command>
```

*e.g.* `\h CREATE TABLE`

**Show Current Search Path:**

```sql
SHOW search_path;
```

---

### Backup & Restore

**Backup a Database (run outside psql):**

```bash
pg_dump -U username -d database_name > backup.sql
```

**Restore a Database (run outside psql):**

```bash
psql -U username -d database_name < backup.sql
```

---

This cheat sheet covers a broad range of common tasks you might perform: connecting to the database, creating and modifying tables, inserting and querying data, and performing routine maintenance. Keep it handy as a quick reference during your work with PostgreSQL.
