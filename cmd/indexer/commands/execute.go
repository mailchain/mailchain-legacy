package commands

// Execute runs the command
func Execute() error {
	root, err := rootCmd()
	if err != nil {
		return err
	}

	return root.Execute()
}