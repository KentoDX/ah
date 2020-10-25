// Package domain defines domain types which are used throughout the application.
// It imports only external packages which compose domain types.
// Use this file as a reference for Javascript objects, enums, etc.  They should
// all match as closely as possible.
// TODO: Might be nice to rename this package "earthwalker", but conflicts with
//       the default executable filename.
package domain

// == Shared Internal Structs ========

// Config holds server-wide settings
type Config struct {
	ConfigPath           string
	StaticPath           string
	DBPath               string
	Port                 string
	TileServerURL        string
	NoLabelTileServerURL string
}

// == Domain Enums ========

// PanoConnectedness is the enum representing that Map option
type PanoConnectedness int

const (
	// ConnectedAny is no restriction
	ConnectedAny PanoConnectedness = iota
	// ConnectedAlways is include only panos with at least one adjacent pano
	ConnectedAlways
	// ConnectedNever is include only panos with no adjacent panos
	ConnectedNever
)

func (c PanoConnectedness) String() string {
	return [...]string{"any", "always", "never"}[c]
}

// PanoCopyright is the enum representing that Map option
type PanoCopyright int

const (
	// CopyrightAny is no restriction
	CopyrightAny PanoCopyright = iota
	// CopyrightGoogle is include only panos with Google in the copyright
	CopyrightGoogle
	// CopyrightThirdParty is include only panos without Google in the copyright
	CopyrightThirdParty
)

func (c PanoCopyright) String() string {
	return [...]string{"any", "Google only", "third party only"}[c]
}

// PanoSource is the enum representing that Map option
type PanoSource int

const (
	// SourceAny specifies to query the streetview API for all panos
	SourceAny PanoSource = iota
	// SourceOutdoors specifies to query the API only for outdoors panos
	SourceOutdoors
)

// == Domain Types and Stores ========
// TODO: consider reducing stutter (Map.MapID, Challenge.ChallengeID, etc.)

// Map is a collection of user provided settings from which a Challenge
// is generated.  Formerly challenge/Settings
type Map struct {
	MapID string
	Name  string
	// TODO: consider adding description field
	Polygon       map[string]interface{} // geoJSON bounding the game area(s)
	Area          float32                // area bounded by Polygon, in square meters
	NumRounds     int
	TimeLimit     int // time limit per round, in seconds
	GraceDistance int // radius in meters within which max points are awarded
	MinDensity    int // minimum population density, 0-100
	MaxDensity    int // maximum population density, 0-100
	Connectedness PanoConnectedness
	Copyright     PanoCopyright
	Source        PanoSource
	ShowLabels    bool // whether to display place labels on the in-game minimap
	// TODO: consider adding CreatedAt (datetime) field
}

// MapStore is implemented by structs which provide access to a database
// containing Maps.
type MapStore interface {
	Insert(Map) error
	Get(mapID string) (Map, error)
}

// Challenge is a list of coordinates of panos.
type Challenge struct {
	ChallengeID string
	MapID       string
	Places      []ChallengePlace
}

// ChallengeStore is implemented by structs which provide access to a database
// containing Challenges.
type ChallengeStore interface {
	Insert(Challenge) error
	Get(challengeID string) (Challenge, error)
}

// ChallengePlace is the location of a pano.
// (May contain FOV, heading, etc. in the future.)
type ChallengePlace struct {
	ChallengeID string
	RoundNum    int
	Location    Coords
}

// ChallengeResult is a player's Guesses in a challenge.
// A ChallengeResult may still be in progress.
// Similar to former player/Session
type ChallengeResult struct {
	ChallengeResultID string
	ChallengeID       string

	// these could potentially be replaced with UserID if we were to
	// implement user auth/accounts
	Nickname string
	Icon     int

	Guesses []Guess
}

// ChallengeResultStore is implemented by structs which provide access to a
// database containing ChallengeResults.
type ChallengeResultStore interface {
	Insert(ChallengeResult) error
	Get(challengeResultID string) (ChallengeResult, error)
	GetAll(challengeID string) ([]ChallengeResult, error)
}

// Guess is a guessed location for one pano in a Challenge.
type Guess struct {
	ChallengeResultID string
	RoundNum          int
	Location          Coords
}

// Coords in degrees
type Coords struct {
	Lat float64
	Lng float64
}
