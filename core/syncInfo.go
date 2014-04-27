package core

type SyncInfo struct {
    Index           int
    SyncedEntities  int
    EntitiesPerHour int
    Target          string
    TaskId          int
    TaskStarted     string
    MinutesActive   int
    ElapsedTime     string
}
