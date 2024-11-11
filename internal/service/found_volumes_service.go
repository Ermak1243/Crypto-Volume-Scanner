package service

import (
	"cvs/internal/models"
	"strconv"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// FoundVolumesService defines the interface for managing found volumes.
// This interface includes methods for updating or inserting found volume data and retrieving all found volumes for a user.
type FoundVolumesService interface {
	UpsertFoundVolume(userData models.UserPairs, foundVolume models.FoundVolume) // Method to update or insert found volume data
	GetAllFoundVolume(userID int) ([]models.FoundVolume, error)                  // Method to retrieve all found volumes for a user
	DeleteFoundVolume(userPairData models.UserPairs)                             // Method to delete found volume data
}

// foundVolumesService is a concrete implementation of FoundVolumesService.
// It holds a concurrent map to store found volumes data.
type foundVolumesService struct {
	//first key - userID
	// second key - pair + exchange + side
	foundVolumesData cmap.ConcurrentMap[string, cmap.ConcurrentMap[string, models.FoundVolume]]
}

// NewFoundVolumesService creates a new instance of foundVolumesService.
// It initializes the concurrent map for storing found volumes data.
func NewFoundVolumesService() FoundVolumesService {
	return &foundVolumesService{
		foundVolumesData: cmap.New[cmap.ConcurrentMap[string, models.FoundVolume]](),
	}
}

// UpsertFoundVolume inserts or updates a found volume for a user in the stored data.
//
// This method retrieves the cached found volumes data for a specific user ID and either inserts
// or updates the found volume identified by a unique key composed of the pair, exchange, and side attributes.
// If the price of the found volume is zero, it will remove the existing entry instead of updating it.
//
// Parameters:
//   - userPairData: A models.UserPairs struct containing information about the user and their trading pair.
//   - foundVolume: A models.FoundVolume struct representing the volume data to be inserted or updated.
//
// This method does not return any values and does not produce errors. If there is no existing data for
// the user ID, it creates a new entry. If the found volume's price is zero, it removes any existing
// entry associated with that unique key.
func (fvs *foundVolumesService) UpsertFoundVolume(userPairData models.UserPairs, foundVolume models.FoundVolume) {
	userID := strconv.Itoa(userPairData.UserID)                                        // Convert UserID to string for use as a key
	foundVolumeUniqueKey := foundVolume.Pair + foundVolume.Exchange + foundVolume.Side // Create a unique key for the found volume

	// Check if user data exists
	userFoundVolumesData, ok := fvs.foundVolumesData.Get(userID) // Retrieve cached data for the user ID
	if !ok {
		var foundVolumesMap = cmap.New[models.FoundVolume]() // Create a new concurrent map for found volumes

		foundVolumesMap.Set(foundVolumeUniqueKey, foundVolume) // Insert found volume data
		fvs.foundVolumesData.Set(userID, foundVolumesMap)      // Store the new map in foundVolumesData

		return // Exit after inserting new data
	}

	if foundVolume.Price != 0 {
		userFoundVolumesData.Set(foundVolumeUniqueKey, foundVolume) // Update existing volume data
	} else {
		userFoundVolumesData.Remove(foundVolumeUniqueKey) // Remove entry if price is zero
	}

	fvs.foundVolumesData.Set(userID, userFoundVolumesData) // Update stored data for the user
}

// DeleteFoundVolume removes a specified found volume for a user from the stored data.
//
// This method retrieves the cached found volumes data for a specific user ID and attempts to remove
// the found volume identified by a unique key composed of the pair and exchange attributes.
// If the user does not have any found volumes stored, no action is taken.
//
// Parameters:
//   - userPairData: A models.UserPairs struct containing information about the user and their trading pair.
//   - foundVolume: A models.FoundVolume struct representing the volume data to be deleted.
//
// This method does not return any values and does not produce errors. However, if the user ID is not found,
// it will simply exit without making any changes.
func (fvs *foundVolumesService) DeleteFoundVolume(userPairData models.UserPairs) {
	userID := strconv.Itoa(userPairData.UserID)            // Convert UserID to string for use as a key
	uniqueKey := userPairData.Pair + userPairData.Exchange // Create a unique key for the found volume
	asksUniqueKey := uniqueKey + "asks"                    // Unique key for asks
	bidsUniqueKey := uniqueKey + "bids"                    // Unique key for bids

	// Retrieve cached data for the user ID
	userFoundVolumesData, _ := fvs.foundVolumesData.Get(userID)

	// Remove both asks and bids using their unique keys
	userFoundVolumesData.Remove(asksUniqueKey)
	userFoundVolumesData.Remove(bidsUniqueKey)
}

// GetAllFoundVolume retrieves all found volumes for a given user ID.
//
// Parameters:
//   - userID: The ID of the user whose found volumes are to be retrieved.
//
// Returns:
//   - A slice of FoundVolume and an error if any occurs during retrieval.
func (fvs *foundVolumesService) GetAllFoundVolume(userID int) ([]models.FoundVolume, error) {
	var volumesToReturn []models.FoundVolume

	userFoundVolumes, ok := fvs.foundVolumesData.Get(strconv.Itoa(userID)) // Retrieve cached data for the user ID
	if !ok {
		err := errGettingFoundVolume // Custom error indicating failure to get found volume

		return volumesToReturn, err // Return empty slice and error if not found
	}

	for _, volume := range userFoundVolumes.Items() { // Iterate over all found volumes
		volumesToReturn = append(volumesToReturn, volume)
	}

	return volumesToReturn, nil // Return all found volumes retrieved
}
