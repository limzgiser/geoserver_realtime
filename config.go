package main

import "github.com/hishamkaram/geoserver"

var (
	GeoHost           = "http://localhost:8088/geoserver/"
	GeoUser           = "dxy123"
	GeoPassword       = "dxy"
	GeoWorkSpace      = "DXY_Vector"
	CoverageWorkSpace = "DXY_Raster"
	DataStoreName     = "db_gis"
	DataStoreConnect  = geoserver.DatastoreConnection{
		Name:   DataStoreName,
		Port:   5432,
		Type:   "postgis",
		Host:   "db",
		DBName: "db_gis",
		DBPass: "postgres",
		DBUser: "postgres",
	}
)
