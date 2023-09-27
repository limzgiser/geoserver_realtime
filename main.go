package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hishamkaram/geoserver"
)

var gsCatalog *geoserver.GeoServer

var workSpace = "golang"
var dataStore = "postgis_datastore"
var tbName = "building" // building,gd,china_polygon

var host = "http://localhost:8088/geoserver/"

func main() {

	gsCatalog = geoserver.GetCatalog(host, "admin", "geoserver")

	// publishPostgisLayer()

	// deleteLayer()

	// r := gin.Default()
	// r.GET("/tilejson", func(c *gin.Context) {
	// 	getTileJSON(c)
	// })
	// r.Run()

}



func getTileJSON(c *gin.Context) {
	layer, _ := gsCatalog.GetLayer(workSpace, tbName)
	gz, minz, maxz := 10, 0, 22
	gt := TsPoint
	attribution := "RealtimeTileEngine"
	minx, miny, maxx, maxy := 60.0, 0.0, 180.0, 60.0
	vtLyrs := []VectorLayer{}
	switch layer.Type {
	case "VECTOR":
		ft, _ := gsCatalog.GetFeatureType(workSpace, dataStore, tbName)

		minx = ft.NativeBoundingBox.Minx
		miny = ft.NativeBoundingBox.Miny
		maxx = ft.NativeBoundingBox.Maxx
		maxy = ft.NativeBoundingBox.Maxy

		fmt.Println(minx)
		fmt.Println(miny)
		fmt.Println(maxx)
		fmt.Println(maxy)

		fields := make(map[string]string)
		for _, av := range ft.Attributes.Attribute {
			// av.Binding
			if av.Name == "geometry" || av.Name == "geom" || av.Name == "the_geom" {
				splits := strings.Split(av.Binding, ".")
				geoType := splits[len(splits)-1]
				switch geoType {
				case "Point", "MultiPoint":
					gt = TsPoint
				case "LineString", "MultiLineString":
					gt = TsLine
				case "Polygon", "MultiPolygon":
					gt = TsPolygon
				default:
					gt = TsUnknown
				}
				continue
			}
			fields[av.Name] = GetGsFieldType(av.Binding[strings.LastIndex(av.Binding, ".")+1:])
		}
		// build our vector layer details
		vlayer := VectorLayer{
			ID:           layer.Name,
			Name:         layer.Name,
			MinZoom:      uint(minz),
			MaxZoom:      uint(maxz),
			GeometryType: gt,
			Fields:       fields,
		}
		// add our Layer to our tile layer response
		vtLyrs = append(vtLyrs, vlayer)
	}
	tileJSON := TileJSON{
		Version:      "2",
		TileJSON:     "2.1.0",
		Attribution:  &attribution,
		Name:         &layer.Name,
		Bounds:       [4]float64{minx, miny, maxx, maxy},
		Center:       [3]float64{(minx + maxx) / 2, (miny + maxy) / 2, float64(gz)},
		MinZoom:      uint(minz),
		MaxZoom:      uint(maxz),
		Format:       "pbf",
		Scheme:       "tms",
		Grids:        make([]string, 0),
		Data:         make([]string, 0),
		VectorLayers: vtLyrs,
	}
	tileUrl := fmt.Sprintf("%s/gwc/service/tms/1.0.0/%s:%s@EPSG:900913@pbf/{z}/{x}/{y}.pbf", host, workSpace, tbName)
	tileJSON.Tiles = append(tileJSON.Tiles, tileUrl)

	//cache control headers (no-cache )
	// c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	c.JSON(http.StatusOK, tileJSON)
}

func deleteLayer() {
	_, err := gsCatalog.DeleteLayer(workSpace, tbName, true)
	if err == nil {
		fmt.Println("删除成功!", workSpace+":"+tbName)
	}
}

func publishPostgisLayer() {
	_, e := gsCatalog.GetWorkspace(workSpace)
	if e != nil {
		_, ee := gsCatalog.CreateWorkspace(workSpace)
		if ee == nil {
			fmt.Println("创建工作空间成功!")
		}
	}

	conn := geoserver.DatastoreConnection{
		Name:   dataStore,
		Port:   5432,
		Type:   "postgis",
		Host:   "db",
		DBName: "db_gis",
		DBPass: "postgres",
		DBUser: "postgres",
	}
	_, serr := gsCatalog.GetDatastoreDetails(workSpace, dataStore)

	if serr != nil {
		_, err := gsCatalog.CreateDatastore(conn, workSpace)
		if err != nil {
			fmt.Println("创建链接失败")
			return
		}
		fmt.Println("创建链接成功！")
	}

	_, dbErr := gsCatalog.PublishPostgisLayer(workSpace, dataStore, tbName, tbName)
	if dbErr != nil {
		fmt.Println("发布图层失败!")
		return
	}
	fmt.Println("发布图层成功!")

}

// 创建workspace
func createWorkSpace(workspaceName string) {

	_, err := gsCatalog.GetWorkspace(workspaceName)

	if err == nil {
		return
	}
	created, err := gsCatalog.CreateWorkspace(workspaceName)
	if err != nil {
		fmt.Printf("\nError:%s\n", err)
	}
	fmt.Println(strconv.FormatBool(created))
}

// 获取图层列表
func getLayers() {
	layers, err := gsCatalog.GetLayers("")

	if err != nil {
		fmt.Printf("\nError:%s\n", err)
	}
	for _, lyr := range layers {
		fmt.Println(lyr)
	}
}

// 获取图层
func getLayer(workspaceName string, layerName string) {
	layer, err := gsCatalog.GetLayer(workspaceName, layerName)
	if err != nil {
		fmt.Printf("\nError:%s\n", err)
	} else {
		fmt.Printf("%+v\n", layer)
	}

}
