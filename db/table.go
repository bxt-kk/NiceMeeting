package db

type Table struct {
    name string
    query string
}

var tables = []Table{
    {"user", `CREATE TABLE "user" (
        "id"	INTEGER,
        "name"	TEXT NOT NULL,
        "email"	TEXT NOT NULL UNIQUE,
        "password"	TEXT NOT NULL,
        "status"	TEXT,
        PRIMARY KEY("id" AUTOINCREMENT)
    );`},

    {"meeting", `CREATE TABLE "meeting" (
        "id"	INTEGER,
        "title"	TEXT NOT NULL DEFAULT 'untitled',
        "owner_id"	INTEGER NOT NULL,
        "content"	TEXT NOT NULL,
        "status"	TEXT,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY("owner_id") REFERENCES "user"("id") ON DELETE RESTRICT
    );`},

    {"feedback", `CREATE TABLE "feedback" (
        "id"	INTEGER,
        "audience_id"	INTEGER NOT NULL,
        "meeting_id"	INTEGER NOT NULL,
        "type"	TEXT NOT NULL,
        "value"	INTEGER NOT NULL DEFAULT 0,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY("audience_id") REFERENCES "user"("id") ON DELETE CASCADE,
        FOREIGN KEY("meeting_id") REFERENCES "meeting"("id") ON DELETE CASCADE
    );`},

    {"Link", `CREATE TABLE "Link" (
        "id"	INTEGER,
        "from_id"	INTEGER NOT NULL,
        "to_id"	INTEGER NOT NULL,
        "description"	TEXT,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY("from_id") REFERENCES "meeting"("id") ON DELETE CASCADE,
        FOREIGN KEY("to_id") REFERENCES "meeting"("id") ON DELETE CASCADE
    );`},

}
