package main

import (
    "encoding/json"
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "github.com/mattbaird/elastigo/api"
    "github.com/mattbaird/elastigo/core"
    "os"
)

type SyncInfo struct {
    Index       int
    Target      string
    TaskId      int
    TaskStarted int
}

func main() {
    api.Domain = os.Getenv("ES_ENDPOINT")
    api.Port = os.Getenv("ES_HTTP_PORT")

    index := "production_v4"

    m := martini.Classic()
    m.Use(render.Renderer(render.Options{
        Layout: "layout",
    }))

    m.Get("/", func(r render.Render) {

        out, _ := core.SearchRequest(index, "info", nil, nil)

        if len(out.Hits.Hits) > 0 {

            infos := []SyncInfo{}

            for _, hit := range out.Hits.Hits {
                var syncInfo SyncInfo

                json.Unmarshal(*hit.Source, &syncInfo)
                syncInfo.Index = len(infos) + 1
                syncInfo.TaskStarted = syncInfo.TaskId
                infos = append(infos, syncInfo)
            }

            r.HTML(200, "syncList", infos)
        }
    })

    // m.Get("/", func() string {

    //     searchJson := `{
    //         "sort" : [
    //             { "reputation" : {"order" : "desc"}}
    //         ],
    //         query: {
    //           match: {
    //               "_ninya_sync_task_id": 1398088041702
    //           }
    //         }
    //     }`

    //     out, _ := core.SearchRequest(index, "user", nil, searchJson)

    //     if len(out.Hits.Hits) >= 1 {
    //         jsonStr, _ := out.Hits.Hits[0].Source.MarshalJSON()
    //         fmt.Println(string(jsonStr))
    //     }

    //     return "foo"
    // })

    m.Run()
}
