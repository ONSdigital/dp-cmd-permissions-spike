# dp-cmd-permissions-spike

### Example
- Register the `identity.Handler` middleware to populate the context with a user or service account identifier .
- Wrap the route `handlerFunc` in the `permission.RequireViewer` or `permission.RequireEditor` as required.
- The permission wrapper will send the identity in the context to zebedee to determined if the request should proceed

Requires:
- A service auth token

Or

- A Florence auth token
- A collection ID (if user is not an editor or admin)
- A dataset ID (if user is not an editor or admin)

```go
func main() {
    ...
    
    router := mux.NewRouter()
    
    healthcheckHandler := healthcheck.NewMiddleware(healthcheck.Do)
    middleware := alice.New(healthcheckHandler)
    
    if cfg.EnablePrivateEnpoints {
        middleware = middleware.Append(identity.Handler(cfg.ZebedeeURL))
    }

    router.HandleFunc("/datasets/{dataset_id}", permissions.RequireViewer(GetDataset())).Methods("GET")
    router.HandleFunc("/datasets/{dataset_id}", permissions.RequireEditor(CreateDataset())).Methods("POST")
    router.HandleFunc("/datasets/{dataset_id}", permissions.RequireEditor(UpdateDataset())).Methods("PUT")
    router.HandleFunc("/datasets/{dataset_id}", permissions.RequireEditor(DeleteDataset())).Methods("DELETE")
    ...
}
```
