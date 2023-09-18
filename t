[1mdiff --git a/cmd/handlers.go b/cmd/handlers.go[m
[1mindex 07c207c..1cc7685 100644[m
[1m--- a/cmd/handlers.go[m
[1m+++ b/cmd/handlers.go[m
[36m@@ -67,18 +67,23 @@[m [mfunc (app *Config) Upload(w http.ResponseWriter, r *http.Request) {[m
 [m
 func (app *Config) Convert(w http.ResponseWriter, r *http.Request) {[m
 	log.Debug("Convert")[m
[31m-	// TODO: send to queue[m
[31m-[m
[31m-	// Proof of concept[m
[31m-	// err := ffmpeg.Input("/tmp/uuid.mp3").[m
[31m-	// 	Output("/tmp/uuid.ogg").[m
[31m-	// 	OverWriteOutput().ErrorToStdOut().Run()[m
[31m-[m
[31m-	// if err != nil {[m
[31m-	// 	log.Error(err)[m
[31m-	// 	w.WriteHeader(http.StatusInternalServerError)[m
[31m-	// 	return[m
[31m-	// }[m
[32m+[m	[32muser := r.Context().Value(userID("user")).(*model.User)[m
[32m+[m	[32merr := r.ParseForm()[m
[32m+[m
[32m+[m	[32mif err != nil {[m
[32m+[m		[32mapp.write(w, err.Error(), http.StatusInternalServerError)[m
[32m+[m		[32mreturn[m
[32m+[m	[32m}[m
[32m+[m
[32m+[m	[32mformat := r.PostForm.Get("format")[m
[32m+[m	[32mkbps := r.PostForm.Get("kbps")[m
[32m+[m	[32mmsg := user.UUID + format + kbps[m
[32m+[m	[32merr = app.QueueRepo.Push(msg)[m
[32m+[m	[32mif err != nil {[m
[32m+[m		[32mapp.write(w, err.Error(), http.StatusInternalServerError)[m
[32m+[m		[32mreturn[m
[32m+[m	[32m}[m
[32m+[m	[32mlog.Debug("Sent message: " + msg)[m
 [m
 	http.Redirect(w, r, "/", http.StatusSeeOther)[m
 }[m
[1mdiff --git a/cmd/main.go b/cmd/main.go[m
[1mindex c048af7..47da70a 100644[m
[1m--- a/cmd/main.go[m
[1m+++ b/cmd/main.go[m
[36m@@ -6,6 +6,7 @@[m [mimport ([m
 	"github.com/gorilla/sessions"[m
 	"github.com/nu12/audio-gonverter/internal/database"[m
 	"github.com/nu12/audio-gonverter/internal/logging"[m
[32m+[m	[32m"github.com/nu12/audio-gonverter/internal/queue"[m
 	"github.com/nu12/audio-gonverter/internal/repository"[m
 )[m
 [m
[36m@@ -19,6 +20,7 @@[m [mtype Config struct {[m
 	StaticFilesPath string[m
 	SessionStore    *sessions.CookieStore[m
 	DatabaseRepo    repository.DatabaseRepository[m
[32m+[m	[32mQueueRepo       repository.QueueRepository[m
 	Env             map[string]string[m
 }[m
 [m
[36m@@ -45,6 +47,7 @@[m [mfunc main() {[m
 	}[m
 [m
 	app.DatabaseRepo = database.NewRedis(app.Env["REDIS_HOST"], app.Env["REDIS_PORT"], "")[m
[32m+[m	[32mapp.QueueRepo = &queue.QueueMock{}[m
 	c := make(chan error, 1)[m
 [m
 	if app.Env["WEB_ENABLED"] == "true" {[m
