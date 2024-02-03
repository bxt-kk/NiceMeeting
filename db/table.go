package db

type Table struct {
    name string
    query string
}

var tables = []Table{
    {"site", `CREATE TABLE "site" (
        "id"	INTEGER,
        "name"	TEXT NOT NULL DEFAULT 'Hello NiceMeeting',
        "description"	TEXT NOT NULL,
        PRIMARY KEY("id" AUTOINCREMENT)
    );`},

    {"user", `CREATE TABLE "user" (
        "id"	INTEGER,
        "name"	TEXT NOT NULL,
        "email"	TEXT NOT NULL UNIQUE,
        "password"	TEXT NOT NULL,
        "status"	TEXT,
        "level"	INTEGER NOT NULL DEFAULT 0,
        "registration_time"	INTEGER NOT NULL,
        "last_login_time"	INTEGER NOT NULL DEFAULT 0,
        PRIMARY KEY("id" AUTOINCREMENT)
    );`},

    {"meeting", `CREATE TABLE "meeting" (
        "id"	INTEGER,
        "owner_id"	INTEGER NOT NULL,
        "title"	TEXT NOT NULL DEFAULT 'untitled',
        "type"	TEXT NOT NULL,
        "content"	TEXT NOT NULL,
        "status"	TEXT,
        "creation_time"	INTEGER NOT NULL,
        "last_edited_time"	INTEGER NOT NULL DEFAULT 0,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY("owner_id") REFERENCES "user"("id") ON DELETE RESTRICT
    );`},

    {"feedback", `CREATE TABLE "feedback" (
        "id"	INTEGER,
        "user_id"	INTEGER NOT NULL,
        "meeting_id"	INTEGER NOT NULL,
        "type"	TEXT NOT NULL,
        "value"	INTEGER NOT NULL DEFAULT 0,
        "time"	INTEGER NOT NULL,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY("user_id") REFERENCES "user"("id") ON DELETE CASCADE,
        FOREIGN KEY("meeting_id") REFERENCES "meeting"("id") ON DELETE CASCADE
    );`},

    {"link", `CREATE TABLE "link" (
        "id"	INTEGER,
        "from_id"	INTEGER NOT NULL,
        "to_id"	INTEGER NOT NULL,
        "description"	TEXT,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY("from_id") REFERENCES "meeting"("id") ON DELETE CASCADE,
        FOREIGN KEY("to_id") REFERENCES "meeting"("id") ON DELETE CASCADE
    );`},

    {"tag", `CREATE TABLE "tag" (
        "id"	INTEGER,
        "meeting_id"	INTEGER NOT NULL,
        "label"	TEXT NOT NULL,
        PRIMARY KEY("id" AUTOINCREMENT),
        FOREIGN KEY("meeting_id") REFERENCES "meeting"("id") ON DELETE CASCADE
    );`},

}
