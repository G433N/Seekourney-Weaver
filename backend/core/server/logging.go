package server

import (
	"log"
	"seekourney/core/indexAPI"
)

// logDispatchErrors checks and reports result of an indexer dispatch attempt.
func logDispatchErrors(errs indexAPI.DispatchErrors) {
	if errs.IndexerWasRunning {
		log.Print("Indexer was running when Core made dispatch attempt")
	} else {
		log.Print(
			"Indexer needed to be started when Core made dispatch attempt")
	}
	if errs.StartupAttempt != nil {
		log.Print("Dispatch attempt failed to start up indexer " +
			"dispatch attempt aborted, failed with error: " +
			errs.StartupAttempt.Error())
	}
	if errs.DispatchAttempt != nil {
		log.Print(
			"Indexer was not able to handle indexing dispatch request" +
				", failed with error: " + errs.DispatchAttempt.Error())
	} else {
		log.Print("successfully dispatched indexing request to indexer")
	}
}
