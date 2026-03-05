package cli

import (
	"fmt"
	"strconv"

	"github.com/icl00ud/goban/pkg/cli/format"
	"github.com/spf13/cobra"
)

var cardsCmd = &cobra.Command{
	Use:   "cards",
	Short: "Manage cards",
}

var (
	cardColumnID uint64
	cardBoardID  uint64
	cardTitle    string
	cardDesc     string
	cardPriority string
	cardToColumn uint64
	cardPosition int
)

var cardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List cards in a board",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		if cardBoardID == 0 {
			return fmt.Errorf("--board is required")
		}
		cards, err := newClient().ListCards(uint(cardBoardID))
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(cards)
			return nil
		}
		rows := make([][]string, len(cards))
		for i, c := range cards {
			rows[i] = []string{
				fmt.Sprint(c.ID),
				c.Title,
				c.Priority,
				fmt.Sprint(c.ColumnID),
			}
		}
		format.Table([]string{"ID", "Title", "Priority", "Column"}, rows)
		return nil
	},
}

var cardsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get card details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		id, _ := strconv.ParseUint(args[0], 10, 64)
		card, err := newClient().GetCard(uint(id))
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(card)
			return nil
		}
		format.Table(
			[]string{"ID", "Title", "Priority", "Column", "Description"},
			[][]string{{
				fmt.Sprint(card.ID),
				card.Title,
				card.Priority,
				fmt.Sprint(card.ColumnID),
				card.Description,
			}},
		)
		return nil
	},
}

var cardsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a card",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		if cardColumnID == 0 {
			return fmt.Errorf("--column is required")
		}
		if cardTitle == "" {
			return fmt.Errorf("--title is required")
		}
		card, err := newClient().CreateCard(uint(cardColumnID), cardTitle, cardDesc, cardPriority)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(card)
			return nil
		}
		format.Msg("Created card #%d: %s [%s]", card.ID, card.Title, card.Priority)
		return nil
	},
}

var cardsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		id, _ := strconv.ParseUint(args[0], 10, 64)
		card, err := newClient().UpdateCard(uint(id), cardTitle, cardDesc, cardPriority)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(card)
			return nil
		}
		format.Msg("Updated card #%d: %s", card.ID, card.Title)
		return nil
	},
}

var cardsMoveCmd = &cobra.Command{
	Use:   "move <id>",
	Short: "Move a card to another column",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		if cardToColumn == 0 {
			return fmt.Errorf("--to-column is required")
		}
		id, _ := strconv.ParseUint(args[0], 10, 64)
		card, err := newClient().MoveCard(uint(id), uint(cardToColumn), cardPosition)
		if err != nil {
			return err
		}
		if flagJSON {
			format.JSON(card)
			return nil
		}
		format.Msg("Moved card #%d to column #%d", card.ID, card.ColumnID)
		return nil
	},
}

var cardsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a card",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireToken(cmd); err != nil {
			return err
		}
		id, _ := strconv.ParseUint(args[0], 10, 64)
		if err := newClient().DeleteCard(uint(id)); err != nil {
			return err
		}
		format.Msg("Deleted card #%d", id)
		return nil
	},
}

func init() {
	cardsListCmd.Flags().Uint64VarP(&cardBoardID, "board", "b", 0, "Board ID (required)")

	cardsCreateCmd.Flags().Uint64VarP(&cardColumnID, "column", "c", 0, "Column ID (required)")
	cardsCreateCmd.Flags().StringVar(&cardTitle, "title", "", "Card title (required)")
	cardsCreateCmd.Flags().StringVar(&cardDesc, "description", "", "Card description")
	cardsCreateCmd.Flags().StringVar(&cardPriority, "priority", "medium", "Priority: low|medium|high")

	cardsUpdateCmd.Flags().StringVar(&cardTitle, "title", "", "New title")
	cardsUpdateCmd.Flags().StringVar(&cardDesc, "description", "", "New description")
	cardsUpdateCmd.Flags().StringVar(&cardPriority, "priority", "", "New priority")

	cardsMoveCmd.Flags().Uint64VarP(&cardToColumn, "to-column", "t", 0, "Target column ID (required)")
	cardsMoveCmd.Flags().IntVar(&cardPosition, "position", 0, "Position in target column")

	cardsCmd.AddCommand(cardsListCmd, cardsGetCmd, cardsCreateCmd, cardsUpdateCmd, cardsMoveCmd, cardsDeleteCmd)
}
