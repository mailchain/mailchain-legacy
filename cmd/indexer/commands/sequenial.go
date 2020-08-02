package commands

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff"
	"github.com/mailchain/mailchain/cmd/indexer/internal/processor"
	"github.com/spf13/cobra"
)

func doSequential(cmd *cobra.Command, p *processor.Sequential, numRetry int) {
	failureCounter := 0
	for {
		operation := func() error {
			err := p.NextBlock(context.Background())
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "%+v\n", err)
			}

			return err
		}

		if err := backoff.Retry(operation, backoff.NewExponentialBackOff()); err != nil {
			failureCounter++
			fmt.Fprintf(cmd.ErrOrStderr(), "%+v\n", err)
		}

		if failureCounter == numRetry {
			fmt.Fprintf(cmd.ErrOrStderr(), "Number of retries has reached to %d. Exiting.\n", numRetry)
			return
		}
	}
}
