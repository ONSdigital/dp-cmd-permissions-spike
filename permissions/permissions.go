package permissions

import (
	"net/http"
	"time"

	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

func RequireViewer(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	cli := http.Client{Timeout: 5 * time.Second}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if common.IsUserPresent(req.Context()) {
			ctx := req.Context()
			vars := mux.Vars(req)
			datasetID := vars["dataset_id"]
			collectionID := req.Header.Get("Collection-ID")
			florenceToken := req.Header.Get(common.FlorenceHeaderKey)

			logD := log.Data{
				"user_identity": common.Caller(ctx),
				"dataset_id":    datasetID,
				"collection_id": collectionID,
			}

			log.Event(ctx, "checking has viewer permissions", logD)

			authReq, _ := http.NewRequest("GET", "http://localhost:8082/canView", nil)
			authReq.Header.Set("X-Florence-Token", florenceToken)
			authReq.Header.Set("Dataset-ID", datasetID)
			authReq.Header.Set("Collection-ID", collectionID)

			resp, err := cli.Do(authReq)
			if err != nil {
				log.Event(ctx, "check permissions request error", log.Error(err), logD)
				w.WriteHeader(500)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				logD["permission_status"] = resp.StatusCode
				log.Event(ctx, "unexpected status code returned for permission check", logD)
				w.WriteHeader(401)
				return
			}
			handler(w, req)
		} else {
			w.WriteHeader(500)
			return
		}
	})
}
