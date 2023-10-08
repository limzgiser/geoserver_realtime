package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hishamkaram/geoserver"
)

var gsCatalog *geoserver.GeoServer

func main() {

	gsCatalog = geoserver.GetCatalog(GeoHost, GeoUser, GeoPassword)

	_init()

	router := gin.Default()

	router.POST("/geoserver/layer", func(ctx *gin.Context) {
		layerName := ctx.PostForm("layerName")
		geoType := ctx.PostForm("type") // pg shp tiff

		if geoType == "pg" {
			publishPostgisLayer(ctx, layerName)
			return
		}

		if geoType == "shp" {
			uploadShpfile(ctx, layerName)
			return
		}

		if geoType == "tiff" {
			publishGeoTiffLayer(ctx, layerName)
			return
		}
		ctx.JSON(500, gin.H{
			"status": 500,
			"data":   "error",
		})
	})

	router.DELETE("/geoserver/layer/:layerName", func(ctx *gin.Context) {
		layerName := ctx.Param("layerName")
		deleteLayer(ctx, layerName)
	})

	router.GET("/geoserver/tilejson", func(ctx *gin.Context) {
		layerName := ctx.Query("layerName")
		getTileJSON(ctx, layerName)
	})

	router.Run()

}

func _init() {

	// 矢量workspace
	_, err1 := gsCatalog.GetWorkspace(GeoWorkSpace)
	if err1 != nil {
		_, err := gsCatalog.CreateWorkspace(GeoWorkSpace)
		if err == nil {
			fmt.Println("创建矢量工作空间成功!")
		}
	}
	// 栅格
	_, err2 := gsCatalog.GetWorkspace(CoverageWorkSpace)
	if err2 != nil {
		_, err := gsCatalog.CreateWorkspace(CoverageWorkSpace)
		if err == nil {
			fmt.Println("创建栅格工作空间成功!")
		}
	}

	_, d_err := gsCatalog.GetDatastoreDetails(GeoWorkSpace, DataStoreName)

	if d_err != nil {
		_, err := gsCatalog.CreateDatastore(DataStoreConnect, GeoWorkSpace)
		if err != nil {
			fmt.Println("创建链接失败")
			return
		}
		fmt.Println("创建链接成功！")
	}
}

func publishPostgisLayer(ctx *gin.Context, layerName string) {

	_, dbErr := gsCatalog.PublishPostgisLayer(GeoWorkSpace, DataStoreName, layerName, layerName)
	if dbErr != nil {
		fmt.Println("发布图层失败!")
		ctx.JSON(500, gin.H{
			"status": 500,
			"data":   "",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data":   layerName,
	})

}

func uploadShpfile(ctx *gin.Context, layerName string) {
	zippedShapefile := filepath.Join(getGoGeoserverPackageDir(), "/statics/shp/", layerName /**"shpp2.zip"*/)
	_, err := gsCatalog.UploadShapeFile(zippedShapefile, GeoWorkSpace, "")
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"data":   "",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data":   layerName,
	})
}

func publishGeoTiffLayer(ctx *gin.Context, layerName string) {

	//1.上传tiff文件到容器目录
	//2. 创建store
	//3. 发布服务
	url := "/data_dir/" + layerName + ".tif"
	coverageStore := geoserver.CoverageStore{

		Name: layerName,
		Type: "GeoTIFF",
		URL:  url,
		Workspace: &geoserver.Resource{
			Name: CoverageWorkSpace,
		},
		Enabled: true,
	}
	_, err := gsCatalog.CreateCoverageStore(CoverageWorkSpace, coverageStore)
	if err != nil && !strings.Contains(err.Error(), "exists") {

		fmt.Println("创建CoverageStore失败")
		return

	}
	_, err2 := gsCatalog.PublishGeoTiffLayer(CoverageWorkSpace, layerName, layerName, layerName)
	if err2 != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"data":   "",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data":   layerName,
	})

}

func deleteLayer(ctx *gin.Context, layerName string) {
	_, err := gsCatalog.DeleteLayer(GeoWorkSpace, layerName, true)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"data":   "",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data":   layerName,
	})
}

func getTileJSON(c *gin.Context, layerName string) {
	layer, _ := gsCatalog.GetLayer(GeoWorkSpace, layerName)
	gz, minz, maxz := 10, 0, 22
	gt := TsPoint
	attribution := "RealtimeTileEngine"
	minx, miny, maxx, maxy := 60.0, 0.0, 180.0, 60.0
	vtLyrs := []VectorLayer{}
	switch layer.Type {
	case "VECTOR":
		// 如果是shp数据元，这个DataStoreName 是动态变的，目前使用的是文件名
		ft, _ := gsCatalog.GetFeatureType(GeoWorkSpace, DataStoreName, layerName)

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
	tileUrl := fmt.Sprintf("%s/gwc/service/tms/1.0.0/%s:%s@EPSG:900913@pbf/{z}/{x}/{y}.pbf", GeoHost, GeoWorkSpace, layerName)
	tileJSON.Tiles = append(tileJSON.Tiles, tileUrl)

	//cache control headers (no-cache )
	// c.Header("Cache-Control", "no-cache, no-store, must-revalidate")

	c.JSON(http.StatusOK, tileJSON)
}

// // 创建workspace
// func createWorkSpace(workspaceName string) {

// 	_, err := gsCatalog.GetWorkspace(workspaceName)

// 	if err == nil {
// 		return
// 	}
// 	created, err := gsCatalog.CreateWorkspace(workspaceName)
// 	if err != nil {
// 		fmt.Printf("\nError:%s\n", err)
// 	}
// 	fmt.Println(strconv.FormatBool(created))
// }

// // 获取图层列表
// func getLayers() {
// 	layers, err := gsCatalog.GetLayers("")

// 	if err != nil {
// 		fmt.Printf("\nError:%s\n", err)
// 	}
// 	for _, lyr := range layers {
// 		fmt.Println(lyr)
// 	}
// }

// // 获取图层
// func getLayer(workspaceName string, layerName string) {
// 	layer, err := gsCatalog.GetLayer(workspaceName, layerName)
// 	if err != nil {
// 		fmt.Printf("\nError:%s\n", err)
// 	} else {
// 		fmt.Printf("%+v\n", layer)
// 	}
// }
