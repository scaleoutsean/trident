// Copyright 2022 NetApp, Inc. All Rights Reserved.

package rest

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/netapp/trident/config"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	Middleware  []mux.MiddlewareFunc
	HandlerFunc http.HandlerFunc
}

type Routes []Route

const (
	// arbitrarily large number to limit maximum routines waiting for global lock
	updateVolumePublicationRateLimit = 10000.0 // requests per second
	updateVolumePublicationBurst     = 10000   // maximum request burst
)

var controllerRoutes = Routes{
	Route{
		"GetVersion",
		"GET",
		strings.Replace(config.VersionURL, "/v"+config.OrchestratorAPIVersion, "", 1),
		nil,
		GetVersion,
	},
	Route{
		"GetVersion",
		"GET",
		config.VersionURL,
		nil,
		GetVersion,
	},
	Route{
		"AddBackend",
		"POST",
		config.BackendURL,
		nil,
		AddBackend,
	},
	Route{
		"UpdateBackend",
		"POST",
		config.BackendURL + "/{backend}",
		nil,
		UpdateBackend,
	},
	Route{
		"UpdateBackendState",
		"POST",
		config.BackendURL + "/{backend}" + "/state",
		nil,
		UpdateBackendState,
	},
	Route{
		"GetBackend",
		"GET",
		config.BackendURL + "/{backend}",
		nil,
		GetBackend,
	},
	Route{
		"ListBackends",
		"GET",
		config.BackendURL,
		nil,
		ListBackends,
	},
	Route{
		"DeleteBackend",
		"DELETE",
		config.BackendURL + "/{backend}",
		nil,
		DeleteBackend,
	},
	Route{
		"AddVolume",
		"POST",
		config.VolumeURL,
		nil,
		AddVolume,
	},
	Route{
		"GetVolume",
		"GET",
		config.VolumeURL + "/{volume}",
		nil,
		GetVolume,
	},
	Route{
		"ListVolumes",
		"GET",
		config.VolumeURL,
		nil,
		ListVolumes,
	},
	Route{
		"DeleteVolume",
		"DELETE",
		config.VolumeURL + "/{volume}",
		nil,
		DeleteVolume,
	},
	Route{
		"UpdateVolume",
		"PUT",
		config.VolumeURL + "/{volume}/luksPassphraseNames",
		nil,
		UpdateVolumeLUKSPassphraseNames,
	},
	Route{
		"ImportVolume",
		"POST",
		config.VolumeURL + "/import",
		nil,
		ImportVolume,
	},
	Route{
		"UpgradeVolume",
		"POST",
		config.VolumeURL + "/{volume}/upgrade",
		nil,
		UpgradeVolume,
	},
	Route{
		"AddStorageClass",
		"POST",
		config.StorageClassURL,
		nil,
		AddStorageClass,
	},
	Route{
		"GetStorageClass",
		"GET",
		config.StorageClassURL + "/{storageClass}",
		nil,
		GetStorageClass,
	},
	Route{
		"ListStorageClasses",
		"GET",
		config.StorageClassURL,
		nil,
		ListStorageClasses,
	},
	Route{
		"DeleteStorageClass",
		"DELETE",
		config.StorageClassURL + "/{storageClass}",
		nil,
		DeleteStorageClass,
	},
	Route{
		"AddOrUpdateNode",
		"PUT",
		config.NodeURL + "/{node}",
		nil,
		AddNode,
	},
	Route{
		"GetNode",
		"GET",
		config.NodeURL + "/{node}",
		nil,
		GetNode,
	},
	Route{
		"ListNodes",
		"GET",
		config.NodeURL,
		nil,
		ListNodes,
	},
	Route{
		"DeleteNode",
		"DELETE",
		config.NodeURL + "/{node}",
		nil,
		DeleteNode,
	},
	Route{
		"GetVolumePublication",
		"GET",
		config.PublicationURL + "/{volume}/{node}",
		nil,
		GetVolumePublication,
	},
	Route{
		"ListVolumePublications",
		"GET",
		config.PublicationURL,
		nil,
		ListVolumePublications,
	},
	Route{
		"ListVolumePublicationsForVolume",
		"GET",
		config.VolumeURL + "/{volume}/publication",
		nil,
		ListVolumePublicationsForVolume,
	},
	Route{
		"ListVolumePublicationsForNode",
		"GET",
		config.NodeURL + "/{node}/publication",
		nil,
		ListVolumePublicationsForNode,
	},
	Route{
		"UpdateVolumePublication",
		"PUT",
		config.PublicationURL + "/{volume}/{node}",
		[]mux.MiddlewareFunc{
			rateLimiterMiddleware(updateVolumePublicationRateLimit, updateVolumePublicationBurst),
		},
		UpdateVolumePublication,
	},
	Route{
		"ListSnapshots",
		"GET",
		config.SnapshotURL,
		nil,
		ListSnapshots,
	},
	Route{
		"ListSnapshotsForVolume",
		"GET",
		config.VolumeURL + "/{volume}/snapshot",
		nil,
		ListSnapshotsForVolume,
	},
	Route{
		"GetSnapshot",
		"GET",
		config.SnapshotURL + "/{volume}/{snapshot}",
		nil,
		GetSnapshot,
	},
	Route{
		"AddSnapshot",
		"POST",
		config.SnapshotURL,
		nil,
		AddSnapshot,
	},
	Route{
		"DeleteSnapshot",
		"DELETE",
		config.SnapshotURL + "/{volume}/{snapshot}",
		nil,
		DeleteSnapshot,
	},
	Route{
		"GetCHAP",
		"GET",
		config.ChapURL + "/{volume}/{node}",
		nil,
		GetCHAP,
	},
}
