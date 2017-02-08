package plugins

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/noonien/nag/bot"
	"github.com/sorcix/irc"
	"github.com/turnage/graw/reddit"
)

type RedditFilth struct {
	PreloadCount int
	Lurker       reddit.Lurker

	bot   *bot.Bot
	close chan bool
}

func (p *RedditFilth) Load(b *bot.Bot) (*bot.PluginInfo, error) {
	p.bot = b
	p.close = make(chan bool)

	if p.PreloadCount < 1 {
		p.PreloadCount = 10
	}

	for i := range mizerii {
		mizerii[i].register(p)
	}
	p.bot.HandleCmdRateLimited("cmd.porn", p.roulette)

	return &bot.PluginInfo{
		Name:        "RedditFilth",
		Author:      "noonien",
		Description: "Reddit RedditFilth",
		Version:     "1.0",
	}, nil
}

func (p *RedditFilth) Unload() error {
	close(p.close)
	return nil
}

func (p *RedditFilth) roulette(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
	mizerie := mizerii[rand.Intn(len(mizerii))]
	cmd = mizerie.Cmds[0]
	return p.bot.Event("cmd."+cmd, source, target, cmd, args)
}

var imageHosts = []string{"imgur.com", "gfycat.com", "giphy.com", "i.redditmedia.com", "i.redd.it", "media.tumblr.com"}

func checkIsImage(post *reddit.Post) bool {
	linkURL, err := url.Parse(post.URL)
	if err != nil {
		return false
	}

	for _, host := range imageHosts {
		if strings.Contains(linkURL.Host, host) {
			return true
		}
	}

	return false
}

func chooseRandStr(opt []string) string {
	return opt[rand.Intn(len(opt))]
}

type mizerie struct {
	Cmds          []string
	Subs          []string
	RedditListTag string
	Ce            []string
	Check         func(*reddit.Post) bool

	posts []*reddit.Post
	mu    sync.Mutex
	close chan bool
}

var mizerii = []mizerie{
	mizerie{
		Cmds: []string{"nsfw"},
		Subs: []string{
			"nsfw", "nsfwhardcore", "nsfw2", "HighResNSFW", "BonerMaterial",
			"porn", "iWantToFuckHer", "NSFW_nospam", "Sexy", "nude",
			"UnrealGirls", "primes", "THEGOLDSTANDARD", "nsfw_hd", "UHDnsfw",
			"BeautifulTitsAndAss", "FuckMarryOrKill", "NSFWCute",
			"badassgirls", "HotGirls", "PornPleasure", "nsfwnonporn",
			"NSFWcringe", "NSFW_PORN_ONLY", "Sex_Games", "BareGirls",
			"lusciousladies", "Babes", "FilthyGirls", "NaturalWomen",
			"ImgurNSFW", "Adultpics", "sexynsfw", "nsfw_sets", "OnlyGoodPorn",
			"TumblrArchives", "HardcoreSex", "PornLovers", "NSFWgaming",
			"Fapucational", "RealBeauties", "fappitt", "exotic_oasis", "TIFT",
			"nakedbabes", "oculusnsfw", "CrossEyedFap", "TitsAssandNoClass",
			"formylover", "Ass_and_Titties", "Ranked_Girls", "fapfactory",
			"NSFW_hardcore", "Sexyness", "debs_and_doxies", "nsfwonly",
			"pornpedia", "lineups", "Nightlysex", "spod", "nsfwnew",
			"pinupstyle", "NoBSNSFW", "awwyea", "nsfwdumps", "FoxyLadies",
			"nsfwcloseups", "NudeBeauty", "SimplyNaked", "fappygood",
			"FaptasticImages", "WhichOneWouldYouPick", "TumblrPorn",
			"SaturdayMorningGirls", "NSFWSector", "GirlsWithBigGuns",
			"QualityNsfw", "nsfwPhotoshopBattles", "hawtness",
			"fapb4momgetshome", "SeaSquared", "SexyButNotPorn", "WoahPoon",
			"Reflections", "Hotness", "Erotic_Galleries", "carnalclass",
			"nsfw_bw", "LaBeauteFeminine", "Sweet_Sexuality", "NSFWart",
			"WomenOfColorRisque",
		},
		Ce:    []string{"nsfw", "imagini de ascuns cand vine sefu"},
		Check: checkIsImage,
	},
	mizerie{
		Cmds: []string{"buci", "booci", "cur"},
		Subs: []string{
			"AssOnTheGlass", "BoltedOnBooty", "BubbleButts",
			"ButtsAndBareFeet", "Cheeking", "HighResASS",
			"LoveToWatchYouLeave", "NoTorso", "SpreadEm", "TheUnderbun",
			"Top_Tier_Asses", "Tushy", "Underbun", "ass", "assgifs",
			"bigasses", "booty", "booty_gifs", "datass", "datbuttfromthefront",
			"hugeass", "juicybooty", "pawg", "twerking", "whooties",
		},
		Ce:    []string{"buci", "booci", "cur", "curuletz", "funduletz"},
		Check: checkIsImage,
	},
	mizerie{
		Cmds: []string{"țț", "țâțe", "tzatze", "boobs"},
		Subs: []string{
			"BeforeAndAfterBoltons", "Bigtitssmalltits", "BoltedOnMaxed",
			"Boobies", "BreastEnvy", "EpicCleavage", "HardBoltOns",
			"JustOneBoob", "OneInOneOut", "PM_ME_YOUR_TITS_GIRL",
			"PerfectTits", "Perky", "Rush_Boobs", "Saggy", "SloMoBoobs",
			"TheHangingBoobs", "TheUnderboob", "Titsgalore", "TittyDrop",
			"bananatits", "boltedontits", "boobbounce", "boobgifs",
			"boobkarma", "boobland", "boobs", "breastplay", "breasts",
			"cleavage", "feelthemup", "handbra", "hanging", "hersheyskisstits",
			"homegrowntits", "knockers", "naturaltitties", "sideboob",
			"tits", "titsagainstglass", "torpedotits", "underboob",
		},
		Ce:    []string{"țț", "țâțe", "tzatze"},
		Check: checkIsImage,
	},
}

func (m *mizerie) get() *reddit.Post {
	for i := 0; i < 5; i++ {
		m.mu.Lock()
		var post *reddit.Post
		if len(m.posts) > 0 {
			post = m.posts[len(m.posts)-1]
			m.posts = m.posts[:len(m.posts)-1]
		}
		m.mu.Unlock()

		if post != nil {
			return post
		}

		select {
		case <-m.close:
			return nil
		case <-time.After(time.Second):
		}
	}

	return nil
}

func (m *mizerie) register(plug *RedditFilth) {
	m.posts = make([]*reddit.Post, 0, plug.PreloadCount)
	m.close = plug.close

	go func() {
		if len(m.RedditListTag) > 0 {
			m.getSubredditList()
		}

		if len(m.Subs) == 0 {
			return
		}

		m.preload(plug.Lurker)
	}()

	handler := func(source *irc.Prefix, target string, cmd string, args []string) (bool, error) {
		niste := chooseRandStr(m.Ce)

		post := m.get()
		if post == nil {
			plug.bot.Message(bot.PrivMsg(target, fmt.Sprintf("%s: n-am %s inca boss", source.Name, niste)))
			return true, nil
		}

		var msg string
		if len(args) > 0 {
			msg = fmt.Sprintf("%s, ia niste %s de la %s: %s NSFW (https://redd.it/%s)", args[0], niste, source.Name, post.URL, post.ID)
		} else {
			msg = fmt.Sprintf("%s, ia niste %s: %s NSFW (https://redd.it/%s)", source.Name, niste, post.URL, post.ID)
		}

		plug.bot.Message(bot.PrivMsg(target, msg))
		return true, nil

	}

	for _, cmd := range m.Cmds {
		plug.bot.HandleCmdRateLimited("cmd."+cmd, handler)
	}
}

func (m *mizerie) getSubredditList() {
	url := "http://redditlist.com/nsfw/category/" + m.RedditListTag
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Println("failed to get reddit list subreddits for ", m.RedditListTag)
		return
	}

	var subs []string
	doc.Find(".result-item-slug a").Each(func(i int, s *goquery.Selection) {
		sub := strings.TrimPrefix(s.Text(), "/r/")
		subs = append(subs, sub)
	})

	m.Subs = append(m.Subs, subs...)
}

func (m *mizerie) preload(lurk reddit.Lurker) {
	for {
		select {
		case <-m.close:
			return
		case <-time.After(2 * time.Second):
			m.mu.Lock()
			full := len(m.posts) == cap(m.posts)
			m.mu.Unlock()
			if full {
				continue
			}

			sub := m.Subs[rand.Intn(len(m.Subs))]
			for {
				post, err := lurk.Thread("/r/" + sub + "/random")
				if err != nil {
					log.Printf("error while getting random post from %s: %v\n", sub, err)
					sub = m.Subs[rand.Intn(len(m.Subs))]
					continue
				}

				if m.Check != nil && !m.Check(post) {
					continue
				}

				m.mu.Lock()
				m.posts = append(m.posts, post)
				m.mu.Unlock()
				break
			}

		}
	}
}
