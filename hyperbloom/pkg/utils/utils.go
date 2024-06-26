package utils

import (
	"fmt"
	"gopds/hyperbloom/internal/database/postgres"
	"gopds/hyperbloom/internal/service"
	"os"
	"sync"
)

// Cleanup handles OS interrupt signals to perform graceful shutdown tasks.
// It waits for a signal on osChan, shuts down the hyperbloom update coroutine,
// closes the PostgreSQL database connection, and then exits the program.
func Cleanup(osChan chan os.Signal, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this goroutine as done when function exits

	// Wait for an OS interrupt signal
	sig := <-osChan

	// Print the received signal
	fmt.Println("Encountered signal:", sig.String())

	// Perform shutdown tasks
	fmt.Println("Shutting down hyperbloom update coroutine and closing DB conn")

	// Send signal to stop async updates
	close(service.StopAsyncBloomUpdate)

	// Close the PostgreSQL database connection
	postgres.DbClient.Close()

	// Close osChan to signal completion of cleanup
	close(osChan)

	// Print final cleanup message
	fmt.Println("Cleaned up, exiting the program")

	// Exit the program with status code 0
	os.Exit(0)
}
