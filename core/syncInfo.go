package core

type SyncInfo struct {
    Index           int
    SyncedEntities  int
    EntitiesPerHour float32
    Target          string
    TaskId          int
    TaskStarted     string
    MinutesActive   int
    ElapsedTime     string
}
