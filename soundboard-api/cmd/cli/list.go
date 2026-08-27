package main

import (
	"context"
	"fmt"
)

func runList(args []string) error {
	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	soundbites, err := st.queries.ListSoundbites(context.Background())
	if err != nil {
		return fmt.Errorf("list soundbites: %w", err)
	}

	if len(soundbites) == 0 {
		fmt.Println("no soundbites stored yet")
		return nil
	}

	for _, sb := range soundbites {
		dateMade := "-"
		if sb.DateMade.Valid {
			dateMade = sb.DateMade.String
		}
		fmt.Printf("#%-4d %-30s %6.1fs  %6d plays  made:%-12s file:%s\n",
			sb.ID, sb.Name, sb.LengthSeconds, sb.PlayCount, dateMade, sb.Filename)
	}
	return nil
}
