package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Dataman-Cloud/seckilling/seckill-proxy/autogen"
	"github.com/Dataman-Cloud/seckilling/seckill-proxy/types"
	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// WebProvider is a provider.Provider implementation that provides the UI.
// FIXME to be handled another way.
type WebProvider struct {
	Address           string
	CertFile, KeyFile string
}

var (
	templatesRenderer = render.New(render.Options{
		Directory: "nowhere",
	})
)

// Provide allows the provider to provide configurations to traefik
// using the given configuration channel.
func (provider *WebProvider) Provide(configurationChan chan<- types.ConfigMessage) error {
	systemRouter := mux.NewRouter()

	// health route
	systemRouter.Methods("GET").Path("/health").HandlerFunc(getHealthHandler)

	// API routes
	systemRouter.Methods("GET").Path("/api").HandlerFunc(getConfigHandler)
	systemRouter.Methods("GET").Path("/api/providers").HandlerFunc(getConfigHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}").HandlerFunc(getProviderHandler)
	systemRouter.Methods("PUT").Path("/api/providers/{provider}").HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		if vars["provider"] != "web" {
			response.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(response, "Only 'web' provider can be updated through the REST API")
			return
		}

		configuration := new(types.Configuration)
		body, _ := ioutil.ReadAll(request.Body)
		err := json.Unmarshal(body, configuration)
		if err == nil {
			configurationChan <- types.ConfigMessage{"web", configuration}
			getConfigHandler(response, request)
		} else {
			log.Errorf("Error parsing configuration %+v", err)
			http.Error(response, fmt.Sprintf("%+v", err), http.StatusBadRequest)
		}
	})
	systemRouter.Methods("GET").Path("/api/providers/{provider}/backends").HandlerFunc(getBackendsHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}/backends/{backend}").HandlerFunc(getBackendHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}/backends/{backend}/servers").HandlerFunc(getServersHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}/backends/{backend}/servers/{server}").HandlerFunc(getServerHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}/frontends").HandlerFunc(getFrontendsHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}/frontends/{frontend}").HandlerFunc(getFrontendHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}/frontends/{frontend}/routes").HandlerFunc(getRoutesHandler)
	systemRouter.Methods("GET").Path("/api/providers/{provider}/frontends/{frontend}/routes/{route}").HandlerFunc(getRouteHandler)

	// Expose dashboard
	systemRouter.Methods("GET").Path("/").HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/dashboard/", 302)
	})
	systemRouter.Methods("GET").PathPrefix("/dashboard/").Handler(http.StripPrefix("/dashboard/", http.FileServer(&assetfs.AssetFS{Asset: autogen.Asset, AssetDir: autogen.AssetDir, Prefix: "static"})))

	go func() {
		if len(provider.CertFile) > 0 && len(provider.KeyFile) > 0 {
			err := http.ListenAndServeTLS(provider.Address, provider.CertFile, provider.KeyFile, systemRouter)
			if err != nil {
				log.Fatal("Error creating server: ", err)
			}
		} else {
			err := http.ListenAndServe(provider.Address, systemRouter)
			if err != nil {
				log.Fatal("Error creating server: ", err)
			}
		}
	}()
	return nil
}

func getHealthHandler(response http.ResponseWriter, request *http.Request) {
	templatesRenderer.JSON(response, http.StatusOK, metrics.Data())
}

func getConfigHandler(response http.ResponseWriter, request *http.Request) {
	templatesRenderer.JSON(response, http.StatusOK, currentConfigurations)
}

func getProviderHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	if provider, ok := currentConfigurations[providerID]; ok {
		templatesRenderer.JSON(response, http.StatusOK, provider)
	} else {
		http.NotFound(response, request)
	}
}

func getBackendsHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	if provider, ok := currentConfigurations[providerID]; ok {
		templatesRenderer.JSON(response, http.StatusOK, provider.Backends)
	} else {
		http.NotFound(response, request)
	}
}

func getBackendHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	backendID := vars["backend"]
	if provider, ok := currentConfigurations[providerID]; ok {
		if backend, ok := provider.Backends[backendID]; ok {
			templatesRenderer.JSON(response, http.StatusOK, backend)
			return
		}
	}
	http.NotFound(response, request)
}

func getServersHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	backendID := vars["backend"]
	if provider, ok := currentConfigurations[providerID]; ok {
		if backend, ok := provider.Backends[backendID]; ok {
			templatesRenderer.JSON(response, http.StatusOK, backend.Servers)
			return
		}
	}
	http.NotFound(response, request)
}

func getServerHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	backendID := vars["backend"]
	serverID := vars["server"]
	if provider, ok := currentConfigurations[providerID]; ok {
		if backend, ok := provider.Backends[backendID]; ok {
			if server, ok := backend.Servers[serverID]; ok {
				templatesRenderer.JSON(response, http.StatusOK, server)
				return
			}
		}
	}
	http.NotFound(response, request)
}

func getFrontendsHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	if provider, ok := currentConfigurations[providerID]; ok {
		templatesRenderer.JSON(response, http.StatusOK, provider.Frontends)
	} else {
		http.NotFound(response, request)
	}
}

func getFrontendHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	frontendID := vars["frontend"]
	if provider, ok := currentConfigurations[providerID]; ok {
		if frontend, ok := provider.Frontends[frontendID]; ok {
			templatesRenderer.JSON(response, http.StatusOK, frontend)
			return
		}
	}
	http.NotFound(response, request)
}

func getRoutesHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	frontendID := vars["frontend"]
	if provider, ok := currentConfigurations[providerID]; ok {
		if frontend, ok := provider.Frontends[frontendID]; ok {
			templatesRenderer.JSON(response, http.StatusOK, frontend.Routes)
			return
		}
	}
	http.NotFound(response, request)
}

func getRouteHandler(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	providerID := vars["provider"]
	frontendID := vars["frontend"]
	routeID := vars["route"]
	if provider, ok := currentConfigurations[providerID]; ok {
		if frontend, ok := provider.Frontends[frontendID]; ok {
			if route, ok := frontend.Routes[routeID]; ok {
				templatesRenderer.JSON(response, http.StatusOK, route)
				return
			}
		}
	}
	http.NotFound(response, request)
}
