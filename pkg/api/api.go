package api

import (
 "encoding/json"
 "go-news/pkg/storage"
 "net/http"

 "github.com/gorilla/mux"
)

type API struct {
 db     storage.Interface
 router *mux.Router
}

func New(db storage.Interface) *API {
 api := API{
  db: db,
 }
 api.router = mux.NewRouter()
 api.endpoints()
 return &api
}

func (api *API) endpoints() {
 api.router.HandleFunc("/posts", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
 api.router.HandleFunc("/posts", api.addPostHandler).Methods(http.MethodPost, http.MethodOptions)
 api.router.HandleFunc("/posts", api.updatePostHandler).Methods(http.MethodPut, http.MethodOptions)
 api.router.HandleFunc("/posts", api.deletePostHandler).Methods(http.MethodDelete, http.MethodOptions)
}

func (api *API) Router() *mux.Router {
 return api.router
}

func (api *API) postsHandler(w http.ResponseWriter, r *http.Request) {
 posts, err := api.db.Tasks()
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 bytes, err := json.Marshal(posts)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 _, err = w.Write(bytes)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
}

func (api *API) addPostHandler(w http.ResponseWriter, r *http.Request) {
 var p storage.Task
 err := json.NewDecoder(r.Body).Decode(&p)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 err = api.db.AddTask(p)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 w.WriteHeader(http.StatusOK)
}

func (api *API) updatePostHandler(w http.ResponseWriter, r *http.Request) {
 var p storage.Task
 err := json.NewDecoder(r.Body).Decode(&p)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 err = api.db.UpdateTask(p)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 w.WriteHeader(http.StatusOK)
}

func (api *API) deletePostHandler(w http.ResponseWriter, r *http.Request) {
 var p storage.Task
 err := json.NewDecoder(r.Body).Decode(&p)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 err = api.db.DeleteTask(p)
 if err != nil {
  http.Error(w, err.Error(), http.StatusInternalServerError)
  return
 }
 w.WriteHeader(http.StatusOK)
}

