package commands

func doSequential(p *processor.Sequential){
	for {
		err := p.NextBlock(context.Background())

		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "%+v", err)
		}
	}
}