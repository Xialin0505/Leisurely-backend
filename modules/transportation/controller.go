package transportation

import (
	_ "github.com/gofiber/fiber/v2"
	_ "leisurely/utils"
	"github.com/rs/zerolog/log"
	Models "leisurely/database/models"
	_ "fmt"
)

func DTOtoModelTag (durations float64, startE int, endE int) (Models.Transportation, error) {
	var trans Models.Transportation

	trans.Duration = durations
	trans.StartEvent = startE
	trans.DestinationEvent = endE

	return trans, nil
}

func InsertTransportation(planID int, durations []float64, events []int) error {
	if len(durations) == 0 {
		return nil
	}

	for i := 0; i < len(events) - 1; i++ {
		// Necessary to maintain synchronous database 
		currentTID, transErr := GetCurrentTransportationID()
	
		if transErr != nil {
			log.Error().Stack().Err(transErr).Msg("TransporationController: Database errors cannot get largest TransporationID")
		}

		oneTrans := Models.Transportation {
			TransportationID: currentTID,
			PlanID: planID,
			Duration: durations[i],
			StartEvent: events[i],
			DestinationEvent: events[i+1],
		}
		err := CreateTransportation(&oneTrans)
		if err != nil {
			log.Error().Stack().Err(err).Msg("TransportationController: Failed to insert transportation")
			return err
		}
	}

	return nil
}