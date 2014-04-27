package main

import (
    "encoding/json"
    "fmt"
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/mattbaird/elastigo/api"
    ecCore "github.com/mattbaird/elastigo/core"
    "github.com/ninya-io/ninya-operations/async"
    "github.com/ninya-io/ninya-operations/core"
    "github.com/ninya-io/ninya-operations/format"
    "os"
    "time"
)

func main() {
    api.Domain = os.Getenv("ES_ENDPOINT")
    api.Port = os.Getenv("ES_HTTP_PORT")

    index := "production_v4"

    m := martini.Classic()
    m.Use(render.Renderer(render.Options{
        Layout: "layout",
    }))

    m.Get("/", func(r render.Render) {

        out, _ := ecCore.SearchRequest(index, "info", nil, nil)

        if len(out.Hits.Hits) > 0 {

            infos := []core.SyncInfo{}

            sem := make(async.Semaphore, out.Hits.Len())

            for _, hit := range out.Hits.Hits {
                go func(hit ecCore.Hit) {

                    var syncInfo core.SyncInfo
                    json.Unmarshal(*hit.Source, &syncInfo)

                    searchJson := fmt.Sprintf(`{
                        "sort" : [
                            { "reputation" : {"order" : "desc"}}
                        ],
                        query: {
                          match: {
                              "_ninya_sync_task_id": %v
                          }
                        }
                    }`, syncInfo.TaskId)

                    taskInfo, _ := ecCore.SearchRequest(index, "user", nil, searchJson)
                    syncInfo.SyncedEntities = taskInfo.Hits.Total
                    syncInfo.Index = len(infos) + 1

                    taskStartedAt := syncInfo.TaskId / 1000

                    const layout = "Jan 2, 2006 at 3:04pm (MST)"
                    syncInfo.TaskStarted = time.Unix(int64(taskStartedAt), 0).Format(layout)

                    secondsElapsed := (int(time.Now().Unix()) - taskStartedAt)
                    syncInfo.ElapsedTime = format.Duration(secondsElapsed)

                    syncInfo.EntitiesPerHour = int(float32(syncInfo.SyncedEntities) / float32(secondsElapsed) * 60 * 60)
                    infos = append(infos, syncInfo)

                    sem.Signal()

                }(hit)
            }

            // FIXME: it's strange that we have to +1 here. Our semaphore seems to have a bug
            sem.Wait(out.Hits.Len() + 1)
            r.HTML(200, "syncList", infos)
        }
    })

    m.Run()
}
