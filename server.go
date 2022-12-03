package main

import (
	"fmt"
	api "k8-api/api"
	apply "k8-api/apply"
	"k8-api/install"
	"net/http"
	"runtime"

	"github.com/distribution/distribution/v3/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sirupsen/logrus"
)

func main() {

	e := echo.New()

	log := logrus.New()

	//ctx:= context.Context
	log.SetReportCaller(true)
	log.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", f.File, f.Line)
		},
	}

	// Middleware"
	//e.Use(uuidMiddleware())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("uuid", "kube-ez-"+uuid.Generate().String()[:8])
			cc := c
			return next(cc)
		}
	})

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"level":"INFO","time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	api.Main()
	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// Root route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Yes! I am alive!\n")
	})

	//route to get the Pods info in a namespace
	e.GET("/pods", func(c echo.Context) error {
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get pods intitiated")
		namespace := c.QueryParam("namespace")
		containerDetails := c.QueryParam("containerDetails") == "True" || c.QueryParam("containerDetails") == "true"
		return c.String(http.StatusOK, api.Pods(namespace, containerDetails, l))
	})

	e.GET("/namespace", func(c echo.Context) error {
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Namespace intitiated")
		return c.String(http.StatusOK, api.NameSpace(l))
	})

	e.GET("/deployments", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Deployments intitiated")
		return c.String(http.StatusOK, api.Deployments(namespace, l))
	})

	e.GET("/configmaps", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Configmaps intitiated")
		return c.String(http.StatusOK, api.Configmaps(namespace, l))
	})

	e.GET("/services", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Services intitiated")
		return c.String(http.StatusOK, api.Services(namespace, l))
	})

	e.GET("/events", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Events intitiated")
		return c.String(http.StatusOK, api.Events(namespace, l))
	})

	e.GET("/secrets", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Secrets intitiated")
		return c.String(http.StatusOK, api.Secrets(namespace, l))
	})

	e.GET("/replicationController", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get RepilicationControllers intitiated")
		return c.String(http.StatusOK, api.ReplicationController(namespace, l))
	})

	e.GET("/daemonset", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Daemaonsets intitiated")
		return c.String(http.StatusOK, api.DaemonSet(namespace, l))
	})

	e.GET("/podLogs", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		pod := c.QueryParam("pod")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Pod's Logs intitiated")
		return c.String(http.StatusOK, api.PodLogs(namespace, pod, l))
	})

	e.GET("/helmRepoUpdate", func(c echo.Context) error {
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Get Helm Repo updates intitiated")
		return c.String(http.StatusOK, install.RepoUpdate(l))
	})

	e.POST("/helmRepoAdd", func(c echo.Context) error {
		url := c.QueryParam("url")
		repoName := c.QueryParam("repoName")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Adding Helm Repo intitiated")
		return c.String(http.StatusOK, install.RepoAdd(repoName, url, l))
	})

	e.POST("/helmInstall", func(c echo.Context) error {
		namespace := c.QueryParam("namespace")
		chartName := c.QueryParam("chartName")
		name := c.QueryParam("name")
		repo := c.QueryParam("repo")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Adding Helm Install intitiated")
		return c.String(http.StatusOK, install.InstallChart(namespace, chartName, name, repo, l))
	})

	e.POST("/createNamespace", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Creating Namespace intitiated")
		return c.String(http.StatusOK, api.CreateNamespace(namespace, l))
	})

	e.POST("/applyFile", func(c echo.Context) error {
		filepath := c.FormValue("filepath")
		l := log.WithFields(logrus.Fields{"uuid": c.Get("uuid")})
		l.Info("Intiating File appliying")
		return c.String(http.StatusOK, apply.Main(filepath, l))
	})

	e.DELETE("/deleteHelm", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		name := c.FormValue("name")
		return c.String(http.StatusOK, install.DeleteChart(name, namespace))
	})

	e.DELETE("/deleteNamespace", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		return c.String(http.StatusOK, api.DeleteNamespace(namespace))
	})

	e.DELETE("/deleteDeployment", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		deployment := c.FormValue("deployment")
		return c.String(http.StatusOK, api.DeleteDeployment(namespace, deployment))
	})

	e.DELETE("/deleteService", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		service := c.FormValue("service")
		return c.String(http.StatusOK, api.DeleteService(namespace, service))
	})

	e.DELETE("/deleteConfigMap", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		configMap := c.FormValue("configMap")
		return c.String(http.StatusOK, api.DeleteConfigMap(namespace, configMap))
	})

	e.DELETE("/deleteSecret", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		secret := c.FormValue("secret")
		return c.String(http.StatusOK, api.DeleteSecret(namespace, secret))
	})

	e.DELETE("/deleteReplicationController", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		replicationController := c.FormValue("replicationController")
		return c.String(http.StatusOK, api.DeleteReplicationController(namespace, replicationController))
	})

	e.DELETE("/deleteDaemonSet", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		daemonSet := c.FormValue("daemonSet")
		return c.String(http.StatusOK, api.DeleteDaemonSet(namespace, daemonSet))
	})

	e.DELETE("/deletePod", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		pod := c.FormValue("pod")
		return c.String(http.StatusOK, api.DeletePod(namespace, pod))
	})

	e.DELETE("/deleteEvent", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		event := c.FormValue("event")
		return c.String(http.StatusOK, api.DeleteEvent(namespace, event))
	})

	e.DELETE("/deleteAll", func(c echo.Context) error {
		namespace := c.FormValue("namespace")
		return c.String(http.StatusOK, api.DeleteAll(namespace))
	})

	// Run Server
	e.Logger.Fatal(e.Start(":8000"))
}
