package autocoins

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/LompeBoer/go-autocoins/internal/discord"
)

type Writer interface {
	WriteResult([]MarketSwing, *QuarantineMessages) error
	WriteError(string) error
}

type OutputWriter struct {
	Writers []Writer
}

// WriteResult outputs the calculated results from AutoCoins.
func (w *OutputWriter) WriteResult(data []SymbolDataObject, lists SymbolLists) error {
	marketSwings := CalculateMarketSwing(data)
	q := w.writeQuarantineMessage(lists)

	for _, wr := range w.Writers {
		wr.WriteResult(marketSwings, q)
	}
	return nil
}

// WriteError outputs the error message.
func (w *OutputWriter) WriteError(message string) error {
	for _, wr := range w.Writers {
		wr.WriteError(message)
	}
	return nil
}

type QuarantineMessages struct {
	NewQuarantined string
	Quarantined    string
	Unquarantined  string
	OpenPositions  string
	Failed         string
}

func (w *OutputWriter) writeQuarantineMessage(lists SymbolLists) *QuarantineMessages {
	return &QuarantineMessages{
		NewQuarantined: strings.Join(lists.QuarantinedNew, ", "),
		Quarantined:    strings.Join(lists.Quarantined, ", "),
		Unquarantined:  strings.Join(lists.QuarantinedRemoved, ", "),
		OpenPositions:  strings.Join(lists.QuarantinedSkipped, ", "),
		Failed:         strings.Join(lists.FailedToProcess, ", "),
	}
}

// ConsoleOutputWriter writes to the console.
type ConsoleOutputWriter struct {
}

func (w *ConsoleOutputWriter) WriteResult(marketSwings []MarketSwing, q *QuarantineMessages) error {
	b := strings.Builder{}
	d := strings.Builder{}
	for _, m := range marketSwings {
		fmt.Fprintf(&b, "MarketSwing - Last %s - %s\n", m.Timeframe, m.SwingMood)
		fmt.Fprintf(&b, "| %.0f%% Long | %d Coins | Avg %.2f%% | Max %.2f%% %s\n", m.Positive.Percent, m.Positive.CoinCount, m.Positive.Average, m.Positive.Max, m.Positive.MaxCoin)
		fmt.Fprintf(&b, "| %.0f%% Short | %d Coins | Avg %.2f%% | Max %.2f%% %s\n", m.Negative.Percent, m.Negative.CoinCount, m.Negative.Average, m.Negative.Max, m.Negative.MaxCoin)
	}

	if len(q.NewQuarantined) > 0 {
		fmt.Fprintf(&d, "NEW QUARANTINED: %s\n", q.NewQuarantined)
	}
	fmt.Fprintf(&d, "QUARANTINED: %s\n", q.Quarantined)
	if len(q.Unquarantined) > 0 {
		fmt.Fprintf(&d, "UNQUARANTINED: %s\n", q.Unquarantined)
	}
	if len(q.OpenPositions) > 0 {
		fmt.Fprintf(&d, "OPEN POSITIONS - NOT QUARANTINED: %s\n", q.OpenPositions)
	}
	if len(q.Failed) > 0 {
		fmt.Fprintf(&d, "FAILED TO PROCESS: %s\n", q.Failed)
	}

	fmt.Println(b.String())
	fmt.Println(d.String())

	// # $longvwap24 = [math]::Round((($settings.longVwapMax - $settings.longVwapMin) * ($negpercent24 / 100)) + $settings.longVwapMin, 1)
	// # $shortvwap24 = [math]::Round((($settings.shortVwapMax - $settings.shortVwapMin) * ($pospercent24 / 100)) + $settings.shortVwapMin, 1)
	// # $longvwap1 = [math]::Round((($settings.longVwapMax - $settings.longVwapMin) * ($negpercent1 / 100)) + $settings.longVwapMin, 1)
	// # $shortvwap1 = [math]::Round((($settings.shortVwapMax - $settings.shortVwapMin) * ($pospercent1 / 100)) + $settings.shortVwapMin, 1)

	// $message = "**MarketSwing - Last 1hr** - $swingmood1`n$pospercent1% Long | $poscoincount1 Coins | Ave $posave1% | Max $posmax1% $posmaxcoin1`n" + "$negpercent1% Short | $negcoincount1 Coins | Ave $negave1% | Max $negmax1% $negmaxcoin1 `n**MarketSwing - Last 4hrs** - $swingmood4`n$pospercent4% Long | $poscoincount4 Coins | Ave $posave4% | Max $posmax4% $posmaxcoin4`n" + "$negpercent4% Short | $negcoincount4 Coins | Ave $negave4% | Max $negmax4% $negmaxcoin4 `n**MarketSwing - Last 24hrs** - $swingmood24`n$pospercent24% Long | $poscoincount24 Coins | Ave $posave24% | Max $posmax24% $posmaxcoin24`n" + "$negpercent24% Short | $negcoincount24 Coins | Ave $negave24% | Max $negmax24% $negmaxcoin24"

	return nil
}

func (w *ConsoleOutputWriter) WriteError(message string) error {
	log.Println(message)
	return nil
}

// DiscordOutputWriter sends output to a Discord WebHook.
type DiscordOutputWriter struct {
	Version        string
	MentionOnError bool                   // TODO: Update when changing config file.
	WebHook        discord.DiscordWebHook // TODO: Update when changing config file.
}

func (w *DiscordOutputWriter) WriteResult(marketSwings []MarketSwing, q *QuarantineMessages) error {
	header := discord.DiscordEmbed{
		Title:       "AutoCoins MarketSwing report",
		Description: fmt.Sprintf("Generated %s (using v%s)", time.Now().Format("2006-01-02 15:04"), w.Version),
	}

	msg := discord.DiscordWebhookMessage{
		Embeds: []discord.DiscordEmbed{
			header,
		},
	}

	for _, market := range marketSwings {
		valueLong := fmt.Sprintf("%.0f%% Long | %d Coins | Avg %.2f%% | Max %.2f%% %s\n", market.Positive.Percent, market.Positive.CoinCount, market.Positive.Average, market.Positive.Max, market.Positive.MaxCoin)
		valueShort := fmt.Sprintf("%.0f%% Short | %d Coins | Avg %.2f%% | Max %.2f%% %s", market.Negative.Percent, market.Negative.CoinCount, market.Negative.Average, market.Negative.Max, market.Negative.MaxCoin)
		value := valueLong + "\n" + valueShort
		color := 3066993
		if market.Swing < 0 {
			color = 15158332
		}
		msg.Embeds = append(msg.Embeds, discord.DiscordEmbed{
			Color: color,
			Fields: []discord.DiscordEmbedField{
				{
					Name:   fmt.Sprintf("Last %s - %s\n", market.Timeframe, market.SwingMood),
					Value:  value,
					Inline: false,
				},
			},
		})
	}

	coins := discord.DiscordEmbed{}
	if len(q.NewQuarantined) > 0 {
		coins.Fields = append(coins.Fields, discord.DiscordEmbedField{
			Name: "New quarantined", Value: q.NewQuarantined, Inline: false,
		})
	}
	coins.Fields = append(coins.Fields, discord.DiscordEmbedField{
		Name: "Quarantined", Value: q.Quarantined, Inline: false,
	})
	if len(q.Unquarantined) > 0 {
		coins.Fields = append(coins.Fields, discord.DiscordEmbedField{
			Name: "Unquarantined", Value: q.Unquarantined, Inline: false,
		})
	}
	if len(q.OpenPositions) > 0 {
		coins.Fields = append(coins.Fields, discord.DiscordEmbedField{
			Name: "Open positions - not quarantined", Value: q.OpenPositions, Inline: false,
		})
	}
	if len(q.Failed) > 0 {
		coins.Fields = append(coins.Fields, discord.DiscordEmbedField{
			Name: "Failed to process", Value: q.Failed, Inline: false,
		})
	}

	msg.Embeds = append(msg.Embeds, coins)

	return w.WebHook.SendMessage(msg)
}

func (w *DiscordOutputWriter) WriteError(message string) error {
	return w.WebHook.SendError(message, w.MentionOnError)
}
