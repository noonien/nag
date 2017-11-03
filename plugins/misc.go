package plugins

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/noonien/nag/bot"
	"github.com/sorcix/irc"
)

type Misc struct {
	bot *bot.Bot

	fuckers map[string]string
}

func (p *Misc) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b

	p.bot.HandleCmd("cmd.burp", p.burp())

	p.bot.HandleCmdRateLimited("cmd.digi", p.digi())

	p.textCmd("cmd.macanache", []string{"nu fi dilimache"})
	p.textCmd("cmd.satraiesti", []string{"satz traiasca familia boss", "sa traiesti boss"})
	p.textCmd("cmd.noroc", []string{"hai noroc"})
	p.textCmd("cmd.birzan", []string{"hurr durr sebyk suge"})
	p.textCmd("cmd.plp", []string{"fac laba"})
	p.textCmd("cmd.gabem", []string{"am suflet dar nu am suflet unde mi-e sufletu"})
	p.textCmd("cmd.florin", []string{"<?php echo 'haha fraerri chiar ma cred ca fac python';"})
	p.textCmd("cmd.jupi", []string{"hiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii _rs"})

	p.bot.HandleCmdRateLimited("cmd.ba", p.ba)
	p.bot.HandleCmdRateLimited("cmd.bă", p.ba)

	ai := "Altă întrebare! (https://soundcloud.com/armies/alta-intrebare-armies-edit)"
	p.textCmd("cmd.next", []string{ai})
	p.textReply("irc.privmsg", ai, func(line string) bool {
		line = strings.ToLower(line)
		return strings.HasSuffix(line, ", next") || strings.HasSuffix(line, " next!")
	})

	p.textReply("irc.privmsg", "mai zi", func(line string) bool {
		line = strings.ToLower(line)
		return strings.HasSuffix(line, "fascinant")
	})

	p.textReply("irc.privmsg", "*voiam", func(line string) bool {
		line = strings.ToLower(line)
		return strings.Contains(line, "vroiam")
	})

	p.bot.HandleIRC("irc.invite", p.invite)
	p.bot.HandleIRC("irc.kick", p.kick)
	p.bot.HandleIRC("irc.join", p.join)

	p.bot.HandleCmdRateLimited("cmd.bullshit", p.bullshit)
	p.bot.HandleCmdRateLimited("cmd.bs", p.bullshit)

	p.fuckers = make(map[string]string)

	return &bot.PluginInfo{
		Name:        "Misc",
		Author:      "noonien",
		Description: "Misc cmds",
		Version:     "1.0",
	}, nil
}

func (p *Misc) Unload() error {
	return nil
}

func (p *Misc) textCmd(cmd string, texts []string) {
	if len(texts) == 0 {
		return
	}

	handler := func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		text := texts[rand.Intn(len(texts))]

		if len(args) > 0 {
			text = args[0] + ": " + text
		}

		p.bot.Message(bot.PrivMsg(target, text))

		return true, nil
	}

	p.bot.HandleCmdRateLimited(cmd, handler)
}

func (p *Misc) textReply(cmd, text string, check func(string) bool) {
	handler := func(msg *irc.Message) (bool, error) {
		if !check(msg.Trailing) {
			return false, nil
		}

		if p.bot.RateLimiter.Limited(msg.Params[0]) {
			return false, nil
		}

		p.bot.Message(bot.PrivMsg(msg.Params[0], text))

		return false, nil
	}

	p.bot.HandleIRC(cmd, handler)
}

func (p *Misc) bullshit(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
	if len(args) == 0 {
		return true, nil
	}

	msg := fmt.Sprintf("%s: Dragnea s-a bullshit pe tine!", args[0])
	p.bot.Message(bot.PrivMsg(target, msg))
	return true, nil
}

func (p *Misc) digi() bot.CmdHandler {
	var mu sync.Mutex
	var shows []digiShow

	getShow := func() (digiShow, error) {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()

		for {
			for i, show := range shows {
				if show.when.Before(now) {
					continue
				}

				if i == 0 {
					break
				}

				return show, nil
			}

			var err error
			shows, err = getTodaysDigiShows()
			if err != nil {
				return digiShow{}, err
			}
		}
	}

	return func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		show, err := getShow()
		if err != nil {
			return false, err
		}

		p.bot.Message(bot.PrivMsg(target, "http://www.digi24.ro/live/digi24 | Emisiune: "+show.title))
		return true, nil
	}
}

type digiShow struct {
	when  time.Time
	title string
}

func getTodaysDigiShows() ([]digiShow, error) {
	today := time.Now().Truncate(24 * time.Hour)

	yesterdayShows, err := getDigiShows(today.Add(-24 * time.Hour))
	if err != nil {
		return nil, err
	}

	todayShows, err := getDigiShows(today)
	if err != nil {
		return nil, err
	}

	return append(yesterdayShows, todayShows...), nil
}

func getDigiShows(day time.Time) ([]digiShow, error) {
	day = day.Truncate(24 * time.Hour)
	today := day.Format("02/01/2006")

	url := fmt.Sprintf("http://www.rcs-rds.ro/asistenta/program-tv?cid=904&data=%d", day.Unix())
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	var shows []digiShow
	var isTomorrow bool
	doc.Find("table.vtable tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")

		whenText := strings.TrimSpace(tds.Eq(0).Text())
		var titleParts []string
		tds.Eq(1).Contents().Each(func(i int, s *goquery.Selection) {
			if goquery.NodeName(s) != "#text" {
				return
			}

			titleParts = append(titleParts, strings.TrimSpace(s.Text()))
		})

		title := strings.Join(titleParts, " ")
		title = strings.Replace(title, "Cu", "cu", 1)
		title = strings.Replace(title, " si", ",", -1)

		when, _ := time.Parse("02/01/2006 15:04", today+" "+whenText)

		if isTomorrow || (len(shows) > 0 && shows[len(shows)-1].when.After(when)) {
			isTomorrow = true
			when = when.Add(24 * time.Hour)
		}

		shows = append(shows, digiShow{when: when, title: title})
	})

	if err != nil {
		return nil, err
	}

	return shows, nil
}

func (p *Misc) ba(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
	if len(args) == 0 {
		return true, nil
	}

	perms, err := p.bot.Auth(source)
	if err != nil {
		return false, err
	}

	if perms == nil || !perms.Can("annoy") {
		return true, nil
	}

	lines := []string{
		"%s",
		"%s!",
		"ba %s",
		"%s %[1]s %[1]s %[1]s %[1]s",
		"%s %[1]s %[1]s %[1]s",
		"%s anplm",
	}

	times := rand.Intn(3) + 3
	for i := 0; i < times; i++ {
		line := lines[rand.Intn(len(lines))]
		msg := fmt.Sprintf(line, args[0])
		p.bot.Message(bot.PrivMsg(target, msg))

		time.Sleep(time.Duration(rand.Intn(300)+300) * time.Millisecond)
	}
	return true, nil
}

func (p *Misc) invite(msg *irc.Message) (bool, error) {
	perms, err := p.bot.Auth(msg.Prefix)
	if err != nil {
		return false, err
	}

	if perms == nil || !perms.Can("invite") {
		return true, nil
	}

	channel := msg.Trailing
	err = p.bot.Message(bot.Join(channel))
	return true, nil
}

func (p *Misc) kick(msg *irc.Message) (bool, error) {
	channel, who := msg.Params[0], msg.Params[1]
	if who != p.bot.Config.Nickname {
		return false, nil
	}

	fucker := msg.Prefix.Name

	if fucker == "ChanServ" {
		parts := strings.Fields(msg.Trailing)
		fucker = strings.Trim(parts[len(parts)-1], "()")
	}

	p.fuckers[channel] = fucker

	return false, nil
}

func (p *Misc) join(msg *irc.Message) (bool, error) {
	if msg.Prefix.Name != p.bot.Config.Nickname {
		return false, nil
	}

	channel := msg.Trailing

	fucker, ok := p.fuckers[channel]
	if !ok {
		return false, nil
	}

	delete(p.fuckers, fucker)

	welcome := fmt.Sprintf("%s: _)_", fucker)
	p.bot.Message(bot.PrivMsg(channel, welcome))
	return false, nil
}

func (p *Misc) burp() bot.CmdHandler {
	burps := map[string]bool{}

	return func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		_, ok := burps[source.Name]
		if ok {
			return true, nil
		}

		burps[source.Name] = true
		go func() {
			time.Sleep(time.Duration(180+rand.Intn(420)) * time.Second)
			delete(burps, source.Name)

			p.bot.Message(bot.PrivMsg("ChanServ", fmt.Sprintf("kick %s %s", target, source.Name)))
		}()

		return true, nil
	}
}
