package main

// Version mbstyle version
const Version = 8

const (
	//SourceTypeVector 矢量
	SourceTypeVector = "vector"
	//SourceTypeRaster 影像
	SourceTypeRaster = "raster"
	//SourceTypeGeoJSON geojson
	SourceTypeGeoJSON = "geojson"
	//SourceTypeImage 图像
	SourceTypeImage = "image"
	//SourceTypeVideo 视频
	SourceTypeVideo = "video"
	//SourceTypeCanvas canvas画布
	SourceTypeCanvas = "canvas"
)

type TsGeoType string

const (
	TsPoint   TsGeoType = "point"
	TsLine    TsGeoType = "line"
	TsPolygon TsGeoType = "polygon"
	TsUnknown TsGeoType = "unknown"
)

const (
	SchemeXYZ = "xyz"
	SchemeTMS = "tms"
)

// Source 数据源
type Source struct {
	Type string `json:"type"`
	// An array of one or more tile source URLs, as in the TileJSON spec.
	Tiles []string `json:"tiles,omitempty"`
	// defaults to 0 if not set
	MinZoom int `json:"minzoom,omitempty"`
	// defaults to 22 if not set
	MaxZoom int `json:"maxzoom,omitempty"`
	// url to TileJSON resource
	URL string `json:"url,omitempty"`
	//Optional enum. One of "xyz", "tms". Defaults to "xyz".
	Scheme string `json:"scheme,omitempty"`
}

// Light light
type Light struct {
	Anchor    string  `json:"anchor"`
	Color     string  `json:"color"`
	Intensity float64 `json:"intensity"`
}

// Transition 变换
type Transition struct{}

// TileJSON
//const Version = "2.1.0"

// https://github.com/mapbox/tilejson-spec
type TileJSON struct {
	// OPTIONAL. Default: null. Contains an attribution to be displayed
	// when the map is shown to a user. Implementations MAY decide to treat this
	// as HTML or literal text. For security reasons, make absolutely sure that
	// this field can't be abused as a vector for XSS or beacon tracking.
	Attribution *string `json:"attribution"`
	// OPTIONAL. Default: [-180, -90, 180, 90].
	// The maximum extent of available map tiles. Bounds MUST define an area
	// covered by all zoom levels. The bounds are represented in WGS:84
	// latitude and longitude values, in the order left, bottom, right, top.
	// Values may be integers or floating point numbers.
	Bounds [4]float64 `json:"bounds"`
	// OPTIONAL. Default: null.
	// The first value is the longitude, the second is latitude (both in
	// WGS:84 values), the third value is the zoom level as an integer.
	// Longitude and latitude MUST be within the specified bounds.
	// The zoom level MUST be between minzoom and maxzoom.
	// Implementations can use this value to set the default location. If the
	// value is null, implementations may use their own algorithm for
	// determining a default location.
	Center [3]float64 `json:"center"`
	// pbf - protocol buffer
	Format string `json:"format"`
	// OPTIONAL. Default: 0. >= 0, <= 22.
	// A positive integer specifying the minimum zoom level.
	MinZoom uint `json:"minzoom"`
	// OPTIONAL. Default: 22. >= 0, <= 22.
	// An positive integer specifying the maximum zoom level. MUST be >= minzoom.
	MaxZoom uint `json:"maxzoom"`
	// OPTIONAL. Default: null. A name describing the tileset. The name can
	// contain any legal character. Implementations SHOULD NOT interpret the
	// name as HTML.
	Name *string `json:"name"`
	// OPTIONAL. Default: null. A text description of the tileset. The
	// description can contain any legal character. Implementations SHOULD NOT
	// interpret the description as HTML.
	Description *string `json:"description"`
	// OPTIONAL. Default: "xyz". Either "xyz" or "tms". Influences the y
	// direction of the tile coordinates.
	// The global-mercator (aka Spherical Mercator) profile is assumed.
	Scheme string `json:"scheme"`
	// REQUIRED. A semver.org style version number. Describes the version of
	// the TileJSON spec that is implemented by this JSON object.
	TileJSON string `json:"tilejson"`
	// REQUIRED. An array of tile endpoints. {z}, {x} and {y}, if present,
	// are replaced with the corresponding integers. If multiple endpoints are specified, clients
	// may use any combination of endpoints. All endpoints MUST return the same
	// content for the same URL. The array MUST contain at least one endpoint.
	Tiles []string `json:"tiles"`
	// OPTIONAL. Default: []. An array of interactivity endpoints. {z}, {x}
	// and {y}, if present, are replaced with the corresponding integers. If multiple
	// endpoints are specified, clients may use any combination of endpoints.
	// All endpoints MUST return the same content for the same URL.
	// If the array doesn't contain any entries, interactivity is not supported
	// for this tileset.
	// See https://github.com/mapbox/utfgrid-spec/tree/master/1.2
	// for the interactivity specification.
	Grids []string `json:"grids,omitempty"`
	// OPTIONAL. Default: []. An array of data files in GeoJSON format.
	// {z}, {x} and {y}, if present,
	// are replaced with the corresponding integers. If multiple
	// endpoints are specified, clients may use any combination of endpoints.
	// All endpoints MUST return the same content for the same URL.
	// If the array doesn't contain any entries, then no data is present in
	// the map.
	Data []string `json:"data,omitempty"`
	// OPTIONAL. Default: "1.0.0". A semver.org style version number. When
	// changes across tiles are introduced, the minor version MUST change.
	// This may lead to cut off labels. Therefore, implementors can decide to
	// clean their cache when the minor version changes. Changes to the patch
	// level MUST only have changes to tiles that are contained within one tile.
	// When tiles change significantly, the major version MUST be increased.
	// Implementations MUST NOT use tiles with different major versions.
	Version string `json:"version"`
	// OPTIONAL. Default: null. Contains a mustache template to be used to
	// format data from grids for interaction.
	// See https://github.com/mapbox/utfgrid-spec/tree/master/1.2
	// for the interactivity specification.
	Template *string `json:"template"`
	// OPTIONAL. Default: null. Contains a legend to be displayed with the map.
	// Implementations MAY decide to treat this as HTML or literal text.
	// For security reasons, make absolutely sure that this field can't be
	// abused as a vector for XSS or beacon tracking.
	Legend *string `json:"legend"`
	// vector layer details. This is not part of the tileJSON spec
	// properties mimiced based on other vector provider implementations
	VectorLayers []VectorLayer `json:"vector_layers"`
	Tilestats    TileStats     `json:"tilestats,omitempty"`
}

// vector layers are not officially part of the tileJSON spec.
type VectorLayer struct {
	// REQUIRED. The name of the layer
	// "name" and "id" are identical
	ID string `json:"id"`
	// REQUIRED. The name of the layer
	// "name" and "id" are identical
	Name string `json:"name,omitempty"`
	// OPTIONAL. Default: []
	// an array of feature tags that MAY be included on each feature
	Description string            `json:"description"`
	Fields      map[string]string `json:"fields,omitempty"`
	// OPTIONAL. Default: null
	// possible values include: "point", "line", "polygon", "unknown"
	GeometryType TsGeoType `json:"geometry_type,omitempty"`
	// OPTIONAL. Default: 0. >= 0, <= 22.
	// A positive integer specifying the minimum zoom level.
	MinZoom uint `json:"minzoom"`
	// OPTIONAL. Default: 22. >= 0, <= 22.
	// A positive integer specifying the maximum zoom level. MUST be >= minzoom.
	MaxZoom uint `json:"maxzoom"`
}

// TileStats for stata
type TileStats struct {
	LayerCount int              `json:"layerCount"`
	Layers     []TileStataLayer `json:"layers"`
}
type TileStataLayer struct {
	Name           string      `json:"layer"`
	Count          int         `json:"count"`
	Geometry       GeoType     `json:"geometry"`
	AttributeCount int         `json:"attributeCount,omitempty"`
	Attributes     []Attribute `json:"attributes,omitempty"`
}
type Attribute struct {
	Attribute      string        `json:"attribute"`
	Count          int           `json:"count"`
	Type           string        `json:"geometry"`
	Max            float64       `json:"max"`
	Min            float64       `json:"min"`
	AttributeCount int           `json:"attributeCount,omitempty"`
	Values         []interface{} `json:"values"`
}
type GeoType string
