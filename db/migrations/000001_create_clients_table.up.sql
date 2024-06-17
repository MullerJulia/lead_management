CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    priority INTEGER NOT NULL,
    leadCapacity INTEGER NOT NULL,
    currentLeadCount INTEGER NOT NULL,
    workingHoursStart TEXT NOT NULL,
    workingHoursEnd TEXT NOT NULL
);
