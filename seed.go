/*
@Author: Alioth
@Date: 2024/9/4
@Description:
*/
package geoserver

import (
	"bytes"
	"encoding/json"
)

// SeedRequest represents the request payload for seeding tiles
type SeedRequest struct {
	Name        string `json:"name"`
	GridSetId   string `json:"gridSetId"`
	ZoomStart   int    `json:"zoomStart"`
	ZoomStop    int    `json:"zoomStop"`
	Type        string `json:"type"` // "reseed", "seed", "truncate"
	ThreadCount int    `json:"threadCount"`
}

// SeedRequestWrapper is the outer structure for the seed request
type SeedRequestWrapper struct {
	SeedRequest SeedRequest `json:"seedRequest"`
}

// SeedResponse represents the response from the seed API
type SeedResponse struct {
	// Define fields according to the response from the GeoServer API
}

// SeedTiles initiates a seed, reseed, or truncate operation for a specific layer
func (g *GeoServer) SeedTiles(workspaceName, layerName, gridSetId string, zoomStart, zoomStop int, seedType string, threadCount int) (success bool, err error) {
	targetURL := g.ParseURL("gwc", "rest", "seed", workspaceName+":"+layerName+".json")

	seedRequest := SeedRequest{
		Name:        workspaceName + ":" + layerName,
		GridSetId:   gridSetId,
		ZoomStart:   zoomStart,
		ZoomStop:    zoomStop,
		Type:        seedType, // "reseed", "seed", or "truncate"
		ThreadCount: threadCount,
	}

	// Wrap the seedRequest
	seedRequestWrapper := SeedRequestWrapper{SeedRequest: seedRequest}

	// Serialize the wrapped request to JSON
	serializedRequest, err := json.Marshal(seedRequestWrapper)
	if err != nil {
		return false, err
	}

	httpRequest := HTTPRequest{
		Method:   postMethod,
		Accept:   jsonType,
		Data:     bytes.NewBuffer(serializedRequest),
		DataType: jsonType,
		URL:      targetURL,
		Query:    nil,
	}
	response, responseCode := g.DoRequest(httpRequest)
	if responseCode != statusOk {
		g.logger.Error(string(response))
		success = false
		err = g.GetError(responseCode, response)
		return
	}
	success = true
	return
}
